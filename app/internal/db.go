package internal

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go-chats/app/internal/models"
	"log"
)

type Db struct {
	config ParamsLocal
	user   models.User
}

func NewDb(config ParamsLocal) *Db {
	return &Db{
		config: config,
		user:   models.User{},
	}
}

func (d *Db) GetDb() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", d.config.Db.User+":"+d.config.Db.Password+"@("+d.config.Db.Serv+")/"+d.config.Db.Table)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return db, nil
}
