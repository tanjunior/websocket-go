// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

type Message struct {
	UserName string `json:"userName"`
	Message  string `json:"message"`
}

func (h *Hub) run() {
	log.Printf("[hub.run] clients %v", h.clients)
	for {
		select {
		case client := <-h.register:
			log.Printf("[h.register] registering client to hub: %v\n%v", client, h.clients)
			h.clients[client] = true
			log.Printf("[h.register] Client registered to hub: %v\n%v", client, h.clients)
		case client := <-h.unregister:
			log.Printf("[h.unregister] client: %v", client)
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			m := Message{Message: "disconnected"}
			b, err := Marshal(m)
			if err == nil {
				for client := range h.clients {
					client.send <- b
				}
			}
		case message := <-h.broadcast:
			log.Printf("[h.broadcast] message: %v", string(message[:]))
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
