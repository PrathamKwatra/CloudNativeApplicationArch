package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"log"

	"drexel.edu/polls/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag string
	cacheURL string
	portFlag uint
)

func processCmdLineFlags() {

	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	// flag.StringVar(&hostFlag, "h", "localhost", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1082, "Default Port")
	flag.StringVar(&cacheURL, "c", "localhost:6379", "Default cache location")
	

	flag.Parse()
}

func envVarOrDefault(envVar string, defaultVal string) string {
	envVal := os.Getenv(envVar)
	if envVal != "" {
		return envVal
	}
	return defaultVal
}

func setupParms() {
	//first process any command line flags
	processCmdLineFlags()

	//now process any environment variables
	cacheURL = envVarOrDefault("REDIS_URL", cacheURL)
	hostFlag = envVarOrDefault("RLAPI_HOST", hostFlag)

	pfNew, err := strconv.Atoi(envVarOrDefault("RLAPI_PORT", fmt.Sprintf("%d", portFlag)))
	//only update the port if we were able to convert the env var to an int, else
	//we will use the default we got from the command line, or command line defaults
	if err == nil {
		portFlag = uint(pfNew)
	}

}

func main() {
	setupParms()
	log.Println("Init/cacheURL: " + cacheURL)
	log.Println("Init/hostFlag: " + hostFlag)
	log.Printf("Init/portFlag: %d", portFlag)

	r := gin.Default()
	r.Use(cors.Default())

	apiHandler, err := api.New(cacheURL, api.API{
		Self: "http://localhost:1082",
		Voters: "http://localhost:1081/voters",
		Votes: "http://localhost:1080/votes",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r.GET("/", apiHandler.GetPolls)
	r.GET("/crash", apiHandler.CrashSim)
	r.GET("/polls/health", apiHandler.HealthCheck)
	r.GET("/polls/:pollId", apiHandler.GetPoll)
	r.GET("/polls/:pollId/results", apiHandler.GetResults)
	r.GET("/polls", apiHandler.GetPolls)
	r.POST("/polls/:pollId", apiHandler.PostPoll)
	r.PUT("/polls/:pollId", apiHandler.UpdatePoll)
	r.DELETE("/polls/:pollId", apiHandler.DeletePoll)
	
	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
