package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"drexel.edu/voters/schema"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

const (
	RedisKeyPrefix = "voters:"
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
	Votes string
	Polls string
	Self  string
}

type VotersAPI struct {
	cache
	health Health
	API    API
}

func (v *VotersAPI) validCall() {
	v.health.mu.Lock()
	v.health.totalValidApiCalls++
	v.health.mu.Unlock()
}

func (v *VotersAPI) invalidCall() {
	v.health.mu.Lock()
	v.health.totalApiCallsWithErrors++
	v.health.mu.Unlock()
}

func New(location string, api API) (*VotersAPI, error) {

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
	return &VotersAPI{
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

func (v *VotersAPI) CrashSim(c *gin.Context) {
	panic("Simulating an unexpected crash")
}

func (v *VotersAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"api":                  "Voter API",
		"status":               "ok",
		"uptime":               time.Since(v.health.startTime).String(),
		"msg":                  "Currently healthy",
		"totalValidCalls":      v.health.totalValidApiCalls,
		"totalCallsWithErrors": v.health.totalApiCallsWithErrors,
		"totalCalls":           v.health.totalApiCallsWithErrors + v.health.totalValidApiCalls,
	})
}

func (v *VotersAPI) GetVoter(c *gin.Context) {
	var voter schema.Voter
	id := c.Param("voterId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No Voter ID provided",
		})
		v.invalidCall()
		return
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid voter id",
		})
		v.invalidCall()
		return
	}

	err = getItemFromRedis(id, v, &voter)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"msg": "No such voter in the cache\n" + err.Error(),
			})
		v.invalidCall()
		return
	}

	genHalJSONResponse(&voter, v)
	v.validCall()
	c.JSON(http.StatusOK, voter)
}
func (v *VotersAPI) GetVoters(c *gin.Context) {
	voterKey := RedisKeyPrefix + "*"
	voters, err := v.helper.JSONGet(voterKey, ".")
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"msg": "Error getting voters from cache",
			})
		v.invalidCall()
		return
	}

	var voterList []schema.Voter
	err = json.Unmarshal(voters.([]byte), &voterList)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"msg": "Error getting voters from cache",
			})
		v.invalidCall()
		return
	}

	for i := range voterList {
		genHalJSONResponse(&voterList[i], v)
	}

	v.validCall()
	c.JSON(http.StatusOK, voterList)
}
func (v *VotersAPI) PostVoter(c *gin.Context) {
	var voter schema.Voter
	id := c.Param("voterId")
	if id == "" {
		v.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No voter ID provided"})
		return
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid voter id",
		})
		v.invalidCall()
		return
	}

	// confirm that the voter does not exist
	err = getItemFromRedis(id, v, &voter)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Voter already exists",
		})
		v.invalidCall()
		return
	}

	// bind the request body into a voter struct
	err = c.BindJSON(&voter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Error unmarshalling voter\n" + err.Error(),
		})
		v.invalidCall()
		return
	}

	voter.VoterPolls = []schema.VoterPoll{}
	voter.Meta.TotalVotes = 0
	voter.Meta.CreatedAt = time.Now()
	voter.Meta.UpdatedAt = time.Now()

	genHalJSONResponse(&voter, v)

	err = v.saveVoter(&voter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Error saving voter to cache",
		})
		v.invalidCall()
		return
	}

	v.validCall()
	c.JSON(http.StatusOK, voter)
}
func (v *VotersAPI) UpdateVoter(c *gin.Context) {
	var voter schema.Voter
	id := c.Param("voterId")
	if id == "" {
		v.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No voter ID provided"})
		return
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid voter id",
		})
		v.invalidCall()
		return
	}

	// get the old voter
	err = getItemFromRedis(id, v, &voter)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"msg": "No such voter in the cache\n" + err.Error(),
			})
		v.invalidCall()
		return
	}

	var newVoter schema.Voter
	err = c.BindJSON(&newVoter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Error unmarshalling voter\n" + err.Error(),
		})
		v.invalidCall()
		return
	}

	// updating the voter should reset the votes.
	newVoter.Meta.TotalVotes = 0
	newVoter.Meta.CreatedAt = voter.Meta.CreatedAt
	newVoter.Meta.UpdatedAt = time.Now()

	genHalJSONResponse(&newVoter, v)

	err = v.saveVoter(&newVoter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Error saving voter to cache",
		})
		v.invalidCall()
		return
	}

	v.validCall()
	c.JSON(http.StatusOK, newVoter)
}
func (v *VotersAPI) DeleteVoter(c *gin.Context) {
	id := c.Param("voterId")
	if id == "" {
		v.invalidCall()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No voter ID provided"})
		return
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid voter id",
		})
		v.invalidCall()
		return
	}

	voterKey := RedisKeyPrefix + id
	_, err = v.helper.JSONDel(voterKey, ".")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Error deleting voter from cache",
		})
		v.invalidCall()
		return
	}

	v.validCall()
	c.JSON(http.StatusOK, gin.H{
		"msg": "Voter deleted",
	})
}

func (v *VotersAPI) saveVoter(voter *schema.Voter) error {
	// save vote in redis with votes:<id> as key
	cacheKey := RedisKeyPrefix + strconv.Itoa(voter.Id)
	_, err := v.helper.JSONSet(cacheKey, ".", voter)
	if err != nil {
		return err
	}

	return nil
}

func getItemFromRedis(id string, p *VotersAPI, voter *schema.Voter) error {
	voterKey := RedisKeyPrefix + id
	voterJSON, err := p.helper.JSONGet(voterKey, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(voterJSON.([]byte), &voter)
	if err != nil {
		return err
	}
	return nil
}

func genHalJSONResponse(voter *schema.Voter, p *VotersAPI) {
	voter.Links.Self.Href = p.API.Self + "/voters/" + strconv.Itoa(voter.Id)
	voter.Links.Polls.Href = p.API.Polls
	voter.Links.Votes.Href = p.API.Votes + "/voters/" + strconv.Itoa(voter.Id)
	voter.Links.Vote.Href = p.API.Votes + "/voters/" + strconv.Itoa(voter.Id)
}
