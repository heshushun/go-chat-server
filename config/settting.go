package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"time"
)

type server struct {
	AppMode  string
	HttpPort string
}

var ServerSetting = &server{}

type database struct {
	DbType     string
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

var DatabaseSetting = &database{}

type socket struct {
	MaxMessageSize  int64
	ReadBufferSize  int
	WriteBufferSize int
	WriteWait       time.Duration
	PongWait        time.Duration
}

var SocketSetting = &socket{}

var cfg *ini.File

// init 读取配置文件数据
func init() {
	var err error
	cfg, err = ini.Load("config/config.ini")

	if err != nil {
		fmt.Println("Open file error", err)
		os.Exit(1)
	}

	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("socket", SocketSetting)

}

func mapTo(s string, i interface{}) {
	err := cfg.Section(s).MapTo(i)

	if err != nil {
		log.Fatalf("%s cfg Load error", s)
	}
}
