package routes

import (
	"server_notify/libs"
	"server_notify/service"

	"github.com/gin-gonic/gin"
)

func Setup(g *gin.Engine, epoll *libs.Epoll, hub *libs.Hub) {
	g.Static("/assets", "./web/assets")
	g.StaticFile("/", "./web/index.html")

	service := service.New(epoll, hub)
	go service.StartReadWebsocket()

	g.GET("/ws", service.WebsocketHandler)
}
