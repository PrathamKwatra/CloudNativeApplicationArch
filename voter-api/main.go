package main

import (
	"flag"
	"fmt"
	"os"

	"drexel.edu/voter/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {

	//Note some networking lingo, some frameworks start the server on localhost
	//this is a local-only interface and is fine for testing but its not accessible
	//from other machines.  To make the server accessible from other machines, we
	//need to listen on an interface, that could be an IP address, but modern
	//cloud servers may have multiple network interfaces for scale.  With TCP/IP
	//the address 0.0.0.0 instructs the network stack to listen on all interfaces
	//We set this up as a flag so that we can overwrite it on the command line if
	//needed
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

	flag.Parse()
}

func main() {
	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	apiHandler, err := api.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// r.DELETE("/voter", apiHandler.DeleteAllVoter)
	// r.DELETE("/voter/:id", apiHandler.DeleteVoter)
	r.GET("/crash", apiHandler.CrashSim)
	r.GET("/voter", apiHandler.ListAllVoter)
	r.GET("/voter/:id", apiHandler.GetVoter)
	// r.POST("/voter", apiHandler.AddVoter)
	// r.PUT("/voter", apiHandler.UpdateVoter)

	// //We will now show a common way to version an API and add a new
	// //version of an API handler under /v2.  This new API will support
	// //a path parameter to search for todos based on a status
	// v2 := r.Group("/v2")
	// v2.GET("/todo", apiHandler.ListSelectTodos)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
