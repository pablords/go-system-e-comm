# Diagramas - Payment Service

## Arquitetura Geral

```
┌─────────────────────────────────────────────────────────────────┐
│                         Orders Service                          │
│                          (Port 8080)                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │   HTTP       │  │   Domain     │  │  Repository  │         │
│  │   Handlers   │→ │   Use Cases  │→ │              │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
│         │                                       ↓                │
│         │                                  ┌─────────┐          │
│         │                                  │Orders DB│          │
│         │                                  └─────────┘          │
│         │                                                        │
│         │  gRPC Call                                            │
│         ↓                                                        │
└─────────┼─────────────────────────────────────────────────────┘
          │
          │ ProcessPayment(orderID, amount, method)
          │
          ↓
┌─────────┼─────────────────────────────────────────────────────┐
│         ↓                                                        │
│  ┌─────────────────┐         Payment Service                   │
│  │  gRPC Handler   │            (Port 50051)                   │
│  │ (payment_handler│                                            │
│  │      .go)       │                                            │
│  └────────┬────────┘                                            │
│           │                                                      │
│           ↓                                                      │
│  ┌──────────────────┐                                          │
│  │    Use Cases     │                                          │
│  │ ┌──────────────┐ │                                          │
│  │ │ProcessPayment│ │                                          │
│  │ │GetPayment    │ │                                          │
│  │ │CancelPayment │ │                                          │
│  │ │ListPayments  │ │                                          │
│  │ └──────────────┘ │                                          │
│  └────────┬─────────┘                                          │
│           │                                                      │
│           ↓                                                      │
│  ┌──────────────────┐                                          │
│  │     Domain       │                                          │
│  │ ┌──────────────┐ │                                          │
│  │ │   Payment    │ │  Entidade com regras de negócio         │
│  │ │   Entity     │ │                                          │
│  │ └──────────────┘ │                                          │
│  └────────┬─────────┘                                          │
│           │                                                      │
│           ↓                                                      │
│  ┌──────────────────┐                                          │
│  │   Repository     │                                          │
│  │  Implementation  │                                          │
│  │ (payment_repo.go)│                                          │
│  └────────┬─────────┘                                          │
│           │                                                      │
│           ↓                                                      │
│      ┌──────────┐                                               │
│      │Payment DB│                                               │
│      │ (MySQL)  │                                               │
│      └──────────┘                                               │
└─────────────────────────────────────────────────────────────────┘
```

## Fluxo de Processamento de Pagamento

```
Cliente → Orders API → Payment Service → Gateway → Banco de Dados

1. Cliente faz requisição HTTP POST /orders
   ↓
2. Orders cria pedido e calcula total
   ↓
3. Orders chama Payment.ProcessPayment via gRPC
   ↓
4. PaymentHandler recebe requisição
   ↓
5. ProcessPaymentUseCase é executado
   ↓
6. Payment Entity é criada (NewPayment)
   ↓
7. Payment.Process() simula gateway
   ↓
8. Payment.Approve() ou Decline()
   ↓
9. PaymentRepository.Create() salva no banco
   ↓
10. Resposta retorna pela cadeia
   ↓
11. Orders atualiza status do pedido
   ↓
12. Cliente recebe resposta HTTP
```

## Estrutura de Camadas (Clean Architecture)

```
┌─────────────────────────────────────────────────────────────┐
│                      External World                         │
│  ┌────────────┐  ┌────────────┐  ┌──────────────┐         │
│  │   gRPC     │  │  Database  │  │  External    │         │
│  │  Clients   │  │            │  │  Gateways    │         │
│  └────────────┘  └────────────┘  └──────────────┘         │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│            Interface Adapters (Infrastructure)              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │    gRPC      │  │  Repository  │  │   Database   │     │
│  │   Handler    │  │     Impl     │  │   Client     │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                Application Business Rules                   │
│                       (Use Cases)                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Process    │  │     Get      │  │    Cancel    │     │
│  │   Payment    │  │   Payment    │  │   Payment    │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│              Enterprise Business Rules                      │
│                      (Domain)                               │
│  ┌──────────────────────────────────────────────┐          │
│  │            Payment Entity                     │          │
│  │  - ID, OrderID, Amount, Status, Method       │          │
│  │  - Process(), Approve(), Decline(), Cancel() │          │
│  │  - Business Validations                      │          │
│  └──────────────────────────────────────────────┘          │
│                                                              │
│  ┌──────────────────────────────────────────────┐          │
│  │         Repository Interface                  │          │
│  │  - Create(), FindByID(), Update(), Delete()  │          │
│  └──────────────────────────────────────────────┘          │
└─────────────────────────────────────────────────────────────┘

Dependências apontam para dentro (Dependency Rule)
```

## Estados e Transições de Pagamento

```
┌─────────┐
│ PENDING │ ────────────────────────────┐
└────┬────┘                              │
     │                                   │
     │ Process()                         │ Cancel()
     ↓                                   │
┌────────────┐                           │
│ PROCESSING │ ────────────────────┐     │
└─────┬──────┘                     │     │
      │                            │     │
      │                            │     │
      │                            │     │
   ┌──▼───────┐               ┌───▼─────▼─┐
   │ APPROVED │               │ DECLINED  │
   └────┬─────┘               └───────────┘
        │
        │ Refund()
        ↓
   ┌──────────┐              ┌──────────┐
   │ REFUNDED │              │ CANCELED │
   └──────────┘              └──────────┘

Regras:
- PENDING pode ir para PROCESSING ou CANCELED
- PROCESSING pode ir para APPROVED, DECLINED ou CANCELED
- APPROVED pode ir para REFUNDED
- DECLINED pode ser cancelado
- APPROVED, CANCELED e REFUNDED são estados finais (não podem ser cancelados)
```

## Fluxo de Dados - ProcessPayment

```
┌──────────────────────────────────────────────────────────────┐
│ 1. gRPC Request (ProcessPaymentRequest)                      │
│    - order_id: "order-123"                                   │
│    - amount: 100.50                                          │
│    - payment_method: CREDIT_CARD                             │
│    - customer_email: "user@example.com"                      │
└────────────────────────┬─────────────────────────────────────┘
                         │
                         ↓
┌──────────────────────────────────────────────────────────────┐
│ 2. PaymentHandler.ProcessPayment()                          │
│    - Converte proto → domain types                          │
│    - Cria ProcessPaymentInput                               │
└────────────────────────┬─────────────────────────────────────┘
                         │
                         ↓
┌──────────────────────────────────────────────────────────────┐
│ 3. ProcessPaymentUseCase.Execute()                          │
│    - Valida input                                           │
│    - Cria Payment entity via NewPayment()                   │
└────────────────────────┬─────────────────────────────────────┘
                         │
                         ↓
┌──────────────────────────────────────────────────────────────┐
│ 4. Payment Entity                                            │
│    - NewPayment() valida e cria entidade                    │
│    - Status inicial: PENDING                                │
└────────────────────────┬─────────────────────────────────────┘
                         │
                         ↓
┌──────────────────────────────────────────────────────────────┐
│ 5. Payment.Process(transactionID)                           │
│    - Muda status para PROCESSING                            │
│    - Associa transaction_id                                 │
└────────────────────────┬─────────────────────────────────────┘
                         │
                         ↓
┌──────────────────────────────────────────────────────────────┐
│ 6. simulatePaymentGateway()                                  │
│    - Simula chamada ao gateway                              │
│    - Retorna approved/declined                              │
└────────────────────────┬─────────────────────────────────────┘
                         │
                    ┌────┴────┐
                    │         │
              ┌─────▼───┐ ┌──▼──────┐
              │Approved │ │Declined │
              └─────┬───┘ └──┬──────┘
                    │        │
                    └────┬───┘
                         │
                         ↓
┌──────────────────────────────────────────────────────────────┐
│ 7. Payment.Approve() ou Payment.Decline()                    │
│    - Atualiza status                                         │
│    - Atualiza updated_at                                     │
└────────────────────────┬─────────────────────────────────────┘
                         │
                         ↓
┌──────────────────────────────────────────────────────────────┐
│ 8. PaymentRepository.Create(payment)                         │
│    - Salva no MySQL                                          │
│    - Retorna sucesso/erro                                    │
└────────────────────────┬─────────────────────────────────────┘
                         │
                         ↓
┌──────────────────────────────────────────────────────────────┐
│ 9. ProcessPaymentOutput                                      │
│    - payment_id, status, message, transaction_id            │
└────────────────────────┬─────────────────────────────────────┘
                         │
                         ↓
┌──────────────────────────────────────────────────────────────┐
│ 10. ProcessPaymentResponse (proto)                           │
│     - Converte domain → proto                                │
│     - Retorna ao cliente gRPC                                │
└──────────────────────────────────────────────────────────────┘
```

## Estrutura de Diretórios Detalhada

```
payments/
│
├── cmd/                          # Application entry points
│   └── grpc/
│       └── main.go              # Inicializa servidor gRPC
│
├── internal/                     # Private application code
│   │
│   ├── domain/                  # Domain layer (business logic)
│   │   ├── entity/
│   │   │   └── payment.go      # Payment entity + business rules
│   │   └── repository/
│   │       └── repository.go   # Repository interfaces
│   │
│   ├── usecase/                 # Application layer (use cases)
│   │   ├── process_payment_usecase.go
│   │   ├── get_payment_usecase.go
│   │   ├── cancel_payment_usecase.go
│   │   └── list_payments_usecase.go
│   │
│   └── infra/                   # Infrastructure layer
│       ├── database/
│       │   └── mysql.go        # Database connection
│       ├── grpc/
│       │   └── handler/
│       │       └── payment_handler.go  # gRPC server implementation
│       └── repository/
│           └── payment_repository.go   # MySQL repository impl
│
├── proto/                        # Protocol Buffers definitions
│   ├── payment.proto            # Service & message definitions
│   ├── payment.pb.go            # Generated: messages
│   └── payment_grpc.pb.go       # Generated: service
│
├── migrations/                   # Database migrations
│   └── 001_create_tables.sql
│
├── tests/                        # Test files
│   ├── internal/
│   │   └── domain/
│   │       └── entity/
│   │           └── payment_test.go
│   └── mocks/                   # Test mocks
│
├── examples/                     # Integration examples
│   ├── client/
│   │   └── payment_client.go   # Example gRPC client
│   └── integration/
│       └── orders_integration_example.go
│
├── docs/                         # Additional documentation
│
├── docker-compose.yml            # Docker services definition
├── Dockerfile                    # Container image
├── Makefile                      # Build commands
├── go.mod                        # Go dependencies
├── go.sum                        # Dependency checksums
├── .env                          # Environment variables
├── .gitignore                    # Git ignore rules
│
└── Documentation Files
    ├── README.md                # Main documentation
    ├── SETUP.md                 # Setup guide
    ├── ARCHITECTURE.md          # Architecture details
    ├── INTEGRATION.md           # Integration guide
    ├── EXAMPLES.md              # Request examples
    ├── SUMMARY.md               # Quick reference
    └── DIAGRAMS.md              # This file
```

## Sequência de Integração Orders ↔ Payments

```
Client     Orders API    Payment Service    Orders DB    Payment DB
  │            │                │               │            │
  │  POST      │                │               │            │
  │  /orders   │                │               │            │
  ├───────────→│                │               │            │
  │            │                │               │            │
  │            │  INSERT order  │               │            │
  │            ├───────────────────────────────→│            │
  │            │                │               │            │
  │            │  gRPC: Process │               │            │
  │            │    Payment     │               │            │
  │            ├───────────────→│               │            │
  │            │                │               │            │
  │            │                │  INSERT       │            │
  │            │                │  payment      │            │
  │            │                ├───────────────────────────→│
  │            │                │               │            │
  │            │                │  SELECT       │            │
  │            │                │  payment      │            │
  │            │                │←──────────────────────────┤
  │            │                │               │            │
  │            │   Response     │               │            │
  │            │   (approved)   │               │            │
  │            │←───────────────┤               │            │
  │            │                │               │            │
  │            │  UPDATE order  │               │            │
  │            │  status=paid   │               │            │
  │            ├───────────────────────────────→│            │
  │            │                │               │            │
  │  Response  │                │               │            │
  │  200 OK    │                │               │            │
  │←───────────┤                │               │            │
  │            │                │               │            │
```

## Modelo de Dados

```
┌─────────────────────────────────────────────────────────┐
│                     payments                            │
├─────────────────────────────────────────────────────────┤
│ id               VARCHAR(36)    PRIMARY KEY             │
│ order_id         VARCHAR(36)    NOT NULL  [indexed]    │
│ amount           DECIMAL(10,2)  NOT NULL                │
│ payment_method   VARCHAR(50)    NOT NULL                │
│ status           VARCHAR(50)    NOT NULL  [indexed]    │
│ transaction_id   VARCHAR(255)                           │
│ customer_email   VARCHAR(255)   NOT NULL                │
│ customer_name    VARCHAR(255)   NOT NULL                │
│ created_at       TIMESTAMP      NOT NULL  [indexed]    │
│ updated_at       TIMESTAMP      NOT NULL                │
│ canceled_at      TIMESTAMP      NULL                    │
│ cancel_reason    TEXT           NULL                    │
└─────────────────────────────────────────────────────────┘

Indexes:
- PRIMARY KEY (id)
- INDEX idx_order_id (order_id)
- INDEX idx_status (status)
- INDEX idx_created_at (created_at)

Relacionamento com Orders:
orders.id ←─────→ payments.order_id (1:N)
```

## Deployment Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Load Balancer                          │
│                     (nginx/traefik)                         │
└────────────────────┬────────────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
        ↓            ↓            ↓
┌──────────┐  ┌──────────┐  ┌──────────┐
│ Orders-1 │  │ Orders-2 │  │ Orders-N │
│  :8080   │  │  :8080   │  │  :8080   │
└─────┬────┘  └─────┬────┘  └─────┬────┘
      │             │             │
      └─────────────┼─────────────┘
                    │ gRPC
        ┌───────────┼───────────┐
        │           │           │
        ↓           ↓           ↓
┌──────────┐  ┌──────────┐  ┌──────────┐
│Payment-1 │  │Payment-2 │  │Payment-N │
│  :50051  │  │  :50051  │  │  :50051  │
└─────┬────┘  └─────┬────┘  └─────┬────┘
      │             │             │
      └─────────────┼─────────────┘
                    │
        ┌───────────┼───────────┐
        ↓                       ↓
┌──────────────┐        ┌──────────────┐
│  Orders DB   │        │  Payments DB │
│   (MySQL)    │        │   (MySQL)    │
└──────────────┘        └──────────────┘

Features:
- Horizontal scaling for both services
- Separate databases (data isolation)
- Load balancing
- Service discovery (optional: Consul, etcd)
- Health checks
```

---

**Nota**: Estes diagramas são representações em ASCII. Para produção, considere usar ferramentas como:
- PlantUML
- Mermaid
- Draw.io
- Lucidchart
