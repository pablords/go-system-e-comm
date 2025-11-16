# Contribuindo para o Payment Service

Obrigado por considerar contribuir para o Payment Service! 🎉

## 📋 Índice

- [Como Contribuir](#como-contribuir)
- [Reportando Bugs](#reportando-bugs)
- [Sugerindo Melhorias](#sugerindo-melhorias)
- [Processo de Pull Request](#processo-de-pull-request)
- [Guia de Estilo](#guia-de-estilo)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Desenvolvimento](#desenvolvimento)

## Como Contribuir

Existem várias formas de contribuir:

1. **Reportar bugs** - Encontrou algo que não funciona? Nos avise!
2. **Sugerir melhorias** - Tem uma ideia para tornar o serviço melhor?
3. **Escrever código** - Implemente novas funcionalidades ou corrija bugs
4. **Melhorar documentação** - Ajude outros desenvolvedores
5. **Revisar código** - Participe de code reviews

## Reportando Bugs

Antes de reportar um bug, verifique se ele já não foi reportado. Se não encontrou, abra uma issue com:

### Template de Bug Report

```markdown
**Descrição do Bug**
Descrição clara e concisa do problema.

**Como Reproduzir**
Passos para reproduzir o comportamento:
1. Execute '...'
2. Chame o método '...'
3. Veja o erro

**Comportamento Esperado**
Descrição do que você esperava que acontecesse.

**Comportamento Atual**
O que realmente aconteceu.

**Logs/Screenshots**
Se aplicável, adicione logs ou screenshots.

**Ambiente**
- OS: [e.g. macOS 14.0]
- Go Version: [e.g. 1.21]
- MySQL Version: [e.g. 8.0]
- Docker Version: [e.g. 24.0]

**Contexto Adicional**
Qualquer outra informação relevante.
```

## Sugerindo Melhorias

### Template de Feature Request

```markdown
**Problema a Resolver**
Descrição clara do problema que esta feature resolveria.

**Solução Proposta**
Como você imagina que a feature deveria funcionar.

**Alternativas Consideradas**
Outras soluções que você pensou.

**Impacto**
Quem se beneficiaria desta feature e como.

**Contexto Adicional**
Qualquer informação adicional relevante.
```

## Processo de Pull Request

### Antes de Começar

1. Fork o repositório
2. Clone seu fork localmente
3. Crie uma branch a partir de `main`
4. Faça suas alterações
5. Teste suas alterações
6. Commit suas mudanças
7. Push para seu fork
8. Abra um Pull Request

### Passos Detalhados

```bash
# 1. Fork no GitHub, depois clone
git clone https://github.com/SEU_USUARIO/payments.git
cd payments

# 2. Adicione o repositório original como upstream
git remote add upstream https://github.com/ORIGINAL_OWNER/payments.git

# 3. Crie uma branch para sua feature
git checkout -b feature/nome-da-feature

# 4. Faça suas alterações e commit
git add .
git commit -m "feat: adiciona nova funcionalidade X"

# 5. Mantenha sua branch atualizada
git fetch upstream
git rebase upstream/main

# 6. Push para seu fork
git push origin feature/nome-da-feature

# 7. Abra um Pull Request no GitHub
```

### Critérios para Aprovação

- [ ] Código segue o guia de estilo
- [ ] Testes adicionados/atualizados
- [ ] Documentação atualizada
- [ ] Commits seguem convenção
- [ ] CI/CD passa
- [ ] Code review aprovado

## Guia de Estilo

### Go

Seguimos as convenções padrão de Go:

```bash
# Formatar código
make fmt

# Lint
golangci-lint run
```

**Boas Práticas:**
- Use `gofmt` para formatar código
- Siga [Effective Go](https://golang.org/doc/effective_go.html)
- Nomes de variáveis: camelCase
- Nomes de constantes: CamelCase ou UPPER_CASE
- Interfaces com -er suffix quando possível
- Comentários em funções exportadas
- Erros sempre retornados, nunca em panic

### Commits

Seguimos [Conventional Commits](https://www.conventionalcommits.org/):

```
<tipo>[escopo opcional]: <descrição>

[corpo opcional]

[rodapé opcional]
```

**Tipos:**
- `feat`: Nova funcionalidade
- `fix`: Correção de bug
- `docs`: Documentação
- `style`: Formatação (não afeta código)
- `refactor`: Refatoração
- `test`: Adiciona/modifica testes
- `chore`: Tarefas de manutenção

**Exemplos:**
```
feat(payment): adiciona suporte a PIX
fix(repository): corrige query de listagem
docs: atualiza README com novas instruções
refactor(usecase): simplifica lógica de validação
test(entity): adiciona testes para Payment.Cancel
```

### Código

**Estrutura de Arquivos:**
```go
package entity

import (
    // Standard library
    "errors"
    "time"
    
    // External packages
    "github.com/google/uuid"
    
    // Internal packages
    "payments/internal/domain/repository"
)

// Constants
const (
    StatusPending = "pending"
)

// Types
type Payment struct {
    ID string
}

// Constructors
func NewPayment() *Payment {
    return &Payment{}
}

// Methods (receiver alphabetically)
func (p *Payment) Cancel() error {
    return nil
}
```

**Tratamento de Erros:**
```go
// ✅ Bom
if err != nil {
    return fmt.Errorf("failed to process payment: %w", err)
}

// ❌ Evitar
if err != nil {
    panic(err)
}
```

**Logging:**
```go
// ✅ Bom
slog.Info("Processing payment", "order_id", orderID, "amount", amount)

// ❌ Evitar
fmt.Println("Processing payment for", orderID)
```

### Testes

```go
func TestPaymentProcess(t *testing.T) {
    // Arrange
    payment := NewPayment()
    
    // Act
    err := payment.Process("txn-123")
    
    // Assert
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if payment.Status != StatusProcessing {
        t.Errorf("Expected status %s, got %s", StatusProcessing, payment.Status)
    }
}
```

**Cobertura de Testes:**
- Domínio: > 90%
- Use Cases: > 80%
- Handlers: > 70%

### Documentação

**Comentários em Código:**
```go
// ProcessPayment processa um pagamento para um pedido específico.
// Retorna erro se o pagamento já foi processado ou se os dados são inválidos.
func (uc *ProcessPaymentUseCase) Execute(input ProcessPaymentInput) error {
    // ...
}
```

**README:**
- Use Markdown
- Inclua exemplos de código
- Mantenha atualizado
- Use emojis com moderação

## Estrutura do Projeto

```
payments/
├── cmd/              # Entry points
├── internal/         # Código privado
│   ├── domain/      # Lógica de negócio
│   ├── usecase/     # Casos de uso
│   └── infra/       # Implementações
├── proto/           # Definições gRPC
├── migrations/      # SQL migrations
├── tests/           # Testes
└── examples/        # Exemplos
```

**Regras de Dependência:**
- `domain` não depende de nada
- `usecase` depende apenas de `domain`
- `infra` implementa interfaces de `domain`
- `cmd` conecta tudo

## Desenvolvimento

### Setup Inicial

```bash
# Instalar dependências
go mod download

# Gerar código proto
make proto

# Iniciar banco de dados
docker-compose up -d payments-db

# Executar testes
make test
```

### Adicionando Nova Funcionalidade

1. **Domínio**: Adicione lógica de negócio em `internal/domain/entity/`
2. **Repository**: Se necessário, adicione método na interface
3. **Use Case**: Crie novo use case em `internal/usecase/`
4. **Handler**: Adicione método no gRPC handler
5. **Proto**: Atualize `proto/payment.proto` se necessário
6. **Migration**: Adicione migration se alterar banco
7. **Testes**: Adicione testes para nova funcionalidade
8. **Docs**: Atualize documentação relevante

### Executando Testes

```bash
# Todos os testes
make test

# Com coverage
make test-coverage

# Teste específico
go test -v ./internal/domain/entity/
```

### Debugging

```bash
# Logs do serviço
docker-compose logs -f payments-service

# Conectar ao banco
docker exec -it payments-mysql mysql -uroot -proot payments_db

# Testar gRPC
grpcurl -plaintext localhost:50051 list
```

## Revisão de Código

Ao revisar PRs, verifique:

### Funcionalidade
- [ ] Código faz o que deveria
- [ ] Não quebra funcionalidades existentes
- [ ] Testes passam

### Qualidade
- [ ] Código é legível e manutenível
- [ ] Segue padrões do projeto
- [ ] Não há duplicação desnecessária
- [ ] Erros são tratados adequadamente

### Segurança
- [ ] Não expõe dados sensíveis
- [ ] Valida entrada de usuário
- [ ] Não tem SQL injection
- [ ] Não loga informações sensíveis

### Performance
- [ ] Não tem problemas óbvios de performance
- [ ] Queries de banco são eficientes
- [ ] Não há memory leaks

### Documentação
- [ ] Código está comentado quando necessário
- [ ] README atualizado se necessário
- [ ] Changelog atualizado

## Perguntas?

Se tiver dúvidas sobre como contribuir:

1. Leia a documentação em `docs/`
2. Verifique issues abertas
3. Abra uma issue com sua pergunta
4. Entre em contato com os mantenedores

## Código de Conduta

Ao contribuir, você concorda em seguir nosso Código de Conduta:

- Seja respeitoso e inclusivo
- Aceite críticas construtivas
- Foque no que é melhor para a comunidade
- Mostre empatia com outros membros

## Licença

Ao contribuir, você concorda que suas contribuições serão licenciadas sob a mesma licença do projeto (MIT).

---

**Obrigado por contribuir! 🚀**
