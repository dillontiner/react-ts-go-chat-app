package main

import (
	"log"
	"net/http"
	"os"
	"server/api"
	"server/entities"
	"server/httpserver"
	"server/websocketserver"

	"github.com/joho/godotenv"
)

func main() {
	// TODO Websocket Server on 4001 https://medium.com/rungo/running-multiple-http-servers-in-go-d15300f4e59f
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger := log.New(os.Stdout, "SERVER: ", log.LstdFlags)

	apiClient, err := api.NewClient()
	if err != nil {
		panic(err)
	}

	voteInMemoryQueue := make(chan entities.Vote)
	go func() {
		logger.Println("Running in memory queue")

		for {
			select {
			case vote := <-voteInMemoryQueue:
				logger.Println("QUEUE PROCESSING VOTE", vote)
				// TODO: process votes
				x, e := apiClient.VoteOnMessage(vote)
				logger.Println("QUEUE PROCESSING VOTE 2", x, e)
				continue
			}
		}
	}()

	go func() {
		logger.Println("Websocket Server running on port 4001")
		// TODO: enqueue
		wsServer := websocketserver.NewServer(apiClient, voteInMemoryQueue)
		wsServer.SetupRoutes("/chat")
		http.ListenAndServe(":4001", nil)
	}()

	logger.Println("HTTP Server running on port 4000")
	httpServer := httpserver.NewServer(apiClient)
	httpServer.ServeHTTP(4000)
}
