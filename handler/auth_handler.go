package handler

import (
	"encoding/json"
	"go-ecommerce-api/dto"
	"go-ecommerce-api/response"
	"go-ecommerce-api/service"
	"net/http"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Bad Request")
		return
	}

	result, err := h.service.Register(req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Bad Request")
		return
	}
	response.JSON(w, true, http.StatusCreated, "User registered successfully", result)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Bad Request")
		return
	}
	result, err := h.service.Login(req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Bad Request")
		return
	}
	response.JSON(w, true, http.StatusOK, "User login successfully", result)
}
