# Integração Orders ↔ Payments via gRPC

## 🏗️ Arquitetura

```
┌─────────────────┐         gRPC          ┌──────────────────┐
│  Orders Service │ ──────────────────────>│ Payments Service │
│   (port 8080)   │  ProcessPayment()      │   (port 50051)   │
│                 │  CancelPayment()       │                  │
└─────────────────┘                        └──────────────────┘
```

## 📋 Fluxo de Criação de Pedido

1. Cliente faz requisição HTTP `POST /api/v1/orders/with-payment`
2. Orders cria pedido no banco de dados
3. **Orders chama Payments via gRPC** para processar pagamento
4. Payments processa e retorna status
5. Orders atualiza status do pedido baseado no pagamento
6. Orders retorna resposta ao cliente com `order_id` e `payment_id`

## 🚀 Como Executar

### Pré-requisitos
```bash
# Os serviços já estão compilados e configurados
```

### Iniciar Serviços

**Opção 1: Script automático (Recomendado)**
```bash
cd /Users/pablosantos/estudos/go-system-e-comm
./start-services.sh
```

**Opção 2: Manual**
```bash
# Terminal 1 - Payments
cd payments
./bin/payment-service

# Terminal 2 - Orders
cd orders
go run cmd/api/main.go
```

## 🧪 Testando a Integração

### 1. Criar Pedido com Pagamento

```bash
curl -X POST http://localhost:8080/api/v1/orders/with-payment \
  -H "Content-Type: application/json" \
  -d '{
    "customer_email": "cliente@example.com",
    "customer_name": "João Silva",
    "payment_method": 1,
    "items": [
      {
        "product_id": "prod-123",
        "quantity": 2,
        "price": 50.00
      },
      {
        "product_id": "prod-456",
        "quantity": 1,
        "price": 100.00
      }
    ]
  }'
```

**Resposta esperada:**
```json
{
  "order_id": "uuid-do-pedido",
  "total": 200.00,
  "status": "paid",
  "payment_id": "uuid-do-pagamento"
}
```

### 2. Cancelar Pedido e Pagamento

```bash
curl -X POST http://localhost:8080/api/v1/orders/{order_id}/cancel \
  -H "Content-Type: application/json" \
  -d '{
    "payment_id": "uuid-do-pagamento"
  }'
```

**Resposta esperada:**
```json
{
  "message": "Order canceled successfully"
}
```

## 🔧 Configuração

### Orders .env
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=orders_user
DB_PASSWORD=orders_pass
DB_NAME=orders_db
SERVER_PORT=8080
PAYMENT_SERVICE_ADDR=localhost:50051
```

### Payments .env
```env
DB_DSN=root:root@tcp(localhost:3306)/payments_db?parseTime=true
GRPC_PORT=50051
```

## 📝 Métodos de Pagamento

```
1 = CREDIT_CARD  (Cartão de Crédito)
2 = DEBIT_CARD   (Cartão de Débito)
3 = PIX          (Pix)
4 = BOLETO       (Boleto Bancário)
5 = PAYPAL       (PayPal)
```

## 🔍 Status de Pagamento

```
PENDING    = Aguardando processamento
PROCESSING = Em processamento
APPROVED   = Aprovado ✅
DECLINED   = Recusado ❌
CANCELED   = Cancelado
REFUNDED   = Reembolsado
```

## 🔍 Status de Pedido

Baseado no status do pagamento:
- `APPROVED` → Order status: `paid`
- `PROCESSING` → Order status: `pending`
- `DECLINED` → Order status: `canceled`

## 📊 Endpoints Disponíveis

### Orders Service (HTTP - Port 8080)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/api/v1/orders/with-payment` | Criar pedido com pagamento |
| POST | `/api/v1/orders/{id}/cancel` | Cancelar pedido e pagamento |
| GET | `/api/v1/orders` | Listar pedidos |
| GET | `/api/v1/orders/{id}` | Buscar pedido |
| GET | `/health` | Health check |
| GET | `/swagger/*` | Documentação Swagger |

### Payments Service (gRPC - Port 50051)

| Método | Descrição |
|--------|-----------|
| `ProcessPayment` | Processar pagamento |
| `GetPayment` | Buscar pagamento |
| `CancelPayment` | Cancelar pagamento |
| `ListPayments` | Listar pagamentos de um pedido |

## 🛡️ Tratamento de Erros

O cliente gRPC implementa:
- ✅ Timeout de 10 segundos para ProcessPayment
- ✅ Timeout de 5 segundos para outras operações
- ✅ Logging estruturado em JSON
- ✅ Graceful degradation (se pagamento falhar, pedido é marcado como `payment_failed`)

## 📊 Logs Estruturados

### Orders Service
```json
{
  "level": "info",
  "msg": "Payment processed successfully",
  "payment_id": "pay_123",
  "order_id": "ord_456",
  "status": "APPROVED"
}
```

### Payments Service
```json
{
  "level": "info",
  "msg": "Payment created successfully",
  "payment_id": "pay_123",
  "order_id": "ord_456",
  "amount": 200.00,
  "method": "CREDIT_CARD"
}
```

## 🔐 Próximos Passos (Melhorias)

- [ ] Adicionar autenticação/autorização
- [ ] Implementar TLS/SSL para gRPC
- [ ] Adicionar retry automático com backoff exponencial
- [ ] Implementar circuit breaker
- [ ] Adicionar rate limiting
- [ ] Implementar idempotência
- [ ] Adicionar tracing distribuído (OpenTelemetry)
- [ ] Implementar saga pattern para compensação de transações

## 📚 Arquivos Criados

### Orders Service
- `internal/infra/grpc/client/payment_client.go` - Cliente gRPC
- `internal/usecase/create_order_with_payment_usecase.go` - Use case de criação
- `internal/usecase/cancel_order_usecase.go` - Use case de cancelamento
- `internal/infra/http/handler/order_with_payment_handler.go` - HTTP handlers
- `proto/payment.proto` - Definições protobuf
- `proto/payment.pb.go` - Código gerado
- `proto/payment_grpc.pb.go` - Código gRPC gerado

### Root
- `start-services.sh` - Script para iniciar ambos serviços

## 🐛 Troubleshooting

### Erro: "Failed to connect to payment service"
- Verifique se o Payment Service está rodando
- Confirme que a porta 50051 está livre
- Verifique `PAYMENT_SERVICE_ADDR` no `.env`

### Erro: "Payment processing failed"
- Verifique logs do Payment Service
- Confirme que o banco de dados do Payments está configurado
- Teste o Payment Service diretamente via gRPC

### Erro de compilação: "undefined: pb"
- Execute: `cd orders && make proto` (ou o comando protoc manualmente)
- Verifique se os arquivos `.pb.go` foram gerados em `proto/`

## 📖 Documentação Adicional

- [gRPC Documentation](https://grpc.io/docs/)
- [Protocol Buffers](https://protobuf.dev/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
