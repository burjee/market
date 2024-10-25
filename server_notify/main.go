package main

import (
	"log"
	_ "server_notify/config"
	"server_notify/libs"
	"server_notify/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	redis_client := libs.NewRedisClient()
	defer redis_client.Close()

	go_pool := libs.NewGoPool()
	defer go_pool.Release()

	epoll := libs.NewEpoll()
	defer epoll.Close()

	hub := libs.NewHub(epoll, go_pool)
	go hub.Run()
	defer hub.Close()

	subscriber := libs.NewSubscriber(redis_client, hub)
	go subscriber.Sub()

	order_list := libs.NewOrderList(redis_client)
	go order_list.Broadcast()

	g := gin.Default()
	routes.Setup(g, epoll, hub)

	log.Printf("start server ws://0.0.0.0:8001")
	log.Fatal(g.Run(":8001"))
}
