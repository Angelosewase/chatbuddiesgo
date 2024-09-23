package tests

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Angelosewase/chatbuddiesgo/helpers"
	"github.com/Angelosewase/chatbuddiesgo/internal/database"
)

func GetParticipatingChats(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct{
			UserId string `json:"UserId"`
		}

		userID :=parameters{}
	

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
