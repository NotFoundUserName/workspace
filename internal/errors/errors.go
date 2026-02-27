package errors

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewErrorResponse(err error, message string, code int) ErrorResponse {
	error := ""
	if err != nil {
		error = err.Error()
	}
	return ErrorResponse{
		Error:   error,
		Message: message,
		Code:    code,
	}
}

func WriteErrorResponse(w http.ResponseWriter, err error, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	response := NewErrorResponse(err, message, code)
	json.NewEncoder(w).Encode(response)
}

func BadRequest(w http.ResponseWriter, err error, message string) {
	WriteErrorResponse(w, err, message, http.StatusBadRequest)
}

func InternalServerError(w http.ResponseWriter, err error, message string) {
	WriteErrorResponse(w, err, message, http.StatusInternalServerError)
}

func MethodNotAllowed(w http.ResponseWriter) {
	WriteErrorResponse(w, nil, "Method not allowed", http.StatusMethodNotAllowed)
}

func NotFound(w http.ResponseWriter, message string) {
	WriteErrorResponse(w, nil, message, http.StatusNotFound)
}
