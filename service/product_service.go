package service

import (
	"errors"
	"go-ecommerce-api/dto"
	"go-ecommerce-api/model"
	"go-ecommerce-api/repository"
	"math"
	"strings"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) GetProducts(req dto.ProductListRequest) ([]dto.ProductResponse, dto.PaginationResponse, error) {
	if req.Pagination.Page <= 0 {
		req.Pagination.Page = 1
	}
	if req.Pagination.Limit <= 0 {
		req.Pagination.Limit = 10
	}
	if req.Pagination.Limit > 100 {
		req.Pagination.Limit = 100
	}
	req.Pagination.Offset = (req.Pagination.Page - 1) * req.Pagination.Limit
	result, err := s.repo.GetProducts(req)
	if err != nil {
		return nil, dto.PaginationResponse{}, err
	}
	totalPages := int(
		math.Ceil(
			float64(result.Total) / float64(req.Pagination.Limit),
		),
	)
	products := make([]dto.ProductResponse, 0, len(result.Products))
	for _, product := range result.Products {
		products = append(products, dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		})
	}

	paginationResponse := dto.PaginationResponse{
		Page:       req.Pagination.Page,
		Limit:      req.Pagination.Limit,
		Total:      result.Total,
		TotalPages: totalPages,
	}

	return products, paginationResponse, nil
}

func (s *ProductService) GetProductByID(id int) (model.Product, error) {
	return s.repo.GetProductById(id)
}

func (s *ProductService) CreateProduct(product model.Product) (model.Product, error) {
	if strings.TrimSpace(product.Name) == "" {
		return model.Product{}, errors.New("product name is empty")
	} else if product.Price <= 0 {
		return model.Product{}, errors.New("product price cannot be empty or in negative value")
	} else if product.Stock < 0 {
		return model.Product{}, errors.New("product stock cannot be negative value")
	}
	return s.repo.CreateProduct(product)
}

func (s *ProductService) UpdateProduct(product model.Product) (model.Product, error) {
	if strings.TrimSpace(product.Name) == "" {
		return model.Product{}, errors.New("product name is empty")
	}
	if product.Price <= 0 {
		return model.Product{}, errors.New("product price cannot be empty or in negative value")
	}
	if product.Stock < 0 {
		return model.Product{}, errors.New("product stock cannot be negative value")
	}
	return s.repo.UpdateProduct(product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.DeleteProduct(id)
}
