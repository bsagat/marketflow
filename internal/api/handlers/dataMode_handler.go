package handlers

import (
	"fmt"
	"log/slog"
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

// Core handler for switching datafetcher mode
func (h *SwitchModeHTTPHandler) SwitchMode(w http.ResponseWriter, r *http.Request) {
	mode := r.PathValue("mode")
	if code, err := h.serv.SwitchMode(mode); err != nil {
		slog.Error("Failed to switch mode", "message", err.Error())
		senders.SendMsg(w, code, err.Error())
		return
	}

	// Sending message to the client
	msg := fmt.Sprintf("Datafetcher mode switched to %s", mode)
	senders.SendMsg(w, http.StatusOK, msg)
	slog.Info(msg)
}
