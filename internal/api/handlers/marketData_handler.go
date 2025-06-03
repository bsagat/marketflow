package handlers

import (
	"fmt"
	"log/slog"
	"marketflow/internal/api/senders"
	"marketflow/internal/domain"
	"net/http"
	"time"
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
		if err := senders.SendMsg(w, http.StatusBadRequest, domain.ErrEmptyMetricVal.Error()); err != nil {
			slog.Error("Failed to send message to the client", "error", err.Error())
		}
		return
	}

	exchange := r.PathValue("exchange")
	if len(exchange) == 0 {
		slog.Error("Failed to get exchange value from path: ", "error", domain.ErrEmptyExchangeVal.Error())
		if err := senders.SendMsg(w, http.StatusBadRequest, domain.ErrEmptyExchangeVal.Error()); err != nil {
			slog.Error("Failed to send message to the client", "error", err.Error())
		}
		return
	}

	symbol := r.PathValue("symbol")
	if len(symbol) == 0 {
		slog.Error("Failed to get symbol value from path: ", "error", domain.ErrEmptyExchangeVal)
		if err := senders.SendMsg(w, http.StatusBadRequest, domain.ErrEmptySymbolVal.Error()); err != nil {
			slog.Error("Failed to send message to the client", "error", err.Error())
		}
		return
	}

	period := 60 // Default 1m
	if p := r.URL.Query().Get("period"); p != "" {
		switch p {
		case "1s", "3s", "5s", "10s", "30s", "1m", "3m", "5m": // like this?
			d, _ := time.ParseDuration(p)
			period = int(d.Seconds())
		default:
			http.Error(w, "invalid period", http.StatusBadRequest)
			return
		}
	}

	var data domain.Data
	var code int
	var err error

	switch metric {
	case "highest":
		data, code, err = h.serv.GetHighestPrice(exchange, symbol, period)
	case "lowest":
		data, code, err = h.serv.GetLowestPrice(exchange, symbol, period)
	case "average":
		data, code, err = h.serv.GetAveragePrice(exchange, symbol, period)
	case "latest":
		data, code, err = h.serv.GetLatestData("All", symbol)

	default:
		slog.Error("Failed to get data by metric: ", "exchange", "All", "symbol", symbol, "metric", metric, "error", domain.ErrInvalidMetricVal.Error())
		http.Error(w, fmt.Sprintf(domain.ErrInvalidMetricVal.Error(), metric), http.StatusBadRequest)
		return
	}

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
