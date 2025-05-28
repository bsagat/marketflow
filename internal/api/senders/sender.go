package senders

import (
	"encoding/json"
	"marketflow/internal/domain"
	"net/http"
	"time"
)

func SendMsg(w http.ResponseWriter, code int, msg string) error {
	data := struct {
		Code int    `json:"Code"`
		Msg  string `json:"Message"`
	}{
		Code: code,
		Msg:  msg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

func SendJSON(w http.ResponseWriter, code int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

func SendMetricData(w http.ResponseWriter, code int, rawdata domain.Data) error {
	data := struct {
		ExchangeName string  `json:"exchange"`
		Symbol       string  `json:"symbol"`
		Price        float64 `json:"price"`
		Timestamp    string  `json:"timestamp"` // Readable time :)
	}{
		ExchangeName: rawdata.ExchangeName,
		Symbol:       rawdata.Symbol,
		Price:        rawdata.Price,
		Timestamp:    time.Unix(0, rawdata.Timestamp*int64(time.Millisecond)).Format(time.ANSIC),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}
