package handlers

import "net/http"

type SystemHealthHTTPHandler struct{}

func NewSystemHealthHandler() *SystemHealthHTTPHandler {
	return &SystemHealthHTTPHandler{}
}

func (h *SystemHealthHTTPHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {

}
