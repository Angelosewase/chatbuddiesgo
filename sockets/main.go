package sockets

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Angelosewase/chatbuddiesgo/Handlers"
	"github.com/Angelosewase/chatbuddiesgo/helpers"
	"github.com/Angelosewase/chatbuddiesgo/internal/database"
	"github.com/google/uuid"
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

// ServeSocketServer handles WebSocket connections
func ServeSocketServer(msgStruct *Handlers.MsgHandlersStruct) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Get query params (UserId)
		queryParams := r.URL.Query()
		userId := queryParams.Get("UserId")

		if userId == "" {
			fmt.Println("missing userId")
			helpers.RespondWithError(w, r, http.StatusBadRequest, errors.New("no userId found"))
			return
		}

		// Upgrade the HTTP connection to a WebSocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			helpers.RespondWithError(w, r, http.StatusInternalServerError, fmt.Errorf("internal server error"))
			return
		}
		defer conn.Close()

		Clients[userId] = conn

		for {
			// Read incoming message
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("error reading message: %s\n", err)
				conn.WriteJSON(map[string]string{"error": "failed to receive message"})
				continue
			}

			// Decode the received message if it's JSON
			var receivedMsg receivedMessage
			if err := json.Unmarshal(msg, &receivedMsg); err != nil {
				fmt.Printf("error unmarshalling received message: %s\n", err)
				conn.WriteJSON(map[string]string{"error": "invalid message format"})
				continue
			}

			// Forward the message to the given client (ReceiverID)
			// fmt.Printf("Forwarding message from user %s to user %s: %+v\n", userId, receivedMsg.ReceiverID, receivedMsg)
			receiverID, err := msgStruct.GetReceiverIdFromChatID(receivedMsg.ChatID, userId)
			if err != nil {
				return
			}
			SendMessage(receiverID, receivedMsg, msgStruct)
		}
	})
}

// RunMainSocketServer starts the WebSocket server
func RunMainSocketServer(msgStruct *Handlers.MsgHandlersStruct) {
	ServeSocketServer(msgStruct)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Printf("error starting the socket server: %s", err)
	}
}

// SendMessage sends a message to a specific client
func SendMessage(userId string, message receivedMessage, msgstruct *Handlers.MsgHandlersStruct) {
	WSconn, ok := Clients[userId]
	if !ok {
		fmt.Printf("no client with userId %s found\n", userId)
		return
	}
	if err := msgstruct.AddTextMessageDB(database.AddTextMessageParams{
		ID:       uuid.NewString(),
		ChatID:   message.ChatID,
		SenderID: userId,
		Content:  message.Content,
	}); err != nil {
		fmt.Println("failed sending message")
	}

	msgTOSend, err := json.Marshal(message)
	if err != nil {
		return
	}

	if err := WSconn.WriteMessage(websocket.TextMessage, msgTOSend); err != nil {
		fmt.Printf("error sending message to client %s: %s\n", userId, err)
	}
}

// BroadCast broadcasts a message to all connected clients
func BroadCast(message receivedMessage, msgStruct *Handlers.MsgHandlersStruct) {
	for key := range Clients {
		SendMessage(key, message, msgStruct)
	}
}

// receivedMessage structure represents a message received from a user
type receivedMessage struct {
	ChatID   string `json:"chat_id"`
	Content  string `json:"content"`
	SenderID string `json:"sender_id"`
}
