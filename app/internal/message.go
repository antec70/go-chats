package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"go-chats/app/internal/config"
	"go-chats/app/internal/db"
	"go-chats/app/models"
	"net/http"
	"strconv"
	"time"
)

type Msg struct {
	config config.ParamsLocal
}

const typeMessage = "1"
const topicMessage = "message-"

func NewMessage(conf config.ParamsLocal) *Msg {
	return &Msg{
		config: conf,
	}
}

type data struct {
	Id        int
	Text      string
	ChatId    int
	ChatName  string
	ChatImage string
	ChatType  int
	UserId    int
	Date      string
	Read      bool
}

func (m *Msg) save(data map[string]interface{}, userId int) (models.Message, error) {

	d := db.NewDb(m.config)
	database, er := d.GetDb()
	if data["text"].(string) == "" {
		fmt.Println("empty text")
		return models.Message{}, errors.New("empty text")
	}

	message := models.Message{
		ChatId:     int(data["chatId"].(float64)),
		UserId:     userId,
		Text:       data["text"].(string),
		Status:     1,
		Read:       false,
		CreateDate: time.Now().UTC().Format("2006-01-02 15:04:05"),
	}
	if er != nil {
		return message, er
	}
	sql, er := database.Prepare("INSERT INTO message (chat_id, user_id, text, status, `read`, create_date) VALUES(?,?,?,?,?,?)")
	if er != nil {
		fmt.Println(er)
		return message, er
	}
	_, err := sql.Exec(message.ChatId, message.UserId, message.Text, message.Status, message.Read, message.CreateDate)

	if err != nil {
		fmt.Println(err)
		return message, err
	}
	return message, nil

}

func (m *Msg) publish(message models.Message, user models.User) error {
	d := data{
		Text:      message.Text,
		ChatId:    message.ChatId,
		ChatName:  user.Name + " " + user.SurName,
		ChatImage: user.Photo,
		ChatType:  0,
		Date:      time.Now().UTC().Format("2006-01-02 15:04:05"),
		Read:      false,
	}

	toUserId, er := m.getUserByChat(message.UserId, message.ChatId)
	if er != nil {
		fmt.Println(er)
		return er
	}
	socketio.NewBroadcast().Send("messages-to-"+strconv.Itoa(toUserId), "new-message", d)

	if user.PushMessage == true {
		n := models.JsonRequest{
			Push: &models.Push{
				Type:  typeMessage,
				Image: user.Photo,
				Title: user.Name + " " + user.SurName,
				Body:  message.Text,
			},
			Topic:  topicMessage + strconv.Itoa(toUserId),
			ChatId: strconv.Itoa(message.ChatId),
		}

		b, er := json.Marshal(n)

		if er != nil {
			fmt.Println(er)
		}
		go m.sendPush(b)

	}

	return nil
}

func (m *Msg) getUserByChat(userId int, chatId int) (int, error) {
	d := db.NewDb(m.config)
	database, er := d.GetDb()
	chat := models.Chat{}
	if er != nil {
		fmt.Println(er)
		return 0, er
	}
	err := database.Get(&chat, "SELECT * FROM chat WHERE id=?", chatId)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	if userId == chat.UserId1 {
		return chat.UserId2, nil
	} else {
		return chat.UserId1, nil
	}

}

func (m *Msg) readMsg(chatId int, userId int) error {
	d := db.NewDb(m.config)
	database, er := d.GetDb()
	sql := "UPDATE message SET `read` = true where chat_id = " + strconv.Itoa(chatId) + " AND user_id != " + strconv.Itoa(userId)
	if er != nil {
		fmt.Println(er)
		return er
	}

	_, err := database.Exec(sql)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *Msg) sendPush(data []byte) {

	req, err := http.NewRequest("POST", m.config.Pusher.Url, bytes.NewBuffer(data))

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("X-Api-Key", m.config.Pusher.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, er := client.Do(req)
	if er != nil {
		fmt.Println(err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Pusher: ", resp.StatusCode)
	}

}
