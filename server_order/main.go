package main

import (
	"log"
	_ "server_order/config"
	"server_order/libs"
	"server_order/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	producer := libs.NewProducer()
	defer producer.Close()

	g := gin.Default()
	routes.Setup(g, producer)

	log.Println("start server http://0.0.0.0:8000")
	log.Fatal(g.Run(":8000"))
}
