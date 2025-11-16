package client

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pb "orders/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentClient struct {
	client pb.PaymentServiceClient
	conn   *grpc.ClientConn
	logger *slog.Logger
}

func NewPaymentClient(paymentServiceAddr string, logger *slog.Logger) (*PaymentClient, error) {
	conn, err := grpc.NewClient(
		paymentServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}

	client := pb.NewPaymentServiceClient(conn)

	return &PaymentClient{
		client: client,
		conn:   conn,
		logger: logger,
	}, nil
}

func (c *PaymentClient) Close() error {
	return c.conn.Close()
}

// ProcessPayment processa um pagamento via gRPC
func (c *PaymentClient) ProcessPayment(ctx context.Context, orderID string, amount float64, paymentMethod int32, customerEmail, customerName string) (*pb.ProcessPaymentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	c.logger.Info("Processing payment via gRPC",
		"order_id", orderID,
		"amount", amount,
		"payment_method", paymentMethod,
	)

	request := &pb.ProcessPaymentRequest{
		OrderId:       orderID,
		Amount:        amount,
		PaymentMethod: pb.PaymentMethod(paymentMethod),
		CustomerEmail: customerEmail,
		CustomerName:  customerName,
	}

	c.logger.Info("About to call ProcessPayment gRPC", "order_id", orderID)
	response, err := c.client.ProcessPayment(ctx, request)
	c.logger.Info("ProcessPayment gRPC returned", "order_id", orderID, "error", err)

	if err != nil {
		c.logger.Error("Failed to process payment",
			"error", err,
			"order_id", orderID,
		)
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	c.logger.Info("Payment processed successfully",
		"payment_id", response.PaymentId,
		"status", response.Status,
		"order_id", orderID,
	)

	return response, nil
}

// GetPayment busca detalhes de um pagamento
func (c *PaymentClient) GetPayment(ctx context.Context, paymentID string) (*pb.GetPaymentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	request := &pb.GetPaymentRequest{
		PaymentId: paymentID,
	}

	response, err := c.client.GetPayment(ctx, request)
	if err != nil {
		c.logger.Error("Failed to get payment",
			"error", err,
			"payment_id", paymentID,
		)
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return response, nil
}

// CancelPayment cancela um pagamento
func (c *PaymentClient) CancelPayment(ctx context.Context, paymentID string) (*pb.CancelPaymentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c.logger.Info("Canceling payment via gRPC", "payment_id", paymentID)

	request := &pb.CancelPaymentRequest{
		PaymentId: paymentID,
	}

	response, err := c.client.CancelPayment(ctx, request)
	if err != nil {
		c.logger.Error("Failed to cancel payment",
			"error", err,
			"payment_id", paymentID,
		)
		return nil, fmt.Errorf("failed to cancel payment: %w", err)
	}

	c.logger.Info("Payment canceled successfully", "payment_id", paymentID)

	return response, nil
}

// ListPayments lista pagamentos de um pedido
func (c *PaymentClient) ListPayments(ctx context.Context, orderID string) (*pb.ListPaymentsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	request := &pb.ListPaymentsRequest{
		OrderId: orderID,
	}

	response, err := c.client.ListPayments(ctx, request)
	if err != nil {
		c.logger.Error("Failed to list payments",
			"error", err,
			"order_id", orderID,
		)
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}

	return response, nil
}
