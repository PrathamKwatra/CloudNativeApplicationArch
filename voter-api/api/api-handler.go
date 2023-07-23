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

func (v *VoterAPI) AddVoter(c *gin.Context)       {}
func (v *VoterAPI) CrashSim(c *gin.Context)       {}
func (v *VoterAPI) DeleteAllVoter(c *gin.Context) {}
func (v *VoterAPI) DeleteVoter(c *gin.Context)    {}
func (v *VoterAPI) GetVoter(c *gin.Context) {

	param := c.Param("id")
	id64, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		log.Println("Error (Get Voter): ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter, err := v.db.GetItem(uint(id64))
	if err != nil {
		log.Println("Error (Get Voter): ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voter)
}

func (v *VoterAPI) ListAllVoter(c *gin.Context) {

	voters, err := v.db.GetAllItems()
	if err != nil {
		log.Println("Error (Get All Items): ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if voters == nil {
		voters = make([]db.Voter, 0)
	}

	c.JSON(http.StatusOK, voters)
}

func (v *VoterAPI) UpdateVoter(c *gin.Context) {}
