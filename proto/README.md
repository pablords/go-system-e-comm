# Proto Repository

Este diretório centraliza todos os arquivos `.proto` (contratos gRPC) do sistema e-commerce.

## 📁 Estrutura

```
proto/
├── payment/
│   ├── payment.proto       # Contrato do serviço de pagamentos
│   └── README.md
├── Makefile                # Comandos make para sincronização
├── sync-protos.sh          # Script de sincronização manual
├── watch-protos.sh         # Script para auto-sincronização
└── README.md
```

## 🚀 Como Usar

### 1. Sincronização Manual

```bash
# Sincronizar todos os protos para os serviços
make sync-all

# Sincronizar apenas payment.proto
make sync-payment

# Ou usar o script diretamente
./sync-protos.sh all
./sync-protos.sh payment
```

### 2. Sincronização Automática (Watch Mode)

```bash
# Observar mudanças e sincronizar automaticamente
./watch-protos.sh
```

**Nota:** Requer `fswatch` instalado:
```bash
# macOS
brew install fswatch

# Linux
apt-get install fswatch
```

### 3. Validar Sincronização

```bash
# Verificar se os protos estão sincronizados entre serviços
make validate

# Ou
./sync-protos.sh validate
```

### 4. Gerar Código Go

```bash
# Gerar código Go para todos os protos
make generate-all

# Gerar apenas para payment
make generate-payment
```

### 5. Limpar Arquivos Gerados

```bash
# Remove todos os arquivos .pb.go
make clean
```

## 📋 Workflow Recomendado

### Para adicionar/modificar um proto:

1. **Editar o proto aqui** (pasta `proto/`)
2. **Sincronizar** para os serviços: `make sync-all`
3. **Validar** a sincronização: `make validate`
4. **Testar** os serviços afetados
5. **Commitar** as mudanças

### Exemplo:

```bash
# 1. Editar proto/payment/payment.proto
vim proto/payment/payment.proto

# 2. Sincronizar
make sync-payment

# 3. Validar
make validate

# 4. Testar
cd ../payments && go test ./...
cd ../orders && go test ./...

# 5. Commitar
git add .
git commit -m "feat: add new field to payment proto"
```

## 🔄 Como Funciona a Sincronização

```
proto/payment/payment.proto (SOURCE OF TRUTH)
            │
            ├─── copia para ──→ payments/proto/payment.proto
            │                         │
            │                         └─→ gera payment.pb.go
            │                             gera payment_grpc.pb.go
            │
            └─── copia para ──→ orders/proto/payment.proto
                                      │
                                      └─→ gera payment.pb.go
                                          gera payment_grpc.pb.go
```

## 📝 Versionamento

### Mudanças Breaking (Major)
- Remover campos
- Renomear campos
- Mudar tipos de campos
- Remover métodos RPC

### Mudanças Compatíveis (Minor)
- Adicionar novos campos (com defaults)
- Adicionar novos métodos RPC
- Adicionar novos enums

### Exemplo de Versionamento

```protobuf
syntax = "proto3";

package payment.v1;  // ← Versão no package
option go_package = "payments/proto;payment";

service PaymentService {
  rpc ProcessPayment(ProcessPaymentRequest) returns (ProcessPaymentResponse);
  
  // v1.1.0: Novo método adicionado
  rpc RefundPayment(RefundPaymentRequest) returns (RefundPaymentResponse);
}
```

## 🛠️ Troubleshooting

### Erro: "protoc-gen-go: program not found"

```bash
# Instalar plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Adicionar ao PATH
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Erro: "make: *** No rule to make target 'proto'"

Os serviços precisam ter um Makefile com target `proto`. Exemplo:

```makefile
# payments/Makefile ou orders/Makefile
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/payment.proto
```

### Protos dessincronizados

```bash
# Forçar sincronização
make sync-all

# Validar
make validate
```

## 📚 Referências

- [Protocol Buffers - Google](https://protobuf.dev/)
- [gRPC Go Tutorial](https://grpc.io/docs/languages/go/quickstart/)
- [Buf Schema Registry](https://buf.build/)
- [API Versioning Best Practices](https://cloud.google.com/apis/design/versioning)

## 🤝 Contribuindo

1. Sempre edite os protos **nesta pasta** (`proto/`), nunca diretamente nos serviços
2. Use `make validate` antes de commitar
3. Documente mudanças breaking no changelog
4. Incremente a versão no proto quando aplicável
