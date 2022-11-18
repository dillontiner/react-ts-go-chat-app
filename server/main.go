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
	logger := log.New(os.Stdout, "SERVER: ", log.LstdFlags)

	err := godotenv.Load()
	if err != nil {
		logger.Fatal(err)
	}

	// API for business logic
	apiClient, err := api.NewClient()
	if err != nil {
		logger.Fatal(err)
	}

	// Queue for processing votes
	voteInMemoryQueue := make(chan entities.Vote)
	go func() {
		logger.Println("Running in memory queue to process votes")

		for {
			select {
			case vote := <-voteInMemoryQueue:
				logger.Println("received vote:", vote)
				_, err := apiClient.VoteOnMessage(vote)
				if err != nil {
					logger.Println(err)
				}
				continue
			}
		}
	}()

	// Websocket Server
	go func() {
		logger.Println("Websocket Server running on port 4001")
		// TODO: enqueue
		wsServer := websocketserver.NewServer(apiClient, voteInMemoryQueue)
		wsServer.SetupRoutes("/chat")
		http.ListenAndServe(":4001", nil)
	}()

	// HTTP Server
	logger.Println("HTTP Server running on port 4000")
	httpServer := httpserver.NewServer(apiClient)
	httpServer.ServeHTTP(4000)
}
