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

const (
	MetricHighest = "highest"
	MetricLowest  = "lowest"
	MetricAverage = "average"
	MetricLatest  = "latest"
)

// Core handler for processing metric-based queries by specific exchange
func (h *MarketDataHTTPHandler) ProcessMetricQueryByExchange(w http.ResponseWriter, r *http.Request) {
	var (
		data domain.Data
		msg  string
		code int = 200
		err  error
	)

	metric := r.PathValue("metric")
	if len(metric) == 0 {
		slog.Error("Failed to get metric value from path: ", "error", domain.ErrEmptyMetricVal.Error())
		senders.SendMsg(w, http.StatusBadRequest, domain.ErrEmptyMetricVal.Error())
		return
	}

	exchange := r.PathValue("exchange")
	if len(exchange) == 0 {
		slog.Error("Failed to get exchange value from path: ", "error", domain.ErrEmptyExchangeVal.Error())
		senders.SendMsg(w, http.StatusBadRequest, domain.ErrEmptyExchangeVal.Error())
		return
	}

	symbol := r.PathValue("symbol")
	if len(symbol) == 0 {
		slog.Error("Failed to get symbol value from path: ", "error", domain.ErrEmptySymbolVal)
		senders.SendMsg(w, http.StatusBadRequest, domain.ErrEmptySymbolVal.Error())
		return
	}

	switch metric {
	case MetricHighest:
		period := r.URL.Query().Get("period")
		if period == "" {
			data, code, err = h.serv.GetHighestPrice(exchange, symbol)
			if err != nil {
				slog.Error("Failed to get highest price: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
				senders.SendMsg(w, code, err.Error())
				return
			}

		} else {
			data, code, err = h.serv.GetHighestPriceWithPeriod(exchange, symbol, period)
			if err != nil {
				slog.Error("Failed to get highest price: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
				senders.SendMsg(w, code, err.Error())
				return
			}
		}
		msg = fmt.Sprintf("Highest price for %s at %s duration {%s}: %.2f", symbol, exchange, period, data.Price)
	case MetricLowest:
		period := r.URL.Query().Get("period")

		if period == "" {
			data, code, err = h.serv.GetLowestPrice(exchange, symbol)
			if err != nil {
				slog.Error("Failed to get lowest price: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
				senders.SendMsg(w, code, err.Error())
				return
			}
		} else {
			data, code, err = h.serv.GetLowestPriceWithPeriod(exchange, symbol, period)
			if err != nil {
				slog.Error("Failed to get lowest price: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
				senders.SendMsg(w, code, err.Error())
				return
			}
		}
		msg = fmt.Sprintf("Lowest price for %s at %s duration {%s}: %.2f", symbol, exchange, period, data.Price)
	case MetricAverage:
		period := r.URL.Query().Get("period")
		if period == "" {
			data, code, err = h.serv.GetAveragePrice(exchange, symbol)
			if err != nil {
				slog.Error("Failed to get average price: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
				senders.SendMsg(w, code, err.Error())
				return
			}

		} else {
			data, code, err = h.serv.GetAveragePriceWithPeriod(exchange, symbol, period)
			if err != nil {
				slog.Error("Failed to get average price with period: ", "exchange", exchange, "symbol", symbol, "period", period, "error", err.Error())
				senders.SendMsg(w, code, err.Error())
				return
			}

		}
		msg = fmt.Sprintf("Average price for %s at %s duration {%s}: %.2f", symbol, exchange, period, data.Price)
	case MetricLatest:
		data, code, err = h.serv.GetLatestData(exchange, symbol)
		if err != nil {
			slog.Error("Failed to get latest data: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
			senders.SendMsg(w, code, err.Error())
			return
		}
		msg = fmt.Sprintf("Latest price for %s at %s: %.2f", symbol, exchange, data.Price)

	default:
		slog.Error("Failed to get data by metric: ", "exchange", "All", "symbol", symbol, "metric", metric, "error", domain.ErrInvalidMetricVal.Error())
		senders.SendMsg(w, code, domain.ErrInvalidMetricVal.Error())
		return
	}

	if err := senders.SendMetricData(w, code, data); err != nil {
		slog.Error("Failed to send JSON message: ", "data", data, "error", err.Error())
		senders.SendMsg(w, code, err.Error())
		return
	}

	slog.Info(msg)
}

// Core handler for processing metric-based queries across all exchanges
func (h *MarketDataHTTPHandler) ProcessMetricQueryByAll(w http.ResponseWriter, r *http.Request) {
	var (
		data     domain.Data
		exchange = "All"
		msg      string
		code     int = 200
		err      error
	)
	metric := r.PathValue("metric")
	if len(metric) == 0 {
		slog.Error("Failed to get metric value from path: ", "error", domain.ErrEmptyMetricVal.Error())
		senders.SendMsg(w, http.StatusBadRequest, domain.ErrEmptyMetricVal.Error())
		return
	}

	symbol := r.PathValue("symbol")
	if len(symbol) == 0 {
		slog.Error("Failed to get symbol value from path: ", "error", domain.ErrEmptyExchangeVal)
		senders.SendMsg(w, http.StatusBadRequest, domain.ErrEmptySymbolVal.Error())
		return
	}

	switch metric {
	case MetricHighest:
		period := r.URL.Query().Get("period")
		if period == "" {
			data, code, err = h.serv.GetHighestPrice(exchange, symbol)
			if err != nil {
				slog.Error("Failed to get highest price: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
				senders.SendMsg(w, code, err.Error())
				return
			}
		} else {
			data, code, err = h.serv.GetHighestPriceByAllExchangesWithPeriod(symbol, period)
			if err != nil {
				slog.Error("Failed to get highest price with period: ", "exchange", exchange, "symbol", symbol, "period", period, "error", err.Error())
				senders.SendMsg(w, code, err.Error())
				return
			}
		}
		msg = fmt.Sprintf("Highest price for %s at %s: %.2f", symbol, exchange, data.Price)
	case MetricLowest:
		period := r.URL.Query().Get("period")
		if period == "" {
			data, code, err = h.serv.GetLowestPrice(exchange, symbol)
			if err != nil {
				slog.Error("Failed to get lowest price: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
				senders.SendMsg(w, code, err.Error())
				return
			}
		} else {
			data, code, err = h.serv.GetLowestPriceByAllExchangesWithPeriod(symbol, period)
			if err != nil {
				slog.Error("Failed to get lowest price: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
				senders.SendMsg(w, code, err.Error())
				return
			}
		}

		msg = fmt.Sprintf("Lowest price for %s at %s: %.2f", symbol, exchange, data.Price)
	case MetricAverage:
		data, code, err = h.serv.GetAveragePrice(exchange, symbol)
		if err != nil {
			slog.Error("Failed to get average price: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
			senders.SendMsg(w, code, err.Error())
			return
		}

		msg = fmt.Sprintf("Average price for %s at %s: %.2f", symbol, exchange, data.Price)
	case MetricLatest:
		data, code, err = h.serv.GetLatestData(exchange, symbol)
		if err != nil {
			slog.Error("Failed to get latest data: ", "exchange", exchange, "symbol", symbol, "error", err.Error())
			senders.SendMsg(w, code, err.Error())
			return
		}

		msg = fmt.Sprintf("Latest price for %s at %s: %.2f", symbol, exchange, data.Price)
	default:
		slog.Error("Failed to get data by metric: ", "exchange", exchange, "symbol", symbol, "metric", metric, "error", domain.ErrInvalidMetricVal.Error())
		senders.SendMsg(w, code, domain.ErrInvalidMetricVal.Error())
		return
	}

	if err := senders.SendMetricData(w, code, data); err != nil {
		slog.Error("Failed to send JSON message: ", "data", data, "error", err.Error())
		senders.SendMsg(w, code, err.Error())
		return
	}
	slog.Info(msg)
}
