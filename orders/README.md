# Orders API

API robusta de gerenciamento de pedidos desenvolvida em Go.

## Arquitetura

- **Clean Architecture** com separação clara de responsabilidades
- **Domain-Driven Design** para modelagem de negócio
- **Repository Pattern** para abstração de dados
- **Use Cases** para lógica de negócio

## Tecnologias

- Go 1.21+
- MySQL 8.0
- Chi Router
- Docker & Docker Compose
- Air (Live Reload)
- Testes Unitários com cobertura de 98.5% nas entidades e 66.3% nos use cases

## Estrutura do Projeto

```
orders-go/
├── cmd/api/              # Ponto de entrada da aplicação
├── internal/
│   ├── domain/           # Camada de domínio
│   │   ├── entity/       # Entidades de negócio
│   │   └── repository/   # Interfaces de repositório
│   ├── usecase/          # Casos de uso
│   └── infra/            # Infraestrutura
│       ├── database/     # Configuração do banco
│       ├── http/         # Handlers HTTP
│       └── repository/   # Implementações de repositório
├── migrations/           # Scripts SQL
├── docker-compose.yml
├── Dockerfile
└── .air.toml            # Configuração live reload
```

## Funcionalidades

### Produtos
- ✅ Criar, listar, atualizar e deletar produtos
- ✅ Controle de estoque
- ✅ Validações de negócio

### Pedidos (Carrinho)
- ✅ Criar carrinho
- ✅ Adicionar itens ao carrinho
- ✅ Remover itens do carrinho
- ✅ Atualizar quantidade de itens
- ✅ Calcular total para pagamento
- ✅ Gerenciar status do pedido

### Items
- ✅ Gestão automática de items no carrinho
- ✅ Cálculo automático de totais

## Como Executar

### Com Docker (Recomendado)

```bash
# Iniciar todos os serviços (MySQL + API com live reload)
docker-compose up --build

# A API estará disponível em http://localhost:8080
```

### Localmente

```bash
# Instalar dependências
go mod download

# Configurar variáveis de ambiente
cp .env .env.local

# Executar com live reload (recomendado para desenvolvimento)
make dev

# Ou executar diretamente
make run
# ou
go run cmd/api/main.go
```

## Endpoints da API

### Documentação Swagger

A API possui documentação interativa via Swagger UI:

```
http://localhost:8080/swagger/index.html
```

A documentação Swagger inclui:
- Todos os endpoints disponíveis
- Schemas de request/response
- Exemplos de requisições
- Testes interativos

### Health Check
```
GET /health
```

### Produtos
```
GET    /api/v1/products          # Listar produtos
POST   /api/v1/products          # Criar produto
GET    /api/v1/products/:id      # Obter produto
PUT    /api/v1/products/:id      # Atualizar produto
DELETE /api/v1/products/:id      # Deletar produto
```

### Pedidos
```
GET    /api/v1/orders            # Listar pedidos
GET    /api/v1/orders/:id        # Obter pedido
DELETE /api/v1/orders/:id        # Deletar pedido
```

### Carrinho
```
POST   /api/v1/cart                          # Criar carrinho
GET    /api/v1/cart/:id                      # Obter carrinho
POST   /api/v1/cart/:id/items                # Adicionar item
DELETE /api/v1/cart/:id/items/:itemId        # Remover item
PUT    /api/v1/cart/:id/items/:itemId        # Atualizar quantidade
GET    /api/v1/cart/:id/calculate            # Calcular total
PUT    /api/v1/cart/:id/status               # Atualizar status
```

## Exemplos de Uso

### Criar Produto
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Notebook",
    "description": "Notebook Dell Inspiron",
    "price": 3500.00,
    "stock": 10
  }'
```

### Criar Carrinho
```bash
curl -X POST http://localhost:8080/api/v1/cart
```

### Adicionar Item ao Carrinho
```bash
curl -X POST http://localhost:8080/api/v1/cart/{order_id}/items \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "{product_id}",
    "quantity": 2
  }'
```

### Calcular Total
```bash
curl http://localhost:8080/api/v1/cart/{order_id}/calculate
```

## Variáveis de Ambiente

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=orders_user
DB_PASSWORD=orders_pass
DB_NAME=orders_db
SERVER_PORT=8080
```

## Live Reload

O projeto está configurado com Air para live reload durante o desenvolvimento. Qualquer alteração nos arquivos `.go` irá recompilar e reiniciar automaticamente a aplicação.

## Banco de Dados

O schema do banco é criado automaticamente através do script de migration em `migrations/001_create_tables.sql` quando o container do MySQL é iniciado pela primeira vez.

## Testes

O projeto possui testes unitários completos para as entidades e use cases.

### Executar todos os testes
```bash
make test
# ou
go test -v ./...
```

### Executar testes com cobertura
```bash
make test-coverage
# ou
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Executar apenas testes unitários
```bash
make test-unit
# ou
go test -v ./internal/domain/entity/... ./internal/usecase/...
```

### Ver cobertura por pacote
```bash
make coverage-by-package
```

### Cobertura Atual
- **Entidades**: 98.5% de cobertura
- **Use Cases**: 66.3% de cobertura
- **Total**: ~25.7% (incluindo infraestrutura)

## Comandos Make Disponíveis

- `make dev` - Executar aplicação com live reload (desenvolvimento)
- `make run` - Executar aplicação localmente
- `make build` - Compilar aplicação
- `make test` - Executar todos os testes
- `make test-coverage` - Gerar relatório de cobertura
- `make test-unit` - Executar apenas testes unitários
- `make swagger` - Gerar documentação Swagger
- `make docker-up` - Iniciar containers Docker
- `make docker-down` - Parar containers Docker
- `make clean` - Limpar arquivos gerados
- `make fmt` - Formatar código
- `make deps` - Instalar dependências

## Licença

MIT
