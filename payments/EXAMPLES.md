# Exemplo de Requisições - Payment Service

Este documento contém exemplos de requisições para testar o Payment Service.

## Usando grpcurl

### 1. Processar Pagamento com Cartão de Crédito

```bash
grpcurl -plaintext -d '{
  "order_id": "order-001",
  "amount": 299.99,
  "payment_method": 1,
  "customer_email": "joao.silva@example.com",
  "customer_name": "João Silva"
}' localhost:50051 payment.PaymentService/ProcessPayment
```

### 2. Processar Pagamento com PIX

```bash
grpcurl -plaintext -d '{
  "order_id": "order-002",
  "amount": 150.00,
  "payment_method": 3,
  "customer_email": "maria.santos@example.com",
  "customer_name": "Maria Santos",
  "pix_details": {
    "pix_key": "maria.santos@example.com"
  }
}' localhost:50051 payment.PaymentService/ProcessPayment
```

### 3. Processar Pagamento com Boleto

```bash
grpcurl -plaintext -d '{
  "order_id": "order-003",
  "amount": 500.00,
  "payment_method": 4,
  "customer_email": "pedro.oliveira@example.com",
  "customer_name": "Pedro Oliveira",
  "boleto_details": {
    "customer_document": "12345678901",
    "due_date": "2025-12-01T00:00:00Z"
  }
}' localhost:50051 payment.PaymentService/ProcessPayment
```

### 4. Buscar Pagamento

```bash
# Substitua PAYMENT_ID pelo ID retornado na criação
grpcurl -plaintext -d '{
  "payment_id": "PAYMENT_ID"
}' localhost:50051 payment.PaymentService/GetPayment
```

### 5. Listar Pagamentos de um Pedido

```bash
grpcurl -plaintext -d '{
  "order_id": "order-001"
}' localhost:50051 payment.PaymentService/ListPayments
```

### 6. Cancelar Pagamento

```bash
# Substitua PAYMENT_ID pelo ID do pagamento a cancelar
grpcurl -plaintext -d '{
  "payment_id": "PAYMENT_ID",
  "reason": "Cliente solicitou cancelamento"
}' localhost:50051 payment.PaymentService/CancelPayment
```

## Métodos de Pagamento (payment_method)

- `1` - CREDIT_CARD (Cartão de Crédito)
- `2` - DEBIT_CARD (Cartão de Débito)
- `3` - PIX
- `4` - BOLETO
- `5` - PAYPAL

## Status de Pagamento

- `1` - PENDING (Pendente)
- `2` - PROCESSING (Processando)
- `3` - APPROVED (Aprovado)
- `4` - DECLINED (Recusado)
- `5` - CANCELED (Cancelado)
- `6` - REFUNDED (Reembolsado)

## Testando Integração com Orders

### 1. Iniciar Payment Service

```bash
cd payments
docker-compose up -d
# ou
make run
```

### 2. No Orders Service, fazer requisição HTTP

```bash
# Criar um pedido no Orders (assumindo que você tenha este endpoint)
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "product_id": "prod-001",
        "quantity": 2,
        "price": 50.00
      }
    ]
  }'

# Processar pagamento para o pedido criado
curl -X POST http://localhost:8080/api/v1/orders/ORDER_ID/payment \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "ORDER_ID",
    "amount": 100.00,
    "payment_method": "credit_card",
    "email": "cliente@example.com",
    "name": "Cliente Teste"
  }'
```

## Cenários de Teste

### Cenário 1: Pagamento Aprovado

```bash
# Pagamento com valor menor que 10000 será aprovado
grpcurl -plaintext -d '{
  "order_id": "order-100",
  "amount": 99.90,
  "payment_method": 1,
  "customer_email": "test@example.com",
  "customer_name": "Test User"
}' localhost:50051 payment.PaymentService/ProcessPayment
```

### Cenário 2: Pagamento Recusado

```bash
# Pagamento com valor maior ou igual a 10000 será recusado
grpcurl -plaintext -d '{
  "order_id": "order-101",
  "amount": 15000.00,
  "payment_method": 1,
  "customer_email": "test@example.com",
  "customer_name": "Test User"
}' localhost:50051 payment.PaymentService/ProcessPayment
```

### Cenário 3: Cancelamento de Pagamento

```bash
# 1. Criar pagamento
PAYMENT_RESPONSE=$(grpcurl -plaintext -d '{
  "order_id": "order-102",
  "amount": 200.00,
  "payment_method": 3,
  "customer_email": "test@example.com",
  "customer_name": "Test User"
}' localhost:50051 payment.PaymentService/ProcessPayment)

# 2. Extrair payment_id (manualmente ou com jq)
# PAYMENT_ID=...

# 3. Cancelar
grpcurl -plaintext -d '{
  "payment_id": "PAYMENT_ID",
  "reason": "Pedido cancelado pelo cliente"
}' localhost:50051 payment.PaymentService/CancelPayment
```

## Verificar Dados no Banco

```bash
# Conectar ao MySQL
docker exec -it payments-mysql mysql -u root -proot payments_db

# Listar todos os pagamentos
SELECT id, order_id, amount, payment_method, status, created_at FROM payments;

# Buscar pagamentos por status
SELECT * FROM payments WHERE status = 'approved';

# Buscar pagamentos de um pedido específico
SELECT * FROM payments WHERE order_id = 'order-001';
```

## Logs

### Ver logs do Payment Service

```bash
# Docker
docker-compose logs -f payments-service

# Local
# Os logs aparecem no terminal onde o serviço está rodando
```

## Troubleshooting

### Erro: "connection refused"

- Verifique se o serviço está rodando: `docker-compose ps`
- Verifique a porta: `netstat -an | grep 50051`

### Erro: "payment not found"

- Verifique se o payment_id está correto
- Liste os pagamentos no banco de dados

### Erro: "payment cannot be canceled"

- Pagamentos aprovados não podem ser cancelados
- Use apenas pagamentos com status pending ou processing
