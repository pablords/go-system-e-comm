package handler

import (
	"context"
	"log/slog"
	"payments-go/internal/domain/entity"
	"payments-go/internal/usecase"
	pb "payments-go/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PaymentServiceServer struct {
	pb.UnimplementedPaymentServiceServer
	processPaymentUC *usecase.ProcessPaymentUseCase
	getPaymentUC     *usecase.GetPaymentUseCase
	cancelPaymentUC  *usecase.CancelPaymentUseCase
	listPaymentsUC   *usecase.ListPaymentsUseCase
}

func NewPaymentServiceServer(
	processPaymentUC *usecase.ProcessPaymentUseCase,
	getPaymentUC *usecase.GetPaymentUseCase,
	cancelPaymentUC *usecase.CancelPaymentUseCase,
	listPaymentsUC *usecase.ListPaymentsUseCase,
) *PaymentServiceServer {
	return &PaymentServiceServer{
		processPaymentUC: processPaymentUC,
		getPaymentUC:     getPaymentUC,
		cancelPaymentUC:  cancelPaymentUC,
		listPaymentsUC:   listPaymentsUC,
	}
}

func (s *PaymentServiceServer) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	slog.Info("Received ProcessPayment request", "order_id", req.OrderId, "amount", req.Amount)

	// Convert proto payment method to entity payment method
	paymentMethod := convertProtoPaymentMethodToEntity(req.PaymentMethod)

	input := usecase.ProcessPaymentInput{
		OrderID:       req.OrderId,
		Amount:        req.Amount,
		PaymentMethod: paymentMethod,
		CustomerEmail: req.CustomerEmail,
		CustomerName:  req.CustomerName,
	}

	output, err := s.processPaymentUC.Execute(ctx, input)
	if err != nil {
		slog.Error("Failed to process payment", "error", err)
		return nil, err
	}

	return &pb.ProcessPaymentResponse{
		PaymentId:     output.PaymentID,
		OrderId:       output.OrderID,
		Status:        convertEntityStatusToProto(output.Status),
		Message:       output.Message,
		TransactionId: output.TransactionID,
		CreatedAt:     timestamppb.Now(),
	}, nil
}

func (s *PaymentServiceServer) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error) {
	slog.Info("Received GetPayment request", "payment_id", req.PaymentId)

	payment, err := s.getPaymentUC.Execute(ctx, req.PaymentId)
	if err != nil {
		slog.Error("Failed to get payment", "error", err)
		return nil, err
	}

	return &pb.GetPaymentResponse{
		PaymentId:     payment.ID,
		OrderId:       payment.OrderID,
		Amount:        payment.Amount,
		PaymentMethod: convertEntityPaymentMethodToProto(payment.PaymentMethod),
		Status:        convertEntityStatusToProto(payment.Status),
		TransactionId: payment.TransactionID,
		CreatedAt:     timestamppb.New(payment.CreatedAt),
		UpdatedAt:     timestamppb.New(payment.UpdatedAt),
	}, nil
}

func (s *PaymentServiceServer) CancelPayment(ctx context.Context, req *pb.CancelPaymentRequest) (*pb.CancelPaymentResponse, error) {
	slog.Info("Received CancelPayment request", "payment_id", req.PaymentId)

	input := usecase.CancelPaymentInput{
		PaymentID: req.PaymentId,
		Reason:    req.Reason,
	}

	err := s.cancelPaymentUC.Execute(ctx, input)
	if err != nil {
		slog.Error("Failed to cancel payment", "error", err)
		return &pb.CancelPaymentResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.CancelPaymentResponse{
		Success:    true,
		Message:    "Payment canceled successfully",
		CanceledAt: timestamppb.Now(),
	}, nil
}

func (s *PaymentServiceServer) ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
	slog.Info("Received ListPayments request", "order_id", req.OrderId)

	payments, err := s.listPaymentsUC.Execute(ctx, req.OrderId)
	if err != nil {
		slog.Error("Failed to list payments", "error", err)
		return nil, err
	}

	var pbPayments []*pb.GetPaymentResponse
	for _, payment := range payments {
		pbPayments = append(pbPayments, &pb.GetPaymentResponse{
			PaymentId:     payment.ID,
			OrderId:       payment.OrderID,
			Amount:        payment.Amount,
			PaymentMethod: convertEntityPaymentMethodToProto(payment.PaymentMethod),
			Status:        convertEntityStatusToProto(payment.Status),
			TransactionId: payment.TransactionID,
			CreatedAt:     timestamppb.New(payment.CreatedAt),
			UpdatedAt:     timestamppb.New(payment.UpdatedAt),
		})
	}

	return &pb.ListPaymentsResponse{
		Payments: pbPayments,
	}, nil
}

// Helper functions to convert between proto and entity types

func convertProtoPaymentMethodToEntity(method pb.PaymentMethod) entity.PaymentMethod {
	switch method {
	case pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return entity.PaymentMethodCreditCard
	case pb.PaymentMethod_PAYMENT_METHOD_DEBIT_CARD:
		return entity.PaymentMethodDebitCard
	case pb.PaymentMethod_PAYMENT_METHOD_PIX:
		return entity.PaymentMethodPix
	case pb.PaymentMethod_PAYMENT_METHOD_BOLETO:
		return entity.PaymentMethodBoleto
	case pb.PaymentMethod_PAYMENT_METHOD_PAYPAL:
		return entity.PaymentMethodPayPal
	default:
		return entity.PaymentMethodCreditCard
	}
}

func convertEntityPaymentMethodToProto(method entity.PaymentMethod) pb.PaymentMethod {
	switch method {
	case entity.PaymentMethodCreditCard:
		return pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case entity.PaymentMethodDebitCard:
		return pb.PaymentMethod_PAYMENT_METHOD_DEBIT_CARD
	case entity.PaymentMethodPix:
		return pb.PaymentMethod_PAYMENT_METHOD_PIX
	case entity.PaymentMethodBoleto:
		return pb.PaymentMethod_PAYMENT_METHOD_BOLETO
	case entity.PaymentMethodPayPal:
		return pb.PaymentMethod_PAYMENT_METHOD_PAYPAL
	default:
		return pb.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func convertEntityStatusToProto(status entity.PaymentStatus) pb.PaymentStatus {
	switch status {
	case entity.PaymentStatusPending:
		return pb.PaymentStatus_PAYMENT_STATUS_PENDING
	case entity.PaymentStatusProcessing:
		return pb.PaymentStatus_PAYMENT_STATUS_PROCESSING
	case entity.PaymentStatusApproved:
		return pb.PaymentStatus_PAYMENT_STATUS_APPROVED
	case entity.PaymentStatusDeclined:
		return pb.PaymentStatus_PAYMENT_STATUS_DECLINED
	case entity.PaymentStatusCanceled:
		return pb.PaymentStatus_PAYMENT_STATUS_CANCELED
	case entity.PaymentStatusRefunded:
		return pb.PaymentStatus_PAYMENT_STATUS_REFUNDED
	default:
		return pb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}
}
