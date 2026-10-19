package service

import (
	"database/sql"
	"errors"
	"go-ecommerce-api/dto"
	"go-ecommerce-api/model"
	"go-ecommerce-api/repository"
)

type CartService struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

func NewCartService(cartRepo *repository.CartRepository, productRepo *repository.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

var ErrInvalidQuantity = errors.New("quantity must be greater than zero")
var ErrInsufficientStock = errors.New("stock is insufficient")
var ErrProductNotFound = errors.New("product not found")
var ErrCartItemNotFound = errors.New("cart item not found")

func (s *CartService) AddCartItem(userID int, req dto.AddCartItemRequest) error {
	if req.Quantity <= 0 {
		return ErrInvalidQuantity
	}

	product, err := s.productRepo.GetProductById(req.ProductID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrProductNotFound
		}
		return err
	}

	cartQuantity, err := s.cartRepo.GetCartQuantity(userID, product.ID)
	if err != nil {
		return err
	}
	if cartQuantity+req.Quantity > product.Stock {
		return ErrInsufficientStock
	}

	err = s.cartRepo.AddCartItem(userID, product.ID, req.Quantity)
	if err != nil {
		return err
	}
	return nil
}

func (s *CartService) GetCart(userID int) ([]model.CartItem, error) {
	cart, err := s.cartRepo.GetCart(userID)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (s *CartService) UpdateCartItem(userID, productID, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	exists, err := s.cartRepo.CartItemExists(userID, productID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrCartItemNotFound
	}
	product, err := s.productRepo.GetProductById(productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrProductNotFound
		}
		return err
	}
	if quantity > product.Stock {
		return ErrInsufficientStock
	}
	err = s.cartRepo.UpdateCartItem(userID, productID, quantity)
	if err != nil {
		return err
	}
	return nil
}

func (s *CartService) DeleteCartItem(userID, productID int) error {
	exists, err := s.cartRepo.CartItemExists(userID, productID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrCartItemNotFound
	}
	err = s.cartRepo.DeleteCartItem(userID, productID)
	if err != nil {
		return err
	}
	return nil
}
