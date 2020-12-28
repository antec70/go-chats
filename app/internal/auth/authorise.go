package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go-chats/app/internal/config"
	"go-chats/app/internal/db"
	"go-chats/app/models"

	"log"
)

type tokenClaims struct {
	Type   string `json:"type"`
	UserId int    `json:"usr"`
	jwt.StandardClaims
}

const Signature = "mpYYjH0pF4pLttlKX9va9iimx9i7QI"

func GetUser(accessToken string, conf config.ParamsLocal) (models.User, error) {
	d := db.NewDb(conf)
	database, er := d.GetDb()
	user := models.User{}
	if er != nil {
		fmt.Println(er)
	}

	if accessToken == "" {
		return user, errors.New("missing jwt-accessToken")
	}

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(Signature), nil
	})

	if err != nil {
		fmt.Println(err)
		return user, err

	}

	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		er := database.Select(&user, "select * from users where id=?", claims.UserId)
		if er != nil {
			log.Fatal(er)
		}
		return user, nil

	} else {
		fmt.Println(err)
		return user, nil
	}

}
