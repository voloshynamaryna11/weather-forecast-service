package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"weather-forecast-service/internal/service"
)

type WeatherHandler struct{ svc service.WeatherService }

func New(svc service.WeatherService) *WeatherHandler { return &WeatherHandler{svc: svc} }

func (h *WeatherHandler) Get(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, `{"error":"city is required"}`, http.StatusBadRequest)
		return
	}
	resp, err := h.svc.Get(r.Context(), city)
	if errors.Is(err, service.ErrNotFound) {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
