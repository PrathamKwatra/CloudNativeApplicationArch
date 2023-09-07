package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"drexel.edu/votes/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag  string
	cacheURL  string
	portFlag  uint
	pollsURL  string
	votersURL string
)

func processCmdLineFlags() {

	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	// flag.StringVar(&hostFlag, "h", "localhost", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port (cannot be changed)")
	flag.StringVar(&cacheURL, "c", "localhost:6379", "Default cache location")

	// flags for internal api
	flag.StringVar(&pollsURL, "polls", "http://localhost:1082/polls", "Default polls location")
	flag.StringVar(&votersURL, "voters", "http://localhost:1081/voters", "Default voters location")
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
	pollsURL = envVarOrDefault("POLL_API_URL", pollsURL)
	votersURL = envVarOrDefault("VOTER_API_URL", votersURL)

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
		Polls:  "http://localhost:1082/polls",
		Voters: "http://localhost:1081/voters",
		Self:   "http://localhost:1080",
	},
		api.API{
			Polls:  pollsURL,
			Voters: votersURL,
		},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// r.DELETE("/voters", apiHandler.DeleteAllVoter)
	r.GET("/", apiHandler.GetVotes)
	r.GET("/crash", apiHandler.CrashSim)
	r.GET("/votes/health", apiHandler.HealthCheck)
	r.GET("/votes/:voteId", apiHandler.GetVote)
	r.GET("/votes", apiHandler.GetVotes)
	r.GET("/votes/voters/:voterId", apiHandler.GetVotesByVoter)
	r.GET("/votes/polls/:pollId", apiHandler.GetVotesByPolls)
	r.POST("/votes/:voteId", apiHandler.PostVote)
	r.DELETE("/votes/:voteId", apiHandler.DeleteVote)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
