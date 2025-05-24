package handlers

import (
	"fmt"
	"log"
	"marketflow/internal/api/senders"
	"marketflow/internal/domain"

	"net/http"
)

type SwitchModeHTTPHandler struct {
	serv domain.DataModeService
}

func NewSwitchModeHandler(serv domain.DataModeService) *SwitchModeHTTPHandler {
	return &SwitchModeHTTPHandler{serv: serv}
}

func (h *SwitchModeHTTPHandler) SwitchMode(w http.ResponseWriter, r *http.Request) {
	// Switch mode service logic call...
	mode := r.PathValue("mode")
	if err := h.serv.SwitchMode(mode); err != nil {
		log.Printf("Failed to switch mode: %s \n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Sending message to the client
	msg := fmt.Sprintf("Datafetcher mode switched to %s \n", mode)
	if err := senders.SendMsg(w, http.StatusOK, msg); err != nil {
		log.Printf("Failed to send message to the client: %s \n", err.Error())
		return
	}

	log.Print(msg)
}
