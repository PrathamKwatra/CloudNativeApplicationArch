package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"drexel.edu/votes/schema"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/nitishm/go-rejson/v4"
)

const (
	RedisKeyPrefix = "votes:"
)

type cache struct {
	client  *redis.Client
	helper  *rejson.Handler
	context context.Context
}

type Health struct {
	startTime               time.Time
	totalApiCallsWithErrors int
	totalValidApiCalls      int
	mu                      sync.Mutex
}

type API struct {
	Polls  string
	Voters string
	Self   string
}

type VotesAPI struct {
	cache
	health      Health
	apiClient   *resty.Client
	API         API
	InternalAPI API
}

type Embedded struct {
	Voter schema.Voter
	Poll  schema.Poll
}

func (v *VotesAPI) validCall() {
	v.health.mu.Lock()
	v.health.totalValidApiCalls++
	v.health.mu.Unlock()
}

func (v *VotesAPI) invalidCall() {
	v.health.mu.Lock()
	v.health.totalApiCallsWithErrors++
	v.health.mu.Unlock()
}

func New(location string, api API, internalAPI API) (*VotesAPI, error) {

	apiClient := resty.New()
	//Connect to redis.  Other options can be provided, but the
	//defaults are OK
	client := redis.NewClient(&redis.Options{
		Addr: location,
	})

	//We use this context to coordinate betwen our go code and
	//the redis operaitons
	ctx := context.Background()

	//This is the reccomended way to ensure that our redis connection
	//is working
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error())
		return nil, err
	}

	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	//Return a pointer to a new ToDo struct
	return &VotesAPI{
		cache: cache{
			client:  client,
			helper:  jsonHelper,
			context: ctx,
		},
		API: api,
		health: Health{
			startTime:               time.Now(),
			totalApiCallsWithErrors: 0,
			totalValidApiCalls:      0,
		},
		InternalAPI: internalAPI,
		apiClient:   apiClient,
	}, nil
}

func (v *VotesAPI) CrashSim(c *gin.Context) {
	panic("Simulating an unexpected crash")
}

func (v *VotesAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"api":                  "Votes API",
		"status":               "ok",
		"uptime":               time.Since(v.health.startTime).String(),
		"msg":                  "Currently healthy",
		"totalValidCalls":      v.health.totalValidApiCalls,
		"totalCallsWithErrors": v.health.totalApiCallsWithErrors,
		"totalCalls":           v.health.totalApiCallsWithErrors + v.health.totalValidApiCalls,
	})
}

func (v *VotesAPI) GetVote(c *gin.Context) {
	id := c.Param("voteId")
	if id == "" {
		v.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No vote ID provided"})
		return
	}

	cacheKey := "votes:" + id
	rawVotes, err := v.helper.JSONGet(cacheKey, ".")
	if err != nil {
		v.invalidCall()
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find vote in cache with id=" + cacheKey})
		return
	}

	var vote schema.Vote
	err = json.Unmarshal(rawVotes.([]byte), &vote)
	if err != nil {
		v.invalidCall()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cached data seems to be wrong type"})
		return
	}

	//generate the latest HAL JSON response
	err = v.generateHALJSONResponse(&vote)
	if err != nil {
		v.invalidCall()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate HAL JSON response"})
		return
	}

	v.validCall()
	c.JSON(http.StatusOK, vote)
}

func (v *VotesAPI) GetVotes(c *gin.Context) {

	var votes []schema.Vote
	var vote schema.Vote

	//Lets query redis for all of the items
	pattern := "votes:*"
	ks, _ := v.client.Keys(v.context, pattern).Result()
	for _, key := range ks {
		err := v.getItemFromRedis(key, &vote)
		if err != nil {
			v.invalidCall()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find vote in cache with id=" + key})
			return
		}
		//generate the latest HAL JSON response
		err = v.generateHALJSONResponse(&vote)
		if err != nil {
			v.invalidCall()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate HAL JSON response\n" + err.Error()})
			return
		}
		votes = append(votes, vote)
	}

	v.validCall()
	c.JSON(http.StatusOK, votes)
}

func (v *VotesAPI) GetVotesByPolls(c *gin.Context) {
	id := c.Param("pollId")
	if id == "" {
		v.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No poll ID provided"})
		return
	}

	pollId, err := strconv.Atoi(id)
	if err != nil {
		v.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid poll id",
		})
		return
	}

	var votes []schema.Vote
	var vote schema.Vote
	// get all votes
	pattern := "votes:*"
	ks, _ := v.client.Keys(v.context, pattern).Result()
	for _, key := range ks {
		err := v.getItemFromRedis(key, &vote)
		if err != nil {
			v.invalidCall()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find vote in cache with id=" + key})
			return
		}
		//generate the latest HAL JSON response
		err = v.generateHALJSONResponse(&vote)
		if err != nil {
			v.invalidCall()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate HAL JSON response\n" + err.Error()})
			return
		}
		votes = append(votes, vote)
	}

	// find all votes with pollId
	var votesWithPollId []schema.Vote
	for _, vote := range votes {
		if vote.PollId == pollId {
			votesWithPollId = append(votesWithPollId, vote)
		}
	}

	v.validCall()
	c.JSON(http.StatusOK, votesWithPollId)
}

func (v *VotesAPI) GetVotesByVoter(c *gin.Context) {
	id := c.Param("voterId")
	if id == "" {
		v.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No voter ID provided"})
		return
	}

	voterId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid voter id",
		})
		v.invalidCall()
		return
	}

	var votes []schema.Vote
	var vote schema.Vote
	// get all votes
	pattern := "votes:*"
	ks, _ := v.client.Keys(v.context, pattern).Result()
	for _, key := range ks {
		err := v.getItemFromRedis(key, &vote)
		if err != nil {
			v.invalidCall()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find vote in cache with id=" + key})
			return
		}
		//generate the latest HAL JSON response
		err = v.generateHALJSONResponse(&vote)
		if err != nil {
			v.invalidCall()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate HAL JSON response\n" + err.Error()})
			return
		}
		votes = append(votes, vote)
	}

	// find all votes with voterId
	var votesWithVoterId []schema.Vote
	for _, vote := range votes {
		if vote.VoterId == voterId {
			votesWithVoterId = append(votesWithVoterId, vote)
		}
	}

	v.validCall()
	c.JSON(http.StatusOK, votesWithVoterId)
}

func (v *VotesAPI) PostVote(c *gin.Context) {
	id := c.Param("voteId")
	if id == "" {
		v.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No vote ID provided"})
		return
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		v.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vote ID provided"})
		return
	}

	cacheKey := "votes:" + id
	_, err = v.helper.JSONGet(cacheKey, ".")
	if err == nil {
		v.invalidCall()
		c.JSON(http.StatusNotFound, gin.H{"error": "Vote Id already exists in cache with id=" + cacheKey})
		return
	}

	var vote schema.Vote
	err = c.ShouldBindJSON(&vote)
	if err != nil {
		v.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse JSON"})
		return
	}

	//confirm voter and poll exist
	var voter schema.Voter
	var poll schema.Poll
	err = getVoterAndPoll(&vote, v, &voter, &poll)
	if err != nil {
		v.invalidCall()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find voter or poll\n" + err.Error()})
		return
	}

	// check if the option exists
	if vote.VoteValue < 0 || vote.VoteValue >= len(poll.Options) {
		v.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vote value"})
		return
	}
	// update the poll results
	poll.Results[vote.VoteValue].Votes++

	// update poll
	err = updatePollCounts(&vote, v, &poll)
	if err != nil {
		v.invalidCall()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update poll results in cache"})
		return
	}

	// update the total votes count on Voter
	voter.Meta.TotalVotes++
	voter.VoterPolls = append(voter.VoterPolls, schema.VoterPoll{
		PollId:  poll.Id,
		VoteId:  vote.Id,
		VotedAt: time.Now(),
	})

	//update voter
	err = updateVoter(&vote, v, &voter)
	if err != nil {
		v.invalidCall()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update voter in cache"})
		return
	}

	// set up links and embedded
	setLinkAndEmbeddedProps(v, &vote, voter, poll)

	vote.Meta.CreatedAt = time.Now()

	//save the vote
	err = v.saveVote(&vote)
	if err != nil {
		v.invalidCall()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save vote to cache"})
		return
	}

	// send it
	v.validCall()
	c.JSON(http.StatusOK, vote)
}

func (v *VotesAPI) DeleteVote(c *gin.Context) {
	// get the vote id
	id := c.Param("voteId")
	if id == "" {
		v.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No vote ID provided"})
		return
	}

	// check if the vote exists
	var vote schema.Vote
	cacheKey := "votes:" + id
	err := v.getItemFromRedis(cacheKey, &vote)
	if err != nil {
		v.invalidCall()
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find vote in cache with id=" + cacheKey})
		return
	}

	// get the voter and poll
	var voter schema.Voter
	var poll schema.Poll
	err = getVoterAndPoll(&vote, v, &voter, &poll)
	if err != nil {
		v.invalidCall()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find voter or poll\n" + err.Error()})
		return
	}

	// update the poll results
	poll.Results[vote.VoteValue].Votes--

	// update in redis
	err = updatePollCounts(&vote, v, &poll)
	if err != nil {
		v.invalidCall()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update poll results in cache"})
		return
	}

	// update the total votes count on Voter
	voter.Meta.TotalVotes--
	for i, vp := range voter.VoterPolls {
		if vp.VoteId == vote.Id {
			voter.VoterPolls = append(voter.VoterPolls[:i], voter.VoterPolls[i+1:]...)
			break
		}
	}

	//update in redis
	err = updateVoter(&vote, v, &voter)
	if err != nil {
		v.invalidCall()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update voter in cache"})
		return
	}

	// delete the vote
	_, err = v.helper.JSONDel(cacheKey, ".")
	if err != nil {
		v.invalidCall()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete vote from cache"})
		return
	}

	v.validCall()
	c.JSON(http.StatusOK, gin.H{"message": "Vote deleted"})
}

func (v *VotesAPI) generateHALJSONResponse(vote *schema.Vote) error {
	var poll schema.Poll
	var voter schema.Voter

	//First, we need to get the poll
	//Now we need to get the voter
	// unmarshal the voter and poll
	err := getVoterAndPoll(vote, v, &voter, &poll)
	if err != nil {
		return err
	}

	// set up links and embedded
	setLinkAndEmbeddedProps(v, vote, voter, poll)

	return nil
}

func getVoterAndPoll(vote *schema.Vote, v *VotesAPI, voter *schema.Voter, poll *schema.Poll) error {
	err := getVoter(vote, v, voter)
	if err != nil {
		return err
	}

	err = getPoll(vote, v, poll)
	if err != nil {
		return err
	}

	return nil
}

func getVoter(vote *schema.Vote, v *VotesAPI, voter *schema.Voter) error {
	voterId := vote.VoterId
	voterUrl := v.InternalAPI.Voters + "/" + fmt.Sprint(voterId)
	voterResp, err := v.apiClient.R().Get(voterUrl)
	if err != nil {
		return err
	}

	if voterResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("could not find voter with id=%d", voterId)
	}

	err = json.Unmarshal(voterResp.Body(), &voter)
	if err != nil {
		return err
	}
	return nil
}
func updateVoter(vote *schema.Vote, v *VotesAPI, voter *schema.Voter) error {
	voterId := vote.VoterId
	voterUrl := v.InternalAPI.Voters + "/" + fmt.Sprint(voterId)
	// update voter
	voterResp, err := v.apiClient.R().SetBody(voter).Put(voterUrl)
	if err != nil {
		return err
	}

	if voterResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("could not update voter with id=%d", voterId)
	}

	err = json.Unmarshal(voterResp.Body(), &voter)
	if err != nil {
		return err
	}

	return nil
}
func updatePollCounts(vote *schema.Vote, v *VotesAPI, poll *schema.Poll) error {
	pollId := vote.PollId
	pollUrl := v.InternalAPI.Polls + "/counts/" + fmt.Sprint(pollId)
	// update poll
	pollResp, err := v.apiClient.R().SetBody(poll).Put(pollUrl)
	if err != nil {
		return err
	}

	if pollResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("could not update poll with id=%d", pollId)
	}

	err = json.Unmarshal(pollResp.Body(), &poll)
	if err != nil {
		return err
	}

	return nil
}
func getPoll(vote *schema.Vote, v *VotesAPI, poll *schema.Poll) error {
	pollId := vote.PollId
	pollUrl := v.InternalAPI.Polls + "/" + fmt.Sprint(pollId)
	pollResp, err := v.apiClient.R().Get(pollUrl)
	if err != nil {
		return err
	}
	if pollResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("could not find voter with id=%d", pollId)
	}
	err = json.Unmarshal(pollResp.Body(), &poll)
	if err != nil {
		return err
	}
	return nil
}

func setLinkAndEmbeddedProps(v *VotesAPI, vote *schema.Vote, voter schema.Voter, poll schema.Poll) {
	var links schema.Links
	links.Self.Href = v.API.Self + "/votes/" + fmt.Sprint(vote.Id)
	links.Voter.Href = v.API.Voters + "/" + fmt.Sprint(vote.VoterId)
	links.Poll.Href = v.API.Polls + "/" + fmt.Sprint(vote.PollId)
	links.Results.Href = v.API.Polls + "/" + fmt.Sprint(vote.PollId) + "/results"
	vote.Links = links

	var embedded Embedded
	embedded.Voter = voter
	embedded.Poll = poll
	vote.Embedded = embedded
}

func (v *VotesAPI) saveVote(vote *schema.Vote) error {
	// save vote in redis with votes:<id> as key
	cacheKey := redisKeyFromId(vote.Id)
	_, err := v.helper.JSONSet(cacheKey, ".", vote)
	if err != nil {
		return err
	}

	return nil
}

func redisKeyFromId(id int) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

func (v *VotesAPI) getItemFromRedis(key string, vote *schema.Vote) error {
	itemObject, err := v.helper.JSONGet(key, ".")
	if err != nil {
		return err
	}
	err = json.Unmarshal(itemObject.([]byte), vote)
	if err != nil {
		return err
	}

	return nil
}
