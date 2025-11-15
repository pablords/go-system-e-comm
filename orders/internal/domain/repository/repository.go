package repository

import "orders-go/internal/domain/entity"

type ProductRepository interface {
	Create(product *entity.Product) error
	FindByID(id string) (*entity.Product, error)
	FindAll() ([]entity.Product, error)
	Update(product *entity.Product) error
	Delete(id string) error
}

type OrderRepository interface {
	Create(order *entity.Order) error
	FindByID(id string) (*entity.Order, error)
	FindAll() ([]entity.Order, error)
	Update(order *entity.Order) error
	Delete(id string) error
}

type ItemRepository interface {
	Create(item *entity.Item) error
	FindByID(id string) (*entity.Item, error)
	FindByOrderID(orderID string) ([]entity.Item, error)
	Update(item *entity.Item) error
	Delete(id string) error
}
