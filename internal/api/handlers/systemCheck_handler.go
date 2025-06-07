package handlers

import (
	"log/slog"
	"marketflow/internal/api/senders"
	"net/http"
)

// Core handler for service health checking
func (h *SwitchModeHTTPHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	res := h.serv.CheckHealth()

	if err := senders.SendJSON(w, http.StatusOK, res); err != nil {
		slog.Error("Failed to send checkhealth data: " + err.Error())
		senders.SendMsg(w, http.StatusInternalServerError, err.Error())
	}
}
