# 📋 Guia Rápido - Gerenciamento de Protos

## 🚀 Comandos Essenciais

```bash
cd proto/

# Sincronizar protos para todos os serviços
make sync-all

# Validar se está sincronizado
make validate

# Ver todos os comandos disponíveis
make help
```

## 📖 Estrutura

```
proto/                          ← Source of Truth (edite AQUI)
├── payment/
│   └── payment.proto          ← Contrato centralizado
├── sync-protos.sh             ← Script de sincronização
├── watch-protos.sh            ← Auto-sync (watch mode)
└── Makefile                   ← Comandos make

payments/proto/payment.proto   ← Cópia sincronizada (NÃO edite)
orders/proto/payment.proto     ← Cópia sincronizada (NÃO edite)
```

## ✏️ Workflow: Editar Proto

```bash
# 1. Editar proto centralizado
vim proto/payment/payment.proto

# 2. Sincronizar para serviços
cd proto && make sync-all

# 3. Validar
make validate

# 4. Testar
cd ../payments && go test ./...
cd ../orders && go test ./...

# 5. Commitar
git add . && git commit -m "feat: update payment proto"
```

## 🔍 Comandos de Diagnóstico

```bash
# Verificar se protos estão sincronizados
cd proto && make validate

# Limpar arquivos gerados
make clean

# Regenerar tudo do zero
make clean && make sync-all
```

## 🎯 Regras de Ouro

1. ✅ **Sempre edite em `proto/`** (nunca direto nos serviços)
2. ✅ **Rode `make sync-all` após editar**
3. ✅ **Valide com `make validate` antes de commitar**
4. ✅ **Teste ambos os serviços após mudanças**

## 🆘 Resolução de Problemas

### protoc-gen-go não encontrado

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Protos dessincronizados

```bash
cd proto
make sync-all
make validate
```

## 📚 Documentação Completa

- `README.md` - Documentação detalhada
- `WORKFLOW.md` - Fluxos e diagramas
