package helpers

import (
	"strings"
)

func ParseDatabaseParticipantsString(participantsString string) []string {

	participantsArray := strings.Split(participantsString, ",")
	return participantsArray
}

func RemoveLoggedInUserFromChatParticipantsArray(participantsArray []string, userId string) []string {
	if userId == "" {
		return []string{} 
	}

	var filteredParticipants []string
	for _, participant := range participantsArray {
		if participant != userId {
			filteredParticipants = append(filteredParticipants, participant)
		}
	}
	return filteredParticipants
}
