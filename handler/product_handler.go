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
		response.Error(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.JSON(w, true, http.StatusOK, "Product Retrieved Successfully", products)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	number, err := strconv.Atoi(id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Bad Request")
		return
	}

	product, err := h.service.GetProductByID(number)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Product not found")
		return
	}
	response.JSON(w, true, http.StatusOK, "Product Retrieved Successfully", product)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		response.Error(w, http.StatusBadRequest, "Bad Request")
		return
	}

	result, err := h.service.CreateProduct(product)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Bad Request")
		return
	}
	response.JSON(w, true, http.StatusCreated, "Product Created Successfully", result)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	number, err := strconv.Atoi(id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Bad Request")
		return
	}
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		response.Error(w, http.StatusBadRequest, "Bad Request")
		return
	}
	product.ID = number
	result, err := h.service.UpdateProduct(product)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Product not found")
		return
	}
	response.JSON(w, true, http.StatusOK, "Product Updated Successfully", result)

}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	number, err := strconv.Atoi(id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Bad Request")
		return
	}
	err = h.service.DeleteProduct(number)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Product not found")
		return
	}
	response.JSON(w, true, http.StatusOK, "Product Deleted Successfully", nil)
}
