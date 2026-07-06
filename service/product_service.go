package service

import (
	"errors"
	"go-ecommerce-api/model"
	"go-ecommerce-api/repository"
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

func (s *ProductService) GetProducts() ([]model.Product, error) {
	return s.repo.GetProducts()
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
