package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"orders-go/internal/domain/entity"
	"orders-go/internal/domain/repository"
	"orders-go/internal/infra/grpc/client"
	pb "orders-go/proto"
	"time"
)

type OrderItemInput struct {
	ProductID string
	Quantity  int
	Price     float64
}

type CreateOrderInput struct {
	CustomerEmail string
	CustomerName  string
	Items         []OrderItemInput
	PaymentMethod int32 // 1=CREDIT_CARD, 2=DEBIT_CARD, 3=PIX, 4=BOLETO, 5=PAYPAL
}

type CreateOrderOutput struct {
	OrderID   string
	Total     float64
	Status    string
	PaymentID string
}

type CreateOrderUseCase struct {
	orderRepo     repository.OrderRepository
	productRepo   repository.ProductRepository
	paymentClient *client.PaymentClient
	logger        *slog.Logger
}

func NewCreateOrderUseCase(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	paymentClient *client.PaymentClient,
	logger *slog.Logger,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:     orderRepo,
		productRepo:   productRepo,
		paymentClient: paymentClient,
		logger:        logger,
	}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, input CreateOrderInput) (*CreateOrderOutput, error) {
	// 1. Criar o pedido
	order := entity.NewOrder()

	// 2. Adicionar items ao pedido e validar/criar produtos
	for _, itemInput := range input.Items {
		// Verificar se o produto existe
		product, err := uc.productRepo.FindByID(itemInput.ProductID)
		if err != nil {
			// Produto não existe, criar um produto temporário para o pedido
			uc.logger.Warn("Product not found, creating temporary product",
				"product_id", itemInput.ProductID,
			)

			now := time.Now()
			product = &entity.Product{
				ID:        itemInput.ProductID,
				Name:      "Product " + itemInput.ProductID,
				Price:     itemInput.Price,
				CreatedAt: now,
				UpdatedAt: now,
			}

			if err := uc.productRepo.Create(product); err != nil {
				uc.logger.Error("Failed to create product", "error", err)
				return nil, fmt.Errorf("failed to create product: %w", err)
			}
		}

		item, err := entity.NewItem(order.ID, itemInput.ProductID, product, itemInput.Quantity)
		if err != nil {
			uc.logger.Error("Failed to create item", "error", err)
			return nil, fmt.Errorf("failed to create item: %w", err)
		}

		order.AddItem(item)
	}

	// Validar pedido (verificar se tem itens)
	if len(order.Items) == 0 {
		uc.logger.Error("Order has no items")
		return nil, entity.ErrEmptyOrder
	}

	// 3. Salvar pedido no banco (isso já salva os items também)
	if err := uc.orderRepo.Create(order); err != nil {
		uc.logger.Error("Failed to save order", "error", err)
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	uc.logger.Info("Order created successfully",
		"order_id", order.ID,
		"total", order.Total,
	)

	// 4. Processar pagamento via gRPC
	paymentResponse, err := uc.paymentClient.ProcessPayment(
		ctx,
		order.ID,
		order.Total,
		input.PaymentMethod,
		input.CustomerEmail,
		input.CustomerName,
	)
	if err != nil {
		// Se falhar, marcar pedido como falha no pagamento
		order.Status = "payment_failed"
		_ = uc.orderRepo.Update(order)

		uc.logger.Error("Payment processing failed",
			"error", err,
			"order_id", order.ID,
		)
		return nil, fmt.Errorf("payment processing failed: %w", err)
	}

	// 5. Atualizar status do pedido baseado no pagamento
	switch paymentResponse.Status {
	case pb.PaymentStatus_PAYMENT_STATUS_APPROVED:
		order.Status = entity.OrderStatusPaid
	case pb.PaymentStatus_PAYMENT_STATUS_PROCESSING:
		order.Status = entity.OrderStatusPending
	case pb.PaymentStatus_PAYMENT_STATUS_DECLINED:
		order.Status = entity.OrderStatusCanceled
	default:
		order.Status = entity.OrderStatusPending
	}

	if err := uc.orderRepo.Update(order); err != nil {
		uc.logger.Error("Failed to update order status", "error", err)
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	uc.logger.Info("Order and payment processed successfully",
		"order_id", order.ID,
		"payment_id", paymentResponse.PaymentId,
		"payment_status", paymentResponse.Status,
	)

	return &CreateOrderOutput{
		OrderID:   order.ID,
		Total:     order.Total,
		Status:    string(order.Status),
		PaymentID: paymentResponse.PaymentId,
	}, nil
}
