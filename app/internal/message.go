package internal

import (
	"fmt"
	"go-chats/app/internal/config"
	"go-chats/app/internal/db"
	"go-chats/app/models"
	"time"
)

func Save(data map[string]interface{}, user int, conf config.ParamsLocal) (models.Message, error) {

	d := db.NewDb(conf)
	database, er := d.GetDb()
	message := models.Message{
		ChatId:     int(data["chatId"].(float64)),
		UserId:     user,
		Text:       data["text"].(string),
		Status:     1,
		Read:       true,
		CreateDate: time.Now().UTC().Format("2006-01-02 15:04:05"),
	}
	if er != nil {
		return message, er
	}
	//todo создается 3 записи в бд иногда да а иногда нет хз че за прикол
	sql, er := database.Prepare("INSERT INTO message (chat_id, user_id, text, status, `read`, create_date) VALUES(?,?,?,?,?,?)")
	if er != nil {
		fmt.Println(er)
	}
	row, err := sql.Exec(message.ChatId, message.UserId, message.Text, message.Status, message.Read, message.CreateDate)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(row)
	return message, nil

}
