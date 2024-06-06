package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var receiver *websocket.Conn
var sender *websocket.Conn

var send_msg []byte = make([]byte, 0)

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", rootHandler)

	panic(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "Ola, amigo")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { 
		return true 
	}

	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	if r.Header.Get("Origin") == "Receiver" {
		receiver = conn
		fmt.Println("Receiver connected")
	} else if r.Header.Get("Origin") == "Sender" {
		sender = conn
		fmt.Println("Sender connected")
		go echo(sender)
	}	
}

func echo(conn *websocket.Conn) {
	fmt.Println("Echo started")
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		

		// fmt.Println("Received message type:", t)
		// fmt.Println("Received message:", string(msg))
		if t == websocket.BinaryMessage {
			// fmt.Println("Binary message received")
			if receiver != nil {
				send_msg = append(send_msg, msg...)
				if len(send_msg) == 500 {
					
					err_send := receiver.WriteMessage(websocket.BinaryMessage, send_msg)
					if err_send != nil {
						fmt.Println(err_send)
					}

					send_msg = make([]byte, 0)
				}
			}
		}
	}
}