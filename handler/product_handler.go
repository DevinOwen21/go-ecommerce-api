package handler

import (
	"encoding/json"
	"go-ecommerce-api/dto"
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

func (h *ProductHandler) GetProducts(
	w http.ResponseWriter,
	r *http.Request,
) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")
	categoryStr := r.URL.Query().Get("category")
	sortBy := r.URL.Query().Get("sort_by")
	order := r.URL.Query().Get("order")
	minPrice := r.URL.Query().Get("min_price")
	maxPrice := r.URL.Query().Get("max_price")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	var minPriceValue *int
	if minPrice != "" {
		parsedMinPrice, err := strconv.Atoi(minPrice)

		if err != nil {
			response.Error(
				w,
				http.StatusBadRequest,
				"invalid min_price",
			)
			return
		}
		minPriceValue = &parsedMinPrice
	}

	var maxPriceValue *int
	if maxPrice != "" {
		parsedMaxPrice, err := strconv.Atoi(maxPrice)

		if err != nil {
			response.Error(
				w,
				http.StatusBadRequest,
				"invalid max_price",
			)
			return
		}
		maxPriceValue = &parsedMaxPrice
	}

	req := dto.ProductListRequest{
		Pagination: dto.PaginationRequest{
			Page:  page,
			Limit: limit,
		},
		Search:   search,
		Category: categoryStr,
		SortBy:   sortBy,
		Order:    order,
		MinPrice: minPriceValue,
		MaxPrice: maxPriceValue,
	}

	products, paginationResp, err :=
		h.service.GetProducts(req)
	if err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	response.JSONPagination(
		w,
		true,
		http.StatusOK,
		"Products retrieved successfully",
		products,
		&paginationResp,
	)
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
	var req dto.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Bad Request")
		return
	}

	result, err := h.service.CreateProduct(req)
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
