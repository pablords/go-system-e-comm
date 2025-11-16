package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"payments/internal/domain/repository"
)

type CancelPaymentUseCase struct {
	paymentRepo repository.PaymentRepository
}

func NewCancelPaymentUseCase(paymentRepo repository.PaymentRepository) *CancelPaymentUseCase {
	return &CancelPaymentUseCase{
		paymentRepo: paymentRepo,
	}
}

type CancelPaymentInput struct {
	PaymentID string
	Reason    string
}

func (uc *CancelPaymentUseCase) Execute(ctx context.Context, input CancelPaymentInput) error {
	if input.PaymentID == "" {
		return fmt.Errorf("payment_id cannot be empty")
	}

	slog.Info("Canceling payment", "payment_id", input.PaymentID, "reason", input.Reason)

	// Find payment
	payment, err := uc.paymentRepo.FindByID(ctx, input.PaymentID)
	if err != nil {
		slog.Error("Failed to find payment", "payment_id", input.PaymentID, "error", err)
		return err
	}

	// Cancel payment
	if err := payment.Cancel(input.Reason); err != nil {
		slog.Error("Failed to cancel payment", "payment_id", input.PaymentID, "error", err)
		return err
	}

	// Update payment in database
	if err := uc.paymentRepo.Update(ctx, payment); err != nil {
		slog.Error("Failed to update payment", "payment_id", input.PaymentID, "error", err)
		return fmt.Errorf("failed to update payment: %w", err)
	}

	slog.Info("Payment canceled successfully", "payment_id", input.PaymentID)

	return nil
}
