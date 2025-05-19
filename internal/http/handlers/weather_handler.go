package handlers

import (
	"errors"
	"net/http"
	"weather-forecast-service/internal/http/model"
	"weather-forecast-service/internal/service"
)

type WeatherHandler struct {
	svc service.WeatherService
}

func NewWeatherHandler(svc service.WeatherService) *WeatherHandler {
	return &WeatherHandler{svc: svc}
}

func (h *WeatherHandler) Get(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" || len(city) < 2 {
		model.WriteError(w, http.StatusBadRequest, "invalid city")
		return
	}

	respData, err := h.svc.Get(r.Context(), city)
	if errors.Is(err, service.ErrNotFound) {
		model.WriteError(w, http.StatusNotFound, "City not found")
		return
	} else if err != nil {
		model.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	model.WriteJSON(w, http.StatusOK, respData)
}
