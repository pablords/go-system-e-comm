#!/bin/bash

# Script de watch para sincronizar automaticamente quando os protos mudarem
# Uso: ./watch-protos.sh

set -e

# Cores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Verificar se fswatch está instalado
if ! command -v fswatch &> /dev/null; then
    warning "fswatch não está instalado!"
    info "Para instalar no macOS: brew install fswatch"
    info "Para instalar no Linux: apt-get install fswatch (ou yum install fswatch)"
    exit 1
fi

echo ""
echo "👀 Observando mudanças nos protos..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
info "Pressione Ctrl+C para parar"
echo ""

# Observar mudanças nos arquivos .proto
fswatch -o "$SCRIPT_DIR/payment/" | while read -r change; do
    success "Detectada mudança nos protos!"
    "$SCRIPT_DIR/sync-protos.sh" all
    echo ""
    info "Aguardando próximas mudanças..."
    echo ""
done
