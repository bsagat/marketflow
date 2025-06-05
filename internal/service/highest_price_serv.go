package service

import (
	"errors"
	"log/slog"
	"marketflow/internal/domain"
	"net/http"
	"time"
)

func (serv *DataModeServiceImp) GetHighestPrice(exchange, symbol string) (domain.Data, int, error) {
	var (
		highest domain.Data
		err     error
	)

	if err := CheckExchangeName(exchange); err != nil {
		return domain.Data{}, http.StatusBadRequest, err
	}

	if err := CheckSymbolName(symbol); err != nil {
		return domain.Data{}, http.StatusBadRequest, err
	}

	switch exchange {
	case "All":
		highest, err = serv.DB.GetExtremePriceByAllExchanges("Max_price", symbol)
		if err != nil {
			slog.Error("Failed to get highest price by all exchanges", "error", err.Error())
			return domain.Data{}, http.StatusInternalServerError, err
		}

	default:
		highest, err = serv.DB.GetExtremePriceByExchange("Max_price", exchange, symbol)
		if err != nil {
			slog.Error("Failed to get highest price from exchange", "error", err.Error())
			return domain.Data{}, http.StatusInternalServerError, err
		}
	}

	serv.mu.Lock()
	merged := MergeAggregatedData(serv.DataBuffer)
	serv.mu.Unlock()

	for _, name := range domain.Exchanges {
		key := name + " " + symbol

		if agg, ok := merged[key]; ok {
			if agg.Max_price != 0 {
				highest.Price = max(highest.Price, agg.Max_price)
			}
		} else {
			slog.Warn("Aggregated data not found for key", "key", key)
		}
		highest.Timestamp = time.Now().UnixMilli()
	}

	return highest, http.StatusOK, nil
}

func (serv *DataModeServiceImp) GetHighestPriceWithPeriod(exchange, symbol string, period string) (domain.Data, int, error) {
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

	highest, err := serv.DB.GetExtremePriceByDuration("Max_price", exchange, symbol, startTime, duration)
	if err != nil {
		slog.Error("Failed to get highest price from Exchange by period", "error", err.Error())
		return domain.Data{}, http.StatusInternalServerError, err
	}

	aggregated := serv.GetAggregatedDataByDuration(exchange, symbol, duration)
	merged := MergeAggregatedData(aggregated)

	key := highest.ExchangeName + " " + symbol
	if agg, ok := merged[key]; ok {
		if agg.Max_price != 0 {
			highest.Price = max(highest.Price, agg.Max_price)
		}
	} else {
		slog.Warn("Aggregated data not found for key", "key", key)
	}
	highest.Timestamp = startTime.Add(-duration).UnixMilli()

	return highest, http.StatusOK, nil
}
