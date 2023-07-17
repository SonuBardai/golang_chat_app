package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

func SayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func WSHandler(ws *websocket.Conn) {
	var message string
	for {
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			fmt.Println("error receiving message")
			return
		}
		fmt.Println("Received from socket: ", message)
		err = websocket.Message.Send(ws, message)
		if err != nil {
			fmt.Println("error sending message")
			return
		}
	}
}

func main() {
	handler := http.NewServeMux()
	handler.HandleFunc("/", SayHello)
	handler.Handle("/ws", websocket.Handler(WSHandler))
	http.ListenAndServe(":3001", handler)
}
