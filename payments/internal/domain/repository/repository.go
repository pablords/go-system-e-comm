package repository

import (
	"context"
	"payments-go/internal/domain/entity"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *entity.Payment) error
	FindByID(ctx context.Context, id string) (*entity.Payment, error)
	FindByOrderID(ctx context.Context, orderID string) ([]*entity.Payment, error)
	Update(ctx context.Context, payment *entity.Payment) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*entity.Payment, error)
}
