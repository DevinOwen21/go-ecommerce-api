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

func (s *ProductService) GetProducts(pagination dto.PaginationRequest) ([]dto.ProductResponse, dto.PaginationResponse, error) {
	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}
	if pagination.Limit > 100 {
		pagination.Limit = 100
	}
	pagination.Offset = (pagination.Page - 1) * pagination.Limit
	result, err := s.repo.GetProducts(pagination)
	if err != nil {
		return nil, dto.PaginationResponse{}, err
	}
	totalPages := int(
		math.Ceil(
			float64(result.Total) / float64(pagination.Limit),
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
		Page:       pagination.Page,
		Limit:      pagination.Limit,
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
