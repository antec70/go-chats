package main

import (
	"encoding/json"
	"go-chats/app/internal"
	"go-chats/app/internal/config"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	jsonFile, err := os.Open(config.DefaultConfig)

	if err != nil {
		log.Fatal(err)
	}

	byteValue, er := ioutil.ReadAll(jsonFile)

	if er != nil {
		log.Fatal(er)
	}
	c := config.ParamsLocal{}
	json.Unmarshal(byteValue, &c)

	w := internal.NewWsServ(c)
	w.NewServer()

}
