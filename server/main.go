package main

import (
	"chat-app-server/entities"
	"chat-app-server/persistence"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// TODO Websocket Server on 4001 https://medium.com/rungo/running-multiple-http-servers-in-go-d15300f4e59f

	logger := log.New(os.Stdout, "chat-app-server: ", log.LstdFlags)

	persistenceClient, err := persistence.NewClient()
	if err != nil {
		panic(err)
	}
	user := entities.User{
		Name:     "dillon",
		Email:    "Dillon@gmail.com",
		Password: "asdsadsa",
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
