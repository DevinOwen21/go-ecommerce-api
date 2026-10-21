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

func (r *OrderRepository) GetExpiredOrder(ctx context.Context) ([]int, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id FROM orders WHERE status = 'pending' AND expires_at <= NOW()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *OrderRepository) LockedOrderForExpiration(ctx context.Context, tx *sql.Tx, orderID int) (model.OrderStatus, time.Time, error) {
	var status model.OrderStatus
	var expiresAt time.Time
	row := tx.QueryRowContext(ctx, "SELECT status, expires_at FROM orders WHERE id = ? FOR UPDATE", orderID)
	if err := row.Scan(&status, &expiresAt); err != nil {
		return status, expiresAt, err
	}
	return status, expiresAt, nil
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, tx *sql.Tx, orderID int, status model.OrderStatus) error {
	row, err := tx.ExecContext(ctx, "UPDATE orders SET status = ? WHERE id = ? AND status = 'pending'", status, orderID)
	if err != nil {
		return err
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("order status was not updated")
	}
	return nil
}

func (r *OrderRepository) GetOrderItems(ctx context.Context, tx *sql.Tx, orderID int) ([]model.OrderItem, error) {
	rows, err := tx.QueryContext(ctx, "SELECT id, order_id, product_id, quantity, unit_price FROM order_items WHERE order_id = ?", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItems []model.OrderItem

	for rows.Next() {
		var item model.OrderItem
		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orderItems, nil
}
