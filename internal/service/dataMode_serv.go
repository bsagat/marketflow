package service

import (
	"errors"
	"log"
	"log/slog"
	datafetcher "marketflow/internal/adapters/dataFetcher"
	"marketflow/internal/domain"
	"sync"
	"time"
)

type DataModeServiceImp struct {
	Datafetcher domain.DataFetcher
	DB          domain.Database
	Cache       domain.CacheMemory
	DataBuffer  []map[string]domain.ExchangeData
	mu          sync.Mutex
}

func NewDataFetcher(dataSource domain.DataFetcher, DataSaver domain.Database, Cache domain.CacheMemory) *DataModeServiceImp {
	serv := &DataModeServiceImp{Datafetcher: dataSource, DB: DataSaver, Cache: Cache, DataBuffer: make([]map[string]domain.ExchangeData, 0)}
	serv.ListenAndSave()

	return serv
}

func (serv *DataModeServiceImp) SwitchMode(mode string) error {
	switch mode {
	case "test":
		serv.Datafetcher.Close()
		serv.Datafetcher = datafetcher.NewTestModeFetcher()
		serv.ListenAndSave()
	case "live":
		serv.Datafetcher.Close()
		serv.Datafetcher = datafetcher.NewLiveModeFetcher()
		serv.ListenAndSave()
	default:
		return errors.New("invalid mode name, must be (test or live)")
	}
	return nil
}

func (serv *DataModeServiceImp) ListenAndSave() {
	log.Println("Started listening and saving...")

	aggregated := serv.Datafetcher.SetupDataFetcher()
	done := make(chan bool)
	t := time.NewTicker(time.Minute)

	go func() {

	mainLoop:
		for {
			select {
			case tick := <-t.C:
				slog.Debug(tick.String())
				serv.mu.Lock()

				merged := serv.MergeAggregatedData()
				err := serv.DB.SaveAggregatedData(merged)
				if err != nil {
					log.Printf("Failed to save aggregated data in database: %s \n", err.Error())
				}

				err = serv.Cache.SaveAggregatedData(merged)
				if err != nil {
					log.Printf("Failed to save aggregated data in cache: %s \n", err.Error())
				}

				serv.DataBuffer = make([]map[string]domain.ExchangeData, 0)
				serv.mu.Unlock()
			case <-done:
				break mainLoop
			}
		}
	}()

	go func() {
		for data := range aggregated {
			serv.mu.Lock()
			serv.DataBuffer = append(serv.DataBuffer, data)
			serv.mu.Unlock()
		}
		slog.Debug("Listen and Save goroutine has been finished...")
		done <- true
		close(done)
		t.Stop()
	}()
}

func (serv *DataModeServiceImp) MergeAggregatedData() map[string]domain.ExchangeData {
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

func (serv *DataModeServiceImp) GetAggregatedData(lastNSeconds int) map[string]domain.ExchangeData {
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
