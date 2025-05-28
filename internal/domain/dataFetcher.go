package domain

import "time"

// Raw Data from Exchanges
type Data struct {
	ExchangeName string  `json:"exchange"`
	Symbol       string  `json:"symbol"`
	Price        float64 `json:"price"`
	Timestamp    int64   `json:"timestamp"`
}

// Aggregated data
type ExchangeData struct {
	Pair_name     string    `json:"pair_name"`
	Exchange      string    `json:"exchange"`
	Timestamp     time.Time `json:"timestamp"`
	Average_price float64   `json:"average_price"`
	Min_price     float64   `json:"min_price"`
	Max_price     float64   `json:"max_price"`
}
