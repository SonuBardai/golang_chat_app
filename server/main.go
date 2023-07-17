package main

import (
	"fmt"
	"net/http"

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
}

func (u *Users) addUser(ws *websocket.Conn, userName string) {
	u.conns = append(u.conns, newUser(ws, userName))
}

func (u *Users) removeUser(d *User) {
	for i, user := range u.conns {
		if user.conn == d.conn {
			u.conns = append(u.conns[:i], u.conns[i+1:]...)
		}
	}
}

func (u *Users) broadcast(message string) {
	fmt.Println("Broadcasting to Users: ", u.conns)
	for _, user := range u.conns {
		if err := websocket.Message.Send(user.conn, message); err != nil {
			fmt.Println("error sending message")
			u.removeUser(&user)
		}
	}
}

type WsServer struct {
	conn  *websocket.Conn
	users *Users
}

func (s *WsServer) readLoop() {
	var message string
	for {
		err := websocket.Message.Receive(s.conn, &message)
		if err != nil {
			fmt.Println("error receiving message")
			return
		}
		fmt.Println("Received from socket: ", message)
		s.users.broadcast(message)
	}
}

var globalUsers = Users{conns: make([]User, 0)}

func WSHandler(ws *websocket.Conn) {
	wsServer := WsServer{conn: ws, users: &globalUsers}
	globalUsers.addUser(ws, faker.FirstName())
	wsServer.readLoop()
}

func main() {
	handler := http.NewServeMux()
	handler.Handle("/ws", websocket.Handler(WSHandler))
	fmt.Println("Listening on port 3001")
	http.ListenAndServe(":3001", handler)
}
