package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-task-api/httpError"
)

func LogToConsole(message string) {
	fmt.Println(message)
}

func ParseIDFromRequest(r *http.Request) (int, *httpError.HTTPError) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return 0, httpError.New(http.StatusBadRequest, "id value cannot be null")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, httpError.New(http.StatusBadRequest, "id must be a valid number")
	}

	return id, nil
}

func RouteLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		fmt.Printf("[REQUEST] %v -> %v took: %vms\n", r.Method, r.URL.Path, time.Since(start).Milliseconds())
	}
}
