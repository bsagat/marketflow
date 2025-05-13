package domain

type Data struct {
	ExchangeName string
	Symbol       string  `json:"symbol"`
	Price        float64 `json:"price"`
	Timestamp    int64   `json:"timestamp"`
}
