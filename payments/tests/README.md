# Tests

Testes unitários para o serviço de pagamentos.

## Executar todos os testes

```bash
make test
```

## Executar testes com coverage

```bash
make test-coverage
```

## Estrutura dos testes

- `internal/domain/entity/` - Testes das entidades de domínio
- `internal/usecase/` - Testes dos casos de uso
- `mocks/` - Mocks para testes

## Adicionar novos testes

Siga o padrão de nomenclatura:
- Arquivo de teste: `*_test.go`
- Função de teste: `Test<FunctionName>`
