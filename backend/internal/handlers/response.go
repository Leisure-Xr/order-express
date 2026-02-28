package handlers

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type PaginatedResult[T any] struct {
	Items    []T `json:"items"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func OK(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, ApiResponse{Code: 200, Data: data, Message: "success"})
}

func Fail(w http.ResponseWriter, code int, message string) {
	writeJSON(w, code, ApiResponse{Code: code, Data: nil, Message: message})
}

func NotFound(w http.ResponseWriter, message string) {
	Fail(w, http.StatusNotFound, message)
}

func Unauthorized(w http.ResponseWriter) {
	Fail(w, http.StatusUnauthorized, "unauthorized")
}

func Forbidden(w http.ResponseWriter) {
	Fail(w, http.StatusForbidden, "forbidden")
}

func BadRequest(w http.ResponseWriter, message string) {
	Fail(w, http.StatusBadRequest, message)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
