# Payment Service - Resumo do Projeto

## 📋 Visão Geral

Serviço de pagamentos desenvolvido em Go que processa transações via gRPC, seguindo Clean Architecture e DDD.

## ✨ Características Principais

- ✅ Comunicação via gRPC (alta performance)
- ✅ Múltiplos métodos de pagamento (Cartão, PIX, Boleto, PayPal)
- ✅ Clean Architecture (fácil manutenção e testes)
- ✅ Domain-Driven Design
- ✅ Persistência em MySQL
- ✅ Logging estruturado (slog)
- ✅ Containerização com Docker
- ✅ Testes unitários

## 🚀 Quick Start

```bash
# 1. Clone e entre no diretório
cd payments

# 2. Instalar dependências
go mod download

# 3. Gerar código gRPC
make proto

# 4. Iniciar banco de dados
docker-compose up -d payments-db

# 5. Aguardar MySQL inicializar
sleep 30

# 6. Executar serviço
make run
```

O serviço estará disponível em `localhost:50051`

## 🧪 Testar Rapidamente

```bash
# Processar um pagamento
grpcurl -plaintext -d '{
  "order_id": "test-001",
  "amount": 100.00,
  "payment_method": 1,
  "customer_email": "test@example.com",
  "customer_name": "Test User"
}' localhost:50051 payment.PaymentService/ProcessPayment
```

## 📁 Estrutura

```
payments/
├── cmd/grpc/main.go              # Entry point
├── internal/
│   ├── domain/                   # Entidades e interfaces
│   ├── usecase/                  # Casos de uso
│   └── infra/                    # Implementações técnicas
├── proto/                        # Definições gRPC
├── migrations/                   # SQL migrations
├── tests/                        # Testes
└── examples/                     # Exemplos de integração
```

## 🔌 Integração com Orders

Veja `INTEGRATION.md` para instruções detalhadas de como integrar com o serviço Orders.

### Resumo
1. Copie o arquivo `proto/payment.proto` para o projeto orders
2. Gere o código gRPC no orders
3. Use o cliente de exemplo em `examples/client/payment_client.go`
4. Chame o Payment Service a partir dos handlers do Orders

## 📚 Documentação Completa

- **README.md** - Documentação principal
- **SETUP.md** - Guia de configuração detalhado
- **ARCHITECTURE.md** - Arquitetura e design
- **INTEGRATION.md** - Como integrar com Orders
- **EXAMPLES.md** - Exemplos de requisições
- **proto/README.md** - Sobre Protocol Buffers

## 🛠️ Comandos Úteis

```bash
make proto          # Gerar código gRPC
make run            # Executar serviço
make test           # Executar testes
make test-coverage  # Testes com coverage
make build          # Build do binário
make docker-build   # Build da imagem Docker
make docker-run     # Executar com Docker
make clean          # Limpar arquivos gerados
```

## 🔧 Configuração

Variáveis de ambiente (`.env`):
```env
DB_HOST=localhost
DB_PORT=3307
DB_USER=root
DB_PASSWORD=root
DB_NAME=payments_db
GRPC_PORT=50051
```

## 📊 Métodos de Pagamento Suportados

| Código | Método          |
|--------|-----------------|
| 1      | Cartão Crédito  |
| 2      | Cartão Débito   |
| 3      | PIX             |
| 4      | Boleto          |
| 5      | PayPal          |

## 📈 Status de Pagamento

| Código | Status       | Descrição                    |
|--------|--------------|------------------------------|
| 1      | PENDING      | Aguardando processamento     |
| 2      | PROCESSING   | Em processamento             |
| 3      | APPROVED     | Aprovado                     |
| 4      | DECLINED     | Recusado                     |
| 5      | CANCELED     | Cancelado                    |
| 6      | REFUNDED     | Reembolsado                  |

## 🧩 API gRPC

### ProcessPayment
Processa um novo pagamento para um pedido.

### GetPayment
Busca detalhes de um pagamento específico.

### CancelPayment
Cancela um pagamento pendente ou em processamento.

### ListPayments
Lista todos os pagamentos de um pedido.

## 🐳 Docker

### Executar tudo com Docker Compose

```bash
docker-compose up -d
```

Isso iniciará:
- MySQL (porta 3307)
- Payment Service (porta 50051)

### Verificar logs

```bash
docker-compose logs -f payments-service
```

### Parar serviços

```bash
docker-compose down
```

## 🧪 Testes

### Executar testes

```bash
make test
```

### Coverage

```bash
make test-coverage
```

Isso gera um arquivo `coverage.html` que pode ser aberto no navegador.

## 🔍 Debugging

### Ver serviços disponíveis

```bash
grpcurl -plaintext localhost:50051 list
```

### Ver métodos de um serviço

```bash
grpcurl -plaintext localhost:50051 list payment.PaymentService
```

### Descrever um método

```bash
grpcurl -plaintext localhost:50051 describe payment.PaymentService.ProcessPayment
```

## 🔐 Segurança (Produção)

Para produção, considere:
- [ ] Adicionar TLS para gRPC
- [ ] Implementar autenticação (JWT)
- [ ] Rate limiting
- [ ] Circuit breaker
- [ ] Criptografia de dados sensíveis
- [ ] Audit logs

## 📊 Monitoramento (Futuro)

- [ ] Prometheus metrics
- [ ] Distributed tracing
- [ ] Health checks
- [ ] APM integration
- [ ] Alertas

## 🤝 Integração com Gateway Real

Atualmente, o serviço simula o gateway de pagamento. Para integrar com gateway real:

1. Criar interface `PaymentGateway` no domínio
2. Implementar adaptadores (Stripe, PayPal, etc.)
3. Injetar no use case
4. Configurar credenciais do gateway

## 📝 TODO

- [ ] Implementar webhook handler para notificações de gateway
- [ ] Adicionar suporte a refund
- [ ] Implementar processamento assíncrono
- [ ] Adicionar eventos de pagamento (Event Sourcing)
- [ ] Integrar com message broker (Kafka/RabbitMQ)
- [ ] Adicionar mais testes (integration, e2e)
- [ ] Implementar retry logic
- [ ] Adicionar cache (Redis)

## 🐛 Troubleshooting

### Erro: "protoc: command not found"
```bash
brew install protobuf
```

### Erro: "connection refused" ao conectar no MySQL
```bash
# Aguarde o MySQL inicializar
docker-compose logs payments-db
```

### Erro ao gerar proto files
```bash
# Instalar plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Adicionar ao PATH
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Porta 50051 já em uso
```bash
# Alterar porta no .env
GRPC_PORT=50052
```

## 📧 Suporte

Para dúvidas ou problemas:
1. Consulte a documentação em `docs/`
2. Verifique exemplos em `examples/`
3. Revise os testes em `tests/`

## 🎯 Conceitos Aplicados

- ✅ Clean Architecture
- ✅ Domain-Driven Design (DDD)
- ✅ Dependency Inversion
- ✅ Repository Pattern
- ✅ Use Case Pattern
- ✅ gRPC / Protocol Buffers
- ✅ Structured Logging
- ✅ Database Connection Pooling
- ✅ Docker & Docker Compose
- ✅ Unit Testing

## 🌟 Boas Práticas Implementadas

- Separação clara de camadas
- Interfaces bem definidas
- Código facilmente testável
- Logs estruturados
- Configuração via variáveis de ambiente
- Tratamento adequado de erros
- Validações de domínio
- Migrations versionadas
- Documentação abrangente

## 🚀 Próximos Passos Sugeridos

1. **Adicionar Autenticação**: Implementar JWT/OAuth
2. **Observabilidade**: Adicionar métricas e tracing
3. **Gateway Real**: Integrar com Stripe/PayPal
4. **Events**: Publicar eventos de pagamento
5. **Cache**: Adicionar Redis para cache
6. **Rate Limiting**: Proteger contra abuso
7. **API Gateway**: Adicionar Kong/Traefik
8. **Service Mesh**: Considerar Istio para produção

---

**Desenvolvido com ❤️ em Go**
