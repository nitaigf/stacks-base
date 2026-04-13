# Stacks Base

## Estado Atual

A baseline atual e real, funcional e sem mocks:

- frontend de referencia em `frontends/solidjs/` com SolidJS, Vite, TanStack Router, Ky e Zod
- backend de referencia em `backends/go-net-http/` com Go `net/http`
- banco principal de desenvolvimento em PostgreSQL local `127.0.0.1:5432`
- design system compartilhado em `shared/design-system/`
- schema SQL canonico em `shared/schema.sql`
- contrato OpenAPI em `specs/openapi.yaml`
- collection Bruno em `specs/bruno/`

## O Que a Baseline Entrega Hoje

O sistema atual cobre o escopo prometido para a baseline:

- pagina publica
- paginas de autenticacao
- area privada
- area administrativa
- paginas de erro 403, 404 e 500
- `health`
- `register`, `login`, `logout`, `users/me`
- `forgot-password`, `reset-password`, `change-password`
- gestao real de usuarios com dados do banco:
  - listar paginado
  - visualizar um
  - criar
  - editar
  - apagar com `soft-delete`
  - restaurar
  - inativar
  - reativar
  - excluir definitivamente
  - exportar CSV
  - exportar XLSX
  - gerar PDF imprimivel
- auditoria real com:
  - quando
  - quem
  - rota
  - metodo
  - IP
  - user agent
  - metadata complementar
- seeds de admin e dados demonstrativos reais

## Validacao Tecnica Atual

- **Estado**: `validated-local` (Consolidação concluída)
- **Cobertura**: Auth completo, CRUD Usuarios, Auditoria real, E2E Playwright exaustivo.
- **Backend**: Go (net/http) + PostgreSQL.
- **Frontend**: SolidJS + TanStack Router (Atomic Design).
- **Specs**: OpenAPI Hand-written + Collection Bruno (Happy & Negative paths).

Controle via chamadas HTTP reais contra PostgreSQL local validando:
  - login admin
  - CRUD administrativo de usuarios
  - soft-delete e restore
  - exportacoes CSV, XLSX e PDF
  - change password
  - forgot password
  - reset password
  - logout
  - listagem de auditoria
  - persistencia dos eventos de auditoria no banco

## Principios

- Paridade funcional maxima entre stacks.
- Mesmo contrato de API para todos os backends.
- Mesmo design system para todos os frontends.
- Sem dependencia compartilhada em runtime entre stacks.
- Sem mocks na baseline de referencia.
- Escritas com previsibilidade transacional e erro coerente.
- Documentacao, OpenAPI e especificacoes atualizadas junto com o codigo.

## Fontes de Verdade

- `README.md`: estado e proposta do projeto
- `CONSTITUTION.md`: regras inegociaveis
- `ARCHITECTURE.md`: layout e fronteiras do monorepo
- `AGENTS.md`: papeis e guard rails operacionais
- `ADR.md`: decisoes estruturais aceitas
- `TODO.md`: estado validado por fase
- `shared/schema.sql`: schema SQL canonico
- `shared/design-system/index.css`: design system compartilhado
- `specs/openapi.yaml`: contrato de API canonico
- `specs/bruno/`: validacao de contrato e colecao executavel

## Estrutura Canonica

```text
stacks-base/
  shared/
    .env.example
    schema.sql
    docker-compose.yml
    design-system/
  specs/
    openapi.yaml
    bruno/
  frontends/
    solidjs/
  backends/
    go-net-http/
  docs/
```

## Desenvolvimento Local

- PostgreSQL local e o padrao de desenvolvimento.
- Docker Compose e opcional e serve apenas como alternativa de apoio.
- O package manager oficial do frontend atual e `npm`, com `package-lock.json` versionado.

## Roadmap Imediato

1. [x] **Consolidação da Baseline**: Fechar débitos de segurança e ampliar E2E (Playwright).
2. [ ] **Multi-Stack**: Iniciar a próxima stack (proposta: React + Fastify) mantendo paridade rigorosa.
3. [ ] **CI/CD**: Automação total de testes via GitHub Actions.

## Proximas Stacks Planejadas

- Frontends: React, Svelte, Vue, Angular
- Backends: Node Fastify, Bun Elysia, Python FastAPI
