package model

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIResponse{
		StatusCode: status,
		Data:       data,
	})
}

func WriteError(w http.ResponseWriter, status int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIResponse{
		StatusCode: status,
		Error:      errMsg,
	})
}
