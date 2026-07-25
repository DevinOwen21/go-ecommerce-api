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

	filter := model.ProductListFilter{
		Page:   req.Pagination.Page,
		Limit:  req.Pagination.Limit,
		Offset: req.Pagination.Offset,
		Search: req.Search,
	}

	if strings.TrimSpace(req.Category) != "" {
		category, err := model.ParseProductCategory(req.Category)
		if err != nil {
			return nil, dto.PaginationResponse{}, err
		}
		filter.Category = &category
	}

	result, err := s.repo.GetProducts(filter)
	if err != nil {
		return nil, dto.PaginationResponse{}, err
	}
	totalPages := int(
		math.Ceil(
			float64(result.Total) / float64(filter.Limit),
		),
	)
	products := make([]dto.ProductResponse, 0, len(result.Products))
	for _, product := range result.Products {
		products = append(products, dto.ToProductResponse(product))
	}

	paginationResponse := dto.PaginationResponse{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      result.Total,
		TotalPages: totalPages,
	}

	return products, paginationResponse, nil
}

func (s *ProductService) GetProductByID(id int) (dto.ProductResponse, error) {
	product, err := s.repo.GetProductById(id)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return dto.ToProductResponse(product), nil
}

func (s *ProductService) CreateProduct(req dto.CreateProductRequest) (dto.ProductResponse, error) {
	if strings.TrimSpace(req.Name) == "" {
		return dto.ProductResponse{}, errors.New("product name is empty")
	}
	if req.Price <= 0 {
		return dto.ProductResponse{}, errors.New("product price cannot be empty or in negative value")
	}
	if req.Stock < 0 {

		return dto.ProductResponse{}, errors.New("product stock cannot be negative value")
	}

	category, err := model.ParseProductCategory(req.Category)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	product := model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    category,
	}

	createdProduct, err := s.repo.CreateProduct(product)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return dto.ToProductResponse(createdProduct), nil
}

func (s *ProductService) UpdateProduct(product model.Product) (dto.ProductResponse, error) {
	if strings.TrimSpace(product.Name) == "" {

		return dto.ProductResponse{}, errors.New("product name is empty")
	}
	if product.Price <= 0 {
		return dto.ProductResponse{}, errors.New("product price cannot be empty or in negative value")
	}
	if product.Stock < 0 {
		return dto.ProductResponse{}, errors.New("product stock cannot be negative value")
	}

	updatedProduct, err := s.repo.UpdateProduct(product)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return dto.ToProductResponse(updatedProduct), nil
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.DeleteProduct(id)
}
