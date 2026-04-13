# TECHNOLOGIES

## Visao Geral

Este documento registra as tecnologias efetivamente adotadas na baseline atual e as stacks planejadas para replicacao futura.

## Baseline Atual

### Infraestrutura Local

- **PostgreSQL 16+**
  - banco principal de desenvolvimento
  - porta padrao `5432`
  - schema canonico em `shared/schema.sql`
- **Docker Compose**
  - opcional
  - arquivo de apoio em `shared/docker-compose.yml`
  - nao e o fluxo padrao de desenvolvimento

### Design System

- **CSS puro**
  - sem biblioteca de componentes
  - tokens e componentes compartilhados em `shared/design-system/index.css`

### Frontend: SolidJS

- **Runtime**: Node.js 20+
- **Package Manager**: npm 11.5.1
- **Build Tool**: Vite 6.2.2
- **Framework**: SolidJS 1.9.5
- **TypeScript**: 5.8.3
- **Roteamento**: TanStack Router 1.168.9
- **Cliente HTTP**: Ky 1.7.5
- **Validacao**: Zod 3.24.2
- **Testes**: Vitest 3.0.8 + SolidJS Testing Library

### Backend: Go net/http

- **Runtime**: Go 1.24.0
- **Framework**: net/http
- **Database Driver**: lib/pq 1.10.9
- **Password Hash**: golang.org/x/crypto (Argon2)
- **Tokens**: aidanwoods.dev/go-paseto 1.6.0
- **Migrations**: implementacao propria com SQL idempotente
- **E-mail**: SMTP configuravel
- **Auditoria**: persistida em PostgreSQL

## Contrato, Especificacoes e Testes

- **OpenAPI 3.1.0**
  - arquivo: `specs/openapi.yaml`
  - base: `/api/v1/`
  - respostas JSON com camelCase
  - autenticacao via Bearer token
- **Bruno**
  - localizacao: `specs/bruno/`
  - execucao via CLI
  - assertions por endpoint
- **Testes Unitarios**
  - Vitest no frontend
  - Go test no backend
- **E2E**
  - Playwright como padrao adotado para navegacao real

## Documentacao

- **Fluxos**: diagramas Mermaid em `docs/flows/`
- **Arquitetura**: decisoes em `ADR.md`
- **API**: especificacao versionada em `specs/openapi.yaml`

## Requisitos Minimos

- **Node.js**: 20+
- **npm**: 11+
- **Go**: 1.22+
- **PostgreSQL**: 16+
- **Docker**: opcional

## Stacks Planejadas

### Frontends

- SolidJS (referência principal)
- React Native
- Angular
- React
- Qwik
- Vue

### Backends

- Go Net/HTTP (referência principal) + sqlc
- tRPC + Drizzle
- NestJS + Fastify + Prisma
- FastAPI + SQLModel
- Elysia + Drizzle
- AdonisJS + Lucid ORM

### Fullstack Frameworks

- SvelteKit (Svelte) + Drizzle
- NextJS (React) + Prisma
- NuxtJS (Vue) + TypeORM

### Database Compartilhado

- PostgreSQL
