package service

import (
	"errors"
	"log/slog"
	"marketflow/internal/domain"
	"net/http"
)

// Latest data validation and service logic
func (serv *DataModeServiceImp) GetLatestData(exchange string, symbol string) (domain.Data, int, error) {
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

	latest, err = serv.Cache.GetLatestData(exchange, symbol)
	if err != nil {
		slog.Debug("Failed to get latest data from cache: ", "error", err.Error())
		latest, err = serv.DB.GetLatestData(exchange, symbol)
		if err != nil {
			slog.Error("Failed to get latest data from Db: ", "error", err.Error())
			return latest, http.StatusInternalServerError, err
		}
	}

	if latest.ExchangeName == "" && latest.Price == 0 && latest.Symbol == "" && latest.Timestamp == 0 {
		return domain.Data{}, http.StatusNotFound, errors.New("latest data is not found")
	}

	return latest, http.StatusOK, nil
}
