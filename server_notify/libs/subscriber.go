package libs

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/redis/rueidis"
)

type Subscriber struct {
	redis_client *RedisClient
	hub          *Hub
}

func NewSubscriber(redis_client *RedisClient, hub *Hub) *Subscriber {
	return &Subscriber{redis_client, hub}
}

func (s *Subscriber) Sub() {
	s.redis_client.Client.Receive(context.Background(), s.redis_client.B().Subscribe().Channel("broadcast").Build(), func(msg rueidis.PubSubMessage) {
		decoded_data, err := base64.StdEncoding.DecodeString(msg.Message)
		if err != nil {
			log.Printf("Failed to decode message: %v", err)
		}

		s.hub.Broadcast <- &BroadcastPack{HandleFunc: func(c *Client) func() {
			return func() { c.Write(decoded_data) }
		}}
	})
}
