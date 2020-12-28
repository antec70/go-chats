package models

type Message struct {
	Id         int    `db:"id"`
	ChatId     int    `db:"chat_id"`
	UserId     int    `db:"user_id_"`
	Text       string `db:"text"`
	Status     int    `db:"status"`
	Read       bool   `db:"read"`
	CreateDate string `db:"create_date"`
}
