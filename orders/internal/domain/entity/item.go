package entity

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidQuantity = errors.New("quantity must be greater than zero")
	ErrInvalidProduct  = errors.New("product is required")
)

type Item struct {
	ID        string   `json:"id"`
	OrderID   string   `json:"order_id"`
	ProductID string   `json:"product_id"`
	Product   *Product `json:"product,omitempty"`
	Quantity  int      `json:"quantity"`
	UnitPrice float64  `json:"unit_price"`
	Total     float64  `json:"total"`
}

func NewItem(orderID, productID string, product *Product, quantity int) (*Item, error) {
	if product == nil {
		return nil, ErrInvalidProduct
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	item := &Item{
		ID:        uuid.New().String(),
		OrderID:   orderID,
		ProductID: productID,
		Product:   product,
		Quantity:  quantity,
		UnitPrice: product.Price,
	}

	item.CalculateTotal()
	return item, nil
}

func (i *Item) CalculateTotal() {
	i.Total = i.UnitPrice * float64(i.Quantity)
}

func (i *Item) UpdateQuantity(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	i.Quantity = quantity
	i.CalculateTotal()
	return nil
}
