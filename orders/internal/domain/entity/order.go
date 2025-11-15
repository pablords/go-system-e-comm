package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCanceled  OrderStatus = "canceled"
	OrderStatusCompleted OrderStatus = "completed"
)

var (
	ErrEmptyOrder         = errors.New("order must have at least one item")
	ErrItemNotFound       = errors.New("item not found in order")
	ErrInvalidOrderStatus = errors.New("invalid order status")
)

type Order struct {
	ID        string      `json:"id"`
	Status    OrderStatus `json:"status"`
	Items     []Item      `json:"items"`
	Total     float64     `json:"total"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func NewOrder() *Order {
	return &Order{
		ID:        uuid.New().String(),
		Status:    OrderStatusPending,
		Items:     []Item{},
		Total:     0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (o *Order) AddItem(item *Item) {
	// Check if item already exists, if so, update quantity
	for i, existingItem := range o.Items {
		if existingItem.ProductID == item.ProductID {
			o.Items[i].Quantity += item.Quantity
			o.Items[i].CalculateTotal()
			o.CalculateTotal()
			o.UpdatedAt = time.Now()
			return
		}
	}

	// Add new item
	o.Items = append(o.Items, *item)
	o.CalculateTotal()
	o.UpdatedAt = time.Now()
}

func (o *Order) RemoveItem(itemID string) error {
	for i, item := range o.Items {
		if item.ID == itemID {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
			o.CalculateTotal()
			o.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrItemNotFound
}

func (o *Order) UpdateItemQuantity(itemID string, quantity int) error {
	for i, item := range o.Items {
		if item.ID == itemID {
			if err := o.Items[i].UpdateQuantity(quantity); err != nil {
				return err
			}
			o.CalculateTotal()
			o.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrItemNotFound
}

func (o *Order) CalculateTotal() {
	total := 0.0
	for _, item := range o.Items {
		total += item.Total
	}
	o.Total = total
}

func (o *Order) PrepareForPayment() error {
	if len(o.Items) == 0 {
		return ErrEmptyOrder
	}
	o.CalculateTotal()
	return nil
}

func (o *Order) UpdateStatus(status OrderStatus) error {
	validStatuses := map[OrderStatus]bool{
		OrderStatusPending:   true,
		OrderStatusPaid:      true,
		OrderStatusCanceled:  true,
		OrderStatusCompleted: true,
	}

	if !validStatuses[status] {
		return ErrInvalidOrderStatus
	}

	o.Status = status
	o.UpdatedAt = time.Now()
	return nil
}
