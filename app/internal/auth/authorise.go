package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func auth(accessToken string) (int, error) {
	if accessToken == "" {
		return -1, errors.New("missing jwt-accessToken")
	}

	tokenParts := strings.Split(accessToken, " ")

	if len(tokenParts) != 2 {
		return -1, errors.New("invalid jwt-accessToken")

	}

	//	jwt.ParseWithClaims()

	return 1, nil
}
