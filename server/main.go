package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

var users []*websocket.Conn

func SayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func WSHandler(ws *websocket.Conn) {
	var message string
	users = append(users, ws)
	for {
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			fmt.Println("error receiving message")
			return
		}
		fmt.Println("Received from socket: ", message)
		for _, user := range users {
			err = websocket.Message.Send(user, message)
			if err != nil {
				fmt.Println("error sending message")
				return
			}
		}
	}
}

func main() {
	handler := http.NewServeMux()
	handler.HandleFunc("/", SayHello)
	handler.Handle("/ws", websocket.Handler(WSHandler))
	fmt.Println("Listening on port 3001")
	http.ListenAndServe(":3001", handler)
}
