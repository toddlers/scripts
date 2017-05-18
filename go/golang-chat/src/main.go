package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

//connected clients
var clients = make(map[*websocket.Conn]bool)

// broadcast channel
var broadcast = make(chan Message)

// Configure the upgrader

var upgrader = websocket.Upgrader{}

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client

	clients[ws] = true

	for {

		var msg Message

		// Read in a new message as JSON and map it to a Message object

		if err := ws.ReadJSON(&msg); err != nil {
			log.Printf("error : %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		log.Println(msg)
		broadcast <- msg

	}
}

func handleMessages() {

	for {

		// Grab the next message from the broadcast channel

		msg := <-broadcast
		log.Println(msg)

		// Send it out to every client that is curretnly connected

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// Create a simple file server
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages

	go handleMessages()

	// Start the server on localhost port 8000 and log any errors

	log.Println("http server started on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
