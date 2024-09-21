package Handlers

import (
	"context"
	"database/sql"
	"fmt"

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

	_,err := msgstruct.DB.UpdateLatestMessage(context.Background(), database.UpdateLatestMessageParams{
		Lastmessage: sql.NullString{
			Valid:  true,
			String: data.Content,
		},
		ID: data.ChatID,
	})

	if err != nil{
   return fmt.Errorf("error updating the latest message : %v",err)
	}

	return nil
}

func (msgStruct *MsgHandlersStruct) DeleteMessage(messageId string) error {
	if _, err := msgStruct.DB.DeleteMessage(context.Background(), messageId); err != nil {
		return fmt.Errorf("failed to delete message %v", err)
	}

	return nil
}


func(msgStruct *MsgHandlersStruct) GetReceiverIdFromChatID(chatId string, senderId string)(string, error){
	if chatId == ""{
		return "", fmt.Errorf("invalid chatId")
	}
    chat,err:= msgStruct.DB.GetChatByChatId(context.Background(),chatId)
	if err != nil{
		return"", fmt.Errorf("error fetching chat : %v",err)
	}
  chatParticipants:=helpers.ParseDatabaseParticipantsString(chat.Participants.String)
  receiverId :=helpers.RemoveLoggedInUserFromChatParticipantsArray(chatParticipants, senderId)
	

	return receiverId[0], nil
}
