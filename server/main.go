package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type User struct {
	userId   int
	userName string
	conn     *websocket.Conn
}

func newUser(ws *websocket.Conn, userName string) User {
	user := User{userId: int(uuid.New().ID()), userName: userName, conn: ws}
	return user
}

type Users struct {
	conns []User
	mu    sync.Mutex
}

func (u *Users) addUser(ws *websocket.Conn, userName string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.conns = append(u.conns, newUser(ws, userName))
}

func (u *Users) removeUser(ws *websocket.Conn) {
	u.mu.Lock()
	defer u.mu.Unlock()
	for i, user := range u.conns {
		if user.conn == ws {
			u.conns = append(u.conns[:i], u.conns[i+1:]...)
		}
	}
}

func (u *Users) broadcast(message string) {
	fmt.Println("Broadcasting to Users: ", u.conns)
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, user := range u.conns {
		if err := websocket.Message.Send(user.conn, message); err != nil {
			u.removeUser(user.conn)
			fmt.Println("error sending message")
		}
	}
}

type WsConnection struct {
	conn  *websocket.Conn
	users *Users
}

var globalUsers = Users{conns: make([]User, 0)}

func WSHandler(ws *websocket.Conn) {
	var message string
	wsConnection := WsConnection{conn: ws, users: &globalUsers}
	globalUsers.addUser(ws, faker.FirstName())
	for {
		if err := websocket.Message.Receive(ws, &message); err != nil {
			fmt.Println("error receiving message")
			globalUsers.removeUser(ws)
			return
		}
		fmt.Println("Received from socket: ", message)
		wsConnection.users.broadcast(message)
		for _, user := range wsConnection.users.conns {
			fmt.Println(user.conn.RemoteAddr())
		}
	}
}

func main() {
	handler := http.NewServeMux()
	handler.Handle("/ws", websocket.Handler(WSHandler))
	fmt.Println("Listening on port 3001")
	http.ListenAndServe(":3001", handler)
}
