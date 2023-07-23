package api

import (
	"log"
	"net/http"
	"strconv"

	"drexel.edu/voter/db"
	"github.com/gin-gonic/gin"
)

type VoterAPI struct {
	db *db.Voters
}

func New() (*VoterAPI, error) {
	dbHandler, err := db.New()
	if err != nil {
		return nil, err
	}

	return &VoterAPI{
		db: dbHandler,
	}, nil
}

func (v *VoterAPI) AddVoter(c *gin.Context) {}

func (v *VoterAPI) CrashSim(c *gin.Context) {
	//panic() is go's version of throwing an exception
	panic("Simulating an unexpected crash")
}

func (v *VoterAPI) DeleteAllVoter(c *gin.Context) {}
func (v *VoterAPI) DeleteVoter(c *gin.Context)    {}

func (v *VoterAPI) GetVoter(c *gin.Context) {

	param := c.Param("id")
	id64, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		log.Println("Error (Get Voter, Parsing): ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter, err := v.db.GetItem(id64)
	if err != nil {
		log.Println("Error (Get Voter): ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voter)
}

func (v *VoterAPI) GetVoterHistory(c *gin.Context) {

	param := c.Param("id")
	id64, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		log.Println("Error (Get Voter History, Parsing): ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter, err := v.db.GetItem(id64)
	if err != nil {
		log.Println("Error (Get Voter History): ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voter.VoteHistory)
}

func (v *VoterAPI) GetPoll(c *gin.Context) {

	param1 := c.Param("id")
	vid, err := strconv.ParseUint(param1, 10, 64)
	if err != nil {
		log.Println("Error (Get Poll, Parsing Voter ID): ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	param2 := c.Param("pollsid")
	pollsid, err := strconv.ParseUint(param2, 10, 64)
	if err != nil {
		log.Println("Error (Get Poll, Parsing Poll ID): ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	poll, err := v.db.GetPoll(vid, pollsid)
	if err != nil {
		log.Println("Error (Get Poll): ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, poll)
}

func (v *VoterAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (v *VoterAPI) ListAllVoter(c *gin.Context) {

	voters, err := v.db.GetAllItems()
	if err != nil {
		log.Println("Error (Get All Items, Parsing): ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if voters == nil {
		voters = make([]db.Voter, 0)
	}

	c.JSON(http.StatusOK, voters)
}

func (v *VoterAPI) UpdateVoter(c *gin.Context) {}
