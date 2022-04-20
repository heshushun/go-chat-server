package router

import (
	"go-chat/config"
	"go-chat/middlerware"
	"go-chat/router/v1/socket"
	"go-chat/router/v1/upload"
	"go-chat/router/v1/user"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化相关路由
func InitRouter() *gin.Engine {
	gin.SetMode(config.ServerSetting.AppMode)

	router := gin.New()

	router.Use(middlerware.Cors())

	router.Use(middlerware.Logger())

	hub := socket.NewHub()

	go hub.Run()

	router.GET("/ws", func(c *gin.Context) {
		socket.ServeWs(hub, c)
	})

	router.POST("/upload", upload.PushFile)

	router.GET("/download", upload.DownLoadFile)

	router.POST("/sign_up", user.Registry)

	router.POST("/sign_in", user.AuthHandler)

	return router
}
