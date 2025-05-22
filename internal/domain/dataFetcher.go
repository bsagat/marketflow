package domain

import "time"

// Raw Data from Exchanges
type Data struct {
	ExchangeName string
	Symbol       string  `json:"symbol"`
	Price        float64 `json:"price"`
	Timestamp    int64   `json:"timestamp"`
}

// Aggregated data
type ExchangeData struct {
	Pair_name     string
	Exchange      string
	Timestamp     time.Time
	Average_price float64
	Min_price     float64
	Max_price     float64
}
