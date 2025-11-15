# Setup Inicial - Payment Service

Este guia ajudará você a configurar e executar o serviço de pagamentos.

## 1. Pré-requisitos

Certifique-se de ter instalado:

- Go 1.21 ou superior
- MySQL 8.0 ou superior (ou use Docker)
- Protocol Buffers compiler (protoc)
- grpcurl (opcional, para testes)

### Instalar protoc e plugins Go (macOS)

```bash
# Instalar protoc
brew install protobuf

# Instalar plugins Go para gRPC
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Adicionar ao PATH (se necessário)
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Instalar grpcurl (opcional)

```bash
brew install grpcurl
```

## 2. Configuração do Projeto

### 2.1. Download das dependências

```bash
cd payments
go mod download
```

### 2.2. Gerar código a partir dos arquivos .proto

```bash
make proto
```

Isso irá gerar os arquivos:
- `proto/payment.pb.go`
- `proto/payment_grpc.pb.go`

## 3. Configurar Banco de Dados

### Opção A: Usando Docker (Recomendado)

```bash
# Iniciar apenas o banco de dados
docker-compose up -d payments-db

# Aguardar o MySQL inicializar (cerca de 30 segundos)
sleep 30

# As migrations serão executadas automaticamente
```

### Opção B: MySQL Local

1. Criar o banco de dados:

```sql
CREATE DATABASE payments_db;
```

2. Configurar o `.env`:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=sua_senha
DB_NAME=payments_db
GRPC_PORT=50051
```

3. Executar migrations:

```bash
make migrate
```

## 4. Executar o Serviço

### Opção A: Desenvolvimento Local

```bash
make run
```

### Opção B: Docker Completo

```bash
docker-compose up -d
```

O serviço estará disponível em `localhost:50051`

## 5. Testar o Serviço

### Usando grpcurl

1. Listar serviços disponíveis:

```bash
grpcurl -plaintext localhost:50051 list
```

Resultado esperado:
```
grpc.reflection.v1alpha.ServerReflection
payment.PaymentService
```

2. Listar métodos do PaymentService:

```bash
grpcurl -plaintext localhost:50051 list payment.PaymentService
```

3. Processar um pagamento:

```bash
grpcurl -plaintext -d '{
  "order_id": "order-123",
  "amount": 150.75,
  "payment_method": 1,
  "customer_email": "cliente@example.com",
  "customer_name": "João Silva"
}' localhost:50051 payment.PaymentService/ProcessPayment
```

4. Buscar um pagamento (use o payment_id retornado acima):

```bash
grpcurl -plaintext -d '{
  "payment_id": "SEU_PAYMENT_ID"
}' localhost:50051 payment.PaymentService/GetPayment
```

5. Cancelar um pagamento:

```bash
grpcurl -plaintext -d '{
  "payment_id": "SEU_PAYMENT_ID",
  "reason": "Cliente solicitou cancelamento"
}' localhost:50051 payment.PaymentService/CancelPayment
```

6. Listar pagamentos de um pedido:

```bash
grpcurl -plaintext -d '{
  "order_id": "order-123"
}' localhost:50051 payment.PaymentService/ListPayments
```

## 6. Métodos de Pagamento

Os seguintes métodos de pagamento são suportados:

- `1` - CREDIT_CARD (Cartão de Crédito)
- `2` - DEBIT_CARD (Cartão de Débito)
- `3` - PIX
- `4` - BOLETO
- `5` - PAYPAL

## 7. Status de Pagamento

Os pagamentos podem ter os seguintes status:

- `1` - PENDING (Pendente)
- `2` - PROCESSING (Processando)
- `3` - APPROVED (Aprovado)
- `4` - DECLINED (Recusado)
- `5` - CANCELED (Cancelado)
- `6` - REFUNDED (Reembolsado)

## 8. Executar Testes

```bash
# Testes unitários
make test

# Testes com coverage
make test-coverage
```

## 9. Verificar Logs

### Docker

```bash
docker-compose logs -f payments-service
```

### Local

Os logs são exibidos no terminal onde o serviço está rodando.

## 10. Parar o Serviço

### Docker

```bash
docker-compose down
```

### Local

Pressione `Ctrl+C` no terminal onde o serviço está rodando.

## Troubleshooting

### Erro: "protoc: command not found"

Instale o Protocol Buffers compiler:
```bash
brew install protobuf
```

### Erro: "connection refused" ao conectar no MySQL

Aguarde alguns segundos para o MySQL inicializar completamente ou verifique se o serviço está rodando:
```bash
docker-compose ps
```

### Erro ao gerar proto files

Certifique-se de que os plugins Go estão instalados e no PATH:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Porta 50051 já em uso

Altere a porta no `.env`:
```env
GRPC_PORT=50052
```

## Próximos Passos

Depois de configurar o Payment Service, consulte o arquivo `INTEGRATION.md` para integrar com o Orders Service.
