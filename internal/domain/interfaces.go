package domain

// For adapters
type DataFetcher interface {
	SetupDataFetcher() chan map[string]ExchangeData
	CheckHealth() error
	Close()
}

type CacheMemory interface {
	SaveAggregatedData(aggregatedData map[string]ExchangeData) error
	CheckHealth() error
}

type Database interface {
	SaveAggregatedData(aggregatedData map[string]ExchangeData) error
	CheckHealth() error
}

// For services
type SystemHealthServ interface {
	CheckHealth() []ConnMsg
}
