package handler

import (
	"encoding/json"
	"errors"
	"go-ecommerce-api/dto"
	"go-ecommerce-api/response"
	"go-ecommerce-api/service"
	"go-ecommerce-api/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CartHandler struct {
	service *service.CartService
}

func NewCartHandler(service *service.CartService) *CartHandler {
	return &CartHandler{service: service}
}

func (h *CartHandler) AddCartItem(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := strconv.Atoi(user.UserID)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.AddCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	err = h.service.AddCartItem(userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidQuantity):
			response.Error(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrInsufficientStock):
			response.Error(w, http.StatusConflict, err.Error())
		case errors.Is(err, service.ErrProductNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	response.JSON(w, true, http.StatusOK, "item added to cart successfully", nil)
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, err := strconv.Atoi(user.UserID)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	cartItems, err := h.service.GetCart(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	response.JSON(w, true, http.StatusOK, "item retrieved successfully", cartItems)
}

func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, err := strconv.Atoi(user.UserID)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	productID := chi.URLParam(r, "productID")
	productIDInt, err := strconv.Atoi(productID)
	if err != nil || productIDInt <= 0 {
		response.Error(w, http.StatusBadRequest, "invalid product id")
		return
	}

	var req dto.UpdateCartItemsRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	err = h.service.UpdateCartItem(userID, productIDInt, req.Quantity)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidQuantity):
			response.Error(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrInsufficientStock):
			response.Error(w, http.StatusConflict, err.Error())
		case errors.Is(err, service.ErrProductNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrCartItemNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	response.JSON(w, true, http.StatusOK, "item updated successfully", nil)
}

func (h *CartHandler) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, err := strconv.Atoi(user.UserID)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	productID := chi.URLParam(r, "productID")
	productIDInt, err := strconv.Atoi(productID)
	if err != nil || productIDInt <= 0 {
		response.Error(w, http.StatusBadRequest, "invalid product id")
		return
	}
	err = h.service.DeleteCartItem(userID, productIDInt)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCartItemNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	response.JSON(w, true, http.StatusOK, "item deleted successfully", nil)
}
