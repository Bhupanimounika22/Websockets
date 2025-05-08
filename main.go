package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow connections from any origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents a single connected WebSocket client
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// ClientManager keeps track of all connected clients
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// NewClientManager creates a new ClientManager
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Start begins the ClientManager's main loop
func (manager *ClientManager) Start() {
	for {
		select {
		case client := <-manager.register:
			manager.clients[client] = true
			fmt.Println("New client connected")
		case client := <-manager.unregister:
			if _, ok := manager.clients[client]; ok {
				delete(manager.clients, client)
				close(client.send)
				fmt.Println("Client disconnected")
			}
		case message := <-manager.broadcast:
			for client := range manager.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(manager.clients, client)
				}
			}
		}
	}
}

// HandleClient handles WebSocket connections
func (manager *ClientManager) HandleClient(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	manager.register <- client

	// Start goroutines for reading and writing
	go client.readPump(manager)
	go client.writePump()
}

// readPump handles messages from the WebSocket connection
func (client *Client) readPump(manager *ClientManager) {
	defer func() {
		manager.unregister <- client
		client.conn.Close()
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			break
		}
		// Broadcast the message to all clients
		manager.broadcast <- message
	}
}

// writePump sends messages to the WebSocket connection
func (client *Client) writePump() {
	defer client.conn.Close()

	for message := range client.send {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}

func main() {
	manager := NewClientManager()
	go manager.Start()

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// WebSocket endpoint
	http.HandleFunc("/ws", manager.HandleClient)

	// Start the server
	fmt.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
