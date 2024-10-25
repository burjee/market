package routes

import (
	"server_order/middleware"
	"server_order/service"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

func Setup(g *gin.Engine, producer sarama.SyncProducer) {
	service := service.New(producer)

	g.Use(middleware.Cors())

	g.Static("/assets", "./web/assets")
	g.StaticFile("/", "./web/index.html")

	api := g.Group("/api")
	api.POST("/buy", service.Buy)
	api.POST("/sell", service.Sell)
	api.POST("/test", service.Test)
}
