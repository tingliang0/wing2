package main

import (
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

type connection struct {
	ws   *websocket.Conn
	recv chan []byte
	send chan []byte
	done chan bool
}

type message struct {
	name    string
	payload interface{}
}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			c.done <- true
			break
		}
		c.recv <- message
	}
}

func (c *connection) writer(msg []byte) {
	err := c.ws.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		c.done <- true
		return
	}
}

func get_info() (string, int) {
	name_list := [3]string{"Tom", "Jack", "Ken"}
	return name_list[rand.Intn(len(name_list))], rand.Intn(100-20) + 20
}

func (c *connection) run() {
	h.register <- c
	defer func() { h.unregister <- c }()
	name, age := get_info()
	agent := Agent{}
	agent.conn = c
	agent.Init(name, age)
	go c.reader()

DONE:
	for {
		select {
		case pkg := <-c.recv:
			agent.Dispatcher(pkg)
		case pkg := <-c.send:
			c.writer(pkg)
		case <-c.done: // EOF
			c.ws.Close()
			break DONE
		}
	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Error.Fatal("upgrade websocket fail", err)
		return
	}
	c := &connection{
		send: make(chan []byte, 256),
		recv: make(chan []byte, 256),
		done: make(chan bool),
		ws:   ws,
	}
	c.run()
}
