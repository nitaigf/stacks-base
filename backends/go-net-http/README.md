# Go net/http Backend

Backend de referencia inicial do Stacks Base.

## Endpoints atuais

- GET /health
- POST /api/v1/auth/register
- POST /api/v1/auth/login
- POST /api/v1/auth/forgot-password
- POST /api/v1/auth/reset-password
- POST /api/v1/auth/logout
- POST /api/v1/auth/change-password
- GET /api/v1/users/me
- GET /api/v1/users
- POST /api/v1/users
- GET /api/v1/users/{userId}
- PATCH /api/v1/users/{userId}
- POST /api/v1/users/{userId}/deactivate
- POST /api/v1/users/{userId}/reactivate
- POST /api/v1/users/{userId}/soft-delete
- POST /api/v1/users/{userId}/restore
- DELETE /api/v1/users/{userId}
- GET /api/v1/users/export.csv
- GET /api/v1/users/export.xlsx
- GET /api/v1/users/print
- GET /api/v1/audit-logs

## Requisitos locais

- Go 1.22+
- PostgreSQL em localhost:5432
- SMTP configurado via `.env` para os fluxos reais de e-mail

## Execucao

1. Exporte as variaveis de ambiente da raiz do repositorio.
2. Crie o banco stacks_base localmente.
3. Execute go run ./cmd/server.

Com `AUTO_MIGRATE=true` o backend aplica as migrations ao iniciar.
Com `ADMIN_SEED_ENABLED=true` e `DEMO_SEED_ENABLED=true` o backend popula admin, usuarios demo e auditoria inicial.
