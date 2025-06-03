package domain

import "time"

// For adapters
type DataFetcher interface {
	SetupDataFetcher() (chan map[string]ExchangeData, chan []Data, error)
	CheckHealth() error
	Close()
}

type CacheMemory interface {
	SaveAggregatedData(aggregatedData map[string]ExchangeData) error
	SaveLatestData(latestData map[string]Data) error
	GetLatestData(exchange, symbol string) (Data, error)
	CheckHealth() error
}

type Database interface {
	SaveAggregatedData(aggregatedData map[string]ExchangeData) error
	SaveLatestData(latestData map[string]Data) error
	GetLatestDataByExchange(exchange, symbol string) (Data, error)
	GetLatestDataByAllExchanges(symbol string) (Data, error)
	GetAveragePriceByExchange(exchange, symbol string) (Data, error)
	GetAveragePriceByAllExchanges(symbol string) (Data, error)
	GetAveragePriceWithDuration(exchange, symbol string, startTime time.Time, duration time.Duration) (Data, error)
	CheckHealth() error
}

// For services
type DataModeService interface {
	GetAggregatedDataByDuration(exchange, symbol string, duration time.Duration) []map[string]ExchangeData
	GetLatestData(exchange string, symbol string) (Data, int, error)
	GetAveragePrice(exchange, symbol string) (Data, int, error)
	GetAveragePriceWithPeriod(exchange, symbol, period string) (Data, int, error)
	SaveLatestData(rawDataCh chan []Data)
	SwitchMode(mode string) (int, error)
	CheckHealth() []ConnMsg
	ListenAndSave() error
	StopListening()
}
