package handlers

import (
	"fmt"
	"log/slog"
	"marketflow/internal/api/senders"
	"marketflow/internal/domain"
	"net/http"
)

type MarketDataHTTPHandler struct {
	serv domain.DataModeService
}

func NewMarketDataHandler(serv domain.DataModeService) *MarketDataHTTPHandler {
	return &MarketDataHTTPHandler{serv: serv}
}

func (h *MarketDataHTTPHandler) ProcessMetricQueryByExchange(w http.ResponseWriter, r *http.Request) {
	metric := r.PathValue("metric")
	if len(metric) == 0 {
		slog.Error("Failed to get metric value from path: ", "error", domain.ErrEmptyMetricVal.Error())
		http.Error(w, domain.ErrEmptyMetricVal.Error(), http.StatusBadRequest)
		return
	}

	exchange := r.PathValue("exchange")
	if len(exchange) == 0 {
		slog.Error("Failed to get exchange value from path: ", "error", domain.ErrEmptyExchangeVal.Error())
		http.Error(w, domain.ErrEmptyExchangeVal.Error(), http.StatusBadRequest)
		return
	}

	symbol := r.PathValue("symbol")
	if len(symbol) == 0 {
		slog.Error("Failed to get symbol value from path: ", "error", domain.ErrEmptySymbolVal.Error())
		http.Error(w, domain.ErrEmptySymbolVal.Error(), http.StatusBadRequest)
		return
	}

	switch metric {
	case "highest":
	case "lowest":
	case "average":
	case "latest":
		data, code, err := h.serv.GetLatestData(exchange, symbol)
		if err != nil {
			slog.Error("Failed to get latest data: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
			http.Error(w, err.Error(), code)
			return
		}
		if err := senders.SendMetricData(w, code, data); err != nil {
			slog.Error("Failed to send JSON message: ", "data", data, "error", err.Error())
			http.Error(w, err.Error(), code)
			return
		}
	default:
		slog.Error("Failed to get data by metric: ", "exchange", exchange, "symbol", symbol, "metric", metric, "error", domain.ErrInvalidMetricVal)
		http.Error(w, fmt.Sprintf(domain.ErrInvalidMetricVal.Error(), metric), http.StatusBadRequest)
		return
	}
}

func (h *MarketDataHTTPHandler) ProcessMetricQueryByAll(w http.ResponseWriter, r *http.Request) {
	metric := r.PathValue("metric")
	if len(metric) == 0 {
		slog.Error("Failed to get metric value from path: ", "error", domain.ErrEmptyMetricVal.Error())
		http.Error(w, domain.ErrEmptyMetricVal.Error(), http.StatusBadRequest)
		return
	}

	symbol := r.PathValue("symbol")
	if len(symbol) == 0 {
		slog.Error("Failed to get symbol value from path: ", "error", domain.ErrEmptySymbolVal.Error())
		http.Error(w, domain.ErrEmptySymbolVal.Error(), http.StatusBadRequest)
		return
	}

	switch metric {
	case "highest":
	case "lowest":
	case "average":
	case "latest":
		data, code, err := h.serv.GetLatestData("All", symbol)
		if err != nil {
			slog.Error("Failed to get latest data: ", "exchange", "All", "symbol", symbol, "error", err.Error())
			http.Error(w, err.Error(), code)
			return
		}
		if err := senders.SendMetricData(w, code, data); err != nil {
			slog.Error("Failed to send JSON message: ", "data", data, "error", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		slog.Error("Failed to get data by metric: ", "exchange", "All", "symbol", symbol, "metric", metric, "error", domain.ErrInvalidMetricVal.Error())
		http.Error(w, fmt.Sprintf(domain.ErrInvalidMetricVal.Error(), metric), http.StatusBadRequest)
		return
	}
}
