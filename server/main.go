package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

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

var ErrUserExists = errors.New("user with the same name already exists")

func (u *Users) addUser(ws *websocket.Conn, userName string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, user := range u.conns {
		if user.userName == userName {
			return ErrUserExists
		}
	}
	u.conns = append(u.conns, newUser(ws, userName))
	fmt.Printf("User added: %v\n", ws.RemoteAddr())
	return nil
}

func (u *Users) removeUser(ws *websocket.Conn) {
	u.mu.Lock()
	defer u.mu.Unlock()
	for i, user := range u.conns {
		if user.conn == ws {
			fmt.Printf("User removed: %v (%v)\n", user.userName, user.conn.RemoteAddr())
			u.conns = append(u.conns[:i], u.conns[i+1:]...)
		}
	}
}

func (u *Users) broadcast(message interface{}) {
	fmt.Println("Broadcasting to Users: ", u.conns)
	u.mu.Lock()
	defer u.mu.Unlock()
	var msg []byte
	var err error
	switch message.(type) {
	case NewUserMessage:
		msg, err = json.Marshal(message)
	case BroadcastMessage:
		msg, err = json.Marshal(message)
	case ErrorMessage:
		msg, err = json.Marshal(message)
	default:
		err = errors.New("unknown message type received")
	}
	if err != nil {
		panic(fmt.Sprintf("failed to broadcast message %v", err))
	}
	for _, user := range u.conns {
		if err := websocket.Message.Send(user.conn, string(msg)); err != nil {
			u.removeUser(user.conn)
			fmt.Printf("error sending message %v", err)
		}
	}
}

type WsConnection struct {
	conn  *websocket.Conn
	users *Users
}

var globalUsers = Users{conns: make([]User, 0)}

type MessageType string

const (
	NewUserType          MessageType = "newUser"
	UserLeftType         MessageType = "userLeft"
	BroadcastMessageType MessageType = "broadcastMessage"
	ErrorType            MessageType = "error"
)

type BaseMessage struct {
	MessageType MessageType `json:"messageType"`
}

type BroadcastMessage struct {
	BaseMessage
	Username string `json:"username"`
	Content  string `json:"content"`
}

type NewUserMessage struct {
	BaseMessage
	Username string `json:"username"`
}

type UserLeftMessage struct {
	BaseMessage
	Username string `json:"username"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

func WSHandler(ws *websocket.Conn) {
	fmt.Printf("New connection: %v\n", ws.RemoteAddr())
	var message string
	wsConnection := WsConnection{conn: ws, users: &globalUsers}
	for {
		if err := websocket.Message.Receive(ws, &message); err != nil {
			fmt.Printf("error receiving message %v\n", err)
			globalUsers.removeUser(ws)
			return
		}
		var parsedMessage BaseMessage
		if err := json.Unmarshal([]byte(message), &parsedMessage); err != nil {
			fmt.Printf("Failed to parse received message %v\n", err)
			globalUsers.removeUser(ws)
			return
		}
		switch parsedMessage.MessageType {
		case NewUserType:
			var newUserMessage NewUserMessage
			if err := json.Unmarshal([]byte(message), &newUserMessage); err != nil {
				fmt.Printf("Failed to parse newUser message %v\n", err)
				globalUsers.removeUser(ws)
				return
			}
			if err := globalUsers.addUser(ws, newUserMessage.Username); err != nil {
				if errors.Is(err, ErrUserExists) {
					fmt.Println("Error adding user:", err)
					errorMessage := ErrorMessage{Error: "Error: " + err.Error()}
					errorJson, _ := json.Marshal(errorMessage)
					if err := websocket.Message.Send(ws, string(errorJson)); err != nil {
						globalUsers.removeUser(ws)
						fmt.Println("error sending message")
					}
					continue
				} else {
					fmt.Printf("Failed to parse newUser message %v\n", err)
					globalUsers.removeUser(ws)
					break
				}
			}
			wsConnection.users.broadcast(newUserMessage)
		case BroadcastMessageType:
			var broadcastMessage BroadcastMessage
			if err := json.Unmarshal([]byte(message), &broadcastMessage); err != nil {
				fmt.Printf("Failed to parse newUser message %v\n", err)
				globalUsers.removeUser(ws)
				return
			}
			wsConnection.users.broadcast(broadcastMessage)
			continue
		}
		fmt.Println("Received from socket: ", message)
		fmt.Println("Parsed message: ", parsedMessage)
	}
}

func main() {
	handler := http.NewServeMux()
	handler.Handle("/ws", websocket.Handler(WSHandler))
	fmt.Println("Listening on port 3001")
	http.ListenAndServe(":3001", handler)
}
