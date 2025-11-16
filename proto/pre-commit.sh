#!/bin/bash

# Git pre-commit hook para validar sincronização dos protos
# Para instalar: cp pre-commit.sh .git/hooks/pre-commit && chmod +x .git/hooks/pre-commit

set -e

PROTO_DIR="$(git rev-parse --show-toplevel)/proto"

# Verificar se estamos commitando mudanças em protos
if git diff --cached --name-only | grep -q "^proto/"; then
    echo "🔍 Detectadas mudanças nos protos, validando sincronização..."
    
    cd "$PROTO_DIR"
    
    if ./sync-protos.sh validate > /dev/null 2>&1; then
        echo "✅ Protos sincronizados!"
    else
        echo ""
        echo "❌ ERRO: Protos não estão sincronizados!"
        echo ""
        echo "Execute: cd proto && make sync-all"
        echo ""
        exit 1
    fi
fi

exit 0
