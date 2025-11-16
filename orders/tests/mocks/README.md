# Test Mocks

Este diretório contém mocks reutilizáveis para os testes do projeto.

## Logger Mock

O `MockLogger` é uma implementação de mock do `slog.Logger` para uso em testes unitários.

### Uso

```go
import (
    "orders/tests/mocks"
    "testing"
)

func TestYourFunction(t *testing.T) {
    // Criar um mock logger
    logger := mocks.NewMockLogger()
    
    // Usar o logger nos seus testes
    yourService := NewYourService(logger)
    
    // O logger não vai imprimir nada no console durante os testes
    // mas vai capturar todas as mensagens de log internamente
}
```

### Características

- **Silencioso**: Não imprime nada no console durante os testes
- **Completo**: Implementa todos os métodos necessários do `slog.Handler`
- **Leve**: Não adiciona overhead significativo aos testes
- **Compatível**: Funciona com qualquer código que use `*slog.Logger`

### Exemplo com Use Case

```go
func TestProductUseCase_CreateProduct(t *testing.T) {
    repo := newMockProductRepository()
    logger := mocks.NewMockLogger()
    uc := usecase.NewProductUseCase(repo, logger)

    product, err := uc.CreateProduct("Laptop", "Dell", 1500.00, 10)
    // ... assertions
}
```

### Extensões Futuras

Se você precisar verificar os logs durante os testes, pode estender o mock para capturar e expor as mensagens de log através dos métodos `GetLogs()` e `Clear()` já implementados.
