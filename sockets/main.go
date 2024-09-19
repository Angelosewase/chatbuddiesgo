package sockets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/Angelosewase/chatbuddiesgo/helpers"
	"github.com/gorilla/websocket"
)

var Clients = make(map[string]*websocket.Conn)
var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeSocketServer() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		type Parameters struct {
			UserId string `json:"UserId"`
		}

		parameters := Parameters{}
		err := json.NewDecoder(r.Body).Decode(&parameters)
		if err != nil {
			helpers.RespondWithError(w, r, http.StatusBadRequest, fmt.Errorf("invalid request body"))
		}

		conn, err := upgrader.Upgrade(w, r, nil)

		defer func() {
			conn.Close()

		}()
		if err != nil {
			helpers.RespondWithError(w, r, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		}
		Clients[parameters.UserId] = conn

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				conn.WriteJSON("failed to receive message")
				continue
			}

			conn.WriteJSON(msg)
		}
	})
}

func RunMainSocketServer() {
	ServeSocketServer()
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Printf("error starting the socket server %d", err)
	}
}


func SendMessage(userId string, message interface{}){
 WSconn := Clients[userId]
 error :=WSconn.WriteJSON(message);
 if error !=nil{
  fmt.Printf("error sending message to client %d error: %d",userId, error)
 }
}

func BroadCast(message interface{}){
  for key, _ := range Clients{
    SendMessage(key, message)
  }
}


//group chat 
