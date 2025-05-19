package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"weather-forecast-service/internal/domain/entity/enum"
	"weather-forecast-service/internal/http/model"
	"weather-forecast-service/internal/service"
)

type SubscribeRequest struct {
	Email     string `json:"email" form:"email"`
	City      string `json:"city"  form:"city"`
	Frequency string `json:"frequency" form:"frequency"`
}

type SubscriptionHandler struct {
	svc service.SubscriptionService
}

func NewSubscribeHandler(svc service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc}
}

func (h *SubscriptionHandler) Post(w http.ResponseWriter, r *http.Request) {
	var req SubscribeRequest

	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			model.WriteError(w, http.StatusBadRequest, "bad form")
			return
		}
		req.Email = r.Form.Get("email")
		req.City = r.Form.Get("city")
		req.Frequency = r.Form.Get("frequency")
		if req.Email == "" {
			model.WriteError(w, http.StatusBadRequest, "email is required")
		}
		if _, err := mail.ParseAddress(req.Email); err != nil {
			model.WriteError(w, http.StatusBadRequest, "email is invalid")
		}

		if req.City == "" || len(req.City) < 2 {
			model.WriteError(w, http.StatusBadRequest, "city is invalid")
		}

		if !enum.Frequency(req.Frequency).IsValid() {
			model.WriteError(w, http.StatusBadRequest, "frequency is invalid")
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			model.WriteError(w, http.StatusBadRequest, "bad json")
			return
		}
	}

	err := h.svc.Subscribe(r.Context(),
		req.Email, req.City, enum.Frequency(req.Frequency),
	)

	switch {
	case errors.Is(err, service.ErrAlreadyExists),
		errors.Is(err, service.ErrConfirmationNeeded):
		model.WriteError(w, http.StatusConflict, err.Error())

	case err != nil:
		model.WriteError(w, http.StatusInternalServerError, "internal error")

	default:
		model.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
