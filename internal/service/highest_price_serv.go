package service

import (
	"log/slog"
	"marketflow/internal/domain"
	"net/http"
	"time"
)

// Fetches the highest price for a specific exchange and symbol
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
		highest, err = serv.DB.GetMaxPriceByAllExchanges(symbol)
		if err != nil {
			slog.Error("Failed to get highest price by all exchanges", "error", err.Error())
			return domain.Data{}, http.StatusInternalServerError, err
		}

	default:
		highest, err = serv.DB.GetMaxPriceByExchange(exchange, symbol)
		if err != nil {
			slog.Error("Failed to get highest price from exchange", "error", err.Error())
			return domain.Data{}, http.StatusInternalServerError, err
		}
	}

	serv.mu.Lock()
	merged := MergeAggregatedData(serv.DataBuffer)
	serv.mu.Unlock()

	key := exchange + " " + symbol
	if agg, ok := merged[key]; ok {
		if agg.Max_price > highest.Price {
			highest.Price = agg.Max_price
			highest.Timestamp = agg.Timestamp.UnixMilli()
		}
	} else {
		slog.Warn("Aggregated data not found for key", "key", key)
	}

	if highest.Price == 0 {
		return domain.Data{}, http.StatusNotFound, domain.ErrHighPriceNotFound
	}

	return highest, http.StatusOK, nil
}

// Fetches the average price for a specific exchange and symbol over a given period
func (serv *DataModeServiceImp) GetHighestPriceWithPeriod(exchange, symbol string, period string) (domain.Data, int, error) {
	if err := CheckExchangeName(exchange); err != nil {
		return domain.Data{}, http.StatusBadRequest, err
	}

	if err := CheckSymbolName(symbol); err != nil {
		return domain.Data{}, http.StatusBadRequest, err
	}

	if exchange == "All" {
		return domain.Data{}, http.StatusBadRequest, domain.ErrAllNotSupported
	}

	duration, err := time.ParseDuration(period)
	if err != nil {
		return domain.Data{}, http.StatusBadRequest, err
	}

	startTime := time.Now()

	highest, err := serv.DB.GetMaxPriceByExchangeWithDuration(exchange, symbol, startTime, duration)
	if err != nil {
		slog.Error("Failed to get highest price from Exchange by period", "error", err.Error())
		return domain.Data{}, http.StatusInternalServerError, err
	}

	aggregated := serv.GetAggregatedDataByDuration(exchange, symbol, duration)
	merged := MergeAggregatedData(aggregated)

	key := exchange + " " + symbol
	if agg, ok := merged[key]; ok {
		if agg.Max_price > highest.Price {
			highest.Price = agg.Max_price
			highest.Timestamp = agg.Timestamp.UnixMilli()
		}
	} else {
		slog.Warn("Aggregated data not found for key", "key", key)
	}
	highest.Timestamp = startTime.Add(-duration).UnixMilli()

	if highest.Price == 0 {
		return domain.Data{}, http.StatusNotFound, domain.ErrHighPriceWithPeriodNotFound
	}

	return highest, http.StatusOK, nil
}

// Fetches the average price across all exchanges for a given symbol over a specified period
func (serv *DataModeServiceImp) GetHighestPriceByAllExchangesWithPeriod(symbol string, period string) (domain.Data, int, error) {
	exchange := "All"
	if err := CheckSymbolName(symbol); err != nil {
		return domain.Data{}, http.StatusBadRequest, err
	}

	duration, err := time.ParseDuration(period)
	if err != nil {
		return domain.Data{}, http.StatusBadRequest, err
	}

	startTime := time.Now()

	highest, err := serv.DB.GetMaxPriceByAllExchangesWithDuration(symbol, startTime, duration)
	if err != nil {
		slog.Error("Failed to get highest price from Exchange by period", "error", err.Error())
		return domain.Data{}, http.StatusInternalServerError, err
	}

	aggregated := serv.GetAggregatedDataByDuration(exchange, symbol, duration)
	merged := MergeAggregatedData(aggregated)

	key := exchange + " " + symbol
	if agg, ok := merged[key]; ok {
		if agg.Max_price > highest.Price {
			highest.Price = agg.Max_price
			highest.Timestamp = agg.Timestamp.UnixMilli()
		}
	} else {
		slog.Warn("Aggregated data not found for key", "key", key)
	}

	if highest.Price == 0 {
		return domain.Data{}, http.StatusNotFound, domain.ErrHighPriceWithPeriodNotFound
	}

	return highest, http.StatusOK, nil
}
