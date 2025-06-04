package service

import (
	"errors"
	"fmt"
	"log/slog"
	"marketflow/internal/domain"
	"net/http"
	"time"
)

func (serv *DataModeServiceImp) GetLowestPrice(exchange, symbol string) (domain.Data, int, error) {
	if err := CheckExchangeName(exchange); err != nil {
		return domain.Data{}, http.StatusBadRequest, err
	}

	if err := CheckSymbolName(symbol); err != nil {
		return domain.Data{}, http.StatusBadRequest, err
	}

	var (
		lowest domain.Data
		err    error
	)
	switch exchange {
	case "All":
		lowest, err = serv.DB.GetExtremePriceByAllExchanges("ASC", symbol)
		if err != nil {
			slog.Error("Failed to get lowest price by all exchanges", "error", err.Error())
			return domain.Data{}, http.StatusInternalServerError, err
		}
	default:
		lowest, err = serv.DB.GetExtremePriceByExchange("ASC", exchange, symbol)
		if err != nil {
			slog.Error("Failed to get lowest price from exchange", "error", err.Error())
			return domain.Data{}, http.StatusInternalServerError, err
		}
	}

	fmt.Println(lowest)

	serv.mu.Lock()
	merged := MergeAggregatedData(serv.DataBuffer)
	serv.mu.Unlock()

	key := lowest.ExchangeName + " " + symbol
	if agg, ok := merged[key]; ok {
		if agg.Min_price != 0 {
			lowest.Price = min(lowest.Price, agg.Min_price)
		}
	} else {
		slog.Warn("Aggregated data not found for key", "key", key)
	}
	lowest.Timestamp = time.Now().UnixMilli()

	return lowest, http.StatusOK, nil
}

func (serv *DataModeServiceImp) GetLowestPriceWithPeriod(exchange, symbol string, period string) (domain.Data, int, error) {
	if err := CheckExchangeName(exchange); err != nil {
		return domain.Data{}, http.StatusBadRequest, err
	}

	if err := CheckSymbolName(symbol); err != nil {
		return domain.Data{}, http.StatusBadRequest, err
	}

	if exchange == "All" {
		return domain.Data{}, http.StatusBadRequest, errors.New(`"All" is not supported for period-based queries`)
	}

	duration, err := time.ParseDuration(period)
	if err != nil {
		return domain.Data{}, http.StatusBadRequest, err
	}

	startTime := time.Now()

	lowest, err := serv.DB.GetExtremePriceByDuration("ASC", exchange, symbol, startTime, duration)
	if err != nil {
		slog.Error("Failed to get lowest price from Exchange by period", "error", err.Error())
		return domain.Data{}, http.StatusInternalServerError, err
	}

	fmt.Println(lowest)

	aggregated := serv.GetAggregatedDataByDuration(exchange, symbol, duration)
	merged := MergeAggregatedData(aggregated)

	key := lowest.ExchangeName + " " + symbol
	if agg, ok := merged[key]; ok {
		if agg.Min_price != 0 {
			lowest.Price = min(lowest.Price, agg.Min_price)
		}
	} else {
		slog.Warn("Aggregated data not found for key", "key", key)
	}
	lowest.Timestamp = time.Now().UnixMilli()

	return lowest, http.StatusOK, nil
}
