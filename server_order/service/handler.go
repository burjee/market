package service

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"server_order/structure"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func (s *service) handleOrder(req *structure.Order) error {
	message, err := json.Marshal(req)
	if err != nil {
		log.Println("SERVER ERROR: json.Marshal")
		return err
	}

	if err := s.produceOrder(message); err != nil {
		log.Println("SERVER ERROR: producer.SendMessage")
		return err
	}

	return nil
}

func (s *service) handleTest() error {
	buy_orders, sell_orders, err := s.randomOrders(100)
	if err != nil {
		return err
	}

	for _, buy_order := range buy_orders {
		if err := s.produceOrder(buy_order); err != nil {
			log.Println("SERVER ERROR: producer.SendMessage")
			return err
		}
	}

	for _, sell_order := range sell_orders {
		if err := s.produceOrder(sell_order); err != nil {
			log.Println("SERVER ERROR: producer.SendMessage")
			return err
		}
	}

	return nil
}

func (s *service) produceOrder(message []byte) error {
	topic := viper.GetString("kafka.topic")
	_, _, err := s.producer.SendMessage(&sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.ByteEncoder(message)})
	return err
}

func (s *service) randomOrders(n int) ([][]byte, [][]byte, error) {
	buy_orders := make([][]byte, 0, n)
	sell_orders := make([][]byte, 0, n)
	for i := 0; i < n; i += 1 {
		buy_order, err1 := s.randomOrder("buy")
		sell_order, err2 := s.randomOrder("sell")
		if err1 != nil {
			return nil, nil, err1
		}
		if err2 != nil {
			return nil, nil, err2
		}

		buy_orders = append(buy_orders, buy_order)
		sell_orders = append(sell_orders, sell_order)
	}
	return buy_orders, sell_orders, nil
}

func (s *service) randomOrder(order_type string) ([]byte, error) {
	price := rand.IntN(2) + 10
	quantity := rand.IntN(10) + 1
	order := &structure.Order{Type: order_type, Price: price, Quantity: quantity}
	return json.Marshal(order)
}
