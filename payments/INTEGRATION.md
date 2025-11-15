# Integração do Payment Service com Orders

Este documento descreve como integrar o serviço de pagamentos via gRPC no serviço de orders.

## 1. Adicionar dependências gRPC no Orders

Adicione as seguintes dependências no `go.mod` do serviço orders:

```go
require (
    google.golang.org/grpc v1.60.1
    google.golang.org/protobuf v1.32.0
)
```

## 2. Copiar arquivo proto

Copie o arquivo `payments/proto/payment.proto` para o projeto orders ou compartilhe via módulo Go.

## 3. Gerar código gRPC no Orders

```bash
cd orders
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/payment.proto
```

## 4. Criar cliente gRPC no Orders

Crie um arquivo `internal/infra/grpc/client/payment_client.go`:

```go
package client

import (
    "context"
    "fmt"
    "log/slog"
    "time"

    pb "orders-go/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type PaymentClient struct {
    client pb.PaymentServiceClient
    conn   *grpc.ClientConn
}

func NewPaymentClient(address string) (*PaymentClient, error) {
    conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to payment service: %w", err)
    }

    client := pb.NewPaymentServiceClient(conn)

    return &PaymentClient{
        client: client,
        conn:   conn,
    }, nil
}

func (c *PaymentClient) Close() error {
    if c.conn != nil {
        return c.conn.Close()
    }
    return nil
}

func (c *PaymentClient) ProcessPayment(orderID string, amount float64, method pb.PaymentMethod, email, name string) (*pb.ProcessPaymentResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    req := &pb.ProcessPaymentRequest{
        OrderId:       orderID,
        Amount:        amount,
        PaymentMethod: method,
        CustomerEmail: email,
        CustomerName:  name,
    }

    slog.Info("Calling payment service", "order_id", orderID, "amount", amount)

    resp, err := c.client.ProcessPayment(ctx, req)
    if err != nil {
        slog.Error("Failed to process payment", "error", err)
        return nil, err
    }

    slog.Info("Payment processed", "payment_id", resp.PaymentId, "status", resp.Status)

    return resp, nil
}

func (c *PaymentClient) GetPayment(paymentID string) (*pb.GetPaymentResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    req := &pb.GetPaymentRequest{
        PaymentId: paymentID,
    }

    return c.client.GetPayment(ctx, req)
}

func (c *PaymentClient) CancelPayment(paymentID, reason string) (*pb.CancelPaymentResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    req := &pb.CancelPaymentRequest{
        PaymentId: paymentID,
        Reason:    reason,
    }

    return c.client.CancelPayment(ctx, req)
}
```

## 5. Usar no Order Handler

Exemplo de uso no `order_handler.go`:

```go
// No main.go, inicialize o cliente:
paymentClient, err := client.NewPaymentClient("localhost:50051")
if err != nil {
    log.Fatal(err)
}
defer paymentClient.Close()

// No handler, processe o pagamento:
func (h *OrderHandler) ProcessOrderPayment(w http.ResponseWriter, r *http.Request) {
    var req struct {
        OrderID       string  `json:"order_id"`
        Amount        float64 `json:"amount"`
        PaymentMethod string  `json:"payment_method"`
        Email         string  `json:"email"`
        Name          string  `json:"name"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Converter método de pagamento
    var method pb.PaymentMethod
    switch req.PaymentMethod {
    case "credit_card":
        method = pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
    case "pix":
        method = pb.PaymentMethod_PAYMENT_METHOD_PIX
    default:
        method = pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
    }

    // Processar pagamento
    payment, err := h.paymentClient.ProcessPayment(
        req.OrderID,
        req.Amount,
        method,
        req.Email,
        req.Name,
    )

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Atualizar status do pedido baseado no pagamento
    if payment.Status == pb.PaymentStatus_PAYMENT_STATUS_APPROVED {
        // Atualizar order para "paid"
    }

    json.NewEncoder(w).Encode(payment)
}
```

## 6. Variáveis de Ambiente

Adicione no `.env` do orders:

```
PAYMENT_SERVICE_ADDRESS=localhost:50051
```

## 7. Executar os serviços

Terminal 1 - Payment Service:
```bash
cd payments
docker-compose up -d payments-db
make proto
go mod download
make run
```

Terminal 2 - Orders Service:
```bash
cd orders
make run
```

## 8. Testar integração

Use grpcurl para testar o payment service:

```bash
# Listar serviços
grpcurl -plaintext localhost:50051 list

# Processar pagamento
grpcurl -plaintext -d '{
  "order_id": "123",
  "amount": 100.50,
  "payment_method": 1,
  "customer_email": "test@example.com",
  "customer_name": "Test User"
}' localhost:50051 payment.PaymentService/ProcessPayment
```

## Fluxo de Integração

1. **Order Service** recebe requisição para criar pedido
2. **Order Service** calcula o total do pedido
3. **Order Service** chama **Payment Service** via gRPC para processar pagamento
4. **Payment Service** processa e retorna status
5. **Order Service** atualiza status do pedido baseado no resultado do pagamento
6. **Order Service** retorna resposta ao cliente
