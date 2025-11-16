package client

// Este arquivo demonstra como criar um cliente gRPC para o Payment Service
// que pode ser usado no serviço Orders.
//
// Para usar este cliente no Orders:
// 1. Copie o arquivo proto/payment.proto para o projeto orders
// 2. Gere o código gRPC no orders: make proto
// 3. Adicione as dependências gRPC no go.mod do orders
// 4. Copie este arquivo para orders/internal/infra/grpc/client/payment_client.go
// 5. Ajuste os imports conforme necessário

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pb "payments/proto" // Ajuste para o import correto no projeto orders

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// PaymentClient é o cliente gRPC para o Payment Service
type PaymentClient struct {
	client pb.PaymentServiceClient
	conn   *grpc.ClientConn
}

// NewPaymentClient cria uma nova instância do cliente de pagamento
func NewPaymentClient(address string) (*PaymentClient, error) {
	// Configurar opções de conexão
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	// Estabelecer conexão com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service at %s: %w", address, err)
	}

	client := pb.NewPaymentServiceClient(conn)

	slog.Info("Connected to payment service", "address", address)

	return &PaymentClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close fecha a conexão com o serviço de pagamento
func (c *PaymentClient) Close() error {
	if c.conn != nil {
		slog.Info("Closing connection to payment service")
		return c.conn.Close()
	}
	return nil
}

// ProcessPayment processa um novo pagamento
func (c *PaymentClient) ProcessPayment(
	ctx context.Context,
	orderID string,
	amount float64,
	method pb.PaymentMethod,
	email, name string,
) (*pb.ProcessPaymentResponse, error) {

	req := &pb.ProcessPaymentRequest{
		OrderId:       orderID,
		Amount:        amount,
		PaymentMethod: method,
		CustomerEmail: email,
		CustomerName:  name,
	}

	slog.Info("Processing payment",
		"order_id", orderID,
		"amount", amount,
		"method", method.String(),
	)

	resp, err := c.client.ProcessPayment(ctx, req)
	if err != nil {
		slog.Error("Failed to process payment", "error", err)
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	slog.Info("Payment processed successfully",
		"payment_id", resp.PaymentId,
		"status", resp.Status.String(),
		"transaction_id", resp.TransactionId,
	)

	return resp, nil
}

// GetPayment busca detalhes de um pagamento específico
func (c *PaymentClient) GetPayment(ctx context.Context, paymentID string) (*pb.GetPaymentResponse, error) {
	req := &pb.GetPaymentRequest{
		PaymentId: paymentID,
	}

	slog.Info("Getting payment details", "payment_id", paymentID)

	resp, err := c.client.GetPayment(ctx, req)
	if err != nil {
		slog.Error("Failed to get payment", "error", err)
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return resp, nil
}

// CancelPayment cancela um pagamento pendente
func (c *PaymentClient) CancelPayment(ctx context.Context, paymentID, reason string) (*pb.CancelPaymentResponse, error) {
	req := &pb.CancelPaymentRequest{
		PaymentId: paymentID,
		Reason:    reason,
	}

	slog.Info("Canceling payment", "payment_id", paymentID, "reason", reason)

	resp, err := c.client.CancelPayment(ctx, req)
	if err != nil {
		slog.Error("Failed to cancel payment", "error", err)
		return nil, fmt.Errorf("failed to cancel payment: %w", err)
	}

	if !resp.Success {
		slog.Warn("Payment cancellation was not successful", "message", resp.Message)
	}

	return resp, nil
}

// ListPayments lista todos os pagamentos de um pedido
func (c *PaymentClient) ListPayments(ctx context.Context, orderID string) (*pb.ListPaymentsResponse, error) {
	req := &pb.ListPaymentsRequest{
		OrderId: orderID,
	}

	slog.Info("Listing payments for order", "order_id", orderID)

	resp, err := c.client.ListPayments(ctx, req)
	if err != nil {
		slog.Error("Failed to list payments", "error", err)
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}

	slog.Info("Payments listed successfully", "count", len(resp.Payments))

	return resp, nil
}

// Helper methods

// IsPaymentApproved verifica se o pagamento foi aprovado
func IsPaymentApproved(status pb.PaymentStatus) bool {
	return status == pb.PaymentStatus_PAYMENT_STATUS_APPROVED
}

// IsPaymentDeclined verifica se o pagamento foi recusado
func IsPaymentDeclined(status pb.PaymentStatus) bool {
	return status == pb.PaymentStatus_PAYMENT_STATUS_DECLINED
}

// IsPaymentPending verifica se o pagamento está pendente
func IsPaymentPending(status pb.PaymentStatus) bool {
	return status == pb.PaymentStatus_PAYMENT_STATUS_PENDING ||
		status == pb.PaymentStatus_PAYMENT_STATUS_PROCESSING
}

// ConvertPaymentMethodString converte string para PaymentMethod
func ConvertPaymentMethodString(method string) pb.PaymentMethod {
	switch method {
	case "credit_card":
		return pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case "debit_card":
		return pb.PaymentMethod_PAYMENT_METHOD_DEBIT_CARD
	case "pix":
		return pb.PaymentMethod_PAYMENT_METHOD_PIX
	case "boleto":
		return pb.PaymentMethod_PAYMENT_METHOD_BOLETO
	case "paypal":
		return pb.PaymentMethod_PAYMENT_METHOD_PAYPAL
	default:
		return pb.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
