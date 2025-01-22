package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// a custom response format you can use in your own handlers. Successful operations can use the `Message` property
// while unsuccessful operations can use the `Error` property
type Response struct {
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// a general-purpose function you can use to package data into a JSON response.
// use this in your resource-specific handlers to package the data you want to send in responses
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// takes a request (r) and attempts to decode it into the shape defined in v
// you can use this in your resource-specific handlers to ensure that you're getting data that matches
// the domain model (or DTO)
func ReadJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// a parser for integer-based query parameters, such as limit, page, etc.
// you can use this as a way to help unpack information from a request URL
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
