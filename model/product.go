package model

import (
	"errors"
	"strings"
	"time"
)

type ProductCategory int

const (
	CategoryElectronics ProductCategory = 1
)

type Product struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       int             `json:"price"`
	Stock       int             `json:"stock"`
	Category    ProductCategory `json:"category"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ProductListFilter struct {
	Page     int
	Limit    int
	Offset   int
	Search   string
	Category *ProductCategory
	SortBy   ProductSortBy
	Order    string
}

type ProductSortBy string

const (
	ProductSortByName      ProductSortBy = "name"
	ProductSortByPrice     ProductSortBy = "price"
	ProductSortByCreatedAt ProductSortBy = "created_at"
)

func ParseProductCategory(category string) (ProductCategory, error) {
	category = strings.TrimSpace(category)

	switch {
	case strings.EqualFold(category, "electronics"):
		return CategoryElectronics, nil
	default:
		return 0, errors.New("invalid product category")
	}
}

func (c ProductCategory) String() string {
	switch c {
	case CategoryElectronics:
		return "Electronics"
	default:
		return "Unknown"
	}
}

func (s ProductSortBy) IsValid() bool {
	switch s {
	case ProductSortByName, ProductSortByPrice, ProductSortByCreatedAt:
		return true
	default:
		return false
	}
}
