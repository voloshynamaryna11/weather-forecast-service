package handlers

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	"weather-forecast-service/internal/http/model"

	"github.com/go-chi/chi/v5"

	"weather-forecast-service/internal/service"
)

type ConfirmHandler struct {
	svc service.SubscriptionService
}

func NewConfirmHandler(svc service.SubscriptionService) *ConfirmHandler {
	return &ConfirmHandler{svc: svc}
}

func (h *ConfirmHandler) Get(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		model.WriteError(w, http.StatusBadRequest, "token is required")
		return
	}

	if _, err := uuid.Parse(token); err != nil {
		model.WriteError(w, http.StatusBadRequest, "invalid token")
		return
	}

	err := h.svc.Confirm(r.Context(), token)
	if errors.Is(err, service.ErrTokenNotFound) {
		model.WriteError(w, http.StatusNotFound, "token not found")
		return
	} else if err != nil {
		model.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	model.WriteJSON(w, http.StatusOK, map[string]string{"status": "confirmed"})
}
