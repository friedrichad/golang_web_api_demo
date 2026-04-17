package main

import (
	config "github.com/friedrichad/golang_web_api_demo/configs"
	db "github.com/friedrichad/golang_web_api_demo/db"
)

func main() {
	cfg, err := config.LoadConfig("configs/config.json")
	if err != nil {
		panic(err)
	}

	db.ConnectDB(cfg)

}
