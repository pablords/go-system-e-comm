# Exemplos de Integração

Esta pasta contém exemplos práticos de como integrar o Payment Service com outros serviços.

## Estrutura

```
examples/
├── client/
│   └── payment_client.go              # Cliente gRPC completo
└── integration/
    └── orders_integration_example.go  # Exemplo de integração com Orders
```

## Cliente gRPC (`client/payment_client.go`)

Implementação completa de um cliente gRPC para o Payment Service que pode ser usado em outros serviços Go.

### Características

- ✅ Gerenciamento de conexão
- ✅ Tratamento de erros
- ✅ Context timeout
- ✅ Logging estruturado
- ✅ Helper functions
- ✅ Type conversions

### Como usar

1. Copie o arquivo para seu projeto (ex: `orders/internal/infra/grpc/client/`)
2. Ajuste os imports para seu projeto
3. Inicialize no main.go:

```go
paymentClient, err := client.NewPaymentClient("localhost:50051")
if err != nil {
    log.Fatal(err)
}
defer paymentClient.Close()
```

4. Use nos handlers:

```go
payment, err := paymentClient.ProcessPayment(
    ctx,
    orderID,
    amount,
    pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
    "user@example.com",
    "User Name",
)
```

## Integração com Orders (`integration/orders_integration_example.go`)

Exemplo completo de como integrar o Payment Service no serviço Orders.

### Inclui

- ✅ Estrutura do handler
- ✅ Endpoints REST para pagamentos
- ✅ Request/Response structures
- ✅ Validações
- ✅ Tratamento de erros
- ✅ Atualização de status do pedido

### Endpoints sugeridos

- `POST /api/v1/orders/{orderId}/payment` - Processar pagamento
- `GET /api/v1/orders/{orderId}/payments` - Listar pagamentos do pedido
- `POST /api/v1/payments/{paymentId}/cancel` - Cancelar pagamento

## Testando os Exemplos

### 1. Iniciar Payment Service

```bash
cd payments
docker-compose up -d
# ou
make run
```

### 2. Testar com grpcurl

```bash
# Processar pagamento
grpcurl -plaintext -d '{
  "order_id": "order-123",
  "amount": 150.00,
  "payment_method": 1,
  "customer_email": "test@example.com",
  "customer_name": "Test User"
}' localhost:50051 payment.PaymentService/ProcessPayment
```

### 3. Integrar no Orders

Siga os passos em `INTEGRATION.md` para integração completa.

## Métodos Disponíveis no Cliente

### ProcessPayment

Processa um novo pagamento.

```go
payment, err := client.ProcessPayment(
    ctx,
    "order-123",
    100.50,
    pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
    "user@example.com",
    "User Name",
)
```

### GetPayment

Busca detalhes de um pagamento.

```go
payment, err := client.GetPayment(ctx, "payment-id")
```

### CancelPayment

Cancela um pagamento pendente.

```go
response, err := client.CancelPayment(ctx, "payment-id", "Customer requested")
```

### ListPayments

Lista todos os pagamentos de um pedido.

```go
payments, err := client.ListPayments(ctx, "order-id")
```

## Helper Functions

### IsPaymentApproved

```go
if client.IsPaymentApproved(payment.Status) {
    // Atualizar pedido para "paid"
}
```

### IsPaymentDeclined

```go
if client.IsPaymentDeclined(payment.Status) {
    // Notificar cliente
}
```

### ConvertPaymentMethodString

```go
method := client.ConvertPaymentMethodString("credit_card")
// Retorna: pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
```

## Fluxo Completo de Integração

```
1. Cliente HTTP → Orders API
   POST /orders
   
2. Orders cria pedido
   
3. Orders → Payment Service (gRPC)
   ProcessPayment(orderID, amount, method)
   
4. Payment Service processa
   
5. Payment Service → Orders (gRPC Response)
   PaymentResponse(status, transaction_id)
   
6. Orders atualiza status do pedido
   
7. Orders → Cliente HTTP
   Response com pedido e pagamento
```

## Exemplo de Request/Response

### HTTP Request para Orders

```json
POST /api/v1/orders/order-123/payment
Content-Type: application/json

{
  "order_id": "order-123",
  "amount": 150.00,
  "payment_method": "credit_card",
  "email": "user@example.com",
  "name": "User Name"
}
```

### HTTP Response do Orders

```json
{
  "payment_id": "pay-456",
  "order_id": "order-123",
  "status": "PAYMENT_STATUS_APPROVED",
  "message": "Payment processed successfully",
  "transaction_id": "txn-789"
}
```

## Tratamento de Erros

### Connection Errors

```go
client, err := NewPaymentClient(address)
if err != nil {
    // Não foi possível conectar ao Payment Service
    // Retry ou retornar erro para cliente
}
```

### Processing Errors

```go
payment, err := client.ProcessPayment(...)
if err != nil {
    // Erro ao processar pagamento
    // Log e retornar erro apropriado
}
```

### Declined Payments

```go
if payment.Status == pb.PaymentStatus_PAYMENT_STATUS_DECLINED {
    // Pagamento recusado
    // Notificar cliente e não atualizar pedido
}
```

## Timeout Configuration

Todos os métodos aceitam `context.Context` para controle de timeout:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

payment, err := client.ProcessPayment(ctx, ...)
```

## Boas Práticas

1. **Sempre fechar a conexão**:
   ```go
   defer client.Close()
   ```

2. **Usar timeouts apropriados**:
   - ProcessPayment: 30 segundos
   - GetPayment: 5 segundos
   - CancelPayment: 5 segundos
   - ListPayments: 10 segundos

3. **Validar antes de chamar**:
   ```go
   if amount <= 0 {
       return errors.New("invalid amount")
   }
   ```

4. **Logar todas as operações**:
   ```go
   slog.Info("Processing payment", "order_id", orderID)
   ```

5. **Tratar todos os status**:
   ```go
   switch payment.Status {
   case pb.PaymentStatus_PAYMENT_STATUS_APPROVED:
       // Atualizar pedido
   case pb.PaymentStatus_PAYMENT_STATUS_DECLINED:
       // Notificar cliente
   case pb.PaymentStatus_PAYMENT_STATUS_PENDING:
       // Aguardar processamento
   }
   ```

## Troubleshooting

### Connection refused

- Verifique se o Payment Service está rodando
- Verifique o endereço e porta
- Verifique firewall

### Context deadline exceeded

- Aumente o timeout
- Verifique performance do Payment Service
- Verifique conectividade de rede

### Invalid payment method

- Use os métodos de conversão fornecidos
- Valide input do usuário antes de enviar

## Próximos Passos

1. Copie os arquivos de exemplo para seu projeto
2. Ajuste imports e packages
3. Adicione ao seu main.go
4. Implemente os endpoints no seu handler
5. Teste a integração
6. Adicione testes unitários e de integração

## Documentação Adicional

- Ver `../INTEGRATION.md` para guia completo de integração
- Ver `../EXAMPLES.md` para mais exemplos de requisições
- Ver `../proto/payment.proto` para definições completas do protocolo
