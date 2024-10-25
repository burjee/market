package structure

type OrderList struct {
	Type         string  `json:"type"`
	BuyOrders    []Order `json:"buy_orders"`
	SellOrders   []Order `json:"sell_orders"`
	BuyQuantity  string  `json:"buy_quantity"`
	SellQuantity string  `json:"sell_quantity"`
}

type Order struct {
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}
