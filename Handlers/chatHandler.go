package Handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Angelosewase/chatbuddiesgo/helpers"
	"github.com/Angelosewase/chatbuddiesgo/internal/database"
	"github.com/google/uuid"
)

func GetChats(res http.ResponseWriter, req *http.Request, db *database.Queries) (int, error) {
	type parameters struct {
		UserId string `json:"userId"`
	}

	Parameters := parameters{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&Parameters)
	if err != nil {
		return 500, fmt.Errorf("error parsing the request body %v", err)
	}

	chats, err := db.GetChatsByuserId(req.Context(), sql.NullString{Valid: true, String: Parameters.UserId})

	if err != nil {
		return 4001, fmt.Errorf("error fetching chats: %v", err)
	}

	helpers.RespondWithJson(res, req, chats, 400)

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
		UserId       string            `json:"id"`
		Participants participantsArray `json:"participants"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		Parameters := parameters{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&Parameters)
		if err != nil {
			helpers.RespondWithError(w, r, 500, errors.New("failed to parse request body"))
			return
		}
		var isgroupChat bool

		//other logic to check the number of participants

		if len(Parameters.Participants) > 2 {
			isgroupChat = true
		} else {
			isgroupChat = false
		}

		_, err = db.CreateChat(r.Context(), database.CreateChatParams{
			ID:          uuid.NewString(),
			Createdby:   sql.NullString{Valid: true, String: Parameters.UserId},
			Lastmessage: sql.NullString{Valid: true, String: ""},
			//this logic for adding participants is not good  , i will update it tomorrow
			Participants: sql.NullString{Valid: true, String: Parameters.Participants[0]},
			IsGroupChat:  sql.NullBool{Valid: true, Bool: isgroupChat},
			CreatedAt:    sql.NullTime{Valid: true, Time: time.Now()},
		})

		if err != nil {
			helpers.RespondWithError(w, r, 500, errors.New(err.Error()))

			return
		}
		//some logic to get the chat created and send it back to the user
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
