package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"payments-go/internal/domain/entity"
	"payments-go/internal/domain/repository"
)

type ListPaymentsUseCase struct {
	paymentRepo repository.PaymentRepository
}

func NewListPaymentsUseCase(paymentRepo repository.PaymentRepository) *ListPaymentsUseCase {
	return &ListPaymentsUseCase{
		paymentRepo: paymentRepo,
	}
}

func (uc *ListPaymentsUseCase) Execute(ctx context.Context, orderID string) ([]*entity.Payment, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order_id cannot be empty")
	}

	slog.Info("Listing payments for order", "order_id", orderID)

	payments, err := uc.paymentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		slog.Error("Failed to list payments", "order_id", orderID, "error", err)
		return nil, err
	}

	return payments, nil
}
