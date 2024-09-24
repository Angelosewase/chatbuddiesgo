package Handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Angelosewase/chatbuddiesgo/helpers"
	"github.com/Angelosewase/chatbuddiesgo/internal/database"
)

type MsgHandlersStruct struct {
	DB *database.Queries
}

func (msgstruct *MsgHandlersStruct) GetMessagesByChatId(chatId string) ([]database.Message, error) {
	return msgstruct.DB.GetMessagesByChatId(context.Background(), chatId)
}

func (msgstruct *MsgHandlersStruct) AddTextMessageDB(data database.AddTextMessageParams) error {
	if _, err := msgstruct.DB.GetUserByUserId(context.Background(), data.SenderID); err != nil {
		return fmt.Errorf("user not found %v", err)
	}

	if _, err := msgstruct.DB.GetChatByChatId(context.Background(), data.ChatID); err != nil {
		return fmt.Errorf("chat Not Found %v", err)
	}

	if data.Content == "" {
		return fmt.Errorf("invalid content")
	}

	if _, err := msgstruct.DB.AddTextMessage(context.Background(), data); err != nil {
		return fmt.Errorf("error inserting message: %v", err)
	}

	_, err := msgstruct.DB.UpdateLatestMessage(context.Background(), database.UpdateLatestMessageParams{
		Lastmessage: sql.NullString{
			Valid:  true,
			String: data.Content,
		},
		ID: data.ChatID,
	})

	if err != nil {
		return fmt.Errorf("error updating the latest message : %v", err)
	}

	return nil
}

func (msgStruct *MsgHandlersStruct) DeleteMessage(messageId string) error {
	if _, err := msgStruct.DB.DeleteMessage(context.Background(), messageId); err != nil {
		return fmt.Errorf("failed to delete message %v", err)
	}

	return nil
}

func (msgStruct *MsgHandlersStruct) GetReceiverIdFromChatID(chatId string, senderId string) (string, error) {
	if chatId == "" {
		return "", fmt.Errorf("invalid chatId")
	}
	chat, err := msgStruct.DB.GetChatByChatId(context.Background(), chatId)
	if err != nil {
		return "", fmt.Errorf("error fetching chat : %v", err)
	}
	chatParticipants := helpers.ParseDatabaseParticipantsString(chat.Participants.String)
	receiverId := helpers.RemoveLoggedInUserFromChatParticipantsArray(chatParticipants, senderId)

	return receiverId[0], nil
}

func GetMessagesHandler(MST MsgHandlersStruct) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Extract chatID from query parameters
		queryParams := r.URL.Query()
		chatID := queryParams.Get("chatID")

		if chatID == "" {
			http.Error(w, "chatID is required", http.StatusBadRequest)
			return
		}

		// Fetch messages by chat ID
		chats, err := MST.GetMessagesByChatId(chatID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		type chatResponse struct {
			Id           string    `json:"id"`
			Chat_id      string    `json:"chat_id"`
			Sender_id    string    `json:"sender_id"`
			Content      string    `json:"content"`
			Content_type string    `json:"content_type"`
			Created_at   time.Time `json:"created_at"`
			Updated_at   time.Time `json:"updated_at"`
			Is_deleted   bool      `json:"is_deleted"`
		}

		res := []chatResponse{}

		for _, value := range chats {
			res = append(res, chatResponse{
				Id:         value.ID,
				Chat_id:    value.ChatID,
				Sender_id:  value.SenderID,
				Content:    value.Content,
				Created_at: value.CreatedAt.Time,
				Updated_at: value.UpdatedAt.Time,
				Is_deleted: value.IsDeleted.Bool,
			})
		}

		// Respond with messages in JSON format
		helpers.RespondWithJson(w, r, res, 200)
	}
}
