package domain

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
	GetLatestData(exchange, symbol string) (Data, error)
	CheckHealth() error
	//
	GetExtremePrice(op, exchange, symbol string, period int) (Data, error)
	GetAveragePrice(exchange, symbol string, period int) (Data, error)
}

// For services
type DataModeService interface {
	GetAggregatedData(lastNSeconds int) map[string]ExchangeData
	GetLatestData(exchange, symbol string) (Data, int, error)
	GetHighestPrice(exchange, symbol string, period int) (Data, int, error)
	GetLowestPrice(exchange, symbol string, period int) (Data, int, error)
	GetAveragePrice(exchange, symbol string, period int) (Data, int, error)
	SaveLatestData(rawDataCh chan []Data)
	MergeAggregatedData() map[string]ExchangeData
	SwitchMode(mode string) (int, error)
	CheckHealth() []ConnMsg
	ListenAndSave() error
	StopListening()
}
