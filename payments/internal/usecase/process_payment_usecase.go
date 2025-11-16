package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"payments-go/internal/domain/entity"
	"payments-go/internal/domain/repository"

	"github.com/google/uuid"
)

type ProcessPaymentUseCase struct {
	paymentRepo repository.PaymentRepository
}

func NewProcessPaymentUseCase(paymentRepo repository.PaymentRepository) *ProcessPaymentUseCase {
	return &ProcessPaymentUseCase{
		paymentRepo: paymentRepo,
	}
}

type ProcessPaymentInput struct {
	OrderID       string
	Amount        float64
	PaymentMethod entity.PaymentMethod
	CustomerEmail string
	CustomerName  string
}

type ProcessPaymentOutput struct {
	PaymentID     string
	OrderID       string
	Status        entity.PaymentStatus
	Message       string
	TransactionID string
}

func (uc *ProcessPaymentUseCase) Execute(ctx context.Context, input ProcessPaymentInput) (*ProcessPaymentOutput, error) {
	slog.Info("Processing payment", "order_id", input.OrderID, "amount", input.Amount, "method", input.PaymentMethod)

	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Create new payment
	payment, err := entity.NewPayment(
		input.OrderID,
		input.Amount,
		input.PaymentMethod,
		input.CustomerEmail,
		input.CustomerName,
	)
	if err != nil {
		slog.Error("Failed to create payment", "error", err)
		return nil, err
	}

	// Simulate payment processing with external gateway
	transactionID := uuid.New().String()
	if err := payment.Process(transactionID); err != nil {
		slog.Error("Failed to process payment", "error", err)
		return nil, err
	}

	// Simulate payment approval (in production, this would be async)
	// Here we could call an external payment gateway
	approved := simulatePaymentGateway(input.PaymentMethod, input.Amount)

	if approved {
		if err := payment.Approve(); err != nil {
			slog.Error("Failed to approve payment", "error", err)
			return nil, err
		}
		slog.Info("Payment approved", "payment_id", payment.ID, "transaction_id", transactionID)
	} else {
		if err := payment.Decline(); err != nil {
			slog.Error("Failed to decline payment", "error", err)
			return nil, err
		}
		slog.Warn("Payment declined", "payment_id", payment.ID)
	}

	// Save payment to database
	slog.Info("About to save payment to database", "payment_id", payment.ID)
	if err := uc.paymentRepo.Create(ctx, payment); err != nil {
		slog.Error("Failed to save payment", "error", err)
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}
	slog.Info("Payment saved successfully", "payment_id", payment.ID)

	message := "Payment processed successfully"
	if payment.Status == entity.PaymentStatusDeclined {
		message = "Payment was declined by the payment gateway"
	}

	return &ProcessPaymentOutput{
		PaymentID:     payment.ID,
		OrderID:       payment.OrderID,
		Status:        payment.Status,
		Message:       message,
		TransactionID: payment.TransactionID,
	}, nil
}

// simulatePaymentGateway simulates a payment gateway response
// In production, this would make actual HTTP calls to payment providers
func simulatePaymentGateway(method entity.PaymentMethod, amount float64) bool {
	// Simulate 95% approval rate for valid payments
	// In production, you would integrate with:
	// - Stripe, PayPal, PagSeguro, etc for credit/debit cards
	// - Brazilian payment gateways for PIX and Boleto
	if amount <= 0 {
		return false
	}

	// Simple simulation: approve all payments under 10000
	return amount < 10000.0
}
