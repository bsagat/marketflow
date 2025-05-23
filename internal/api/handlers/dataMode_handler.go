package handlers

import (
	"marketflow/internal/domain"
	"net/http"
)

type SwitchModeHTTPHandler struct {
	Datafetcher domain.DataFetcher
}

func NewSwitchModeHandler(Datafetcher domain.DataFetcher) *SwitchModeHTTPHandler {
	return &SwitchModeHTTPHandler{Datafetcher: Datafetcher}
}

func (h *SwitchModeHTTPHandler) SwitchMode(w http.ResponseWriter, r *http.Request) {
	mode := r.PathValue("mode")
	switch mode {
	case "test":
		h.Datafetcher.Close()

	case "live":
	default:
		// error mazafaka
	}
}
