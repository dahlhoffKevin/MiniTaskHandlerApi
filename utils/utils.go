package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-task-api/httpError"

	"github.com/google/uuid"
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

func ParseUUIDFromRequest(r *http.Request) (uuid.UUID, *httpError.HTTPError) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return uuid.UUID{}, httpError.New(http.StatusBadRequest, "id value cannot be null")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.UUID{}, httpError.New(http.StatusBadRequest, "id is not a valid uuid")
	}

	return id, nil
}

func ParseAndValidateUUID(id string) (uuid.UUID, *httpError.HTTPError) {
	if id == "" {
		return uuid.UUID{}, httpError.New(http.StatusBadRequest, "id value cannot be null")
	}

	errValidate := uuid.Validate(id)
	if errValidate != nil {
		return uuid.UUID{}, httpError.New(http.StatusBadRequest, "id ist not a valid uuid: "+errValidate.Error())
	}

	parsedUUID, errParse := uuid.Parse(id)
	if errParse != nil {
		return uuid.UUID{}, httpError.New(http.StatusInternalServerError, "failed to parse id to uuid: "+errParse.Error())
	}

	return parsedUUID, nil
}

func RouteLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		fmt.Printf("[REQUEST] %v -> %v took: %vms\n", r.Method, r.URL.Path, time.Since(start).Milliseconds())
	}
}
