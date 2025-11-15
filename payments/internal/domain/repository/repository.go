package repository

import "payments-go/internal/domain/entity"

type PaymentRepository interface {
	Create(payment *entity.Payment) error
	FindByID(id string) (*entity.Payment, error)
	FindByOrderID(orderID string) ([]*entity.Payment, error)
	Update(payment *entity.Payment) error
	Delete(id string) error
	List() ([]*entity.Payment, error)
}
