package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"drexel.edu/voter/db"
	"github.com/gin-gonic/gin"
)

type VoterAPI struct {
	db                      *db.Voters
	startTime               time.Time
	totalValidApiCalls      int
	totalApiCallsWithErrors int
}

func New() (*VoterAPI, error) {
	dbHandler, err := db.New()
	if err != nil {
		return nil, err
	}

	return &VoterAPI{
		db:        dbHandler,
		startTime: time.Now(),
	}, nil
}

func (v *VoterAPI) abort(c *gin.Context, status int) {
	v.totalApiCallsWithErrors++
	c.AbortWithStatus(status)
}

func (v *VoterAPI) AddVoter(c *gin.Context) {

	//With HTTP based APIs, a POST request will usually
	//have a body that contains the data to be added
	//to the database.  The body is usually JSON, so
	//we need to bind the JSON to a struct that we
	//can use in our code.
	//This framework exposes the raw body via c.Request.Body
	//but it also provides a helper function ShouldBindJSON()
	//that will extract the body, convert it to JSON and
	//bind it to a struct for us.  It will also report an error
	//if the body is not JSON or if the JSON does not match
	//the struct we are binding to.
	var voter db.Voter

	if err := c.ShouldBindJSON(&voter); err != nil {
		log.Println("Error (Add Voter, JSON Binding): ", err)
		v.abort(c, http.StatusBadRequest)
		return
	}

	if err := v.db.AddVoter(voter); err != nil {
		log.Println("Error (Add Voter): ", err)
		v.abort(c, http.StatusInternalServerError)
		return
	}

	v.totalValidApiCalls++
	c.JSON(http.StatusOK, voter)

}

func (v *VoterAPI) CrashSim(c *gin.Context) {
	//panic() is go's version of throwing an exception
	panic("Simulating an unexpected crash")
}

func (v *VoterAPI) DeleteAllVoter(c *gin.Context) {

	if err := v.db.DeleteAll(); err != nil {
		log.Println("Error (Delete All Voter): ", err)
		v.abort(c, http.StatusInternalServerError)
		return
	}

	v.totalValidApiCalls++
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "all items deleted",
	})
}

func (v *VoterAPI) DeleteVoter(c *gin.Context) {

	param := c.Param("id")
	id64, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		log.Println("Error (Delete Voter, Parsing): ", err)
		v.totalApiCallsWithErrors++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := v.db.DeleteVoter(id64); err != nil {
		log.Println("Error (Delete Voter): ", err)
		v.totalApiCallsWithErrors++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	v.totalValidApiCalls++
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "deleted item " + param,
	})

}

func (v *VoterAPI) GetVoter(c *gin.Context) {

	param := c.Param("id")
	id64, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		log.Println("Error (Get Voter, Parsing): ", err)
		v.abort(c, http.StatusBadRequest)
		return
	}

	voter, err := v.db.GetVoter(id64)
	if err != nil {
		log.Println("Error (Get Voter): ", err)
		v.abort(c, http.StatusNotFound)
		return
	}

	v.totalValidApiCalls++
	c.JSON(http.StatusOK, voter)
}

func (v *VoterAPI) GetVoterHistory(c *gin.Context) {

	param := c.Param("id")
	id64, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		log.Println("Error (Get Voter History, Parsing): ", err)
		v.abort(c, http.StatusBadRequest)
		return
	}

	voter, err := v.db.GetVoter(id64)
	if err != nil {
		log.Println("Error (Get Voter History): ", err)
		v.abort(c, http.StatusNotFound)
		return
	}

	v.totalValidApiCalls++
	c.JSON(http.StatusOK, voter.VoteHistory)
}

func (v *VoterAPI) GetPoll(c *gin.Context) {

	param1 := c.Param("id")
	vid, err := strconv.ParseUint(param1, 10, 64)
	if err != nil {
		log.Println("Error (Get Poll, Parsing Voter ID): ", err)
		v.abort(c, http.StatusBadRequest)
		return
	}

	param2 := c.Param("pollsid")
	pollsid, err := strconv.ParseUint(param2, 10, 64)
	if err != nil {
		log.Println("Error (Get Poll, Parsing Poll ID): ", err)
		v.abort(c, http.StatusBadRequest)
		return
	}

	poll, err := v.db.GetPoll(vid, pollsid)
	if err != nil {
		log.Println("Error (Get Poll): ", err)
		v.abort(c, http.StatusNotFound)
		return
	}

	v.totalValidApiCalls++
	c.JSON(http.StatusOK, poll)
}

func (v *VoterAPI) AddPoll(c *gin.Context) {

	param1 := c.Param("id")
	vid, err := strconv.ParseUint(param1, 10, 64)
	if err != nil {
		log.Println("Error (Add Poll, Parsing): ", err)
		v.abort(c, http.StatusBadRequest)
		return
	}

	param2 := c.Param("pollsid")
	pollsid, err := strconv.ParseUint(param2, 10, 64)
	if err != nil {
		log.Println("Error (Add Poll, Parsing): ", err)
		v.abort(c, http.StatusBadRequest)
		return
	}

	if err := v.db.AddPoll(vid, pollsid); err != nil {
		log.Println("Error (Add Poll): ", err)
		v.abort(c, http.StatusInternalServerError)
		return
	}

	v.totalValidApiCalls++
	c.JSON(http.StatusOK, pollsid)
}

func (v *VoterAPI) UpdatePoll(c *gin.Context) {
	param1 := c.Param("id")
	vid, err := strconv.ParseUint(param1, 10, 64)
	if err != nil {
		log.Println("Error (Update Poll, Parsing): ", err)
		v.abort(c, http.StatusBadRequest)
		return
	}

	param2 := c.Param("pollsid")
	pollsid, err := strconv.ParseUint(param2, 10, 64)
	if err != nil {
		log.Println("Error (Update Poll, Parsing): ", err)
		v.abort(c, http.StatusBadRequest)
		return
	}

	if err := v.db.UpdatePoll(vid, pollsid); err != nil {
		log.Println("Error (Update Poll): ", err)
		v.abort(c, http.StatusInternalServerError)
		return
	}

	v.totalValidApiCalls++
	c.JSON(http.StatusOK, pollsid)
}

func (v *VoterAPI) DeletePoll(c *gin.Context) {
	param1 := c.Param("id")
	vid, err := strconv.ParseUint(param1, 10, 64)
	if err != nil {
		log.Println("Error (Delete Poll, Parsing): ", err)
		v.abort(c, http.StatusBadRequest)
		return
	}

	param2 := c.Param("pollsid")
	pollsid, err := strconv.ParseUint(param2, 10, 64)
	if err != nil {
		log.Println("Error (Delete Poll, Parsing): ", err)
		v.abort(c, http.StatusBadRequest)
		return
	}

	if err := v.db.DeletePoll(vid, pollsid); err != nil {
		log.Println("Error (Delete Poll): ", err)
		v.abort(c, http.StatusInternalServerError)
		return
	}

	v.totalValidApiCalls++
	c.JSON(http.StatusOK, pollsid)
}

func (v *VoterAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":               "ok",
		"uptime":               time.Since(v.startTime).String(),
		"msg":                  "Currently healthy",
		"totalValidCalls":      v.totalValidApiCalls,
		"totalCallsWithErrors": v.totalApiCallsWithErrors,
		"totalCalls":           v.totalApiCallsWithErrors + v.totalValidApiCalls,
	})
}

func (v *VoterAPI) ListAllVoter(c *gin.Context) {

	voters, err := v.db.GetAllItems()
	if err != nil {
		log.Println("Error (Get All Items, Parsing): ", err)
		v.abort(c, http.StatusInternalServerError)
		return
	}

	if voters == nil {
		voters = make([]db.Voter, 0)
	}

	v.totalValidApiCalls++
	c.JSON(http.StatusOK, voters)
}

func (v *VoterAPI) UpdateVoter(c *gin.Context) {

	var voter db.Voter

	if err := c.ShouldBindJSON(&voter); err != nil {
		log.Println("Error (Update Voter, JSON Binding): ", err)
		v.abort(c, http.StatusBadRequest)
		return
	}

	if err := v.db.UpdateVoter(voter); err != nil {
		log.Println("Error (Update Voter): ", err)
		v.abort(c, http.StatusInternalServerError)
		return
	}

	v.totalValidApiCalls++
	c.JSON(http.StatusOK, voter)
}
