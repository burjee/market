package libs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/rcrowley/go-metrics"
	"github.com/redis/rueidis"
)

type request struct {
	Type     string `json:"type"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type brokenRequest struct {
	request
	Remain int `json:"remain"`
}

type Matchmaker struct {
	redis_client    *RedisClient
	matched_metrics metrics.Meter
}

func NewMatchmaker(redis_client *RedisClient, matched_metrics metrics.Meter) *Matchmaker {
	return &Matchmaker{redis_client, matched_metrics}
}

func (m *Matchmaker) match(message *sarama.ConsumerMessage) error {
	var req request
	if err := json.Unmarshal(message.Value, &req); err != nil {
		return err
	}

	id, err := m.newId(req.Type)
	if err != nil {
		return err
	}

	remain := req.Quantity
	max_retries := 3
	retries := 0
	for retries = 0; retries < max_retries; retries += 1 {
		if err := m.tryMatch(&req, id, &remain); err == nil {
			break
		} else {
			log.Printf("Failed to match %v (attempt %d)", err, retries+1)
			time.Sleep(time.Millisecond * 5)
		}
	}

	if retries == 3 {
		return m.appendBrokenOrder(&req, remain)
	}

	return nil
}

func (m *Matchmaker) tryMatch(req *request, id int64, remain *int) error {
	matched_type := m.getMatchedType(req.Type)

Outer:
	for {
		total, docs, err := m.getMatchedOrders(matched_type, req.Price)
		if err != nil {
			return err
		}
		if total == 0 {
			break
		}

		for _, doc := range docs {
			order_key, price, matched_quantity, err := m.getMatchedValue(doc, *remain)
			if err != nil {
				return err
			}

			if err := m.matchOrder(matched_type, order_key, price, matched_quantity); err != nil {
				return err
			}

			m.matched_metrics.Mark(1)

			*remain -= matched_quantity
			if *remain == 0 {
				break Outer
			}
		}
	}

	return m.handleReqOrder(req, id, *remain)
}

func (m *Matchmaker) matchOrder(matched_type, order_key, price string, matched_quantity int) error {
	price_quantity_key := m.getPriceQuantityKey(matched_type, price)
	total_quantity_key := m.getTotalQuantityKey(matched_type)

	script := rueidis.NewLuaScript(match_order_script)
	return script.Exec(context.Background(), m.redis_client.Client, []string{order_key, price_quantity_key, total_quantity_key}, []string{strconv.Itoa(matched_quantity)}).Error()
}

func (m *Matchmaker) handleReqOrder(req *request, id int64, remain int) error {
	m.matched_metrics.Mark(1)

	if remain == 0 {
		return nil
	}

	created_at := strconv.FormatInt(time.Now().Unix(), 10)
	order_key := fmt.Sprintf("%s:order:%d", req.Type, id)
	price_quantity_key := m.getPriceQuantityKey(req.Type, strconv.Itoa(req.Price))
	total_quantity_key := m.getTotalQuantityKey(req.Type)

	script := rueidis.NewLuaScript(req_order_script)
	return script.Exec(context.Background(), m.redis_client.Client, []string{order_key, price_quantity_key, total_quantity_key}, []string{strconv.Itoa(req.Price), strconv.Itoa(remain), created_at}).Error()
}
