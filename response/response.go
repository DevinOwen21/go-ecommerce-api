package response

import (
	"encoding/json"
	"go-ecommerce-api/dto"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginationJSONResponse struct {
	Success    bool                    `json:"success"`
	Message    string                  `json:"message"`
	Data       interface{}             `json:"data,omitempty"`
	Pagination *dto.PaginationResponse `json:"pagination,omitempty"`
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

func JSONPagination(w http.ResponseWriter, success bool, code int, message string, data interface{}, pagination *dto.PaginationResponse) {
	response := PaginationJSONResponse{
		Success:    success,
		Message:    message,
		Pagination: pagination,
		Data:       data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Error(w http.ResponseWriter, code int, message string) {
	JSON(w, false, code, message, nil)
}
