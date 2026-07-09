package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func newResponse(success bool, message string, data interface{}) *Response {
	return &Response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

func JSON(w http.ResponseWriter, success bool, code int, message string, data interface{}) {
	response := newResponse(success, message, data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Error(w http.ResponseWriter, code int, message string) {
	JSON(w, false, code, message, nil)
}
