package handlers

import (
	"fmt"
	"net/http"
)

type SwitchModeHTTPHandler struct{}

func NewSwitchModeHandler() *SwitchModeHTTPHandler {
	return &SwitchModeHTTPHandler{}
}

func (h *SwitchModeHTTPHandler) SwitchMode(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.PathValue("mode"))
}
