package main

import (
	"go-chat/config"
	"go-chat/database"
	"go-chat/router"
)

func main() {
	database.InitDatabase()

	r := router.InitRouter()

	err := r.Run(config.ServerSetting.HttpPort)

	if err != nil {
		panic(err)
	}
}
