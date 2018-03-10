package main

import "fmt"

var online int

type hub struct {
	// Registered clients
	connections map[*connection]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *connection

	// Unregister requests from clients
	unregister chan *connection
}

var h = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			online++
			fmt.Printf("online: %d\n", online)
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
				online--
				fmt.Printf("online: %d\n", online)
			}
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					delete(h.connections, c)
					close(c.send)
				}
			}
		}
	}
}
