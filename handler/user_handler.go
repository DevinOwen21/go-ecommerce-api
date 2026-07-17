package handler

import (
	"go-ecommerce-api/response"
	"go-ecommerce-api/service"
	"go-ecommerce-api/utils"
	"net/http"
	"strconv"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	number, err := strconv.Atoi(userID)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	user, err := h.service.GetProfile(number)
	if err != nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}
	response.JSON(w, true, http.StatusOK, "user Retrieved Successfully", user)

}
