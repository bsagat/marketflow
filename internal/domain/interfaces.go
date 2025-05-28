package domain

// For adapters
type DataFetcher interface {
	SetupDataFetcher() (chan map[string]ExchangeData, chan []Data)
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
	CheckHealth() error
}

// For services
type DataModeService interface {
	GetAggregatedData(lastNSeconds int) map[string]ExchangeData
	GetLatestData(exchange string, symbol string) (Data, int, error)
	SaveLatestData(rawDataCh chan []Data)
	MergeAggregatedData() map[string]ExchangeData
	SwitchMode(mode string) error
	CheckHealth() []ConnMsg
	ListenAndSave()
}
