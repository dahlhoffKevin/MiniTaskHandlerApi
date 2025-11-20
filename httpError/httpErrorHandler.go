package httpError

import (
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func New(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

func Write(w http.ResponseWriter, err *HTTPError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)

	_ = json.NewEncoder(w).Encode(err)
}

func NotFound(message string) *HTTPError {
	return New(http.StatusNotFound, message)
}
