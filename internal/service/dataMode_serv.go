package service

import (
	"log"
	"log/slog"
	"marketflow/internal/domain"
	"sync"
	"time"
)

type DataModeService struct {
	Datafetcher domain.DataFetcher
	DBSaver     domain.Database
	CacheSaver  domain.CacheMemory
	DataBuffer  []map[string]domain.ExchangeData
	mu          sync.Mutex
}

func NewDataFetcher(dataSource domain.DataFetcher, DataSaver domain.Database, CacheSaver domain.CacheMemory) *DataModeService {
	serv := &DataModeService{Datafetcher: dataSource, DBSaver: DataSaver, CacheSaver: CacheSaver, DataBuffer: make([]map[string]domain.ExchangeData, 0)}
	serv.ListenAndSave()

	return serv
}

func (serv *DataModeService) ListenAndSave() {
	aggregated := serv.Datafetcher.SetupDataFetcher()

	go func() {
		t := time.NewTicker(time.Minute)
		for tick := range t.C {
			slog.Debug(tick.String())
			serv.mu.Lock()

			merged := serv.MergeAggregatedData()
			err := serv.DBSaver.SaveAggregatedData(merged)
			if err != nil {
				log.Printf("Failed to save aggregated data in database: %s \n", err.Error())
			}

			err = serv.CacheSaver.SaveAggregatedData(merged)
			if err != nil {
				log.Printf("Failed to save aggregated data in cache: %s \n", err.Error())
			}

			serv.DataBuffer = make([]map[string]domain.ExchangeData, 0)
			serv.mu.Unlock()
		}
	}()

	go func() {
		for data := range aggregated {
			serv.DataBuffer = append(serv.DataBuffer, data)
		}
	}()
}

func (serv *DataModeService) MergeAggregatedData() map[string]domain.ExchangeData {
	result := make(map[string]domain.ExchangeData)
	sums := make(map[string]float64)
	counts := make(map[string]int)

	for _, dataMap := range serv.DataBuffer {
		for key, val := range dataMap {
			agg, exists := result[key]
			if !exists {
				agg = domain.ExchangeData{
					Pair_name: val.Pair_name,
					Exchange:  val.Exchange,
					Min_price: val.Min_price,
					Max_price: val.Max_price,
					Timestamp: val.Timestamp,
				}
			}

			if val.Min_price < agg.Min_price {
				agg.Min_price = val.Min_price
			}
			if val.Max_price > agg.Max_price {
				agg.Max_price = val.Max_price
			}

			sums[key] += val.Average_price
			counts[key]++

			if val.Timestamp.After(agg.Timestamp) {
				agg.Timestamp = val.Timestamp
			}

			result[key] = agg
		}
	}

	// Считаем среднее
	for key, item := range result {
		if count := counts[key]; count > 0 {
			item.Average_price = sums[key] / float64(count)
			result[key] = item
		}
	}

	return result
}

func (serv *DataModeService) GetAggregatedData(lastNSeconds int) map[string]domain.ExchangeData {
	cutoff := time.Now().Add(-time.Duration(lastNSeconds) * time.Second)

	serv.mu.Lock()
	defer serv.mu.Unlock()

	var latest map[string]domain.ExchangeData
	var latestTime time.Time

	for _, dataMap := range serv.DataBuffer {

		for _, data := range dataMap {
			if data.Timestamp.After(cutoff) {
				if latest == nil || data.Timestamp.After(latestTime) {
					latest = dataMap
					latestTime = data.Timestamp
				}
				break
			}
		}
	}

	return latest
}
