package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"server/entities"
	"server/persistence"
	"strings"

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

	// HTTP Server on 4000
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		// TODO: abstract decoding
		authHeader := r.Header.Get("Authorization")
		splitAuthHeader := strings.Split(authHeader, "Basic ")
		encodedEmailPassword := splitAuthHeader[1]
		decodedEmailPassword, err := base64.StdEncoding.DecodeString(encodedEmailPassword)

		if err != nil {
			// failed parsing auth
			http.Error(w, http.StatusText(500), 500)
			return
		}

		emailPassword := strings.SplitN(string(decodedEmailPassword), ":", 2)
		if len(emailPassword) != 2 {
			// failed parsing auth
			http.Error(w, http.StatusText(500), 500)
			return
		}

		user := entities.User{
			Email:    emailPassword[0],
			Password: emailPassword[1],
		}
		createdUser, err := persistenceClient.CreateUser(user)
		if err != nil {
			if err.Error() == "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)" { // hacky error handling
				http.Error(w, http.StatusText(400), 400) // TODO: pass this to FE in interpretable way
				return
			}
			http.Error(w, http.StatusText(500), 500)
			return
		}

		response := entities.LoginResponse{
			UUID: createdUser.UUID,
		}

		json.NewEncoder(w).Encode(response)
	})

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		// TODO: abstract decoding
		authHeader := r.Header.Get("Authorization")
		splitAuthHeader := strings.Split(authHeader, "Basic ")
		encodedEmailPassword := splitAuthHeader[1]
		decodedEmailPassword, err := base64.StdEncoding.DecodeString(encodedEmailPassword)

		if err != nil {
			// failed parsing auth
			http.Error(w, http.StatusText(500), 500)
			return
		}

		emailPassword := strings.SplitN(string(decodedEmailPassword), ":", 2)
		if len(emailPassword) != 2 {
			// failed parsing auth
			http.Error(w, http.StatusText(500), 500)
			return
		}

		email := emailPassword[0]
		password := emailPassword[1]

		userUUID, err := persistenceClient.AuthorizeUser(email, password)
		if err != nil {
			if err.Error() == "unauthorized" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			} else {
				http.Error(w, http.StatusText(500), 500)
			}

			return
		}

		response := entities.LoginResponse{
			UUID: *userUUID,
		}

		json.NewEncoder(w).Encode(response)
	})

	logger.Println("HTTP Server running on port 4000")
	http.ListenAndServe(":4000", r)
}
