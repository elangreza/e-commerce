package rest

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	errs "github.com/elangreza/e-commerce/api/internal/error"
)

func sendSuccessResponse(w http.ResponseWriter, status int, res any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if val, ok := res.(string); ok {
		json.NewEncoder(w).Encode(map[string]any{"message": val})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"data": res})
}

type APIError struct {
	Message string `json:"error"`
}

func sendErrorResponse(w http.ResponseWriter, status int, err error) {
	var apiErr APIError
	switch {
	case errors.As(err, &errs.InvalidCredential{}):
		slog.Error("handler", "service", err.Error())
		status = errs.InvalidCredential{}.HttpStatusCode()
		apiErr.Message = err.Error()
	case errors.As(err, &errs.AlreadyExist{}):
		slog.Error("handler", "service", err.Error())
		status = errs.AlreadyExist{}.HttpStatusCode()
		apiErr.Message = err.Error()
	case errors.As(err, &errs.NotFound{}):
		slog.Error("handler", "service", err.Error())
		status = errs.NotFound{}.HttpStatusCode()
		apiErr.Message = err.Error()
	case errors.As(err, &errs.ValidationError{}):
		slog.Error("handler", "request", err.Error())
		status = errs.ValidationError{}.HttpStatusCode()
		apiErr.Message = err.Error()
	case errors.As(err, &errs.MethodNotAllowedError{}):
		slog.Error("handler", "request", err.Error())
		status = errs.MethodNotAllowedError{}.HttpStatusCode()
		apiErr.Message = err.Error()
	case status == http.StatusBadRequest, status == http.StatusUnauthorized, status == http.StatusForbidden:
		slog.Error("handler", "request", err.Error())
		apiErr.Message = err.Error()
	default:
		slog.Error("handler", "service", err.Error())
		status = http.StatusInternalServerError
		apiErr.Message = "server error"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(apiErr)
}
