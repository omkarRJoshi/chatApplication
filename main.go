package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Home page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// determines whether or not an incoming request from a different domain is allowed to connect
	// and if it isn’t they’ll be hit with a CORS error.
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// upgrade http connection to websocketsc
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client connected")

	reader(ws)
}

// define a reader which will listen for new messages
// being send to our websocket endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in the message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// print out the mesage for clarity
		fmt.Println("message from client", string(p))
		reply := []byte("ack from server for msg " + string(p))

		if err := conn.WriteMessage(1, reply); err != nil {
			log.Println(err)
			return
		}
	}
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Printf("Hello World from main")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
