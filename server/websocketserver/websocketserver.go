package websocketserver

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// var addr = flag.String("addr", "localhost:8080", "http service address")

func checkOrigin(r *http.Request) bool {
	// origin := r.Header.Get("Origin")
	// origin == "http://localhost:3000" || origin == "http://localhost:4001" // TODO: env var this
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: checkOrigin,
} // use default options

func Echo(w http.ResponseWriter, r *http.Request) {
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
			log.Println("WEBSOCKET SERVER: read:", err)
			break
		}
		log.Printf("WEBSOCKET SERVER: received: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
