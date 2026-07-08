package handler

import (
	"encoding/json"
	"go-ecommerce-api/model"
	"go-ecommerce-api/response"
	"go-ecommerce-api/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetProducts()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, products)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	number, err := strconv.Atoi(id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	product, err := h.service.GetProductByID(number)
	if err != nil {
		response.Error(w, http.StatusNotFound, err)
		return
	}
	response.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.CreateProduct(product)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	response.JSON(w, http.StatusCreated, result)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	number, err := strconv.Atoi(id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	product.ID = number
	result, err := h.service.UpdateProduct(product)
	if err != nil {
		response.Error(w, http.StatusNotFound, err)
		return
	}
	response.JSON(w, http.StatusOK, result)

}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	number, err := strconv.Atoi(id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	err = h.service.DeleteProduct(number)
	if err != nil {
		response.Error(w, http.StatusNotFound, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
