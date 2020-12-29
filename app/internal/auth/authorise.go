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

func GetUser(accessToken string, userId int, conf config.ParamsLocal) (models.User, error) {
	d := db.NewDb(conf)
	database, er := d.GetDb()
	user := models.User{}
	if er != nil {
		fmt.Println(er)
		return models.User{}, er
	}

	if userId != 0 {
		er := database.Get(&user, "select * from user where id=?", userId)
		if er != nil {
			log.Fatal(er)
		}
		return user, nil
	}

	if accessToken == "" {
		return user, errors.New("missing jwt-accessToken")
	}

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(conf.Signature), nil
	})

	if err != nil {
		fmt.Println(err)
		return user, err

	}

	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		er := database.Get(&user, "select * from user where id=?", claims.UserId)
		if er != nil {
			log.Fatal(er)
		}
		return user, nil

	} else {
		fmt.Println(err)
		return user, nil
	}

}
