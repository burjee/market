package libs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"server_notify/structure"
	"time"

	"github.com/redis/rueidis"
)

type OrderList struct {
	redis_client *RedisClient
}

func NewOrderList(redis_client *RedisClient) *OrderList {
	return &OrderList{
		redis_client: redis_client,
	}
}

func (o *OrderList) getQuery() (rueidis.Completed, rueidis.Completed) {
	search_buy_builder := o.redis_client.B().
		FtSearch().
		Index("idx:buy:price").
		Query("*").
		Sortby("price").Desc().
		Limit().OffsetNum(0, 6).
		Build()

	search_sell_builder := o.redis_client.B().
		FtSearch().
		Index("idx:sell:price").
		Query("*").
		Sortby("price").Asc().
		Limit().OffsetNum(0, 6).
		Build()

	return search_buy_builder, search_sell_builder
}

func (o *OrderList) Broadcast() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		// order list
		search_buy_builder, search_sell_builder := o.getQuery()
		_, buy_docs, err1 := o.redis_client.Do(context.Background(), search_buy_builder).AsFtSearch()
		_, sell_docs, err2 := o.redis_client.Do(context.Background(), search_sell_builder).AsFtSearch()
		if err1 != nil || err2 != nil {
			log.Println("SERVER ERROR: get list")
			continue
		}

		order_list, ok := o.getOrderList(buy_docs, sell_docs)
		if !ok {
			log.Println("SERVER ERROR: getOrderList")
			continue
		}

		// price quantity
		buy_quantity, err1 := o.redis_client.Do(context.Background(), o.redis_client.B().Get().Key("buy:quantity:counter").Build()).ToString()
		sell_quantity, err2 := o.redis_client.Do(context.Background(), o.redis_client.B().Get().Key("sell:quantity:counter").Build()).ToString()
		if err1 != nil || err2 != nil {
			log.Println("SERVER ERROR: get quantity")
			continue
		}

		order_list.BuyQuantity = buy_quantity
		order_list.SellQuantity = sell_quantity

		// publish
		message, err := json.Marshal(order_list)
		if err != nil {
			log.Println("SERVER ERROR: json.Marshal")
			continue
		}

		encoded_data := base64.StdEncoding.EncodeToString(message)
		if err := o.redis_client.Do(context.Background(), o.redis_client.B().Publish().Channel("broadcast").Message(encoded_data).Build()).Error(); err != nil {
			log.Println("SERVER ERROR: redis.Publish")
		}
	}
}

func (o *OrderList) getOrderList(buy_docs []rueidis.FtSearchDoc, sale_docs []rueidis.FtSearchDoc) (*structure.OrderList, bool) {
	buy_orders := make([]structure.Order, 0, len(buy_docs))
	sell_orders := make([]structure.Order, 0, len(sale_docs))

	if ok := o.appendOrders(buy_docs, &buy_orders); !ok {
		return nil, false
	}

	if ok := o.appendOrders(sale_docs, &sell_orders); !ok {
		return nil, false
	}

	list := &structure.OrderList{Type: "order_list", BuyOrders: buy_orders, SellOrders: sell_orders}
	return list, true
}

func (o *OrderList) appendOrders(docs []rueidis.FtSearchDoc, list *[]structure.Order) bool {
	for _, doc := range docs {
		price, ok1 := doc.Doc["price"]
		quantity, ok2 := doc.Doc["quantity"]
		if !ok1 || !ok2 {
			return false
		}

		order := structure.Order{Price: price, Quantity: quantity}
		*list = append(*list, order)
	}
	return true
}
