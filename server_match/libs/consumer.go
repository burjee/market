package libs

import (
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	redis_client *RedisClient
	matchmaker   *Matchmaker
}

func (Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}

			err := h.matchmaker.match(message)
			session.MarkMessage(message, "")

			if err != nil {
				log.Printf("message error: %s, ignore this message: %s", err.Error(), string(message.Value))
			}

		case <-session.Context().Done():
			return nil
		}
	}
}
