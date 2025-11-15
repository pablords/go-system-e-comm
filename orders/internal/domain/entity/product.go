package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidProductName  = errors.New("product name is required")
	ErrInvalidProductPrice = errors.New("product price must be greater than zero")
)

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewProduct(name, description string, price float64, stock int) (*Product, error) {
	product := &Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := product.Validate(); err != nil {
		return nil, err
	}

	return product, nil
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return ErrInvalidProductName
	}
	if p.Price <= 0 {
		return ErrInvalidProductPrice
	}
	return nil
}

func (p *Product) UpdateStock(quantity int) error {
	if p.Stock+quantity < 0 {
		return errors.New("insufficient stock")
	}
	p.Stock += quantity
	p.UpdatedAt = time.Now()
	return nil
}
