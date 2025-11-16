# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
e este projeto adere ao [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planejado
- Integração com gateway de pagamento real (Stripe)
- Suporte a webhooks
- Processamento assíncrono de pagamentos
- Implementação de event sourcing
- Adição de métricas Prometheus
- Distributed tracing com Jaeger
- Suporte a refund
- Cache com Redis
- Rate limiting

## [1.0.0] - 2025-11-15

### Adicionado
- Serviço gRPC completo para processamento de pagamentos
- Suporte a múltiplos métodos de pagamento:
  - Cartão de Crédito
  - Cartão de Débito
  - PIX
  - Boleto
  - PayPal
- Entidade Payment com regras de negócio
- Estados de pagamento: Pending, Processing, Approved, Declined, Canceled, Refunded
- Repository pattern com implementação MySQL
- Use cases para operações de pagamento:
  - ProcessPayment
  - GetPayment
  - CancelPayment
  - ListPayments
- Migrations de banco de dados
- Logging estruturado com slog
- Dockerização completa com Docker Compose
- Testes unitários para entidades
- Documentação completa:
  - README.md
  - SETUP.md
  - ARCHITECTURE.md
  - INTEGRATION.md
  - EXAMPLES.md
  - DIAGRAMS.md
  - SUMMARY.md
  - PRODUCTION-CHECKLIST.md
- Exemplos de integração com serviço Orders
- Cliente gRPC de exemplo
- Script de quick start
- Makefile com comandos úteis

### Características Técnicas
- Go 1.21
- gRPC / Protocol Buffers
- MySQL 8.0
- Clean Architecture
- Domain-Driven Design (DDD)
- Connection pooling para banco de dados
- Graceful shutdown
- Structured logging

### Decisões de Arquitetura
- Separação clara de camadas (Domain, Use Case, Infrastructure)
- Dependency Inversion Principle aplicado
- Repository pattern para abstração de dados
- Simulação de gateway de pagamento (95% approval rate)
- Validações no domínio
- Status imutáveis após finalização

---

## Formato de Versão

### Major (X.0.0)
- Mudanças incompatíveis com versões anteriores
- Alterações significativas na arquitetura
- Remoção de funcionalidades

### Minor (0.X.0)
- Novas funcionalidades compatíveis com versões anteriores
- Melhorias significativas
- Novas integrações

### Patch (0.0.X)
- Correções de bugs
- Melhorias de performance
- Atualizações de documentação
- Refatorações internas

## Categorias de Mudanças

- **Adicionado** - Novas funcionalidades
- **Alterado** - Mudanças em funcionalidades existentes
- **Descontinuado** - Funcionalidades que serão removidas
- **Removido** - Funcionalidades removidas
- **Corrigido** - Correções de bugs
- **Segurança** - Correções de vulnerabilidades

---

[Unreleased]: https://github.com/seu-usuario/payments/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/seu-usuario/payments/releases/tag/v1.0.0
