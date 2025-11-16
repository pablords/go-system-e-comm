#!/bin/bash

# Script para sincronizar protos do repositório centralizado para os serviços
# Uso: ./sync-protos.sh [payment|order|all]

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Diretórios
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PAYMENTS_DIR="$SCRIPT_DIR/../payments"
ORDERS_DIR="$SCRIPT_DIR/../orders"

# Função para imprimir mensagens coloridas
info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Função para sincronizar payment.proto
sync_payment() {
    info "Sincronizando payment.proto..."
    
    # Criar diretórios se não existirem
    mkdir -p "$PAYMENTS_DIR/proto"
    mkdir -p "$ORDERS_DIR/proto"
    
    # Copiar proto
    cp "$SCRIPT_DIR/payment/payment.proto" "$PAYMENTS_DIR/proto/"
    cp "$SCRIPT_DIR/payment/payment.proto" "$ORDERS_DIR/proto/"
    
    success "payment.proto copiado para payments e orders"
    
    # Gerar código Go
    info "Gerando código Go no payments service..."
    if cd "$PAYMENTS_DIR" && make proto; then
        success "Código gerado no payments service"
    else
        warning "Erro ao gerar código no payments service (verifique se 'make proto' existe)"
    fi
    
    info "Gerando código Go no orders service..."
    if cd "$ORDERS_DIR" && make proto; then
        success "Código gerado no orders service"
    else
        warning "Erro ao gerar código no orders service (verifique se 'make proto' existe)"
    fi
}

# Função para validar sincronização
validate_sync() {
    info "Validando sincronização..."
    
    local errors=0
    
    # Verificar payments
    if ! diff -q "$SCRIPT_DIR/payment/payment.proto" "$PAYMENTS_DIR/proto/payment.proto" > /dev/null 2>&1; then
        error "payment.proto está dessincronizado no payments service!"
        errors=$((errors + 1))
    else
        success "payment.proto sincronizado no payments service"
    fi
    
    # Verificar orders
    if ! diff -q "$SCRIPT_DIR/payment/payment.proto" "$ORDERS_DIR/proto/payment.proto" > /dev/null 2>&1; then
        error "payment.proto está dessincronizado no orders service!"
        errors=$((errors + 1))
    else
        success "payment.proto sincronizado no orders service"
    fi
    
    if [ $errors -eq 0 ]; then
        success "Todos os protos estão sincronizados!"
        return 0
    else
        error "Encontrados $errors erro(s) de sincronização"
        return 1
    fi
}

# Função principal
main() {
    local command="${1:-all}"
    
    echo ""
    echo "🔄 Sincronizador de Protos"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    case "$command" in
        payment)
            sync_payment
            ;;
        all)
            sync_payment
            ;;
        validate)
            validate_sync
            ;;
        *)
            error "Comando desconhecido: $command"
            echo ""
            echo "Uso: $0 [payment|all|validate]"
            echo ""
            echo "Comandos:"
            echo "  payment   - Sincroniza payment.proto"
            echo "  all       - Sincroniza todos os protos (padrão)"
            echo "  validate  - Valida se os protos estão sincronizados"
            echo ""
            exit 1
            ;;
    esac
    
    echo ""
    success "Sincronização completa!"
    echo ""
}

# Executar
main "$@"
