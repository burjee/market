package libs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/rueidis"
	"github.com/spf13/viper"
)

func initClient() rueidis.Client {
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")

	addr := fmt.Sprintf("%s:%s", host, port)

	retries := 0
	max_retries := 15

	client, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{addr}})
	for err != nil && retries < max_retries {
		retries += 1
		log.Printf("Failed to connect to redis (attempt %d/%d): %v", retries, max_retries, err)

		time.Sleep(12 * time.Second)
		client, err = rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{addr}})
	}

	if err != nil {
		panic(err)
	}

	return client
}

type RedisClient struct {
	Client rueidis.Client
}

func NewRedisClient() *RedisClient {
	client := initClient()
	return &RedisClient{client}
}

func (r *RedisClient) Do(ctx context.Context, cmd rueidis.Completed) rueidis.RedisResult {
	return r.Client.Do(ctx, cmd)
}

func (r *RedisClient) B() rueidis.Builder {
	return r.Client.B()
}

func (r *RedisClient) Reconnect() {
	r.Close()
	r.Client = initClient()
}

func (r *RedisClient) Close() {
	r.Client.Close()
}
