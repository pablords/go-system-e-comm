# 🔄 Fluxo de Sincronização de Protos

## Arquitetura

```
┌─────────────────────────────────────────────────────────────────┐
│                    proto/ (Source of Truth)                      │
│                                                                  │
│  ┌────────────────────────────────────────────────────────┐   │
│  │  payment/payment.proto                                  │   │
│  │  - Define mensagens (Request/Response)                  │   │
│  │  - Define serviços (RPC methods)                        │   │
│  │  - Versionamento centralizado                           │   │
│  └────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                           │
                           │ make sync-all
                           │ ou ./sync-protos.sh
                           │
         ┌─────────────────┴─────────────────┐
         │                                   │
         ▼                                   ▼
┌──────────────────────┐           ┌──────────────────────┐
│  payments/           │           │  orders/             │
│  proto/payment.proto │           │  proto/payment.proto │
│         │            │           │         │            │
│         │            │           │         │            │
│         │ make proto │           │         │ make proto │
│         │            │           │         │            │
│         ▼            │           │         ▼            │
│  payment.pb.go       │           │  payment.pb.go       │
│  payment_grpc.pb.go  │           │  payment_grpc.pb.go  │
│         │            │           │         │            │
│         │            │           │         │            │
│         ▼            │           │         ▼            │
│  ┌──────────────┐   │           │  ┌──────────────┐   │
│  │   Server     │   │           │  │   Client     │   │
│  │ Implementa   │   │           │  │   Usa        │   │
│  │ PaymentSvc   │   │           │  │ PaymentSvc   │   │
│  └──────────────┘   │           │  └──────────────┘   │
└──────────────────────┘           └──────────────────────┘
         │                                   │
         │          gRPC Connection          │
         └───────────────────────────────────┘
```

## Comandos Disponíveis

### 📋 Via Makefile

```bash
cd proto/

# Sincronizar todos os protos
make sync-all

# Sincronizar apenas payment
make sync-payment

# Validar sincronização
make validate

# Gerar código Go
make generate-all

# Limpar arquivos gerados
make clean

# Ver todos os comandos
make help
```

### 🔧 Via Scripts

```bash
cd proto/

# Sincronização manual
./sync-protos.sh all        # Sincroniza tudo
./sync-protos.sh payment    # Sincroniza payment
./sync-protos.sh validate   # Valida sincronização

# Sincronização automática (watch mode)
./watch-protos.sh           # Observa mudanças e sincroniza
```

## Fluxo de Desenvolvimento

### Cenário 1: Adicionar novo campo no proto

```bash
# 1. Editar proto centralizado
vim proto/payment/payment.proto

# Adicionar campo:
# message ProcessPaymentRequest {
#   string order_id = 1;
#   double amount = 2;
#   string notes = 6;  // ← NOVO CAMPO
# }

# 2. Sincronizar para os serviços
make sync-all

# 3. Validar
make validate

# 4. Testar serviços
cd ../payments && go test ./...
cd ../orders && go test ./...

# 5. Commitar
git add .
git commit -m "feat(proto): add notes field to payment"
```

### Cenário 2: Adicionar novo método RPC

```bash
# 1. Editar proto
vim proto/payment/payment.proto

# Adicionar método:
# service PaymentService {
#   rpc ProcessPayment(...) returns (...);
#   rpc RefundPayment(...) returns (...);  // ← NOVO MÉTODO
# }

# 2. Sincronizar
make sync-all

# 3. Implementar no servidor (payments)
# Editar: payments/internal/infra/grpc/handler/payment_handler.go
# Implementar: func (s *PaymentServiceServer) RefundPayment(...)

# 4. Usar no cliente (orders)
# Editar: orders/internal/infra/grpc/client/payment_client.go
# Adicionar: func (c *PaymentClient) RefundPayment(...)

# 5. Testar e commitar
```

### Cenário 3: Detectar proto dessincronizado

```bash
# Alguém editou proto diretamente no serviço
vim payments/proto/payment.proto  # ❌ ERRADO!

# Validar detecta inconsistência
make validate
# ❌ payment.proto está dessincronizado no payments service!

# Corrigir: editar no local correto
vim proto/payment/payment.proto

# Sincronizar novamente
make sync-all
```

## Integração com CI/CD

### GitHub Actions

```yaml
# .github/workflows/proto-validation.yml
name: Validate Protos

on: [pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install protoc
        run: |
          sudo apt-get update
          sudo apt-get install -y protobuf-compiler
      
      - name: Install protoc plugins
        run: |
          go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
          go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
      
      - name: Validate proto sync
        run: |
          cd proto
          make validate
```

## Troubleshooting

### Problema: Protos dessincronizados

**Sintoma:** `make validate` falha

**Solução:**
```bash
cd proto
make sync-all
make validate
```

### Problema: Código Go não regenerado

**Sintoma:** Mudanças no proto não aparecem no código Go

**Solução:**
```bash
cd proto
make clean          # Remove .pb.go antigos
make sync-all       # Copia proto e regenera
```

### Problema: protoc-gen-go não encontrado

**Sintoma:** `protoc-gen-go: program not found`

**Solução:**
```bash
# Instalar plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Adicionar ao PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Persistir no shell (zsh)
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

## Boas Práticas

### ✅ DO

- Sempre edite protos na pasta `proto/`
- Use `make sync-all` após editar
- Valide com `make validate` antes de commitar
- Documente mudanças breaking
- Use versionamento semântico

### ❌ DON'T

- Nunca edite protos diretamente em `payments/proto/` ou `orders/proto/`
- Não commite sem validar sincronização
- Não faça mudanças breaking sem planejamento
- Não delete campos (deprecie com `deprecated = true`)

## Versionamento de Protos

### Semantic Versioning

```
v1.2.3
│ │ └── PATCH: bug fixes, documentation
│ └──── MINOR: new features (backward compatible)
└────── MAJOR: breaking changes
```

### Exemplo

```protobuf
syntax = "proto3";

package payment.v1;  // Major version no package

option go_package = "payments/proto;payment";

// v1.0.0: Initial release
// v1.1.0: Added RefundPayment method
// v1.2.0: Added notes field
// v2.0.0: Changed amount from double to int64 (BREAKING)

service PaymentService {
  rpc ProcessPayment(ProcessPaymentRequest) returns (ProcessPaymentResponse);
  
  // @since v1.1.0
  rpc RefundPayment(RefundPaymentRequest) returns (RefundPaymentResponse);
}

message ProcessPaymentRequest {
  string order_id = 1;
  double amount = 2;
  
  // @since v1.2.0
  string notes = 6;
}
```

## Recursos

- [Protocol Buffers Style Guide](https://protobuf.dev/programming-guides/style/)
- [gRPC Best Practices](https://grpc.io/docs/guides/performance/)
- [API Versioning](https://cloud.google.com/apis/design/versioning)
