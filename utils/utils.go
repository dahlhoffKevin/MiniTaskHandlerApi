package utils

import (
	"fmt"
	"net/http"
	"strconv"

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
		return 0, httpError.New(http.StatusBadRequest, "task id must be a valid number")
	}

	return id, nil
}
