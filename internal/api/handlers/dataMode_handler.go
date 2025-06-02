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

func (h *SwitchModeHTTPHandler) SwitchMode(w http.ResponseWriter, r *http.Request) {
	mode := r.PathValue("mode")
	if code, err := h.serv.SwitchMode(mode); err != nil {
		slog.Error("Failed to switch mode", "message", err.Error())
		if err := senders.SendMsg(w, code, err.Error()); err != nil {
			slog.Error("Failed to send message to the client", "error", err.Error())
		}
		return
	}

	// Sending message to the client
	msg := fmt.Sprintf("Datafetcher mode switched to %s", mode)
	if err := senders.SendMsg(w, http.StatusOK, msg); err != nil {
		slog.Error("Failed to send message to the client", "error", err.Error())
		return
	}
	slog.Info(msg)
}
