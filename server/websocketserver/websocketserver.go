package websocketserver

import (
	"encoding/json"
	"log"
	"net/http"
	"server/entities"
	"server/persistence"

	"github.com/gorilla/websocket"
)

// var addr = flag.String("addr", "localhost:8080", "http service address")
type Server struct {
	api *persistence.Client
}

func NewServer(api *persistence.Client) Server {
	return Server{
		api: api,
	}
}

func checkOrigin(r *http.Request) bool {
	// origin := r.Header.Get("Origin")
	// origin == "http://localhost:3000" || origin == "http://localhost:4001" // TODO: env var this
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: checkOrigin,
}

func (s *Server) HandleLiveChat(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(0)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("WEBSOCKET SERVER: upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Print("WEBSOCKET SERVER: read:", err)
			break
		}

		chatMessage := entities.Message{}
		err = json.Unmarshal(message, &chatMessage)
		if err != nil {
			log.Print("WEBSOCKET SERVER: ERROR: failed to unmarshal message", err)
		}

		createdMessage, err := s.api.CreateMessage(chatMessage)
		if err != nil {
			log.Print("WEBSOCKET SERVER: ERROR: failed to persist message", err)
			break
		}

		createdMessageBytes, err := json.Marshal(createdMessage)
		if err != nil {
			log.Print("WEBSOCKET SERVER: ERROR: failed to persist message", err)
			break
		}

		err = c.WriteMessage(mt, createdMessageBytes)
		if err != nil {
			log.Print("WEBSOCKET SERVER: ERROR: failed to write message", err)
			break
		}
	}
}
