package Handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Angelosewase/chatbuddiesgo/helpers"
	"github.com/Angelosewase/chatbuddiesgo/internal/database"
	"github.com/google/uuid"
)

func GetChats(res http.ResponseWriter, req *http.Request, db *database.Queries) (int, error) {
	userId, err := helpers.GetUserIdFromToken(req)
	if err != nil {
		return http.StatusUnauthorized, fmt.Errorf("unauthorised %v", err)
	}

	chats, err := db.GetChatsByuserId(req.Context(), sql.NullString{Valid: true, String: userId})

	if err != nil {
		return 401, fmt.Errorf("error fetching chats: %v", err)
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

	helpers.RespondWithJson(res, req, chatsResponse, http.StatusAccepted)

	return 0, nil
}

func GetChatHandler(db *database.Queries) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		code, err := GetChats(res, req, db)
		if code == 0 {
			return
		}
		if err != nil {
			helpers.RespondWithError(res, req, code, fmt.Errorf("err:%v", err))
		}
	}
}

func CreateChatHandler(db *database.Queries) http.HandlerFunc {

	type participantsArray []string
	type parameters struct {
		Participants participantsArray `json:"participants"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		userId, err := helpers.GetUserIdFromToken(r)
		if err != nil {
			return
		}
		Parameters := parameters{}
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&Parameters)
		if err != nil {
			helpers.RespondWithError(w, r, 500, errors.New("failed to parse request body"))
			return
		}

		// Ensure the participants array is not empty
		if len(Parameters.Participants) == 0 {
			helpers.RespondWithError(w, r, 400, errors.New("participants cannot be empty"))
			return
		}

		var isgroupChat bool
		if len(Parameters.Participants) > 2 {
			isgroupChat = true
		} else {
			isgroupChat = false
		}

		// Logic for storing participants in a consistent format
		participantsString := helpers.ChatParticipantsToDatabaseChatParticipants(Parameters.Participants)

		_, err = db.CreateChat(r.Context(), database.CreateChatParams{
			ID:           uuid.NewString(),
			Createdby:    sql.NullString{Valid: true, String: userId},
			Lastmessage:  sql.NullString{Valid: true, String: ""},
			Participants: sql.NullString{Valid: true, String: participantsString},
			IsGroupChat:  sql.NullBool{Valid: true, Bool: isgroupChat},
			CreatedAt:    time.Now(),
		})

		if err != nil {
			// Check for unique constraint violation
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				helpers.RespondWithError(w, r, 409, errors.New("a chat with these participants already exists"))
			} else {
				helpers.RespondWithError(w, r, 500, errors.New("failed to create chat: "+err.Error()))
			}
			return
		}

		w.WriteHeader(200)
		w.Write([]byte("chat created successfully"))
	}
}

func DeleteChatHandler(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			ChatId string `json:"chatId"`
		}
		Parameters := parameters{}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&Parameters)
		err := db.DeleteChat(r.Context(), Parameters.ChatId)
		if err != nil {
			helpers.RespondWithError(w, r, 500, err)
			return
		}

		w.Write([]byte("chat deleted successfully"))
	}

}


func GetParticipatingChats(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId, err := helpers.GetUserIdFromToken(r)
	
		if err != nil {

			return
		}

		if userId == "" {
			return
		}

		chats, err := db.GetChatNotCreatedByTheUser(r.Context(), database.GetChatNotCreatedByTheUserParams{
			Createdby: sql.NullString{
				Valid:  true,
				String: userId,
			},
			CONCAT: userId,
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
