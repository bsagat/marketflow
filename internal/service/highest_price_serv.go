package service

import (
	"log/slog"
	"marketflow/internal/domain"
	"net/http"
	"time"
)

func (serv *DataModeServiceImp) GetHighestPrice(exchange, symbol string, period string) (domain.Data, int, error) {

}

func (serv *DataModeServiceImp) GetHighestPriceWithPeriod(exchange, symbol string, period string) (domain.Data, int, error) {
	var (
		highest domain.Data
		err     error
	)

	switch exchange {
	case "Exchange1", "Exchange2", "Exchange3", "All":
	default:
		return highest, http.StatusBadRequest, domain.ErrInvalidExchangeVal
	}

	switch symbol {
	case domain.BTCUSDT, domain.DOGEUSDT, domain.ETHUSDT, domain.SOLUSDT, domain.TONUSDT:
	default:
		return highest, http.StatusBadRequest, domain.ErrInvalidSymbolVal
	}

	duration, err := time.ParseDuration(period)
	if err != nil {
		return highest, http.StatusBadRequest, err
	}

	startTime := time.Now()

	result, err := serv.DB.GetExtremePrice("MAX", exchange, symbol, period)
	if err != nil {
		slog.Error("Failed to get highest price from DB", "error", err.Error())
		return domain.Data{}, http.StatusInternalServerError, err
	}

	aggregated := serv.GetAggregatedDataByDuration(exchange, symbol, duration)
	merged := MergeAggregatedData(aggregated)

	result.Price += merged[result.ExchangeName+" "+symbol].Max_price
	result.Timestamp = startTime.Add(-duration).UnixMilli()

	return result, http.StatusOK, nil
}
