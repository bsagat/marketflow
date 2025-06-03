package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	datafetcher "marketflow/internal/adapters/dataFetcher"
	"marketflow/internal/domain"
	"net/http"
	"sync"
	"time"
)

type DataModeServiceImp struct {
	Datafetcher domain.DataFetcher
	DB          domain.Database
	Cache       domain.CacheMemory
	DataBuffer  []map[string]domain.ExchangeData
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	mu          sync.Mutex
}

func NewDataFetcher(dataSource domain.DataFetcher, DataSaver domain.Database, Cache domain.CacheMemory) *DataModeServiceImp {
	ctx, cancel := context.WithCancel(context.Background())
	return &DataModeServiceImp{
		Datafetcher: dataSource,
		DB:          DataSaver,
		Cache:       Cache,
		DataBuffer:  make([]map[string]domain.ExchangeData, 0),
		ctx:         ctx,
		cancel:      cancel,
	}
}

var _ (domain.DataModeService) = (*DataModeServiceImp)(nil)

func (serv *DataModeServiceImp) SwitchMode(mode string) (int, error) {
	// Check if is current datafetcher mode equal to changing mode
	if _, ok := serv.Datafetcher.(*datafetcher.LiveMode); (ok && mode == "live") || (!ok && mode == "test") {
		return http.StatusBadRequest, fmt.Errorf("data mode is already switched to %s", mode)
	}

	switch mode {
	case "test":
		serv.Datafetcher.Close()
		serv.Datafetcher = datafetcher.NewTestModeFetcher()
		if err := serv.ListenAndSave(); err != nil {
			return http.StatusInternalServerError, err
		}
	case "live":
		serv.Datafetcher.Close()
		serv.Datafetcher = datafetcher.NewLiveModeFetcher()
		if err := serv.ListenAndSave(); err != nil {
			return http.StatusInternalServerError, err
		}
	default:
		return http.StatusBadRequest, domain.ErrInvalidModeVal
	}
	return http.StatusOK, nil
}

func (serv *DataModeServiceImp) StopListening() {
	serv.cancel()
	serv.Datafetcher.Close()
	serv.wg.Wait()
	slog.Info("Listen and save goroutine has been finished...")
}

func (serv *DataModeServiceImp) ListenAndSave() error {
	aggregated, rawDataCh, err := serv.Datafetcher.SetupDataFetcher()
	if err != nil {
		return err
	}
	serv.wg.Add(3)

	go func() {
		defer serv.wg.Done()
		serv.SaveLatestData(rawDataCh)
	}()

	go func() {
		defer serv.wg.Done()
		t := time.NewTicker(time.Minute)
		defer t.Stop()

		for {
			select {
			case <-serv.ctx.Done():
				return
			case <-t.C:
				serv.mu.Lock()
				merged := serv.MergeAggregatedData()
				serv.DB.SaveAggregatedData(merged)
				serv.Cache.SaveAggregatedData(merged)
				serv.DataBuffer = nil
				serv.mu.Unlock()
			}
		}
	}()

	go func() {
		defer serv.wg.Done()
		for {
			select {
			case <-serv.ctx.Done():
				return
			case data, ok := <-aggregated:
				if !ok {
					return
				}
				serv.mu.Lock()
				serv.DataBuffer = append(serv.DataBuffer, data)
				slog.Debug("Received data", "buffer_size", len(serv.DataBuffer)) // Tick log
				serv.mu.Unlock()
			}
		}
	}()

	return nil
}

func (serv *DataModeServiceImp) SaveLatestData(rawDataCh chan []domain.Data) {
	for rawData := range rawDataCh {
		latestData := make(map[string]domain.Data)
		for i := len(rawData) - 1; i >= 0; i-- {
			if rawData[i].ExchangeName == "" || rawData[i].Symbol == "" {
				continue
			}

			exchKey := "latest " + rawData[i].ExchangeName + " " + rawData[i].Symbol
			allKey := "latest " + "All" + " " + rawData[i].Symbol

			if _, exist := latestData[exchKey]; !exist {
				latestData[exchKey] = rawData[i]
			}

			if _, exist := latestData[allKey]; !exist {
				latestData[allKey] = rawData[i]
			}

			// Break loop if we find all latest prices
			if len(latestData) == 20 {
				break
			}
		}

		if err := serv.Cache.SaveLatestData(latestData); err != nil {
			slog.Debug("Failed to save latest data to cache: " + err.Error())

			if err := serv.DB.SaveLatestData(latestData); err != nil {
				slog.Error("Failed to save latest data to Db: " + err.Error())
			}
		}

	}
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

// ////////////
func (serv *DataModeServiceImp) GetHighestPrice(exchange, symbol string, period int) (domain.Data, int, error) {
	var (
		latest domain.Data
		err    error
	)

	switch exchange {
	case "Exchange1", "Exchange2", "Exchange3", "All":
	default:
		return latest, http.StatusBadRequest, domain.ErrInvalidExchangeVal
	}

	switch symbol {
	case domain.BTCUSDT, domain.DOGEUSDT, domain.ETHUSDT, domain.SOLUSDT, domain.TONUSDT:
	default:
		return latest, http.StatusBadRequest, domain.ErrInvalidSymbolVal
	}

	result, err := serv.DB.GetExtremePrice("MAX", exchange, symbol, period)
	if err != nil {
		slog.Error("Failed to get highest price from DB", "error", err.Error())
		return domain.Data{}, http.StatusInternalServerError, err
	}

	if result.Price == 0 {
		return domain.Data{}, http.StatusNotFound, errors.New("highest price not found")
	}

	return result, http.StatusOK, nil
}

func (serv *DataModeServiceImp) GetLowestPrice(exchange, symbol string, period int) (domain.Data, int, error) {
	var (
		latest domain.Data
		err    error
	)

	switch exchange {
	case "Exchange1", "Exchange2", "Exchange3", "All":
	default:
		return latest, http.StatusBadRequest, domain.ErrInvalidExchangeVal
	}

	switch symbol {
	case domain.BTCUSDT, domain.DOGEUSDT, domain.ETHUSDT, domain.SOLUSDT, domain.TONUSDT:
	default:
		return latest, http.StatusBadRequest, domain.ErrInvalidSymbolVal
	}

	result, err := serv.DB.GetExtremePrice("MIN", exchange, symbol, period)
	if err != nil {
		slog.Error("Failed to get lowest price from DB", "error", err.Error())
		return domain.Data{}, http.StatusInternalServerError, err
	}

	if result.Price == 0 {
		return domain.Data{}, http.StatusNotFound, errors.New("lowest price not found")
	}

	return result, http.StatusOK, nil
}

func (serv *DataModeServiceImp) GetAveragePrice(exchange, symbol string, period int) (domain.Data, int, error) {
	var (
		latest domain.Data
		err    error
	)

	switch exchange {
	case "Exchange1", "Exchange2", "Exchange3", "All":
	default:
		return latest, http.StatusBadRequest, domain.ErrInvalidExchangeVal
	}

	switch symbol {
	case domain.BTCUSDT, domain.DOGEUSDT, domain.ETHUSDT, domain.SOLUSDT, domain.TONUSDT:
	default:
		return latest, http.StatusBadRequest, domain.ErrInvalidSymbolVal
	}

	result, err := serv.DB.GetAveragePrice(exchange, symbol, period)
	if err != nil {
		slog.Error("Failed to get average price from DB", "error", err.Error())
		return domain.Data{}, http.StatusInternalServerError, err
	}

	if result.Price == 0 {
		return domain.Data{}, http.StatusNotFound, errors.New("average price not found")
	}

	return result, http.StatusOK, nil
}
