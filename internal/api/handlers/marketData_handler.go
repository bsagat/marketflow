package handlers

import (
	"fmt"
	"net/http"
)

type MarketDataHTTPHandler struct{}

func NewMarketDataHandler() *MarketDataHTTPHandler {
	return &MarketDataHTTPHandler{}
}

// GET /prices/{metric}/{exchange}/{symbol}
func (h *MarketDataHTTPHandler) ProcessMetricQueryByExchange(w http.ResponseWriter, r *http.Request) {
	metric := r.PathValue("metric")
	if len(metric) == 0 {
		http.Error(w, "metric value is empty", http.StatusBadRequest)
		return
	}

	switch metric {
	case "highest":
	case "lowest":
	case "average":
	case "latest":
	default:
		http.Error(w, fmt.Sprintf("metric value is invalid %s , must be (highest, lowest, average)", metric), http.StatusBadRequest)
		return
	}

	exchange := r.PathValue("exchange")
	if len(exchange) == 0 {
		http.Error(w, "exchange value is empty", http.StatusBadRequest)
		return
	}

	switch exchange {
	case "exchange1":
	case "exchange2":
	case "exchange3":
	default:
		http.Error(w, fmt.Sprintf("exchange value is invalid %s , must be (exchange1, exchange2, exchange3)", exchange), http.StatusBadRequest)
		return
	}

	symbol := r.PathValue("symbol")
	if len(symbol) == 0 {
		http.Error(w, "symbol value is empty", http.StatusBadRequest)
		return
	}

	fmt.Println(metric, exchange, symbol)
	// Валидность данных нужно проверять в бизнес логике
}

func (h *MarketDataHTTPHandler) ProcessMetricQueryByAll(w http.ResponseWriter, r *http.Request) {
	metric := r.PathValue("metric")
	if len(metric) == 0 {
		http.Error(w, "metric value is empty", http.StatusBadRequest)
		return
	}

	switch metric {
	case "highest":
	case "lowest":
	case "average":
	case "latest":
	default:
		http.Error(w, fmt.Sprintf("metric value is invalid %s , must be (highest, lowest, average)", metric), http.StatusBadRequest)
		return
	}

	symbol := r.PathValue("symbol")
	if len(symbol) == 0 {
		http.Error(w, "symbol value is empty", http.StatusBadRequest)
		return
	}

	fmt.Println(metric, symbol)
	// Валидность данных нужно проверять в бизнес логике
}
