# Go net/http Backend

Backend de referencia inicial do Stacks Base.

## Endpoints iniciais

- GET /health
- POST /api/v1/auth/register
- POST /api/v1/auth/login
- POST /api/v1/auth/logout
- GET /api/v1/users/me

## Requisitos locais

- Go 1.22+
- PostgreSQL em localhost:5432

## Execucao

1. Exporte as variaveis de ambiente de shared/.env.example.
2. Crie o banco stacks_base localmente.
3. Execute go run ./cmd/server.

Com AUTO_MIGRATE=true o backend aplica a migration inicial ao iniciar.