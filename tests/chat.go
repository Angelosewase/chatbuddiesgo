package tests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Angelosewase/chatbuddiesgo/Handlers"
	"github.com/Angelosewase/chatbuddiesgo/helpers"
	"github.com/Angelosewase/chatbuddiesgo/internal/database"
)

func GetParticipatingChats(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			UserId string `json:"UserId"`
		}

		userID := parameters{}

		err := json.NewDecoder(r.Body).Decode(&userID)
		if err != nil {
			w.Write([]byte("not user found"))
			return
		}

		if userID.UserId == "" {
			return
		}

		chats, err := db.GetChatNotCreatedByTheUser(r.Context(), database.GetChatNotCreatedByTheUserParams{
			Createdby: sql.NullString{
				Valid:  true,
				String: userID.UserId,
			},
			CONCAT: userID.UserId,
		})

		if err != nil {
			return
		}

		type ChatResponse struct {
			ID           string    `json:"id"`
			CreatedBy    string    `json:"createdby"`
			LastMessage  string    `json:"lastMessage"`
			Participants string    `json:"participants"`
			CreatedAt    time.Time `json:"created_at"`
			IsGroupChat  bool      `json:"is_group_chat"`
		}

		chatsResponse := []ChatResponse{}

		for _, value := range chats {
			chatsResponse = append(chatsResponse, ChatResponse{
				ID:           value.ID,
				CreatedBy:    value.Createdby.String,
				LastMessage:  value.Lastmessage.String,
				Participants: value.Participants.String,
				CreatedAt:    value.CreatedAt,
				IsGroupChat:  value.IsGroupChat.Bool,
			})
		}

		helpers.RespondWithJson(w, r, chatsResponse, http.StatusAccepted)
	}
}

func TestGetReceiverIdFromChatID(mst Handlers.MsgHandlersStruct) http.HandlerFunc {

	type parameters struct {
		Chat_id   string ` json:"chat_id"`
		Sender_id string `json:"sender_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		Parameters := parameters{}
		if err := json.NewDecoder(r.Body).Decode(&Parameters); err != nil {
			fmt.Println("error decoding the r body")
		}

		_, err := mst.GetReceiverIdFromChatID(Parameters.Chat_id, Parameters.Sender_id)

		if err != nil {
			fmt.Printf("error getting the receiver id from the chat id: %v", err)
			return
		}

		// fmt.Printf("id: %v", id);
	}
}
