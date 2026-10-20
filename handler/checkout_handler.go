package handler

import (
	"go-ecommerce-api/response"
	"go-ecommerce-api/service"
	"go-ecommerce-api/utils"
	"net/http"
	"strconv"
)

type CheckoutHandler struct {
	checkoutService *service.CheckoutService
}

func NewCheckoutHandler(checkoutService *service.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{
		checkoutService: checkoutService,
	}
}

func (handler *CheckoutHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, err := strconv.Atoi(user.UserID)
	if err != nil || userID < 0 {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := handler.checkoutService.Checkout(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	response.JSON(w, true, http.StatusCreated, "order created successfully", result)
}
