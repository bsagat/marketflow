package domain

import "time"

type Data struct {
	ExchangeName string
	Symbol       string  `json:"symbol"`
	Price        float64 `json:"price"`
	Timestamp    int64   `json:"timestamp"`
}

type ExchangeData struct {
	Pair_name     string
	Exchange      string
	Timestamp     time.Time
	Average_price float64
	Min_price     float64
	Max_price     float64
}

type ExchangePrices struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}
