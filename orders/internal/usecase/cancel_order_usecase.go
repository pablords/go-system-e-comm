package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"orders/internal/domain/entity"
	"orders/internal/domain/repository"
	"orders/internal/infra/grpc/client"
)

type CancelOrderUseCase struct {
	orderRepo     repository.OrderRepository
	paymentClient *client.PaymentClient
	logger        *slog.Logger
}

func NewCancelOrderUseCase(
	orderRepo repository.OrderRepository,
	paymentClient *client.PaymentClient,
	logger *slog.Logger,
) *CancelOrderUseCase {
	return &CancelOrderUseCase{
		orderRepo:     orderRepo,
		paymentClient: paymentClient,
		logger:        logger,
	}
}

func (uc *CancelOrderUseCase) Execute(ctx context.Context, orderID string, paymentID string) error {
	// 1. Buscar pedido
	order, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		uc.logger.Error("Failed to find order", "error", err, "order_id", orderID)
		return fmt.Errorf("failed to find order: %w", err)
	}

	// 2. Cancelar pagamento via gRPC
	if paymentID != "" {
		_, err := uc.paymentClient.CancelPayment(ctx, paymentID)
		if err != nil {
			uc.logger.Error("Failed to cancel payment",
				"error", err,
				"payment_id", paymentID,
				"order_id", orderID,
			)
			// Continuar mesmo se falhar o cancelamento do pagamento
		} else {
			uc.logger.Info("Payment canceled successfully",
				"payment_id", paymentID,
				"order_id", orderID,
			)
		}
	}

	// 3. Atualizar status do pedido
	order.Status = entity.OrderStatusCanceled
	if err := uc.orderRepo.Update(order); err != nil {
		uc.logger.Error("Failed to update order status", "error", err)
		return fmt.Errorf("failed to update order status: %w", err)
	}

	uc.logger.Info("Order canceled successfully", "order_id", orderID)

	return nil
}
