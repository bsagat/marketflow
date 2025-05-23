package handlers

import (
	"encoding/json"
	"log"
	"marketflow/internal/domain"
	"net/http"
)

type SystemHealthHTTPHandler struct {
	serv domain.SystemHealthServ
}

func NewSystemHealthHandler(serv domain.SystemHealthServ) *SystemHealthHTTPHandler {
	return &SystemHealthHTTPHandler{serv: serv}
}

func (h *SystemHealthHTTPHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	res := h.serv.CheckHealth()

	// Marshalling result to JSON
	jsonData, err := json.Marshal(res)
	if err != nil {
		log.Println("Failed to marshal checkhealth data: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sending JSON data to client
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonData); err != nil {
		log.Println("Failed to send checkhealth data: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
