package service

import (
	"errors"
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
	merged := MergeAggregatedData(serv.DataBuffer)[exchange+" "+symbol]
	serv.mu.Unlock()

	data.Timestamp = merged.Timestamp.UnixMilli()
	data.Price = (merged.Average_price + data.Price) / 2

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
		return data, http.StatusBadRequest, errors.New("invalid exchange name")
	}

	// Period parse logic

	duration, err := time.ParseDuration(period)
	if err != nil {
		return data, http.StatusBadRequest, err
	}
	startTime := time.Now()

	data, err = serv.DB.GetAveragePriceWithDuration(exchange, symbol, startTime, duration)
	if err != nil {
		return data, http.StatusInternalServerError, err
	}

	aggregated := serv.GetAggregatedDataByDuration(exchange, symbol, duration)
	merged := MergeAggregatedData(aggregated)

	data.Price += merged[data.ExchangeName+" "+symbol].Average_price
	data.Timestamp = startTime.Add(-duration).UnixMilli()

	return data, http.StatusOK, nil
}

/*

GET /prices/average/{exchange}/{symbol}?period={duration}


if time < 60 seconds:
	search in redis
	if not
		search in dataBuffer
		then save it in redis

elif time >= 1 min:

first search in redis
then give it to client

OR
search in postgres
then save it in redis
then give it to client

*/
