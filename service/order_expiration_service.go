package service

import (
	"context"
	"database/sql"
	"errors"
	"go-ecommerce-api/model"
	"go-ecommerce-api/repository"
	"time"
)

type OrderExpirationService struct {
	db          *sql.DB
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
}

func NewOrderExpirationService(db *sql.DB, orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository) *OrderExpirationService {
	return &OrderExpirationService{
		db:          db,
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *OrderExpirationService) ExpireOrder(ctx context.Context, orderID int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	status, expiresAt, err := s.orderRepo.LockedOrderForExpiration(ctx, tx, orderID)
	if err != nil {
		return err
	}

	if status != model.OrderStatusPending {
		return nil
	}

	if time.Now().Before(expiresAt) {
		return nil
	}

	items, err := s.orderRepo.GetOrderItems(ctx, tx, orderID)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return errors.New("order has no items")
	}

	err = s.orderRepo.UpdateOrderStatus(ctx, tx, orderID, model.OrderStatusExpired)
	if err != nil {
		return err
	}

	for _, item := range items {
		err := s.productRepo.ReleaseReserveStock(ctx, tx, item.Quantity, item.ProductID)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *OrderExpirationService) GetExpiredOrderIds(ctx context.Context) ([]int, error) {
	expiredOrder, err := s.orderRepo.GetExpiredOrder(ctx)
	if err != nil {
		return nil, err
	}
	return expiredOrder, nil
}
