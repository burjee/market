package main

import (
	"log"
	"os"
	_ "server_match/config"
	"server_match/libs"
	"time"

	"github.com/rcrowley/go-metrics"
)

func main() {
	redis_client := libs.NewRedisClient()
	defer redis_client.Close()

	createOrderIndex(redis_client)
	createPriceQuantityIndex(redis_client)
	generateTestOrder(redis_client)

	consumer_group := libs.NewConsumerGroup()
	defer consumer_group.Close()

	matched_metrics := metrics.NewMeter()
	metrics.Register("matched-meter", matched_metrics)
	go metrics.Log(metrics.DefaultRegistry, time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	matchmaker := libs.NewMatchmaker(redis_client, matched_metrics)

	log.Println("start matchmaker")
	libs.StartConsume(consumer_group, redis_client, matchmaker)
}
