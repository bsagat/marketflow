package service

import (
	"marketflow/internal/domain"
	"net/http"
)

// Latest data validation and service logic
func (serv *DataModeServiceImp) GetLatestData(exchange string, symbol string) (domain.Data, int, error) {
	var latest domain.Data

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

	latest, err := serv.Cache.GetLatestData(exchange, symbol)
	if err != nil {
		return latest, http.StatusInternalServerError, err
	}
	return latest, http.StatusOK, nil
}
