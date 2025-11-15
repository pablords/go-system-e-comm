#!/bin/bash

# Quick Start Script for Payment Service
# Este script automatiza o setup inicial do serviço

set -e  # Exit on error

echo "🚀 Payment Service - Quick Start"
echo "================================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go não está instalado. Por favor, instale Go 1.21 ou superior."
    exit 1
fi

echo "✅ Go está instalado: $(go version)"
echo ""

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "⚠️  protoc não está instalado."
    echo "📦 Instalando protoc..."
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        brew install protobuf
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        sudo apt-get update
        sudo apt-get install -y protobuf-compiler
    else
        echo "❌ Por favor, instale manualmente o protoc: https://grpc.io/docs/protoc-installation/"
        exit 1
    fi
fi

echo "✅ protoc está instalado: $(protoc --version)"
echo ""

# Install Go plugins for protoc
echo "📦 Instalando plugins Go para protoc..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Add GOPATH/bin to PATH if not already there
export PATH="$PATH:$(go env GOPATH)/bin"

echo "✅ Plugins instalados"
echo ""

# Download dependencies
echo "📦 Baixando dependências..."
go mod download
go mod tidy

echo "✅ Dependências baixadas"
echo ""

# Generate proto files
echo "🔧 Gerando código gRPC..."
make proto

echo "✅ Código gRPC gerado"
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker não está instalado."
    echo "   Por favor, instale Docker para executar o banco de dados."
    echo "   Você pode continuar sem Docker se já tiver MySQL instalado."
    echo ""
    read -p "Continuar sem Docker? (y/n) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
else
    echo "✅ Docker está instalado: $(docker --version)"
    echo ""
    
    # Start MySQL with Docker Compose
    echo "🐳 Iniciando MySQL com Docker Compose..."
    docker-compose up -d payments-db
    
    echo "⏳ Aguardando MySQL inicializar (30 segundos)..."
    sleep 30
    
    echo "✅ MySQL iniciado"
    echo ""
fi

# Create .env if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Criando arquivo .env..."
    cat > .env << EOF
DB_HOST=localhost
DB_PORT=3307
DB_USER=root
DB_PASSWORD=root
DB_NAME=payments_db
GRPC_PORT=50051
EOF
    echo "✅ Arquivo .env criado"
    echo ""
fi

echo ""
echo "🎉 Setup concluído com sucesso!"
echo ""
echo "Próximos passos:"
echo "  1. Execute o serviço:        make run"
echo "  2. Em outro terminal, teste: grpcurl -plaintext localhost:50051 list"
echo ""
echo "Documentação:"
echo "  - README.md       - Documentação principal"
echo "  - SETUP.md        - Guia de configuração"
echo "  - EXAMPLES.md     - Exemplos de requisições"
echo "  - INTEGRATION.md  - Como integrar com Orders"
echo ""
echo "Comandos úteis:"
echo "  make test         - Executar testes"
echo "  make proto        - Regenerar código gRPC"
echo "  make docker-run   - Executar tudo com Docker"
echo ""

# Ask if user wants to start the service
read -p "Deseja iniciar o serviço agora? (y/n) " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🚀 Iniciando Payment Service..."
    echo ""
    make run
fi
