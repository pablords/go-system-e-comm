# Geração de código Go a partir de arquivos .proto

## Instalação das ferramentas necessárias

```bash
# Instalar protoc compiler
brew install protobuf

# Instalar plugins Go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Gerar código

```bash
make proto
```

Os arquivos gerados serão:
- `payment.pb.go` - contém as definições de mensagens
- `payment_grpc.pb.go` - contém as definições do serviço gRPC
