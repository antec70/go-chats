package models

type Chat struct {
	Id         int    `db:"id"`
	Type       int    `db:"type"`
	CreateDate string `db:"create_date"`
	UserId1    int    `db:"user_id_1"`
	UserId2    int    `db:"user_id_2"`
	Status     int    `db:"status"`
}
