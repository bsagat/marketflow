package domain

type DataFetcher interface {
	SetupDataFetcher() chan map[string]ExchangeData
}

type CacheMemory interface {
}

type Database interface {
}
