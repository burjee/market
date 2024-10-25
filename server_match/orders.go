package main

import (
	"context"
	"fmt"
	"server_match/libs"
	"strconv"
)

func createOrderIndex(redis_client *libs.RedisClient) {
	if err := redis_client.Do(context.Background(), redis_client.B().Flushall().Build()).Error(); err != nil {
		panic(err)
	}

	idx_buy_builder := redis_client.B().FtCreate().Index("idx:buy:order").
		OnHash().
		Prefix(1).Prefix("buy:order:").
		Schema().
		FieldName("price").Numeric().Sortable().
		FieldName("quantity").Numeric().Noindex().
		FieldName("created_at").Numeric().Sortable().
		Build()

	idx_sell_builder := redis_client.B().FtCreate().Index("idx:sell:order").
		OnHash().
		Prefix(1).Prefix("sell:order:").
		Schema().
		FieldName("price").Numeric().Sortable().
		FieldName("quantity").Numeric().Noindex().
		FieldName("created_at").Numeric().Sortable().
		Build()

	err := redis_client.Do(context.Background(), idx_buy_builder).Error()
	if err != nil {
		panic(err)
	}

	err = redis_client.Do(context.Background(), idx_sell_builder).Error()
	if err != nil {
		panic(err)
	}
}

func createPriceQuantityIndex(redis_client *libs.RedisClient) {
	idx_buy_price_builder := redis_client.B().FtCreate().Index("idx:buy:price").
		OnHash().
		Prefix(1).Prefix("buy:price:").
		Schema().
		FieldName("price").Numeric().Sortable().
		FieldName("quantity").Numeric().Noindex().
		Build()

	idx_sell_price_builder := redis_client.B().FtCreate().Index("idx:sell:price").
		OnHash().
		Prefix(1).Prefix("sell:price:").
		Schema().
		FieldName("price").Numeric().Sortable().
		FieldName("quantity").Numeric().Noindex().
		Build()

	err := redis_client.Do(context.Background(), idx_buy_price_builder).Error()
	if err != nil {
		panic(err)
	}

	err = redis_client.Do(context.Background(), idx_sell_price_builder).Error()
	if err != nil {
		panic(err)
	}
}

func generateTestOrder(redis_client *libs.RedisClient) {
	// json
	if err := redis_client.Do(context.Background(), redis_client.B().JsonSet().Key("order:broken").Path("$").Value("[]").Build()).Error(); err != nil {
		panic(err)
	}

	price_buy_count := make(map[int]int)
	price_sell_count := make(map[int]int)
	total_buy_count := 0
	total_sell_count := 0

	// insert orders
	for i, buy_order := range buy_orders {
		price_buy_count[buy_order.Price] += buy_order.Quantity
		total_buy_count += buy_order.Quantity

		key := fmt.Sprintf("buy:order:%d", i)
		hset_buy_builder := redis_client.B().
			Hset().Key(key).
			FieldValue().
			FieldValue("price", strconv.Itoa(buy_order.Price)).
			FieldValue("quantity", strconv.Itoa(buy_order.Quantity)).
			FieldValue("created_at", buy_order.CreatedAt).
			Build()
		if err := redis_client.Do(context.Background(), hset_buy_builder).Error(); err != nil {
			panic(err)
		}
	}
	for i, sell_order := range sell_orders {
		price_sell_count[sell_order.Price] += sell_order.Quantity
		total_sell_count += sell_order.Quantity

		key := fmt.Sprintf("sell:order:%d", i)
		hset_sell_builder := redis_client.B().
			Hset().Key(key).
			FieldValue().
			FieldValue("price", strconv.Itoa(sell_order.Price)).
			FieldValue("quantity", strconv.Itoa(sell_order.Quantity)).
			FieldValue("created_at", sell_order.CreatedAt).
			Build()
		if err := redis_client.Do(context.Background(), hset_sell_builder).Error(); err != nil {
			panic(err)
		}
	}

	// set total quantity counter
	if err := redis_client.Do(context.Background(), redis_client.B().Set().Key("buy:quantity:counter").Value(strconv.Itoa(total_buy_count)).Build()).Error(); err != nil {
		panic(err)
	}
	if err := redis_client.Do(context.Background(), redis_client.B().Set().Key("sell:quantity:counter").Value(strconv.Itoa(total_sell_count)).Build()).Error(); err != nil {
		panic(err)
	}

	// set id counter
	if err := redis_client.Do(context.Background(), redis_client.B().Set().Key("buy:id:counter").Value(strconv.Itoa(len(buy_orders))).Build()).Error(); err != nil {
		panic(err)
	}
	if err := redis_client.Do(context.Background(), redis_client.B().Set().Key("sell:id:counter").Value(strconv.Itoa(len(sell_orders))).Build()).Error(); err != nil {
		panic(err)
	}

	// set price quantity counter
	for price, quantity := range price_buy_count {
		key := fmt.Sprintf("buy:price:%d", price)
		hset_buy_price_builder := redis_client.B().
			Hset().Key(key).
			FieldValue().
			FieldValue("price", strconv.Itoa(price)).
			FieldValue("quantity", strconv.Itoa(quantity)).
			Build()
		if err := redis_client.Do(context.Background(), hset_buy_price_builder).Error(); err != nil {
			panic(err)
		}
	}
	for price, quantity := range price_sell_count {
		key := fmt.Sprintf("sell:price:%d", price)
		hset_sell_price_builder := redis_client.B().
			Hset().Key(key).
			FieldValue().
			FieldValue("price", strconv.Itoa(price)).
			FieldValue("quantity", strconv.Itoa(quantity)).
			Build()
		if err := redis_client.Do(context.Background(), hset_sell_price_builder).Error(); err != nil {
			panic(err)
		}
	}
}
