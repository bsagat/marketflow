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
	var (
		data domain.Data
		msg  string
		code int
		err  error
	)

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

	switch metric {
	case "highest":
	case "lowest":
	case "average":
		period := r.URL.Query().Get("period")

		if period == "" {
			data, code, err = h.serv.GetAveragePrice(exchange, symbol)
			if err != nil {
				slog.Error("Failed to get average price: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
				if err := senders.SendMsg(w, code, err.Error()); err != nil {
					slog.Error("Failed to send message to the client", "error", err.Error())
				}
				return
			}

		} else {
			data, code, err = h.serv.GetAveragePriceWithPeriod(exchange, symbol, period)
			if err != nil {
				slog.Error("Failed to get average price with period: ", "exchange", exchange, "symbol", symbol, "period", period, "error", err.Error())
				if err := senders.SendMsg(w, code, err.Error()); err != nil {
					slog.Error("Failed to send message to the client", "error", err.Error())
				}
				return
			}

		}

		if err := senders.SendMetricData(w, code, data); err != nil {
			slog.Error("Failed to send JSON message: ", "data", data, "error", err.Error())
			return
		}

		msg = fmt.Sprintf("Average price for %s at %s duration {%s}: %.2f", symbol, exchange, period, data.Price)
	case "latest":
		data, code, err := h.serv.GetLatestData(exchange, symbol)
		if err != nil {
			slog.Error("Failed to get latest data: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
			if err := senders.SendMsg(w, code, err.Error()); err != nil {
				slog.Error("Failed to send message to the client", "error", err.Error())
			}
			return
		}
		if err := senders.SendMetricData(w, code, data); err != nil {
			slog.Error("Failed to send JSON message: ", "data", data, "error", err.Error())
			return
		}
		msg = fmt.Sprintf("Latest price for %s at %s: %.2f", symbol, exchange, data.Price)
	default:
		slog.Error("Failed to get data by metric: ", "exchange", exchange, "symbol", symbol, "metric", metric, "error", domain.ErrInvalidMetricVal)
		if err := senders.SendMsg(w, http.StatusBadRequest, domain.ErrInvalidMetricVal.Error()); err != nil {
			slog.Error("Failed to send message to the client", "error", err.Error())
		}
		return
	}

	slog.Info(msg)
}

func (h *MarketDataHTTPHandler) ProcessMetricQueryByAll(w http.ResponseWriter, r *http.Request) {
	var (
		data     domain.Data
		exchange = "All"
		msg      string
		code     int
		err      error
	)

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
		data, code, err = h.serv.GetAveragePrice(exchange, symbol)
		if err != nil {
			slog.Error("Failed to get average price: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
			if err := senders.SendMsg(w, code, err.Error()); err != nil {
				slog.Error("Failed to send message to the client", "error", err.Error())
			}
			return
		}

		msg = fmt.Sprintf("Average price for %s at %s: %.2f", symbol, exchange, data.Price)
	case "latest":
		data, code, err = h.serv.GetLatestData(exchange, symbol)
		if err != nil {
			slog.Error("Failed to get latest data: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
			if err := senders.SendMsg(w, code, err.Error()); err != nil {
				slog.Error("Failed to send message to the client", "error", err.Error())
			}
			return
		}

		msg = fmt.Sprintf("Latest price for %s at %s: %.2f", symbol, exchange, data.Price)
	default:
		slog.Error("Failed to get data by metric: ", "exchange", exchange, "symbol", symbol, "metric", metric, "error", domain.ErrInvalidMetricVal.Error())
		if err := senders.SendMsg(w, code, msg); err != nil {
			slog.Error("Failed to send JSON message: ", "data", msg, "error", err.Error())
		}
		return
	}

	if err := senders.SendMetricData(w, code, data); err != nil {
		slog.Error("Failed to send JSON message: ", "data", data, "error", err.Error())
		return
	}
	slog.Info(msg)
}
