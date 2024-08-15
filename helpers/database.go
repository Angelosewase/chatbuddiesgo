package helpers

import "strings"

func DatabaseChatParticipantstoChatParticipants(chats string) []string {
	chatsSlice := strings.Split(chats, ",")
	if len(chatsSlice) > 0 {
		return chatsSlice
	}

	return make([]string, 0)
}

func ChatParticipantsToDatabaseChatParticipants(chats []string) string {

	if len(chats) == 0 {
		return ""
	}
	var returnString string

	for _, value := range chats {
		returnString += ","
		returnString += value
	}

	return returnString
}
