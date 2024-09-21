package helpers

import "strings"

// parseDatabaseParticipantsString splits a comma-separated participants string into a slice of strings.
func ParseDatabaseParticipantsString(participantsString string) []string {
	// Split the string by commas and return the resulting slice
	participantsArray := strings.Split(participantsString, ",")
	return participantsArray
}

// removeLoggedInUserFromChatParticipantsArray removes the logged-in user from the participants slice.
func RemoveLoggedInUserFromChatParticipantsArray(participantsArray []string, userId string) []string {
	if userId == "" {
		return []string{} // Return an empty slice if the userId is empty
	}

	// Filter out the userId from the participantsArray
	var filteredParticipants []string
	for _, participant := range participantsArray {
		if participant != userId {
			filteredParticipants = append(filteredParticipants, participant)
		}
	}
	return filteredParticipants
}
