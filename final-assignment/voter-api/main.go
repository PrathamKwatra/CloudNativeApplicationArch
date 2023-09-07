package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"drexel.edu/voters/api"
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
	flag.UintVar(&portFlag, "p", 1081, "Default Port (cannot be changed)")
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

	// pfNew, err := strconv.Atoi(envVarOrDefault("RLAPI_PORT", fmt.Sprintf("%d", portFlag)))
	// //only update the port if we were able to convert the env var to an int, else
	// //we will use the default we got from the command line, or command line defaults
	// if err == nil {
	// 	portFlag = uint(pfNew)
	// }

}

func main() {
	setupParms()
	log.Println("Init/cacheURL: " + cacheURL)
	log.Println("Init/hostFlag: " + hostFlag)
	log.Printf("Init/portFlag: %d", portFlag)

	r := gin.Default()
	r.Use(cors.Default())

	apiHandler, err := api.New(cacheURL, api.API{
		Self:  "http://localhost:1081",
		Polls: "http://localhost:1082/polls",
		Votes: "http://localhost:1080/votes",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r.GET("/", apiHandler.GetVoters)
	r.GET("/crash", apiHandler.CrashSim)
	r.GET("/voters/health", apiHandler.HealthCheck)
	r.GET("/voters/:voterId", apiHandler.GetVoter)
	// r.GET("/voters/:voterId/polls/:pollsId", apiHandler.GetPoll)
	r.GET("/voters", apiHandler.GetVoters)
	r.POST("/voters/:voterId", apiHandler.PostVoter)
	r.PUT("/voters/:voterId", apiHandler.UpdateVoter)
	r.DELETE("/voters/:voterId", apiHandler.DeleteVoter)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
