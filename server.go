package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http"
)

type Message struct {
	Id   int
	Type string
	Body string
}

type Client struct {
	Connection *websocket.Conn
	Id         int
}

var clients = make([]*Client, 0)
var idx = 0

func TrackClient(client *Client) {
	client.Id = idx
	clients = append(clients, client)
	idx = idx + 1
}

func EchoServer(ws *websocket.Conn) {
	client := new(Client)
	client.Connection = ws

	TrackClient(client)

	fmt.Printf("Client %d connected\n", client.Id)

	var msg string
	for {
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			break
		}

		fmt.Printf("Client %d sent: %s\n", client.Id, msg)
		Broadcast(msg)
	}

	fmt.Printf("Client %d disconnected\n", client.Id)
}

func Broadcast(message string) {
	for _, val := range clients {
		websocket.Message.Send(val.Connection, message)
	}
}

func main() {
	http.Handle("/echo", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

	for {
		var s string
		fmt.Scan(s)
		Broadcast(s)
	}

}
