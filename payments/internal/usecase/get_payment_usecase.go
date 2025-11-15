package usecase

import (
	"fmt"
	"log/slog"
	"payments-go/internal/domain/entity"
	"payments-go/internal/domain/repository"
)

type GetPaymentUseCase struct {
	paymentRepo repository.PaymentRepository
}

func NewGetPaymentUseCase(paymentRepo repository.PaymentRepository) *GetPaymentUseCase {
	return &GetPaymentUseCase{
		paymentRepo: paymentRepo,
	}
}

func (uc *GetPaymentUseCase) Execute(paymentID string) (*entity.Payment, error) {
	if paymentID == "" {
		return nil, fmt.Errorf("payment_id cannot be empty")
	}

	slog.Info("Getting payment", "payment_id", paymentID)

	payment, err := uc.paymentRepo.FindByID(paymentID)
	if err != nil {
		slog.Error("Failed to get payment", "payment_id", paymentID, "error", err)
		return nil, err
	}

	return payment, nil
}
