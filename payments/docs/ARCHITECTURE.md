# Arquitetura do Payment Service

## Visão Geral

O Payment Service é um microserviço responsável por processar pagamentos via gRPC. Ele segue os princípios de Clean Architecture e Domain-Driven Design (DDD).

## Estrutura do Projeto

```
payments/
├── cmd/
│   └── grpc/
│       └── main.go                    # Entry point do servidor gRPC
├── internal/
│   ├── domain/                        # Camada de domínio
│   │   ├── entity/                    # Entidades de negócio
│   │   │   └── payment.go            # Entidade Payment
│   │   └── repository/                # Interfaces dos repositórios
│   │       └── repository.go         # PaymentRepository interface
│   ├── infra/                         # Camada de infraestrutura
│   │   ├── database/                  # Conexões de banco de dados
│   │   │   └── mysql.go              # Cliente MySQL
│   │   ├── grpc/
│   │   │   └── handler/              # Handlers gRPC
│   │   │       └── payment_handler.go # Implementação do servidor gRPC
│   │   └── repository/               # Implementações dos repositórios
│   │       └── payment_repository.go # Repositório MySQL
│   └── usecase/                       # Casos de uso (regras de negócio)
│       ├── process_payment_usecase.go
│       ├── get_payment_usecase.go
│       ├── cancel_payment_usecase.go
│       └── list_payments_usecase.go
├── proto/
│   ├── payment.proto                  # Definição do serviço gRPC
│   ├── payment.pb.go                  # Código gerado (messages)
│   └── payment_grpc.pb.go            # Código gerado (service)
├── migrations/
│   └── 001_create_tables.sql         # Scripts de migração do banco
├── tests/                             # Testes
│   └── internal/
│       └── domain/
│           └── entity/
│               └── payment_test.go
├── examples/                          # Exemplos de integração
│   ├── client/
│   │   └── payment_client.go         # Cliente gRPC de exemplo
│   └── integration/
│       └── orders_integration_example.go
├── docker-compose.yml                 # Configuração Docker
├── Dockerfile                         # Imagem Docker
├── Makefile                          # Comandos úteis
├── go.mod                            # Dependências Go
├── .env                              # Variáveis de ambiente
└── README.md                         # Documentação principal
```

## Camadas da Arquitetura

### 1. Domain Layer (Camada de Domínio)

**Responsabilidade**: Contém as regras de negócio puras, independentes de frameworks e infraestrutura.

#### Entities (`internal/domain/entity/`)
- **payment.go**: Define a entidade Payment com:
  - Atributos: ID, OrderID, Amount, PaymentMethod, Status, etc.
  - Métodos de negócio: Process(), Approve(), Decline(), Cancel(), Refund()
  - Validações de domínio
  - Estados permitidos (PaymentStatus, PaymentMethod)

#### Repository Interfaces (`internal/domain/repository/`)
- **repository.go**: Define contratos para acesso a dados
  - `PaymentRepository`: Interface com métodos CRUD

**Características**:
- Não depende de outras camadas
- Contém apenas lógica de negócio pura
- Fácil de testar (testes unitários puros)

### 2. Use Case Layer (Camada de Casos de Uso)

**Responsabilidade**: Orquestra o fluxo de dados entre as camadas, implementa regras de aplicação.

#### Use Cases (`internal/usecase/`)
- **process_payment_usecase.go**: Processa novos pagamentos
  - Valida dados de entrada
  - Cria entidade Payment
  - Simula gateway de pagamento
  - Persiste no repositório
  
- **get_payment_usecase.go**: Busca detalhes de pagamento
- **cancel_payment_usecase.go**: Cancela pagamentos pendentes
- **list_payments_usecase.go**: Lista pagamentos por order_id

**Características**:
- Depende apenas da camada de domínio
- Coordena operações entre entidades e repositórios
- Implementa regras de aplicação (não de negócio)

### 3. Infrastructure Layer (Camada de Infraestrutura)

**Responsabilidade**: Implementa detalhes técnicos e integrações externas.

#### Database (`internal/infra/database/`)
- **mysql.go**: Gerencia conexão com MySQL
  - Pool de conexões
  - Health checks
  - Configurações de timeout

#### gRPC Handlers (`internal/infra/grpc/handler/`)
- **payment_handler.go**: Implementa o servidor gRPC
  - Converte proto messages para entidades
  - Chama use cases apropriados
  - Trata erros e retorna respostas
  - Implementa todos os métodos do PaymentService:
    - ProcessPayment
    - GetPayment
    - CancelPayment
    - ListPayments

#### Repository Implementations (`internal/infra/repository/`)
- **payment_repository.go**: Implementação MySQL do PaymentRepository
  - Operações CRUD no banco
  - Tratamento de erros SQL
  - Conversões entre entidades e modelos de banco

**Características**:
- Depende das interfaces definidas no domínio
- Implementa detalhes técnicos
- Facilmente substituível (ex: trocar MySQL por PostgreSQL)

### 4. Interface Layer (gRPC)

**Responsabilidade**: Define o contrato de comunicação com clientes externos.

#### Proto Definitions (`proto/`)
- **payment.proto**: Define:
  - Messages (request/response)
  - Enums (PaymentMethod, PaymentStatus)
  - Service (PaymentService)
  - RPCs (ProcessPayment, GetPayment, etc.)

**Características**:
- Contract-first design
- Type-safe communication
- Language-agnostic (pode ser usado por clientes em outras linguagens)

## Fluxo de Dados

### Exemplo: Processar Pagamento

```
1. Cliente gRPC (Orders Service)
   ↓
2. PaymentServiceServer.ProcessPayment (Handler)
   ↓ [converte proto → entity]
3. ProcessPaymentUseCase.Execute
   ↓ [cria entidade Payment]
4. Payment.Process() (Entity)
   ↓ [valida e atualiza estado]
5. PaymentRepository.Create (Interface)
   ↓
6. PaymentRepositoryMySQL.Create (Implementation)
   ↓
7. MySQL Database
   ↓ [retorna sucesso]
8. ← Resposta sobe pelas camadas
   ↓
9. Cliente recebe ProcessPaymentResponse
```

## Princípios Aplicados

### 1. Clean Architecture
- **Independência de Frameworks**: Lógica de negócio não depende de frameworks específicos
- **Testabilidade**: Cada camada pode ser testada isoladamente
- **Independência de UI**: Pode expor REST, gRPC, GraphQL sem alterar domínio
- **Independência de Banco**: Pode trocar MySQL por outro banco

### 2. Domain-Driven Design (DDD)
- **Entities**: Payment com identidade e comportamento
- **Value Objects**: PaymentMethod, PaymentStatus
- **Repository Pattern**: Abstração de acesso a dados
- **Use Cases**: Representam intenções de negócio

### 3. Dependency Inversion
- Camadas de alto nível não dependem de detalhes
- Abstrações (interfaces) no domínio
- Implementações concretas na infraestrutura

### 4. Single Responsibility
- Cada camada tem uma responsabilidade única
- Cada use case resolve um problema específico
- Handlers apenas convertem formatos

## Tecnologias Utilizadas

### Core
- **Go 1.21+**: Linguagem de programação
- **gRPC**: Framework de comunicação
- **Protocol Buffers**: Serialização de dados

### Banco de Dados
- **MySQL 8.0**: Banco de dados relacional
- **database/sql**: Driver Go para MySQL

### Infraestrutura
- **Docker & Docker Compose**: Containerização
- **godotenv**: Gerenciamento de variáveis de ambiente

### Observabilidade
- **log/slog**: Logging estruturado

## Segurança

### Atual (Desenvolvimento)
- Conexão gRPC sem TLS (insecure)
- Sem autenticação/autorização

### Recomendações para Produção
- [ ] Adicionar TLS para gRPC
- [ ] Implementar autenticação via tokens (JWT)
- [ ] Adicionar rate limiting
- [ ] Implementar circuit breaker
- [ ] Criptografar dados sensíveis (card details)
- [ ] Adicionar audit logs
- [ ] Implementar RBAC para operações sensíveis

## Escalabilidade

### Horizontal Scaling
- Serviço stateless (pode escalar horizontalmente)
- Múltiplas instâncias com load balancer

### Performance
- Connection pooling no MySQL
- Context timeouts em todas as operações
- Índices no banco de dados (order_id, status, created_at)

### Monitoring (Futuro)
- [ ] Prometheus metrics
- [ ] Distributed tracing (Jaeger/Zipkin)
- [ ] Health checks endpoints
- [ ] APM integration

## Integração com Orders Service

```
┌─────────────────┐         gRPC          ┌──────────────────┐
│  Orders Service │ ────────────────────→ │ Payment Service  │
│   (HTTP/REST)   │                       │     (gRPC)       │
└─────────────────┘                       └──────────────────┘
        │                                           │
        │                                           │
        ↓                                           ↓
  ┌──────────┐                              ┌──────────┐
  │ Orders DB│                              │Payment DB│
  └──────────┘                              └──────────┘
```

### Comunicação
1. Orders Service recebe requisição HTTP do cliente
2. Orders Service chama Payment Service via gRPC
3. Payment Service processa e retorna status
4. Orders Service atualiza status do pedido
5. Orders Service retorna resposta HTTP ao cliente

## Extensibilidade

### Adicionar Novo Método de Pagamento
1. Adicionar enum em `payment.proto`
2. Adicionar constante em `entity/payment.go`
3. Atualizar validações
4. Implementar lógica específica no use case

### Adicionar Novo Status
1. Adicionar enum em `payment.proto`
2. Adicionar constante em `entity/payment.go`
3. Implementar transições de estado na entidade
4. Atualizar handlers

### Integrar Gateway Real
1. Criar interface `PaymentGateway` no domínio
2. Implementar adaptadores para gateways específicos (Stripe, PayPal, etc.)
3. Injetar no use case via dependency injection

## Testing Strategy

### Unit Tests
- Entidades: Testar regras de negócio
- Use Cases: Testar lógica de aplicação (com mocks)

### Integration Tests
- Repositórios: Testar com banco real (testcontainers)
- gRPC: Testar handlers com use cases mockados

### E2E Tests
- Testar fluxo completo com banco e gRPC
- Usar grpcurl ou cliente de teste

## Próximos Passos

1. **Segurança**: Adicionar TLS e autenticação
2. **Monitoring**: Implementar observabilidade completa
3. **Gateway Integration**: Integrar com gateway real (Stripe, etc.)
4. **Event Driven**: Publicar eventos de pagamento (Kafka, RabbitMQ)
5. **Async Processing**: Processar pagamentos de forma assíncrona
6. **Webhook Handler**: Receber notificações de gateways
7. **Retry Logic**: Implementar retentativas para falhas transitórias
