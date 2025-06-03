package service

import (
	"errors"
	"log/slog"
	"marketflow/internal/domain"
	"net/http"
)

func (serv *DataModeServiceImp) GetHighestPrice(exchange, symbol string, period string) (domain.Data, int, error) {
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
