package httpserver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"server/api"
	"server/entities"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	API *api.Client
}

func NewServer(api *api.Client) Server {
	return Server{
		API: api,
	}
}

func (s *Server) ServeHTTP(port int) {
	//
	// HTTP Server and MiddleWare
	//
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // improvement: env var client url
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	//
	// HTTP Routes
	//
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		splitAuthHeader := strings.Split(authHeader, "Basic ")
		encodedEmailPassword := splitAuthHeader[1]
		decodedEmailPassword, err := base64.StdEncoding.DecodeString(encodedEmailPassword)

		if err != nil {
			// failed parsing auth
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		emailPassword := strings.SplitN(string(decodedEmailPassword), ":", 2)
		if len(emailPassword) != 2 {
			// failed parsing auth
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		user := entities.User{
			Email:    emailPassword[0],
			Password: emailPassword[1],
		}
		createdUser, err := s.API.CreateUser(user)
		if err != nil {
			if err.Error() == "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)" { // hacky error handling
				http.Error(w, http.StatusText(400), 400)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		response := entities.LoginResponse{
			UUID: createdUser.UUID,
		}

		json.NewEncoder(w).Encode(response)
	})

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		splitAuthHeader := strings.Split(authHeader, "Basic ")
		encodedEmailPassword := splitAuthHeader[1]
		decodedEmailPassword, err := base64.StdEncoding.DecodeString(encodedEmailPassword)

		if err != nil {
			// failed parsing auth
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		emailPassword := strings.SplitN(string(decodedEmailPassword), ":", 2)
		if len(emailPassword) != 2 {
			// failed parsing auth
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		email := emailPassword[0]
		password := emailPassword[1]

		userUUID, err := s.API.AuthorizeUser(email, password)
		if err != nil {
			if err.Error() == "unauthorized" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			return
		}

		response := entities.LoginResponse{
			UUID: *userUUID,
		}

		json.NewEncoder(w).Encode(response)
	})

	r.Get("/chat", func(w http.ResponseWriter, r *http.Request) {
		messages, err := s.API.GetMessages()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		response := entities.GetChatResponse{
			Messages: *messages,
		}

		json.NewEncoder(w).Encode(response)
	})

	formattedPort := fmt.Sprintf(":%v", port)
	http.ListenAndServe(formattedPort, r)
}
