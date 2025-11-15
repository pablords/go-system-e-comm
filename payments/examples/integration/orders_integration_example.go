package integration

// Exemplo de uso do Payment Client no serviço Orders
// Este arquivo demonstra como integrar o cliente de pagamento nos handlers do Orders

// Ajuste os imports conforme sua estrutura de projeto
// "orders-go/internal/infra/grpc/client"
// pb "orders-go/proto"

/*
// 1. No main.go do Orders, inicialize o cliente:

func main() {
    // ... código existente ...

    // Conectar ao Payment Service
    paymentServiceAddr := os.Getenv("PAYMENT_SERVICE_ADDRESS")
    if paymentServiceAddr == "" {
        paymentServiceAddr = "localhost:50051"
    }

    paymentClient, err := client.NewPaymentClient(paymentServiceAddr)
    if err != nil {
        slog.Error("Failed to connect to payment service", "error", err)
        os.Exit(1)
    }
    defer paymentClient.Close()

    slog.Info("Connected to payment service", "address", paymentServiceAddr)

    // Passar o paymentClient para os handlers
    orderHandler := handler.NewOrderHandler(orderUseCase, paymentClient)

    // ... resto do código ...
}
*/

/*
// 2. Atualizar o OrderHandler para incluir o PaymentClient:

type OrderHandler struct {
    orderUseCase  *usecase.OrderUseCase
    paymentClient *client.PaymentClient
}

func NewOrderHandler(orderUseCase *usecase.OrderUseCase, paymentClient *client.PaymentClient) *OrderHandler {
    return &OrderHandler{
        orderUseCase:  orderUseCase,
        paymentClient: paymentClient,
    }
}
*/

/*
// 3. Criar endpoint para processar pagamento do pedido:

// ProcessOrderPayment processa o pagamento de um pedido
// @Summary Process order payment
// @Description Process payment for an existing order
// @Tags payments
// @Accept json
// @Produce json
// @Param request body ProcessPaymentRequest true "Payment details"
// @Success 200 {object} ProcessPaymentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{orderId}/payment [post]
func (h *OrderHandler) ProcessOrderPayment(w http.ResponseWriter, r *http.Request) {
    var req ProcessPaymentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validar campos obrigatórios
    if req.OrderID == "" || req.Amount <= 0 || req.Email == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    // Buscar o pedido para verificar se existe
    order, err := h.orderUseCase.GetOrder(req.OrderID)
    if err != nil {
        http.Error(w, "Order not found", http.StatusNotFound)
        return
    }

    // Verificar se o valor está correto
    if order.Total != req.Amount {
        http.Error(w, "Payment amount does not match order total", http.StatusBadRequest)
        return
    }

    // Converter método de pagamento
    paymentMethod := client.ConvertPaymentMethodString(req.PaymentMethod)
    if paymentMethod == pb.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
        http.Error(w, "Invalid payment method", http.StatusBadRequest)
        return
    }

    // Criar contexto com timeout
    ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
    defer cancel()

    // Processar pagamento via gRPC
    payment, err := h.paymentClient.ProcessPayment(
        ctx,
        req.OrderID,
        req.Amount,
        paymentMethod,
        req.Email,
        req.Name,
    )

    if err != nil {
        slog.Error("Failed to process payment", "error", err)
        http.Error(w, "Failed to process payment", http.StatusInternalServerError)
        return
    }

    // Atualizar status do pedido baseado no resultado do pagamento
    if client.IsPaymentApproved(payment.Status) {
        if err := h.orderUseCase.UpdateOrderStatus(req.OrderID, "paid"); err != nil {
            slog.Error("Failed to update order status", "error", err)
        }
    } else if client.IsPaymentDeclined(payment.Status) {
        slog.Warn("Payment declined", "order_id", req.OrderID, "payment_id", payment.PaymentId)
    }

    // Retornar resposta
    response := ProcessPaymentResponse{
        PaymentID:     payment.PaymentId,
        OrderID:       payment.OrderId,
        Status:        payment.Status.String(),
        Message:       payment.Message,
        TransactionID: payment.TransactionId,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
*/

/*
// 4. Criar endpoint para consultar status do pagamento:

// GetOrderPayment busca os pagamentos de um pedido
// @Summary Get order payments
// @Description Get all payments for a specific order
// @Tags payments
// @Produce json
// @Param orderId path string true "Order ID"
// @Success 200 {object} ListPaymentsResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{orderId}/payments [get]
func (h *OrderHandler) GetOrderPayments(w http.ResponseWriter, r *http.Request) {
    orderID := chi.URLParam(r, "orderId")
    if orderID == "" {
        http.Error(w, "Order ID is required", http.StatusBadRequest)
        return
    }

    // Criar contexto com timeout
    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
    defer cancel()

    // Buscar pagamentos via gRPC
    payments, err := h.paymentClient.ListPayments(ctx, orderID)
    if err != nil {
        slog.Error("Failed to get order payments", "error", err)
        http.Error(w, "Failed to get payments", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(payments)
}
*/

/*
// 5. Criar endpoint para cancelar pagamento:

// CancelOrderPayment cancela um pagamento
// @Summary Cancel payment
// @Description Cancel a pending payment
// @Tags payments
// @Accept json
// @Produce json
// @Param paymentId path string true "Payment ID"
// @Param request body CancelPaymentRequest true "Cancel reason"
// @Success 200 {object} CancelPaymentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /payments/{paymentId}/cancel [post]
func (h *OrderHandler) CancelOrderPayment(w http.ResponseWriter, r *http.Request) {
    paymentID := chi.URLParam(r, "paymentId")
    if paymentID == "" {
        http.Error(w, "Payment ID is required", http.StatusBadRequest)
        return
    }

    var req CancelPaymentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Criar contexto com timeout
    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
    defer cancel()

    // Cancelar pagamento via gRPC
    response, err := h.paymentClient.CancelPayment(ctx, paymentID, req.Reason)
    if err != nil {
        slog.Error("Failed to cancel payment", "error", err)
        http.Error(w, "Failed to cancel payment", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
*/

/*
// 6. Adicionar rotas no router:

func setupRoutes(r *chi.Mux, orderHandler *handler.OrderHandler) {
    r.Route("/api/v1", func(r chi.Router) {
        r.Route("/orders", func(r chi.Router) {
            // Rotas existentes...

            // Novas rotas de pagamento
            r.Post("/{orderId}/payment", orderHandler.ProcessOrderPayment)
            r.Get("/{orderId}/payments", orderHandler.GetOrderPayments)
        })

        r.Route("/payments", func(r chi.Router) {
            r.Post("/{paymentId}/cancel", orderHandler.CancelOrderPayment)
        })
    })
}
*/

// Request/Response structures

type ProcessPaymentRequest struct {
	OrderID       string  `json:"order_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"` // "credit_card", "pix", etc.
	Email         string  `json:"email"`
	Name          string  `json:"name"`
}

type ProcessPaymentResponse struct {
	PaymentID     string `json:"payment_id"`
	OrderID       string `json:"order_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
	TransactionID string `json:"transaction_id"`
}

type CancelPaymentRequest struct {
	Reason string `json:"reason"`
}

type CancelPaymentResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	CanceledAt string `json:"canceled_at,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
