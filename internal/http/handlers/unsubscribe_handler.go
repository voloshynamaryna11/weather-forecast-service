package handlers

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	"weather-forecast-service/internal/http/model"

	"github.com/go-chi/chi/v5"

	"weather-forecast-service/internal/service"
)

type UnsubscribeHandler struct{ svc service.SubscriptionService }

func NewUnsubscribeHandler(svc service.SubscriptionService) *UnsubscribeHandler {
	return &UnsubscribeHandler{svc: svc}
}

func (h *UnsubscribeHandler) Get(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		model.WriteError(w, http.StatusBadRequest, "token is required")
		return
	}

	if _, err := uuid.Parse(token); err != nil {
		model.WriteError(w, http.StatusBadRequest, "invalid token")
		return
	}
	if err := h.svc.Unsubscribe(r.Context(), token); errors.Is(err, service.ErrTokenNotFound) {
		model.WriteError(w, http.StatusNotFound, "token not found")
		return
	} else if err != nil {
		model.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	model.WriteJSON(w, http.StatusOK, map[string]string{"status": "unsubscribed"})
}
