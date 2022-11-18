package websocketserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/entities"
	"server/persistence"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// var addr = flag.String("addr", "localhost:8080", "http service address")
type Server struct {
	API       *persistence.Client
	VoteQueue chan entities.Vote
}

func NewServer(api *persistence.Client, voteQueue chan entities.Vote) Server {
	return Server{
		API:       api,
		VoteQueue: voteQueue,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

type Client struct {
	ID        string
	Conn      *websocket.Conn
	Pool      *Pool
	API       *persistence.Client
	VoteQueue chan entities.Vote
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		responseBody := ""

		chatMessage := entities.Message{}
		err = json.Unmarshal(p, &chatMessage)
		if err != nil {
			log.Print("WEBSOCKET SERVER: ERROR: failed to unmarshal chat message", err)
		} else if chatMessage.Body != "" { // MVP for handling different objects
			fmt.Println("Persisting chat message")

			createdMessage, err := c.API.CreateMessage(chatMessage)
			if err != nil {
				log.Print("WEBSOCKET SERVER: ERROR: failed to persist message", err)
				break
			}

			createdMessageBytes, err := json.Marshal(createdMessage)
			if err != nil {
				log.Print("WEBSOCKET SERVER: ERROR: failed to persist message", err)
				break
			}

			responseBody = string(createdMessageBytes)
		}

		vote := entities.Vote{}
		err = json.Unmarshal(p, &vote)
		if err != nil {
			log.Print("WEBSOCKET SERVER: ERROR: failed to unmarshal vote", err)
		} else if vote.VoterUUID != uuid.Nil { // MVP for handling different objects
			log.Print("WEBSOCKET SERVER: got vote", vote)
			// enqueues message to be persisted
			c.VoteQueue <- vote
			responseBody = string(p)
		}

		if responseBody == "" {
			// failed to parse objects
			log.Print("WEBSOCKET SERVER: ERROR: failed to unmarshal object", err)
		}

		message := Message{Type: messageType, Body: responseBody}
		c.Pool.Broadcast <- message
		log.Printf("Message Received: %+v\n", message)
	}
}

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

func (s *Server) ServeWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &Client{
		Conn:      conn,
		Pool:      pool,
		API:       s.API,
		VoteQueue: s.VoteQueue,
	}

	pool.Register <- client
	client.Read()
}

func (s *Server) SetupRoutes(path string) {
	pool := NewPool()
	go pool.Start()

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		s.ServeWs(pool, w, r)
	})
}
