package senders

import (
	"encoding/json"
	"net/http"
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
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}
