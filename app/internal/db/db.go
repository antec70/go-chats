package db

import (
	"github.com/jmoiron/sqlx"
	"go-chats/app/internal/config"
	"log"
)

type Db struct {
	config config.ParamsLocal
}

func NewDb(config config.ParamsLocal) *Db {
	return &Db{
		config: config,
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
