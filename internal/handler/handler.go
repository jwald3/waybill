package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Status  int    `json:"status"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ReadJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func GetIDParam(r *http.Request) (int64, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		return 0, errors.New("invalid path")
	}
	return strconv.ParseInt(parts[len(parts)-1], 10, 64)
}

func GetIDFromPath(path string) (int64, error) {
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return 0, errors.New("invalid path")
	}

	idStr := parts[len(parts)-1]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid ID format")
	}

	return id, nil
}

func getQueryIntParam(r *http.Request, key string, defaultValue int) int {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}

	return val
}
