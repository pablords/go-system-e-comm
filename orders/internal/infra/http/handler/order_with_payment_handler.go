package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"orders-go/internal/usecase"
)

type OrderWithPaymentHandler struct {
	createOrderUseCase *usecase.CreateOrderUseCase
	cancelOrderUseCase *usecase.CancelOrderUseCase
	logger             *slog.Logger
}

func NewOrderWithPaymentHandler(
	createOrderUseCase *usecase.CreateOrderUseCase,
	cancelOrderUseCase *usecase.CancelOrderUseCase,
	logger *slog.Logger,
) *OrderWithPaymentHandler {
	return &OrderWithPaymentHandler{
		createOrderUseCase: createOrderUseCase,
		cancelOrderUseCase: cancelOrderUseCase,
		logger:             logger,
	}
}

type CreateOrderWithPaymentRequest struct {
	CustomerEmail string             `json:"customer_email"`
	CustomerName  string             `json:"customer_name"`
	Items         []OrderItemRequest `json:"items"`
	PaymentMethod int32              `json:"payment_method"` // 1=CREDIT_CARD, 2=DEBIT_CARD, 3=PIX, 4=BOLETO, 5=PAYPAL
}

type OrderItemRequest struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CreateOrderWithPaymentResponse struct {
	OrderID   string  `json:"order_id"`
	Total     float64 `json:"total"`
	Status    string  `json:"status"`
	PaymentID string  `json:"payment_id"`
}

// CreateOrderWithPayment godoc
// @Summary Create order with payment processing
// @Description Creates a new order and processes payment via gRPC
// @Tags orders
// @Accept json
// @Produce json
// @Param request body CreateOrderWithPaymentRequest true "Order and Payment Info"
// @Success 201 {object} CreateOrderWithPaymentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/with-payment [post]
func (h *OrderWithPaymentHandler) CreateOrderWithPayment(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderWithPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validar request
	if req.CustomerEmail == "" {
		respondWithError(w, http.StatusBadRequest, "Customer email is required")
		return
	}
	if req.CustomerName == "" {
		respondWithError(w, http.StatusBadRequest, "Customer name is required")
		return
	}
	if len(req.Items) == 0 {
		respondWithError(w, http.StatusBadRequest, "At least one item is required")
		return
	}
	if req.PaymentMethod < 1 || req.PaymentMethod > 5 {
		respondWithError(w, http.StatusBadRequest, "Invalid payment method (1-5)")
		return
	}

	// Converter items
	items := make([]usecase.OrderItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = usecase.OrderItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	// Executar use case
	input := usecase.CreateOrderInput{
		CustomerEmail: req.CustomerEmail,
		CustomerName:  req.CustomerName,
		Items:         items,
		PaymentMethod: req.PaymentMethod,
	}

	output, err := h.createOrderUseCase.Execute(r.Context(), input)
	if err != nil {
		h.logger.Error("Failed to create order with payment", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create order: "+err.Error())
		return
	}

	response := CreateOrderWithPaymentResponse{
		OrderID:   output.OrderID,
		Total:     output.Total,
		Status:    output.Status,
		PaymentID: output.PaymentID,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

type CancelOrderRequest struct {
	PaymentID string `json:"payment_id"`
}

// CancelOrder godoc
// @Summary Cancel order and payment
// @Description Cancels an order and its associated payment via gRPC
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body CancelOrderRequest true "Payment ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id}/cancel [post]
func (h *OrderWithPaymentHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("id")
	if orderID == "" {
		respondWithError(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	var req CancelOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.cancelOrderUseCase.Execute(r.Context(), orderID, req.PaymentID); err != nil {
		h.logger.Error("Failed to cancel order", "error", err, "order_id", orderID)
		respondWithError(w, http.StatusInternalServerError, "Failed to cancel order: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, MessageResponse{Message: "Order canceled successfully"})
}

type MessageResponse struct {
	Message string `json:"message"`
}
