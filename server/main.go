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
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		var createLoginRequest entities.CreateLoginRequest
		err := dec.Decode(&createLoginRequest)
		if err != nil {
			// TODO: handle error
			panic(err)
		}

		user := entities.User{
			Name:     createLoginRequest.Name,
			Email:    createLoginRequest.Email,
			Password: createLoginRequest.Password,
		}
		_, err = persistenceClient.CreateUser(user)
		if err != nil {
			if err.Error() == "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)" { // hacky error handling
				http.Error(w, http.StatusText(400), 400) // TODO: pass this to FE in interpretable way
				return
			}
			http.Error(w, http.StatusText(500), 500)
			return
		}

		// TODO: return user uuid

		response := entities.LoginResponse{
			UUID: user.UUID,
		}

		json.NewEncoder(w).Encode(response)
	})

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		// TODO: abstract decoding
		authHeader := r.Header.Get("Authorization")
		splitAuthHeader := strings.Split(authHeader, "Basic ")
		encodedUsernamePassword := splitAuthHeader[1]
		decodedUsernamePassword, err := base64.StdEncoding.DecodeString(encodedUsernamePassword)

		if err != nil {
			// failed parsing auth
			http.Error(w, http.StatusText(500), 500)
			return
		}

		usernamePassword := strings.SplitN(string(decodedUsernamePassword), ":", 2)
		if len(usernamePassword) != 2 {
			// failed parsing auth
			http.Error(w, http.StatusText(500), 500)
			return
		}

		username := usernamePassword[0]
		password := usernamePassword[1]

		userUUID, err := persistenceClient.AuthorizeUser(username, password)
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
