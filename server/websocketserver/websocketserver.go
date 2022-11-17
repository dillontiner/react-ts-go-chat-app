package websocketserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/entities"
	"server/persistence"

	"github.com/gorilla/websocket"
)

// var addr = flag.String("addr", "localhost:8080", "http service address")
type Server struct {
	API *persistence.Client
}

func NewServer(api *persistence.Client) Server {
	return Server{
		API: api,
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
	ID   string
	Conn *websocket.Conn
	Pool *Pool
	API  *persistence.Client
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
		message := Message{Type: messageType, Body: string(p)}
		fmt.Println("Persisting message")
		chatMessage := entities.Message{}
		err = json.Unmarshal([]byte(message.Body), &chatMessage)
		if err != nil {
			log.Print("WEBSOCKET SERVER: ERROR: failed to unmarshal message", err)
		}

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

		c.Pool.Broadcast <- Message{Type: message.Type, Body: string(createdMessageBytes)}
		fmt.Printf("Message Received: %+v\n", message)
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
			for client, _ := range pool.Clients {
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
		Conn: conn,
		Pool: pool,
		API:  s.API,
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
