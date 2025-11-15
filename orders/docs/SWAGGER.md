# Swagger Documentation

## Acessar a Documentação

Após iniciar a aplicação, acesse:

```
http://localhost:8080/swagger/index.html
```

## Regenerar Documentação

Sempre que modificar as anotações nos handlers, execute:

```bash
make swagger
# ou
swag init -g cmd/api/main.go -o docs
```

## Estrutura da Documentação

A documentação Swagger é gerada automaticamente a partir das anotações nos arquivos:

- `cmd/api/main.go` - Informações gerais da API
- `internal/infra/http/handler/*.go` - Documentação dos endpoints

## Anotações Disponíveis

### Informações Gerais (main.go)
```go
// @title Orders API
// @version 1.0
// @description API de gerenciamento de pedidos
// @host localhost:8080
// @BasePath /api/v1
```

### Endpoints (handlers)
```go
// @Summary Resumo do endpoint
// @Description Descrição detalhada
// @Tags categoria
// @Accept json
// @Produce json
// @Param id path string true "ID do recurso"
// @Param body body RequestType true "Corpo da requisição"
// @Success 200 {object} ResponseType
// @Failure 400 {object} ErrorResponse
// @Router /endpoint [method]
```

## Testando a API via Swagger

1. Acesse http://localhost:8080/swagger/index.html
2. Expanda o endpoint desejado
3. Clique em "Try it out"
4. Preencha os parâmetros
5. Clique em "Execute"
6. Veja a resposta abaixo

## Schemas Documentados

- **Product** - Entidade de produto
- **Order** - Entidade de pedido
- **Item** - Item do pedido
- **CreateProductRequest** - Request para criar produto
- **UpdateProductRequest** - Request para atualizar produto
- **AddItemRequest** - Request para adicionar item ao carrinho
- **UpdateItemRequest** - Request para atualizar quantidade
- **UpdateOrderStatusRequest** - Request para atualizar status do pedido
- **ErrorResponse** - Response de erro padrão
