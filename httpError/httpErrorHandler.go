package httpError

import (
	"errors"
	"net/http"
)

type HTTPError struct {
	Error error
	Code  int
}

func Write(w http.ResponseWriter, err *HTTPError) {
	http.Error(w, err.Error.Error(), err.Code)
}

func New(code int, message string) *HTTPError {
	return &HTTPError{
		Error: errors.New(message),
		Code:  code,
	}
}

func NotFound(message string) *HTTPError {
	return New(http.StatusNotFound, message)
}
