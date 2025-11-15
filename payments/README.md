# Payment Service 💳

Serviço de pagamentos em Go que processa transações via gRPC, seguindo Clean Architecture e Domain-Driven Design.

## 📚 Documentação

- **[SETUP.md](SETUP.md)** - Guia completo de configuração e instalação
- **[ARCHITECTURE.md](ARCHITECTURE.md)** - Detalhes da arquitetura e design patterns
- **[INTEGRATION.md](INTEGRATION.md)** - Como integrar com o serviço Orders
- **[EXAMPLES.md](EXAMPLES.md)** - Exemplos práticos de requisições
- **[DIAGRAMS.md](DIAGRAMS.md)** - Diagramas da arquitetura e fluxos
- **[SUMMARY.md](SUMMARY.md)** - Resumo rápido do projeto
- **[examples/](examples/)** - Código de exemplo para integração

## 🚀 Quick Start

```bash
# 1. Executar script de setup automático
chmod +x quick-start.sh
./quick-start.sh

# OU seguir manualmente:

# 2. Baixar dependências
go mod download

# 3. Gerar código gRPC
make proto

# 4. Iniciar banco de dados
docker-compose up -d payments-db

# 5. Executar serviço
make run
```

Serviço disponível em `localhost:50051`

## ✨ Características

- ✅ Comunicação via gRPC (alta performance)
- ✅ Clean Architecture (testável e manutenível)
- ✅ Domain-Driven Design (DDD)
- ✅ Múltiplos métodos de pagamento (Cartão, PIX, Boleto, PayPal)
- ✅ Persistência em MySQL
- ✅ Logging estruturado
- ✅ Containerização com Docker
- ✅ Documentação completa

## 📋 Requisitos

- Go 1.21+
- MySQL 8.0+
- Docker e Docker Compose (opcional)
- Protocol Buffers compiler (protoc)

## Instalação

### Instalar dependências

```bash
go mod download
```

### Gerar código gRPC

```bash
make proto
```

### Configurar variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=payments_db
GRPC_PORT=50051
```

### Executar migrations

```bash
make migrate
```

### Executar o serviço

```bash
make run
```

## Docker

### Build da imagem

```bash
docker build -t payments-service .
```

### Executar com Docker Compose

```bash
docker-compose up -d
```

## Desenvolvimento

### Executar testes

```bash
make test
```

### Executar testes com coverage

```bash
make test-coverage
```

### Clean

```bash
make clean
```

## API gRPC

O serviço expõe os seguintes métodos via gRPC:

- `ProcessPayment`: Processa um novo pagamento
- `GetPayment`: Busca detalhes de um pagamento
- `CancelPayment`: Cancela um pagamento pendente

## Estrutura do Projeto

```
payments/
├── cmd/
│   └── grpc/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   └── repository/
│   ├── infra/
│   │   ├── database/
│   │   ├── grpc/
│   │   └── repository/
│   └── usecase/
├── migrations/
├── proto/
└── tests/
```
