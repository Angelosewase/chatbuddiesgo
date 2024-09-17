package Handlers

// import (
// 	"encoding/json"
// 	"net/http"
// 	"time"

// 	"github.com/Angelosewase/chatbuddiesgo/helpers"
// 	"github.com/Angelosewase/chatbuddiesgo/socket"
// )

// func SendMessage(userId string, srvConfig socket.Server, message interface{}) {
// 	srvConfig.EmitSocket(userId, "newMessage", message)
// }

// func NewMessageHandler(res http.ResponseWriter, req *http.Request) {
// 	_, err := helpers.GetUserIdFromToken(req)
// 	if err != nil {
// 		return
// 	}
// 	type Parameters struct {
// 		ReceiverId string    `json:"receiverId"`
// 		Message    string    `json:"message"`
// 		CreatedAt  time.Time `json:"createdAt"`
// 	}

// 	parameters := Parameters{}
// 	decoder := json.NewDecoder(req.Body)
// 	err = decoder.Decode(&parameters)
// 	if err != nil {
// 		return
// 	}
// }

// func DeleteMessageHandler(r http.ResponseWriter, w *http.Request) {
// 	_, err := helpers.GetUserIdFromToken(w)
// 	if err != nil {
// 		return
// 	}
// 	type Parameters struct {
// 		MessageId string `json:"messageId"`
// 	}
// 	parameters := Parameters{}
// 	decoder := json.NewDecoder(w.Body)
// 	err = decoder.Decode(&parameters)
// 	if err != nil {
// 		return
// 	}

// }
