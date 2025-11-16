.PHONY: help generate-all generate-payment sync-all sync-payment sync-orders clean

# Diretórios dos serviços
PAYMENTS_DIR := ../payments
ORDERS_DIR := ../orders

# Diretórios de destino dos protos
PAYMENTS_PROTO_DIR := $(PAYMENTS_DIR)/proto
ORDERS_PROTO_DIR := $(ORDERS_DIR)/proto

help:
	@echo "📋 Comandos disponíveis:"
	@echo ""
	@echo "  make generate-all     - Gera código Go para todos os protos"
	@echo "  make generate-payment - Gera código Go apenas para payment"
	@echo ""
	@echo "  make sync-all         - Sincroniza protos para todos os serviços"
	@echo "  make sync-payment     - Sincroniza payment.proto para payments e orders"
	@echo ""
	@echo "  make clean            - Remove arquivos gerados (.pb.go)"
	@echo ""

# Gerar código Go a partir dos protos
generate-all: generate-payment

generate-payment:
	@echo "🔨 Gerando código Go para payment.proto..."
	@protoc --go_out=$(PAYMENTS_PROTO_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PAYMENTS_PROTO_DIR) --go-grpc_opt=paths=source_relative \
		payment/payment.proto
	@echo "✅ Código gerado em $(PAYMENTS_PROTO_DIR)"

# Sincronizar protos para os serviços
sync-all: sync-payment

sync-payment:
	@echo "📦 Sincronizando payment.proto..."
	@mkdir -p $(PAYMENTS_PROTO_DIR)
	@mkdir -p $(ORDERS_PROTO_DIR)
	@cp -v payment/payment.proto $(PAYMENTS_PROTO_DIR)/
	@cp -v payment/payment.proto $(ORDERS_PROTO_DIR)/
	@echo "✅ payment.proto sincronizado para payments e orders"
	@echo ""
	@echo "🔨 Gerando código nos serviços..."
	@cd $(PAYMENTS_DIR) && make proto || true
	@cd $(ORDERS_DIR) && make proto || true
	@echo "✅ Sincronização completa!"

# Limpar arquivos gerados
clean:
	@echo "🧹 Limpando arquivos gerados..."
	@find $(PAYMENTS_PROTO_DIR) -name "*.pb.go" -delete 2>/dev/null || true
	@find $(ORDERS_PROTO_DIR) -name "*.pb.go" -delete 2>/dev/null || true
	@echo "✅ Arquivos .pb.go removidos"

# Validar que os protos estão sincronizados
validate:
	@echo "🔍 Validando sincronização dos protos..."
	@if ! diff -q payment/payment.proto $(PAYMENTS_PROTO_DIR)/payment.proto > /dev/null 2>&1; then \
		echo "❌ payment.proto está dessincronizado no payments service!"; \
		exit 1; \
	fi
	@if ! diff -q payment/payment.proto $(ORDERS_PROTO_DIR)/payment.proto > /dev/null 2>&1; then \
		echo "❌ payment.proto está dessincronizado no orders service!"; \
		exit 1; \
	fi
	@echo "✅ Todos os protos estão sincronizados"
