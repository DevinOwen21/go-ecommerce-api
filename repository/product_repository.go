package repository

import (
	"database/sql"
	"errors"
	"go-ecommerce-api/dto"
	"go-ecommerce-api/model"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProducts(pagination dto.PaginationRequest) (model.ProductListResult, error) {
	var total int
	count := r.db.QueryRow("SELECT COUNT(*) FROM products")
	err := count.Scan(&total)
	if err != nil {
		return model.ProductListResult{}, err
	}
	rows, err := r.db.Query("SELECT id, name, description, price, stock, created_at, updated_at FROM products LIMIT ? OFFSET ?", pagination.Limit, pagination.Offset)
	if err != nil {
		return model.ProductListResult{}, err
	}
	defer rows.Close()

	products := []model.Product{}
	for rows.Next() {
		var product model.Product
		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
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
	result := model.ProductListResult{
		Products: products,
		Total:    total,
	}
	return result, nil
}

func (r *ProductRepository) GetProductById(id int) (model.Product, error) {
	var product model.Product
	row := r.db.QueryRow("SELECT id, name, description, price, stock, created_at, updated_at FROM products WHERE id = ?", id)
	err := row.Scan(
		&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (r *ProductRepository) CreateProduct(product model.Product) (model.Product, error) {
	result, err := r.db.Exec("INSERT INTO products (name, description, price, stock) VALUES (?, ?, ?, ?)", product.Name, product.Description, product.Price, product.Stock)
	if err != nil {
		return model.Product{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return model.Product{}, err
	}
	product.ID = int(id)
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
