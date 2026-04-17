package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"plexus-bff-service-go/internal/app/apperrors"
	"plexus-bff-service-go/internal/app/dto"
)

func parseIntDefault(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func decodeOpenAPIJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		writeOpenAPIJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON request body"})
		return false
	}
	return true
}

func writeOpenAPIJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if value != nil {
		_ = json.NewEncoder(w).Encode(value)
	}
}

func writeOpenAPIError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case apperrors.IsValidation(err):
		status = http.StatusBadRequest
	case apperrors.IsNotFound(err):
		status = http.StatusNotFound
	case apperrors.IsConflict(err):
		status = http.StatusConflict
	}
	writeOpenAPIJSON(w, status, dto.ErrorResponse{Error: err.Error()})
}
