package main

import (
	"fmt"
	"log"

	"github.com/rgamba/evtwebsocket"
)

func main() {
	for i := 0; i < 1000; i++ {
		go func() {
			c := evtwebsocket.Conn{
				OnConnected: func(w *evtwebsocket.Conn) {
					fmt.Println("Connected")
				},

				OnMessage: func(msg []byte, w *evtwebsocket.Conn) {
					log.Printf("Received uncatched message: %s\n", msg)
				},

				OnError: func(err error) {
					fmt.Printf("** ERROR **\n%s\n", err.Error())
				},

				MatchMsg: func(req, resp []byte) bool {
					return string(req) == string(resp)
				},
			}

			// Connect
			if err := c.Dial("ws://localhost:8080/ws", ""); err != nil {
				log.Fatal(err)
			}

			msg := evtwebsocket.Msg{
				Body: nil,
				Callback: func(resp []byte, w *evtwebsocket.Conn) {
					fmt.Printf("Got back: %s\n", resp)
				},
			}

			log.Printf("%s\n", msg.Body)
		}()
	}
	select {}
}
