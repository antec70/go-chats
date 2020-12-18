package main

import (
	"encoding/json"
	"go-chats/app/internal"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	jsonFile, err := os.Open(internal.DefaultConfig)

	if err != nil {
		log.Fatal(err)
	}

	byteValue, er := ioutil.ReadAll(jsonFile)

	if er != nil {
		log.Fatal(er)
	}
	c := internal.ParamsLocal{}
	json.Unmarshal(byteValue, &c)

	internal.NewServer()

}
