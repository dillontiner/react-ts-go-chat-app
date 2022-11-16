package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"server/entities"
	"server/persistence"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// TODO Websocket Server on 4001 https://medium.com/rungo/running-multiple-http-servers-in-go-d15300f4e59f
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger := log.New(os.Stdout, "server: ", log.LstdFlags)

	persistenceClient, err := persistence.NewClient()
	if err != nil {
		panic(err)
	}
	user := entities.User{
		Name:     "dillon",
		Email:    "dillontiner@gmail.com",
		Password: "dsfkjhkdsf",
	}
	x, e := persistenceClient.CreateUser(user)
	fmt.Println(x, e)

	// HTTP Server on 4000
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	logger.Println("HTTP Server running on port 4000")
	http.ListenAndServe(":4000", r)
}
