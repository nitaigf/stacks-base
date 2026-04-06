# TECHNOLOGIES

## Visão Geral

Este documento registra todas as tecnologias utilizadas no projeto Stacks Base, incluindo as stacks atuais, planejadas e a infraestrutura compartilhada.

## Infraestrutura Compartilhada

### Banco de Dados
- **PostgreSQL 16** - Banco de dados principal
  - Porta padrão: 5432
  - Schema definido em `shared/schema.sql`
  - Extensões: `pgcrypto` para geração de UUID

### Containerização
- **Docker** - Orquestração de serviços
- **Docker Compose** - Definição de infraestrutura local
  - Arquivo: `shared/docker-compose.yml`
  - Serviços: PostgreSQL com volumes persistidos

### Design System
- **CSS Puro** - Sem bibliotecas de componentes
  - Tokens de design em CSS Custom Properties
  - Componentes compartilhados em `shared/design-system/index.css`
  - Tema claro/escuro via variável CSS

## Stack Atual: Baseline de Referência

### Frontend: SolidJS
- **Runtime**: Node.js 20+
- **Build Tool**: Vite 6.2.2
- **Framework**: SolidJS 1.9.5
- **TypeScript**: 5.8.3
- **Cliente HTTP**: Ky 1.7.5
- **Validação**: Zod 3.24.2
- **Testes**: Vitest 3.0.8 + SolidJS Testing Library
- **Roteamento**: Navegação SPA nativa

### Backend: Go net/http
- **Runtime**: Go 1.24.0
- **Framework**: net/http (stdlib)
- **Database Driver**: lib/pq 1.10.9
- **Password Hash**: golang.org/x/crypto (Argon2)
- **Tokens**: aidanwoods.dev/go-paseto 1.6.0
- **Migrations**: Implementação própria com SQL idempotente

## Stacks Planejadas

### Frontends

#### React
- **Runtime**: Node.js + Vite + TypeScript
- **Roteamento**: TanStack Router
- **Estado Servidor**: TanStack Query
- **Formulários**: React Hook Form + Zod
- **Cliente HTTP**: Ky
- **Testes**: Vitest + Testing Library

#### Svelte
- **Runtime**: Node.js + Vite + TypeScript (Svelte 5, SPA)
- **Roteamento**: TanStack Router (adaptador Svelte)
- **Estado Servidor**: TanStack Query (adaptador Svelte)
- **Formulários**: Superforms + Zod
- **Cliente HTTP**: Ky
- **Testes**: Vitest

#### Vue
- **Runtime**: Node.js + Vite + TypeScript
- **Roteamento**: Vue Router
- **Estado**: Pinia + TanStack Query
- **Formulários**: VeeValidate + Zod
- **Cliente HTTP**: Ky
- **Testes**: Vitest + Vue Testing Library

#### Angular
- **Runtime**: Node.js + Angular CLI + TypeScript (Angular 17+)
- **Build**: esbuild nativo
- **Roteamento**: Angular Router
- **Estado Servidor**: TanStack Query (adaptador Angular)
- **Formulários**: Angular Reactive Forms + Zod
- **Cliente HTTP**: Angular HttpClient
- **Testes**: Jest + Angular Testing Library

### Backends

#### Node.js + Fastify
- **Runtime**: Node.js 20+ + TypeScript
- **Build**: tsup
- **Framework**: Fastify
- **ORM**: Drizzle ORM
- **Cache/Sessão**: Redis (ioredis)
- **Hash**: argon2
- **Tokens**: @panva/paseto
- **Validação**: Zod
- **Testes**: Vitest

#### Bun + Elysia
- **Runtime**: Bun 1.x + TypeScript
- **Framework**: Elysia
- **ORM**: Drizzle ORM
- **Banco**: bun:postgres ou postgres.js
- **Cache/Sessão**: Redis (ioredis)
- **Hash**: argon2
- **Tokens**: @panva/paseto
- **Validação**: TypeBox + Zod
- **Testes**: Bun test

#### Python + FastAPI
- **Runtime**: Python 3.12+ + Uvicorn
- **Gerenciamento**: uv
- **Framework**: FastAPI
- **ORM**: SQLAlchemy 2.0 (async) + asyncpg
- **Migrations**: Alembic
- **Cache/Sessão**: Redis (redis-py async)
- **Hash**: argon2-cffi
- **Tokens**: python-paseto
- **Validação**: Pydantic v2
- **Testes**: pytest + pytest-asyncio + httpx

## Contrato e Validação

### API
- **OpenAPI 3.1.0** - Contrato de API único
  - Arquivo: `shared/openapi.yaml`
  - Base: `/api/v1/`
  - Respostas JSON com camelCase
  - Autenticação via Bearer token (PASETO v4)

### Testes de Contrato
- **Bruno** - Collection de testes de API
  - Localização: `shared/bruno/`
  - Execução via CLI
  - Assertions por endpoint

## Comunicação e E-mail

### SMTP
- **Configuração**: Variáveis de ambiente
- **Provedores**: Resend, Mailgun, Mailtrap (dev)
- **Headers**: From name e address configuráveis

## Segurança

### Autenticação
- **Access Token**: PASETO v4 local (15m TTL)
- **Refresh Token**: PASETO v4 local via cookie httpOnly (168h TTL)
- **Password Hash**: Argon2id com custo calibrado

### Headers de Segurança
- **CORS**: Configurável por origem
- **HSTS**: Via Caddy reverse proxy
- **CSP**: Content-Security-Policy via Caddy

## Desenvolvimento e Testes

### Testes Automatizados
- **Unitários**: Vitest (frontend), Go test (backend)
- **E2E**: Playwright com navegador real
- **Contrato**: Bruno CLI collection

### Documentação
- **Fluxos**: Diagramas Mermaid em `docs/flows/`
- **Arquitetura**: Decisões em ADR.md
- **API**: OpenAPI servida em `/docs`

## Variáveis de Ambiente

### Prefixos
- Frontends Vite: `VITE_`
- Backend: sem prefixo
- Compartilhadas: definidas em `shared/.env.example`

### Categorias
- **App**: APP_ENV, APP_NAME
- **Database**: DATABASE_*
- **Backend**: BACKEND_*
- **Frontend**: FRONTEND_*, VITE_*
- **Tokens**: ACCESS_TOKEN_*, REFRESH_TOKEN_*
- **SMTP**: SMTP_*
- **Admin Seed**: ADMIN_*

## Versionamento e Compatibilidade

### Requisitos Mínimos
- **Node.js**: 20+ (para todos os frontends)
- **Go**: 1.22+ (para backends Go)
- **Python**: 3.12+ (para backends Python)
- **Bun**: 1.x+ (para backends Bun)
- **PostgreSQL**: 16+
- **Docker**: 20.10+

### Compatibilidade
- Todas as stacks implementam o mesmo contrato OpenAPI
- Design system CSS é compatível com todos os browsers modernos
- Schema SQL é idêntico entre todas as implementações
