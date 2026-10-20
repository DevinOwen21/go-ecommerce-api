package service

import (
	"context"
	"database/sql"
	"errors"
	"go-ecommerce-api/repository"
	"time"
)

type CheckoutService struct {
	db          *sql.DB
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
	orderRepo   *repository.OrderRepository
	userRepo    *repository.UserRepository
}

func NewCheckoutService(db *sql.DB, cartRepo *repository.CartRepository, productRepo *repository.ProductRepository, orderRepo *repository.OrderRepository, userRepo *repository.UserRepository) *CheckoutService {
	return &CheckoutService{
		db:          db,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		orderRepo:   orderRepo,
		userRepo:    userRepo,
	}
}

func (s *CheckoutService) Checkout(ctx context.Context, userID int) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	err = s.userRepo.LockUserForCheckout(ctx, tx, userID)
	if err != nil {
		return 0, err
	}

	items, err := s.cartRepo.GetCartForCheckout(ctx, tx, userID)
	if err != nil {
		return 0, err
	}

	if len(items) == 0 {
		return 0, errors.New("cart is empty")
	}

	totalAmount := 0
	const maxAmount = 2147483647

	for _, item := range items {
		if item.Quantity <= 0 {
			return 0, errors.New("invalid quantity")
		}

		if item.Price <= 0 {
			return 0, errors.New("invalid price")
		}

		if item.Quantity > maxAmount/item.Price {
			return 0, errors.New("item subtotal exceeds maximum amount")
		}

		subTotal := item.Quantity * item.Price

		if subTotal > maxAmount-totalAmount {
			return 0, errors.New("order total exceeds maximum amount")
		}

		totalAmount += subTotal
	}

	for _, item := range items {
		err := s.productRepo.ReserveStock(ctx, tx, item.Quantity, item.ProductID)
		if err != nil {
			return 0, err
		}
	}

	expiresAt := time.Now().Add(10 * time.Minute)

	orderID, err := s.orderRepo.CreateOrder(ctx, tx, userID, totalAmount, expiresAt)
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		err := s.orderRepo.CreateOrderItem(ctx, tx, orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return 0, err
		}
	}

	err = s.cartRepo.ClearCart(ctx, tx, userID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return orderID, nil
}
