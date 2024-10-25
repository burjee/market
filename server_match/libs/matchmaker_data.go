package libs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

var match_order_script = `
local order_key = KEYS[1]
local price_quantity_key = KEYS[2]
local total_quantity_key = KEYS[3]

local matched_quantity = ARGV[1]

if redis.call("HINCRBY", order_key, "quantity", -matched_quantity) == 0 then
	redis.call("DEL", order_key)
end

if redis.call("HINCRBY", price_quantity_key, "quantity", -matched_quantity) == 0 then
	redis.call("DEL", price_quantity_key)
end

redis.call("DECRBY", total_quantity_key, matched_quantity)

return 0
`

var req_order_script = `
local order_key = KEYS[1]
local price_quantity_key = KEYS[2]
local total_quantity_key = KEYS[3]

local price = ARGV[1]
local remain_quantity = ARGV[2]
local created_at = ARGV[3]

redis.call("HSET", order_key, "price", price, "quantity", remain_quantity, "created_at", created_at)
redis.call("HINCRBY", price_quantity_key, "quantity", remain_quantity)
redis.call("HMSET", price_quantity_key, "price", price)
redis.call("INCRBY", total_quantity_key, remain_quantity)

return 0
`

func (m *Matchmaker) newId(req_type string) (int64, error) {
	key := fmt.Sprintf("%s:id:counter", req_type)
	return m.redis_client.Do(context.Background(), m.redis_client.B().Incr().Key(key).Build()).AsInt64()
}

func (m *Matchmaker) getMatchedType(req_type string) string {
	if req_type == "buy" {
		return "sell"
	} else {
		return "buy"
	}
}

func (m *Matchmaker) getBuilderArgs(matched_type string, price int) (string, string, string) {
	var index string
	var query string
	var sort string
	if matched_type == "buy" {
		index = "idx:buy:order"
		query = fmt.Sprintf("@price:[%d +inf]", price)
		sort = "DESC"
	} else {
		index = "idx:sell:order"
		query = fmt.Sprintf("@price:[-inf %d]", price)
		sort = "ASC"
	}

	return index, query, sort
}

func (m *Matchmaker) getMatchedOrders(matched_type string, price int) (int64, []map[string]string, error) {
	index, query, sort := m.getBuilderArgs(matched_type, price)
	ftaggregate_builder := m.redis_client.B().
		FtAggregate().
		Index(index).
		Query(query).
		Load(2).Field("__key", "quantity").
		Sortby(4).Property("@price").Property(sort).Property("@created_at").Property("ASC").
		Limit().OffsetNum(0, 10).
		Build()

	return m.redis_client.Do(context.Background(), ftaggregate_builder).AsFtAggregate()
}

func (m *Matchmaker) getOrderKey(doc map[string]string) (string, error) {
	key, ok := doc["__key"]
	if !ok {
		return "", errors.New("doc error")
	}

	return key, nil
}

func (m *Matchmaker) getPrice(doc map[string]string) (string, error) {
	price, ok := doc["price"]
	if !ok {
		return "", errors.New("doc error")
	}

	return price, nil
}

func (m *Matchmaker) getQuantity(remain int, doc map[string]string) (int, error) {
	_quantity, ok := doc["quantity"]
	if !ok {
		return 0, errors.New("doc error")
	}

	quantity, err := strconv.Atoi(_quantity)
	if err != nil {
		return 0, err
	}

	matched_quantity := min(remain, quantity)
	return matched_quantity, nil
}

func (m *Matchmaker) getMatchedValue(doc map[string]string, remain int) (string, string, int, error) {
	var order_key, price string
	var matched_quantity int
	var err error

	order_key, err1 := m.getOrderKey(doc)
	price, err2 := m.getPrice(doc)
	matched_quantity, err3 := m.getQuantity(remain, doc)

	if err1 != nil {
		err = err1
	}
	if err2 != nil {
		err = err2
	}
	if err3 != nil {
		err = err3
	}

	return order_key, price, matched_quantity, err
}

func (m *Matchmaker) getPriceQuantityKey(order_type string, price string) string {
	return fmt.Sprintf("%s:price:%s", order_type, price)
}

func (m *Matchmaker) getTotalQuantityKey(order_type string) string {
	return fmt.Sprintf("%s:quantity:counter", order_type)
}

func (m *Matchmaker) appendBrokenOrder(req *request, remain int) error {
	broken_request := &brokenRequest{*req, remain}
	b, err := json.Marshal(broken_request)
	if err != nil {
		return err
	}

	return m.redis_client.Do(context.Background(), m.redis_client.B().JsonArrappend().Key("order:broken").Path("$").Value(string(b)).Build()).Error()
}
