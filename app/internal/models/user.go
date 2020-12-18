package models

import "database/sql"

type User struct {
	ID            int            `db:"id"`
	Name          string         `db:"name"`
	SurName       string         `db:"surname"`
	Phone         string         `db:"phone"`
	Patronymic    sql.NullString `db:"patronymic"`
	Status        sql.NullBool   `db:"status"`
	PasswordHash  sql.NullString `db:"password_hash"`
	CreateDate    sql.NullString `db:"create_date"`
	IsReg         sql.NullBool   `db:"is_reg"`
	SmsHash       sql.NullString `db:"sms_hash"`
	Birthday      sql.NullString `db:"birthday"`
	Gender        sql.NullString `db:"gender"`
	Nickname      sql.NullString `db:"nickname"`
	PhotoOriginal sql.NullString `db:"photo_original"`
	Approved      sql.NullBool   `db:"approved"`
	ActiveDate    sql.NullString `db:"active_date"`
	Timezone      sql.NullString `db:"timezone"`
	NewPhone      sql.NullString `db:"new_phone"`
	PushCenter    sql.NullBool   `db:"push_center"`
	PushMessage   sql.NullBool   `db:"push_message"`
	Photo         sql.NullString `db:"photo"`
}
