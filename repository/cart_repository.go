package repository

import (
	"database/sql"
	"go-ecommerce-api/model"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetCartQuantity(userID, productID int) (int, error) {
	var userCart model.UserCart
	row := r.db.QueryRow("SELECT quantity FROM user_carts where user_id = ? AND product_id = ?", userID, productID)
	err := row.Scan(&userCart.Quantity)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return userCart.Quantity, nil
}

func (r *CartRepository) AddCartItem(userID, productID, quantity int) error {
	_, err := r.db.Exec("INSERT INTO user_carts (user_id, product_id, quantity) VALUES (?,?,?) ON DUPLICATE KEY UPDATE quantity = quantity + ?", userID, productID, quantity, quantity)
	if err != nil {
		return err
	}
	return nil
}

func (r *CartRepository) GetCart(userID int) ([]model.CartItem, error) {
	rows, err := r.db.Query("SELECT uc.product_id, p.name, p.price, uc.quantity FROM user_carts uc INNER JOIN products p ON uc.product_id = p.id WHERE uc.user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.CartItem{}
	for rows.Next() {
		var item model.CartItem
		if err := rows.Scan(&item.ProductID, &item.Name, &item.Price, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CartRepository) UpdateCartItem(userID, productID, quantity int) error {
	_, err := r.db.Exec("UPDATE user_carts SET quantity = ? WHERE user_id = ? AND product_id = ?", quantity, userID, productID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CartRepository) CartItemExists(userID, productID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS( SELECT 1 FROM user_carts WHERE user_id = ? and product_id = ?)", userID, productID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *CartRepository) DeleteCartItem(userID, productID int) error {
	_, err := r.db.Exec("DELETE FROM user_carts WHERE user_id = ? AND product_id = ?", userID, productID)
	if err != nil {
		return err
	}
	return nil
}
