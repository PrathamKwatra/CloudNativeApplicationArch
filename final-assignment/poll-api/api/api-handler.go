package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"drexel.edu/polls/schema"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

const (
	RedisKeyPrefix      = "polls:"
	votesRedisKeyPrefix = "votes:"
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
	Votes  string
	Voters string
	Self   string
}

type PollsAPI struct {
	cache
	health    Health
	API       API
}

func (v *PollsAPI) validCall() {
	v.health.mu.Lock()
	v.health.totalValidApiCalls++
	v.health.mu.Unlock()
}

func (v *PollsAPI) invalidCall() {
	v.health.mu.Lock()
	v.health.totalApiCallsWithErrors++
	v.health.mu.Unlock()
}

func New(location string, api API) (*PollsAPI, error) {

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
	return &PollsAPI{
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
	}, nil
}

func (v *PollsAPI) CrashSim(c *gin.Context) {
	panic("Simulating an unexpected crash")
}

func (v *PollsAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"api":                  "Polls API",
		"status":               "ok",
		"uptime":               time.Since(v.health.startTime).String(),
		"msg":                  "Currently healthy",
		"totalValidCalls":      v.health.totalValidApiCalls,
		"totalCallsWithErrors": v.health.totalApiCallsWithErrors,
		"totalCalls":           v.health.totalApiCallsWithErrors + v.health.totalValidApiCalls,
	})
}

func (p *PollsAPI) GetPoll(c *gin.Context) {
	var poll schema.Poll
	id := c.Param("pollId")
	if id == "" {
		p.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No poll ID provided"})
		return
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid poll id",
		})
		p.invalidCall()
		return
	}

	err = getItemFromRedis(id, p, poll)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"msg": "Unable to retrieve poll.\n" + err.Error(),
			})
		p.invalidCall()
		return
	}

	genHalJSONResponse(poll, p)

	p.validCall()
	c.JSON(http.StatusOK, poll)
}

func getItemFromRedis(id string, p *PollsAPI, poll schema.Poll) error {
	pollKey := RedisKeyPrefix + id
	pollJSON, err := p.helper.JSONGet(pollKey, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(pollJSON.(string)), &poll)
	if err != nil {
		return err
	}
	return nil
}

func (p *PollsAPI) GetResults(c *gin.Context) {
	id := c.Param("pollId")
	if id == "" {
		p.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No poll ID provided"})
		return
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid poll id",
		})
		p.invalidCall()
		return
	}

	pollKey := RedisKeyPrefix + id
	pollJSON, err := p.helper.JSONGet(pollKey, ".")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Either the poll does not exist or there was an error retrieving it.",
		})
		p.invalidCall()
		return
	}

	var poll schema.Poll
	err = json.Unmarshal([]byte(pollJSON.(string)), &poll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Error unmarshalling poll",
		})
		p.invalidCall()
		return
	}
	var result struct {
		Results []schema.Results `json:"results"`
		Meta    schema.Meta      `json:"_meta"`
		Links   schema.Links     `json:"_links"`
	}
	result.Results = poll.Results
	result.Meta = poll.Meta
	result.Links = poll.Links
	p.validCall()
	c.JSON(http.StatusOK, result)
}

func (p *PollsAPI) GetPolls(c *gin.Context) {
	pollKey := RedisKeyPrefix + "*"
	polls, err := p.helper.JSONGet(pollKey, ".")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Error retrieving polls",
		})
		p.invalidCall()
		return
	}

	var pollList []schema.Poll
	err = json.Unmarshal([]byte(polls.(string)), &pollList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Error unmarshalling polls",
		})
		p.invalidCall()
		return
	}

	p.validCall()
	c.JSON(http.StatusOK, pollList)
}

func (p *PollsAPI) PostPoll(c *gin.Context) {
	var poll schema.Poll
	id := c.Param("pollId")
	if id == "" {
		p.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No poll ID provided"})
		return
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid poll id",
		})
		p.invalidCall()
		return
	}

	// bind the request body into a poll struct
	err = c.BindJSON(&poll)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Error unmarshalling poll",
		})
		p.invalidCall()
		return
	}

	poll.Results = make([]schema.Results, len(poll.Options))
	for i, option := range poll.Options {
		poll.Results[i].OptionId = option.Id
		poll.Results[i].Votes = 0
	}

	poll.Meta.TotalVotes = 0
	poll.Meta.CreatedAt = time.Now()

	// generate the links and embedded
	genHalJSONResponse(poll, p)

	err = p.savePoll(&poll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Error saving poll",
		})
		p.invalidCall()
		return
	}

	p.validCall()
	c.JSON(http.StatusOK, poll)
}

func genHalJSONResponse(poll schema.Poll, p *PollsAPI) {
	poll.Links.Self.Href = p.API.Self + "/polls/" + strconv.Itoa(poll.Id)
	poll.Links.Vote.Href = p.API.Votes
	poll.Links.Voters.Href = p.API.Voters
	poll.Links.Results.Href = p.API.Self + "/polls/" + strconv.Itoa(poll.Id) + "/results"
}

func (p *PollsAPI) UpdatePoll(c *gin.Context) {
	var poll schema.Poll
	id := c.Param("pollId")
	if id == "" {
		p.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No poll ID provided"})
		return
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid poll id",
		})
		p.invalidCall()
		return
	}

	err = getItemFromRedis(id, p, poll)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"msg": "Unable to retrieve poll.\n" + err.Error(),
			})
		p.invalidCall()
		return
	}

	// updating the poll should reset the votes.
	var newPoll schema.Poll
	err = c.BindJSON(&newPoll)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Error unmarshalling poll",
		})
		p.invalidCall()
		return
	}

	newPoll.Results = make([]schema.Results, len(newPoll.Options))
	for i, option := range newPoll.Options {
		newPoll.Results[i].OptionId = option.Id
		newPoll.Results[i].Votes = 0
	}

	newPoll.Meta.CreatedAt = poll.Meta.CreatedAt
	newPoll.Meta.UpdatedAt = time.Now()

	// generate the links and embedded
	genHalJSONResponse(newPoll, p)

	err = p.savePoll(&newPoll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Error saving poll",
		})
		p.invalidCall()
		return
	}

	p.validCall()
	c.JSON(http.StatusOK, poll)
}

func (p *PollsAPI) DeletePoll(c *gin.Context) {
	// get the poll id
	id := c.Param("pollId")
	if id == "" {
		p.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No poll ID provided"})
		return
	}

	// check if the poll exists
	cacheKey := RedisKeyPrefix + id
	_, err := p.helper.JSONGet(cacheKey, ".")
	if err != nil {
		p.invalidCall()
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find poll in cache with id=" + cacheKey})
		return
	}

	// recursive delete of votes will be tough, unless the ids are known...
	// value search for a key in redis is currently not understood. :(
	// delete the poll
	_, err = p.helper.JSONDel(cacheKey, ".")
	if err != nil {
		p.invalidCall()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete poll from cache"})
		return
	}

	p.validCall()
	c.JSON(http.StatusOK, gin.H{"message": "Poll deleted"})
}

func (p *PollsAPI) savePoll(poll *schema.Poll) error {
	// save vote in redis with votes:<id> as key
	cacheKey := RedisKeyPrefix + strconv.Itoa(poll.Id)
	_, err := p.helper.JSONSet(cacheKey, ".", poll)
	if err != nil {
		return err
	}

	return nil
}
