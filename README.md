# Sistema E-commerce com Microserviços

Sistema de e-commerce construído com arquitetura de microserviços usando Go, gRPC e Clean Architecture.

## 🏗️ Arquitetura

```
┌─────────────────────┐
│   Orders Service    │
│    (HTTP/REST)      │
│    Port: 8080       │
└──────────┬──────────┘
           │ gRPC
           │
           ▼
┌─────────────────────┐
│  Payments Service   │
│      (gRPC)         │
│    Port: 50051      │
└─────────────────────┘
```

## 📦 Serviços

### Orders Service
- **Tecnologia**: Go + HTTP/REST + Chi Router
- **Porta**: 8080
- **Banco de dados**: MySQL
- **Responsabilidades**:
  - Gerenciar produtos
  - Gerenciar pedidos
  - Gerenciar carrinho de compras
  - Integrar com Payment Service via gRPC

### Payments Service
- **Tecnologia**: Go + gRPC
- **Porta**: 50051
- **Banco de dados**: MySQL
- **Responsabilidades**:
  - Processar pagamentos
  - Gerenciar status de pagamentos
  - Cancelar/Reembolsar pagamentos
  - Listar pagamentos por pedido

## 🚀 Como Executar

### Pré-requisitos

- Go 1.24+
- MySQL 8.0+
- Protocol Buffers compiler (protoc)
- Make (opcional)

### 1. Configurar Bancos de Dados

```sql
-- Banco Orders
CREATE DATABASE orders_db;

-- Banco Payments
CREATE DATABASE payments_db;
```

Execute as migrations:
```bash
# Orders
mysql -u root -p orders_db < orders/migrations/001_create_tables.sql

# Payments
mysql -u root -p payments_db < payments/migrations/001_create_tables.sql
```

### 2. Configurar Variáveis de Ambiente

**Orders (.env)**
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=orders_user
DB_PASSWORD=orders_pass
DB_NAME=orders_db
SERVER_PORT=8080
PAYMENT_SERVICE_ADDR=localhost:50051
```

**Payments (.env)**
```env
DB_DSN=root:root@tcp(localhost:3306)/payments_db?parseTime=true
GRPC_PORT=50051
```

### 3. Compilar os Serviços

```bash
# Payments
cd payments
go mod download
go build -o bin/payment-service ./cmd/grpc

# Orders
cd ../orders
go mod download
go build ./cmd/api
```

### 4. Iniciar os Serviços

**Opção 1: Script Automático (Recomendado)**
```bash
./start-services.sh
```

**Opção 2: Manual**
```bash
# Terminal 1 - Payment Service
cd payments
./bin/payment-service

# Terminal 2 - Orders Service
cd orders
go run cmd/api/main.go
```

### 5. Testar a Integração

```bash
./test-integration.sh
```

Ou manualmente:
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
      }
    ]
  }'
```

## 📚 Documentação

- [INTEGRATION.md](./INTEGRATION.md) - Documentação detalhada da integração gRPC
- [Orders Swagger](http://localhost:8080/swagger/) - Documentação da API REST
- [Payments Proto](./payments/proto/README.md) - Documentação do contrato gRPC

## 🧪 Testes

### Testes Unitários

```bash
# Orders
cd orders
go test ./...

# Payments
cd payments
go test ./...
```

### Testes de Integração

```bash
./test-integration.sh
```

## 📖 Endpoints

### Orders Service (HTTP)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/health` | Health check |
| GET | `/swagger/*` | Documentação Swagger |
| POST | `/api/v1/orders/with-payment` | Criar pedido com pagamento |
| POST | `/api/v1/orders/{id}/cancel` | Cancelar pedido e pagamento |
| GET | `/api/v1/orders` | Listar pedidos |
| GET | `/api/v1/orders/{id}` | Buscar pedido |
| GET | `/api/v1/products` | Listar produtos |
| POST | `/api/v1/products` | Criar produto |
| POST | `/api/v1/cart` | Criar carrinho |

### Payments Service (gRPC)

| RPC Method | Descrição |
|------------|-----------|
| `ProcessPayment` | Processar pagamento |
| `GetPayment` | Buscar pagamento |
| `CancelPayment` | Cancelar pagamento |
| `ListPayments` | Listar pagamentos de um pedido |

## 🔧 Métodos de Pagamento

```
1 = CREDIT_CARD  (Cartão de Crédito)
2 = DEBIT_CARD   (Cartão de Débito)
3 = PIX          (Pix)
4 = BOLETO       (Boleto Bancário)
5 = PAYPAL       (PayPal)
```

## 📊 Status

### Pagamento
- `PENDING` - Aguardando processamento
- `PROCESSING` - Em processamento
- `APPROVED` - Aprovado
- `DECLINED` - Recusado
- `CANCELED` - Cancelado
- `REFUNDED` - Reembolsado

### Pedido
- `pending` - Pendente
- `paid` - Pago
- `canceled` - Cancelado
- `completed` - Completo

## 🛠️ Tecnologias

### Backend
- **Go 1.24**: Linguagem principal
- **gRPC**: Comunicação entre microserviços
- **Protocol Buffers**: Serialização de dados
- **Chi Router**: HTTP router para REST API
- **MySQL**: Banco de dados
- **Swagger**: Documentação da API

### Ferramentas
- **protoc**: Compiler de Protocol Buffers
- **protoc-gen-go**: Plugin Go para protoc
- **protoc-gen-go-grpc**: Plugin gRPC para protoc
- **swag**: Gerador de documentação Swagger

## 📁 Estrutura do Projeto

```
.
├── orders/                      # Serviço de Pedidos
│   ├── cmd/api/                # Entrypoint HTTP
│   ├── internal/
│   │   ├── domain/             # Entidades e interfaces
│   │   ├── usecase/            # Casos de uso
│   │   └── infra/              # Infraestrutura
│   │       ├── http/           # Handlers HTTP
│   │       ├── grpc/client/    # Cliente gRPC
│   │       ├── repository/     # Repositórios
│   │       └── database/       # Conexão DB
│   ├── proto/                  # Arquivos proto copiados
│   └── migrations/             # Migrations do banco
│
├── payments/                    # Serviço de Pagamentos
│   ├── cmd/grpc/               # Entrypoint gRPC
│   ├── internal/
│   │   ├── domain/             # Entidades e interfaces
│   │   ├── usecase/            # Casos de uso
│   │   └── infra/              # Infraestrutura
│   │       ├── grpc/handler/   # Handlers gRPC
│   │       ├── repository/     # Repositórios
│   │       └── database/       # Conexão DB
│   ├── proto/                  # Definições protobuf
│   ├── migrations/             # Migrations do banco
│   └── bin/                    # Binário compilado
│
├── start-services.sh           # Script para iniciar serviços
├── test-integration.sh         # Script de testes
├── INTEGRATION.md              # Documentação de integração
└── README.md                   # Este arquivo
```

## 🎯 Padrões Utilizados

- **Clean Architecture**: Separação em camadas (domain, usecase, infra)
- **Domain-Driven Design**: Entidades ricas com regras de negócio
- **Repository Pattern**: Abstração do acesso a dados
- **Dependency Injection**: Injeção via construtores
- **SOLID Principles**: Código limpo e manutenível

## 🔐 Segurança

**Implementado:**
- ✅ Validação de entrada
- ✅ Structured logging
- ✅ Timeouts em operações gRPC
- ✅ Graceful shutdown

**Próximos Passos:**
- [ ] Autenticação JWT
- [ ] TLS/SSL para gRPC
- [ ] Rate limiting
- [ ] Circuit breaker
- [ ] Tracing distribuído

## 📝 Licença

MIT

## 👥 Contribuindo

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Add nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## 📧 Contato

Para dúvidas ou sugestões, abra uma issue no GitHub.
