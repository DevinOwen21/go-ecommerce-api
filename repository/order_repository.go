package repository

import (
	"context"
	"database/sql"
	"errors"
	"go-ecommerce-api/model"
	"time"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, tx *sql.Tx, userID int, totalAmount int, expiresAt time.Time) (int, error) {
	result, err := tx.ExecContext(ctx, "INSERT INTO orders(user_id, total_amount, status, expires_at) VALUES (?, ?, ?, ?)", userID, totalAmount, model.OrderStatusPending, expiresAt)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *OrderRepository) CreateOrderItem(ctx context.Context, tx *sql.Tx, orderID int, productID int, quantity int, unitPrice int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	_, err := tx.ExecContext(ctx, "INSERT INTO order_items(order_id, product_id, quantity, unit_price) VALUES (?, ?, ?, ?)", orderID, productID, quantity, unitPrice)
	if err != nil {
		return err
	}
	return nil
}
