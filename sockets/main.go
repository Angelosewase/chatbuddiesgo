package sockets

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/Angelosewase/chatbuddiesgo/Handlers"
	"github.com/Angelosewase/chatbuddiesgo/helpers"
	"github.com/Angelosewase/chatbuddiesgo/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var Clients sync.Map // Thread-safe map

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Improve for production security
	},
}

// ServeSocketServer handles WebSocket connections
func ServeSocketServer(msgStruct *Handlers.MsgHandlersStruct) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		userId := queryParams.Get("UserId")

		if userId == "" {
			fmt.Println("missing userId")
			helpers.RespondWithError(w, r, http.StatusBadRequest, errors.New("no userId found"))
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			helpers.RespondWithError(w, r, http.StatusInternalServerError, fmt.Errorf("internal server error: %v", err))
			return
		}
		defer func() {
			Clients.Delete(userId)
			conn.Close()
		}()

		Clients.Store(userId, conn)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("error reading message: %s\n", err)
				conn.WriteJSON(map[string]string{"error": "failed to receive message"})
				break
			}

			var receivedMsg receivedMessage
			if err := json.Unmarshal(msg, &receivedMsg); err != nil {
				fmt.Printf("error unmarshalling received message: %s\n", err)
				conn.WriteJSON(map[string]string{"error": "invalid message format"})
				continue
			}

			receiverID, err := msgStruct.GetReceiverIdFromChatID(receivedMsg.ChatID, userId)
			if err != nil {
				conn.WriteJSON(map[string]string{"error": "failed to find receiver"})
				continue
			}

			SendMessage(receiverID, receivedMsg, msgStruct, receivedMsg.SenderID)
		}
	})
}

func RunMainSocketServer(msgStruct *Handlers.MsgHandlersStruct) {
	ServeSocketServer(msgStruct)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Printf("error starting the socket server: %s", err)
	}
}

func SendMessage(userId string, message receivedMessage, msgStruct *Handlers.MsgHandlersStruct, senderId string) {
	if wsConn, ok := Clients.Load(userId); ok {
		conn := wsConn.(*websocket.Conn)

		if err := msgStruct.AddTextMessageDB(database.AddTextMessageParams{
			ID:       uuid.NewString(),
			ChatID:   message.ChatID,
			SenderID: senderId,
			Content:  message.Content,
		}); err != nil {
			fmt.Println("failed to save message to DB")
			return
		}

		msgToSend, err := json.Marshal(message)
		if err != nil {
			fmt.Printf("error marshalling message: %v\n", err)
			return
		}

		if err := conn.WriteMessage(websocket.TextMessage, msgToSend); err != nil {
			fmt.Printf("error sending message to client %s: %s\n", userId, err)
		}
	} else {
		fmt.Printf("no client with userId %s found\n", userId)
	}
}

func BroadCast(message receivedMessage, msgStruct *Handlers.MsgHandlersStruct, userID string) {
	Clients.Range(func(key, value interface{}) bool {
		SendMessage(key.(string), message, msgStruct, userID)
		return true
	})
}

type receivedMessage struct {
	ChatID   string `json:"chat_id"`
	Content  string `json:"content"`
	SenderID string `json:"sender_id"`
}
