package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	// os.Setenv("PORT", "8080")
	setupRoutes()

	port := os.Getenv("PORT")
	addr := ":" + port

	log.Println("server is listening on port", port)
	log.Fatal(http.ListenAndServe(addr, nil))

}
