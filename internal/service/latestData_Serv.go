package service

import (
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

	if err := CheckExchangeName(exchange); err != nil {
		slog.Error("Failed to get latest data: ", "error", err.Error())
		return latest, http.StatusBadRequest, err
	}

	if err := CheckSymbolName(symbol); err != nil {
		slog.Error("Failed to get latest data: ", "error", err.Error())
		return latest, http.StatusBadRequest, err
	}

	// first we look for data in the cache
	latest, err = serv.Cache.GetLatestData(exchange, symbol)
	if err != nil {
		// If Redis is not available, se look for data in the DB
		slog.Debug("Failed to get latest data from cache: ", "error", err.Error())
		if exchange == "All" {
			latest, err = serv.DB.GetLatestDataByAllExchanges(symbol)
			if err != nil {
				slog.Error("Failed to get latest data by all exchanges from Db: ", "error", err.Error())
				return latest, http.StatusInternalServerError, err
			}
		} else {
			latest, err = serv.DB.GetLatestDataByExchange(exchange, symbol)
			if err != nil {
				slog.Error("Failed to get latest data by exchange from Db: ", "error", err.Error())
				return latest, http.StatusInternalServerError, err
			}
		}
	}

	if latest.Price == 0 {
		return domain.Data{}, http.StatusNotFound, domain.ErrLatestPriceNotFound
	}

	return latest, http.StatusOK, nil
}
