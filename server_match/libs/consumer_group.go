package libs

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func getConfig() *sarama.Config {
	version, err := sarama.ParseKafkaVersion(sarama.DefaultVersion.String())
	if err != nil {
		panic(err)
	}

	config := sarama.NewConfig()
	config.Version = version
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = time.Second
	return config
}

func NewConsumerGroup() sarama.ConsumerGroup {
	brokers := viper.GetStringSlice("kafka.brokers")
	group := viper.GetString("kafka.group")
	config := getConfig()

	retries := 0
	max_retries := 15

	consumer_group, err := sarama.NewConsumerGroup(brokers, group, config)
	for err != nil && retries < max_retries {
		retries += 1
		log.Printf("Failed to connect to Kafka (attempt %d/%d): %v", retries, max_retries, err)

		time.Sleep(12 * time.Second)
		consumer_group, err = sarama.NewConsumerGroup(brokers, group, config)
	}

	if err != nil {
		panic(err)
	}

	return consumer_group
}

func StartConsume(consumer_group sarama.ConsumerGroup, redis_client *RedisClient, matchmaker *Matchmaker) {
	topics := viper.GetStringSlice("kafka.topics")
	consumer := Consumer{redis_client, matchmaker}

	for {
		if err := consumer_group.Consume(context.Background(), topics, &consumer); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return
			}
			log.Printf("Error from consumer: %v", err)
		}
	}
}
