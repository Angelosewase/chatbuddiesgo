package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJwtToken() (string, error) {

	godotenv.Load()
	SECRET_JWT_KEY := os.Getenv("SECRET_JWT_KEY")

	if SECRET_JWT_KEY == "" {
		return "", errors.New("failed to load the env variables")
	}

	claims := Claims{
		Username: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "angel",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(SECRET_JWT_KEY))

	if err != nil {
		return "", fmt.Errorf("error generating the token %v", err)
	}
	return tokenString, nil
}

func ValidateToken(tokenstring string) (string, error) {
	godotenv.Load()

	token, err := jwt.ParseWithClaims(tokenstring, &Claims{}, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("SECRET_JWT_KEY")), nil

	})

	if err != nil {
		return "", fmt.Errorf("failed parsing jwt token %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Username, nil
	} else {
		return "", fmt.Errorf("invalide claims")
	}

}
