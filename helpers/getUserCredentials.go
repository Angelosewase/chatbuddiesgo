package helpers

import (
	"fmt"
	"net/http"

	"github.com/Angelosewase/chatbuddiesgo/internal/auth"
)

func GetUserIdFromToken(req *http.Request) (string, error) {
	jwtCookie, err := req.Cookie("jwt")
	if err != nil {
		return "", fmt.Errorf("user not authorised %v", err)
	}
	userId, err := auth.ValidateToken(jwtCookie.Value)

	if err != nil {
		return "", fmt.Errorf("invalid token %v", err)
	}

	return userId, nil
}
