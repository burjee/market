package main

import (
	"strconv"
	"time"
)

type Order struct {
	Price     int
	Quantity  int
	CreatedAt string
}

var t = time.Now()

var buy_orders = []Order{
	{9, 1, nextSecond(&t)},
	{9, 2, nextSecond(&t)},
	{8, 2, nextSecond(&t)},
	{8, 1, nextSecond(&t)},
	{8, 1, nextSecond(&t)},
	{8, 3, nextSecond(&t)},
	{8, 2, nextSecond(&t)},
	{7, 3, nextSecond(&t)},
	{7, 3, nextSecond(&t)},
	{7, 3, nextSecond(&t)},
	{7, 2, nextSecond(&t)},
	{7, 1, nextSecond(&t)},
	{7, 2, nextSecond(&t)},
	{7, 4, nextSecond(&t)},
	{7, 4, nextSecond(&t)},
	{6, 10, nextSecond(&t)},
	{6, 10, nextSecond(&t)},
	{6, 10, nextSecond(&t)},
	{5, 20, nextSecond(&t)},
	{5, 20, nextSecond(&t)},
	{5, 20, nextSecond(&t)},
	{4, 30, nextSecond(&t)},
	{4, 30, nextSecond(&t)},
	{4, 30, nextSecond(&t)},
	{4, 30, nextSecond(&t)},
	{3, 50, nextSecond(&t)},
	{3, 50, nextSecond(&t)},
	{3, 50, nextSecond(&t)},
	{3, 50, nextSecond(&t)},
}

var sell_orders = []Order{
	{10, 2, nextSecond(&t)},
	{10, 1, nextSecond(&t)},
	{11, 2, nextSecond(&t)},
	{11, 3, nextSecond(&t)},
	{11, 2, nextSecond(&t)},
	{11, 3, nextSecond(&t)},
	{12, 2, nextSecond(&t)},
	{12, 3, nextSecond(&t)},
	{12, 1, nextSecond(&t)},
	{12, 1, nextSecond(&t)},
	{12, 2, nextSecond(&t)},
	{13, 1, nextSecond(&t)},
	{13, 2, nextSecond(&t)},
	{13, 4, nextSecond(&t)},
	{13, 4, nextSecond(&t)},
	{13, 10, nextSecond(&t)},
	{14, 5, nextSecond(&t)},
	{14, 13, nextSecond(&t)},
	{14, 23, nextSecond(&t)},
	{14, 20, nextSecond(&t)},
	{15, 10, nextSecond(&t)},
	{15, 35, nextSecond(&t)},
	{15, 33, nextSecond(&t)},
	{15, 20, nextSecond(&t)},
	{16, 32, nextSecond(&t)},
	{16, 57, nextSecond(&t)},
	{16, 52, nextSecond(&t)},
	{16, 51, nextSecond(&t)},
	{16, 55, nextSecond(&t)},
}

func nextSecond(t *time.Time) string {
	*t = t.Add(time.Second)
	return strconv.FormatInt(t.Unix(), 10)
}
