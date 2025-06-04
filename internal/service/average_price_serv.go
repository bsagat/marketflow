package service

import (
	"errors"
	"log/slog"
	"marketflow/internal/domain"
	"net/http"
	"time"
)

func (serv *DataModeServiceImp) GetAveragePrice(exchange, symbol string) (domain.Data, int, error) {
	var (
		data domain.Data
		err  error
	)

	if err := CheckExchangeName(exchange); err != nil {
		return data, http.StatusBadRequest, err
	}

	if err := CheckSymbolName(symbol); err != nil {
		return data, http.StatusBadRequest, err
	}

	switch exchange {
	case "All":
		data, err = serv.DB.GetAveragePriceByAllExchanges(symbol)
		if err != nil {
			return data, http.StatusInternalServerError, err
		}
	default:
		data, err = serv.DB.GetAveragePriceByExchange(exchange, symbol)
		if err != nil {
			return data, http.StatusInternalServerError, err
		}
	}

	// we also search it in the DataBuffer
	serv.mu.Lock()
	merged := MergeAggregatedData(serv.DataBuffer)
	serv.mu.Unlock()

	data.Timestamp = time.Now().UnixMilli()
	key := exchange + " " + symbol
	if avg, ok := merged[key]; ok {
		if avg.Average_price != 0 {
			data.Price = (avg.Average_price + data.Price) / 2
		}
	} else {
		slog.Warn("Aggregated data not found for key", "key", key)
	}

	return data, http.StatusOK, nil
}

func (serv *DataModeServiceImp) GetAveragePriceWithPeriod(exchange, symbol, period string) (domain.Data, int, error) {
	var (
		data domain.Data
		err  error
	)

	if err := CheckExchangeName(exchange); err != nil {
		return data, http.StatusBadRequest, err
	}

	if err := CheckSymbolName(symbol); err != nil {
		return data, http.StatusBadRequest, err
	}

	if exchange == "All" {
		return data, http.StatusBadRequest, errors.New(`"All" is not supported for period-based queries`)
	}

	duration, err := time.ParseDuration(period)
	if err != nil {
		return data, http.StatusBadRequest, err
	}
	startTime := time.Now()

	data, err = serv.DB.GetAveragePriceWithDuration(exchange, symbol, startTime, duration)
	if err != nil {
		return data, http.StatusInternalServerError, err
	}

	data.Timestamp = startTime.Add(-duration).UnixMilli()

	aggregated := serv.GetAggregatedDataByDuration(exchange, symbol, duration)
	merged := MergeAggregatedData(aggregated)

	key := data.ExchangeName + " " + symbol
	if agg, ok := merged[key]; ok {
		if agg.Average_price != 0 {
			data.Price = (agg.Average_price + data.Price) / 2
		}
	} else {
		slog.Warn("Aggregated data not found for key", "key", key)
	}

	return data, http.StatusOK, nil
}
