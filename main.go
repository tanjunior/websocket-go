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
}

func index(rw http.ResponseWriter, r *http.Request) {
	// http.ServeFile(rw, r, "index.html")
	fmt.Fprintf(rw, "index")
}

func reader(conn *websocket.Conn) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		log.Println(string(message))

		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			break
		}
	}
}

func wsEndpoint(rw http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// fmt.Fprintf(rw, "wsEndpoint")

	ws, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	log.Println("Client Connected")

	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", index)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	// setupAPI()
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))

}
