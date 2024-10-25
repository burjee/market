package libs

import (
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
	config.Producer.Idempotent = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 50
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Net.MaxOpenRequests = 1
	return config
}

func NewProducer() sarama.SyncProducer {
	brokers := viper.GetStringSlice("kafka.brokers")
	config := getConfig()

	retries := 0
	max_retries := 15

	producer, err := sarama.NewSyncProducer(brokers, config)
	for err != nil && retries < max_retries {
		retries += 1
		log.Printf("Failed to connect to Kafka (attempt %d/%d): %v", retries, max_retries, err)

		time.Sleep(12 * time.Second)
		producer, err = sarama.NewSyncProducer(brokers, config)
	}

	if err != nil {
		panic(err)
	}

	return producer
}
