package repository

import (
	"database/sql"
	"errors"
	"go-ecommerce-api/model"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProducts(filter model.ProductListFilter) (model.ProductListResult, error) {
	whereQuery := " WHERE 1 = 1"
	countQuery := `SELECT COUNT(*) FROM products`

	productQuery := `SELECT id,name,description,price,stock,category,created_at,updated_at FROM products`

	countArgs := []interface{}{}
	productArgs := []interface{}{}

	if filter.Search != "" {
		whereQuery += ` AND name LIKE ?`

		searchValue := "%" + filter.Search + "%"

		countArgs = append(countArgs, searchValue)
		productArgs = append(productArgs, searchValue)
	}
	if filter.Category != nil {
		whereQuery += ` AND category = ?`

		countArgs = append(countArgs, *filter.Category)
		productArgs = append(productArgs, *filter.Category)
	}
	countQuery += whereQuery
	productQuery += whereQuery

	var total int

	err := r.db.QueryRow(
		countQuery,
		countArgs...,
	).Scan(&total)
	if err != nil {
		return model.ProductListResult{}, err
	}

	productQuery += ` LIMIT ? OFFSET ?`

	productArgs = append(
		productArgs,
		filter.Limit,
		filter.Offset,
	)

	rows, err := r.db.Query(
		productQuery,
		productArgs...,
	)
	if err != nil {
		return model.ProductListResult{}, err
	}
	defer rows.Close()

	products := make(
		[]model.Product,
		0,
		filter.Limit,
	)

	for rows.Next() {
		var product model.Product

		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.Category,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return model.ProductListResult{}, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return model.ProductListResult{}, err
	}

	return model.ProductListResult{
		Products: products,
		Total:    total,
	}, nil
}

func (r *ProductRepository) GetProductById(id int) (model.Product, error) {
	var product model.Product
	row := r.db.QueryRow("SELECT id, name, description, price, stock, category ,created_at, updated_at FROM products WHERE id = ?", id)
	err := row.Scan(
		&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.Category, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (r *ProductRepository) CreateProduct(product model.Product) (model.Product, error) {
	result, err := r.db.Exec("INSERT INTO products (name, description, price, stock, category) VALUES (?, ?, ?, ?, ?)", product.Name, product.Description, product.Price, product.Stock, product.Category)
	if err != nil {
		return model.Product{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return model.Product{}, err
	}
	return r.GetProductById(int(id))
}

func (r *ProductRepository) UpdateProduct(product model.Product) (model.Product, error) {
	result, err := r.db.Exec(
		"UPDATE products SET name = ?, description = ?, price = ?, stock = ? WHERE id = ?",
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ID,
	)
	if err != nil {
		return model.Product{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Product{}, err
	}

	if rowsAffected == 0 {
		existingProduct, err := r.GetProductById(product.ID)
		if err != nil {
			return model.Product{}, err
		}
		return existingProduct, nil
	}

	return r.GetProductById(product.ID)
}

func (r *ProductRepository) DeleteProduct(id int) error {
	result, err := r.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}
