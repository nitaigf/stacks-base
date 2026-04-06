# Stacks Base — Proposta Refinada

## Estado Atual da Baseline

O repositorio ja possui uma baseline inicial funcional com SolidJS no frontend, Go com net/http no backend e PostgreSQL local em `localhost:5432` como banco principal de desenvolvimento. O corte vertical atualmente implementado e validado localmente cobre `health`, `register`, `login`, `users/me` e `logout`, com design system compartilhado, schema SQL canonico, contrato OpenAPI inicial, integracao SMTP via Mailtrap e paginas publicas, de autenticacao, privadas, administrativas e de erro.

Os artefatos compartilhados ficam em `shared/`, o frontend de referencia em `frontends/solidjs/` e o backend de referencia em `backends/go-net-http/`. A baseline tambem ja passou por validacao de build, testes unitarios e fluxo E2E local em navegador real.

---

## 1. Visão Geral e Propósito

Stacks Base é um projeto open source que oferece implementações completas, funcionais e prontas para execução de uma aplicação web de referência, em diferentes combinações de frontend e backend. Não é um framework, nem um gerador de código — é uma base real, testada e documentada que qualquer desenvolvedor pode clonar, configurar e continuar.

O objetivo central é que alguém possa escolher a stack com a qual já trabalha (ou quer aprender), clonar o repositório correspondente, ajustar as variáveis de ambiente e o `docker-compose.yml`, e já ter um sistema funcionando com autenticação, controle de acesso, gestão de usuários, auditoria e envio de e-mail — sem escrever uma linha sequer dessas fundações.

---

## 2. Princípios Fundamentais

**Paridade total de funcionalidades.** Toda stack implementa exatamente as mesmas telas, os mesmos fluxos e as mesmas regras de negócio. Nenhuma stack é "mais completa" que outra.

**Contrato de API único.** Todos os backends implementam o mesmo contrato definido em OpenAPI. Um frontend React pode apontar para um backend Go sem qualquer adaptação.

**Design system compartilhado.** Todos os frontends compartilham o mesmo arquivo CSS e a mesma estrutura HTML de cada componente. O visual é idêntico entre todas as implementações.

**Estrutura de pastas padronizada.** As convenções de nomenclatura e organização de arquivos são as mesmas entre stacks, respeitando apenas o que o idioma ou framework tornam obrigatoriamente diferente.

**Configuração unificada.** Mesmas variáveis de ambiente, mesmo `docker-compose.yml` base, mesmos serviços de infraestrutura.

**Independência total entre stacks.** Cada combinação funciona de forma autônoma. Não há código ou dependência compartilhada em tempo de execução entre implementações.

**Qualidade mínima garantida em todas.** Linting, formatação, testes com cobertura mínima definida e CI básico são requisitos — não opcionais.

**Previsibilidade transacional e de erro.** Escritas em banco devem acontecer com transação e rollback quando necessário; leituras e escritas devem produzir mensagens e códigos de erro previsíveis, sem estados silenciosos ou tratamentos implícitos frágeis.

**Ordem de construção obrigatória.** A filosofia de trabalho do projeto respeita a seguinte sequência: testes unitários de frontend, frontend e e2e; database único do projeto; testes unitários de backend, backend e e2e quando aplicável; por fim atualização de documentação, OpenAPI e especificações compartilhadas.

---

## 3. O que é Igual Entre Todas as Stacks

### 3.1 Design System

Esta é a decisão mais importante do projeto: **não haverá biblioteca de componentes de UI em nenhum frontend**. Em vez disso, um único design system em CSS puro será compartilhado entre todas as implementações.

O design system consistirá de:

- **Tokens de design** definidos como CSS Custom Properties: cores, tipografia, espaçamento, bordas, sombras, breakpoints e transições. Troca de tema claro/escuro via uma única variável no `:root`.
- **Componentes CSS** para cada elemento visual recorrente: botão (variantes primário, secundário, destrutivo, ghost), campo de formulário, label, mensagem de erro, card, modal, toast, tabela, badge de status, spinner, avatar, dropdown, menu lateral e topbar.
- **Layout base** com classes utilitárias mínimas para grid, flex, espaçamento e responsividade — sem depender de Tailwind ou UnoCSS.
- **Estrutura HTML semântica definida** para cada componente. Cada framework renderiza exatamente essa estrutura — a diferença é só a sintaxe de template.

O resultado visual será pixel-perfect entre React, SolidJS, Svelte, Vue e Angular.

### 3.2 Contrato de API

Todos os backends expõem os mesmos endpoints, com os mesmos métodos HTTP, os mesmos formatos de request e response, e os mesmos códigos de status. O contrato é formalizado em um arquivo OpenAPI (YAML) que vive na raiz do repositório e serve como fonte de verdade.

Convenções do contrato:
- Base: `/api/v1/`
- Recursos em kebab-case e plural: `/api/v1/users`, `/api/v1/audit-logs`
- Respostas em JSON com camelCase
- Erros sempre no formato `{ error: { code, message, details? } }`
- Paginação sempre no formato `{ data: [], meta: { page, perPage, total, totalPages } }`
- Autenticação via Bearer token (PASETO v4 local) no header Authorization
- Refresh token via cookie httpOnly

### 3.3 Funcionalidades Implementadas

A especificação funcional é idêntica ao documento original — todas as telas, todos os fluxos de autenticação, CRUD de usuários, logs de auditoria, bloqueio de conta, confirmações por e-mail e aprovação de novos registros. Nenhuma stack omite ou simplifica qualquer parte.

Na baseline atual, ja existem explicitamente:
- pagina publica
- paginas de autenticacao
- pagina privada
- pagina administrativa
- paginas de erro para 403, 404 e 500

### 3.4 Variáveis de Ambiente

Um único `.env.example` de referência define todas as variáveis necessárias, com nomes idênticos em todos os projetos. Frontends usam prefixo de framework quando necessário (ex: `VITE_` no Vite), mas os nomes base são os mesmos.

Na baseline inicial, o envio de e-mail e configurado por SMTP e a referencia local de desenvolvimento usa Mailtrap.

### 3.5 Docker Compose

Um único `docker-compose.yml` base com perfis `dev` e `prod`. Todos os projetos utilizam os mesmos serviços de infraestrutura: PostgreSQL 16, Redis 7 e Caddy como reverse proxy com TLS automático. Healthchecks configurados em todos os serviços.

### 3.6 Banco de Dados

O schema do banco de dados é definido uma vez, como referência, em SQL puro. Todas as implementações de backend reproduzem esse schema exatamente, cada uma com sua própria ferramenta de migrations. As tabelas, colunas, tipos, constraints e índices são idênticos.

### 3.7 Convenções Gerais

- **Commits:** Conventional Commits (`feat:`, `fix:`, `docs:`, `chore:`, `test:`)
- **Branches:** `feat/`, `fix/`, `docs/`, `chore/`
- **Endpoints REST:** kebab-case, plural, versionados
- **JSON:** camelCase em requests e responses
- **Variáveis de ambiente:** UPPER\_SNAKE\_CASE

### 3.8 Testes

Cada stack define os mesmos cenários de teste: registro, login, logout, recuperação de senha, edição de perfil, CRUD de usuários, auditoria. A cobertura mínima exigida é a mesma em todos. A abordagem (unitário, integração, e2e) é adaptada ao que é idiomático em cada linguagem.

- **Testes unitários**: Vitest para frontends TypeScript, Go test para backends Go.
- **Testes de contrato**: collection Bruno versionada em `shared/bruno/` com assertions por endpoint, executável via CLI.
- **Testes E2E**: Playwright para navegação em browser real.
- **Documentação de fluxo**: diagramas Mermaid em `docs/flows/` para que agentes de IA e desenvolvedores entendam rapidamente a navegação do frontend e as sequências de API.

---

## 4. O que é Diferente Entre as Stacks

O que varia entre implementações é exclusivamente o que pertence ao ecossistema de cada linguagem ou framework:

- Sintaxe de template e reatividade (JSX, Single-File Components, diretivas)
- Gerenciamento de estado local e servidor
- Roteamento e proteção de rotas
- Validação de formulários
- Cliente HTTP e interceptors
- ORM, query builder ou acesso direto ao banco
- Ferramenta de migrations
- Framework de testes e utilitários
- Idiomas específicos da linguagem (middleware vs interceptor vs plugin)

O que **nunca** varia: o comportamento visível ao usuário, o design, os endpoints consumidos, as regras de negócio e a estrutura de pastas (dentro do que o idioma permite).

---

## 5. Frontends

### Princípios comuns a todos os frontends

- Design system CSS compartilhado, sem biblioteca de componentes de UI
- Cliente HTTP com interceptors para injeção de token e tratamento de erros
- Token PASETO armazenado em memória (nunca em localStorage); refresh token em cookie httpOnly
- Zod para validação de formulários e parse de responses da API
- Rotas protegidas com verificação de papel (público, autenticado, admin)
- Internacionalização preparada mas não obrigatória na v1

---

### React
**Runtime:** Node.js + Vite + TypeScript

| Categoria | Ferramenta |
|---|---|
| Roteamento | TanStack Router |
| Estado servidor | TanStack Query |
| Formulários | React Hook Form + Zod |
| Cliente HTTP | Ky |
| Testes | Vitest + Testing Library |

**Removidos da proposta original:** Argon2 e Paseto (pertencem ao backend), tRPC (REST puro; variante tRPC é separada), Shadcn/Radix e Tailwind (substituídos pelo design system compartilhado).

---

### SolidJS
**Runtime:** Node.js + Vite + TypeScript

| Categoria | Ferramenta |
|---|---|
| Roteamento | TanStack Router |
| Estado servidor | TanStack Query (adaptador Solid) |
| Formulários | Lógica nativa com signals + Zod |
| Cliente HTTP | Ky |
| Testes | Vitest |

**Removidos:** Argon2, Paseto, tRPC, Kobalte, Corvu, UnoCSS.

---

### Svelte
**Runtime:** Node.js + Vite + TypeScript (Svelte 5, SPA pura — não SvelteKit, para manter consistência com as outras stacks)

| Categoria | Ferramenta |
|---|---|
| Roteamento | TanStack Router (adaptador Svelte) |
| Estado servidor | TanStack Query (adaptador Svelte) |
| Formulários | Superforms + Zod |
| Cliente HTTP | Ky |
| Testes | Vitest |

**Removidos:** Argon2, Paseto, tRPC, Svelte Material UI, UnoCSS.

**Nota importante:** escolher Svelte puro + Vite (SPA) em vez de SvelteKit mantém a paridade arquitetural com as demais stacks. SvelteKit (SSR/SSG) seria uma variante separada futura.

---

### Vue
**Runtime:** Node.js + Vite + TypeScript

| Categoria | Ferramenta |
|---|---|
| Roteamento | Vue Router (padrão do ecossistema) |
| Estado servidor + cliente | Pinia + TanStack Query |
| Formulários | VeeValidate + Zod |
| Cliente HTTP | Ky |
| Testes | Vitest + Vue Testing Library |

**Removidos:** Argon2, Paseto, tRPC, Vuetify, UnoCSS, TanStack Router (Vue Router é o padrão idiomático).

---

### Angular
**Runtime:** Node.js + Angular CLI + TypeScript (Angular 17+, build com esbuild nativo)

| Categoria | Ferramenta |
|---|---|
| Roteamento | Angular Router (padrão do framework) |
| Estado servidor | TanStack Query (adaptador Angular, experimental) |
| Formulários | Angular Reactive Forms + Zod (adaptador) |
| Cliente HTTP | Angular HttpClient (padrão do framework) |
| Testes | Jest + Angular Testing Library |

**Removidos:** Vite (Angular usa seu próprio build system), TanStack Router, Ky, Argon2, Paseto, tRPC, Angular Material, UnoCSS.

**Nota importante:** Angular tem o ecossistema mais diferente dos demais. O que muda é a sintaxe e as ferramentas — o comportamento, o design e o contrato de API são idênticos. Injetores, serviços e signals são os idiomas corretos aqui.

---

## 6. Backends

### Princípios comuns a todos os backends

- Mesmo contrato OpenAPI
- PASETO v4 local para tokens de acesso e refresh
- Argon2id para hash de senha (custo calibrado e documentado)
- Rate limiting em rotas de autenticação
- Headers de segurança (CORS, HSTS, Content-Security-Policy via Caddy)
- Auditoria registrada de forma assíncrona (não bloqueia a resposta)
- Envio de e-mail via SMTP configurável (Resend, Mailgun, Mailtrap para dev)
- Migrations versionadas e reversíveis
- Health check endpoint (`/health`)
- Documentação OpenAPI servida em `/docs`
- Escritas em banco encapsuladas em transações com rollback em caso de falha
- Tratamento explícito e previsível de erros de leitura, escrita e autenticação

---

### Node.js + Fastify
**Runtime:** Node.js 20+ + TypeScript + tsup

| Categoria | Ferramenta |
|---|---|
| Framework HTTP | Fastify |
| ORM | Drizzle ORM |
| Banco | PostgreSQL (pg / postgres.js) |
| Cache / sessão | Redis (ioredis) |
| Hash de senha | argon2 |
| Tokens | @panva/paseto |
| Validação | Zod |
| Testes | Vitest |

**Removidos:** tRPC (REST puro; tRPC seria uma variante separada TS-to-TS).

---

### Bun + Elysia
**Runtime:** Bun 1.x + TypeScript

| Categoria | Ferramenta |
|---|---|
| Framework HTTP | Elysia (type-safe nativo, usa TypeBox) |
| ORM | Drizzle ORM |
| Banco | PostgreSQL (bun:postgres ou postgres.js) |
| Cache / sessão | Redis (ioredis) |
| Hash de senha | argon2 |
| Tokens | @panva/paseto |
| Validação | TypeBox (nativo do Elysia) + Zod nos limites da API |
| Testes | Bun test |

---

### Python + FastAPI
**Runtime:** Python 3.12+ + Uvicorn + uv (gerenciamento de dependências)

| Categoria | Ferramenta |
|---|---|
| Framework HTTP | FastAPI |
| ORM | SQLAlchemy 2.0 (async) + asyncpg |
| Migrations | Alembic |
| Cache / sessão | Redis (redis-py async) |
| Hash de senha | argon2-cffi |
| Tokens | python-paseto |
| Validação | Pydantic v2 |
| Testes | pytest + pytest-asyncio + httpx |

**Removidos:** tRPC (exclusivo do ecossistema JS/TS).

---

### Go
**Runtime:** Go 1.22+

| Categoria | Ferramenta |
|---|---|
| Framework HTTP | Chi (router leve e idiomático) |
| Acesso a dados | sqlc + pgx (geração de código a partir de SQL — mais idiomático que GORM) |
| Migrations | golang-migrate |
| Cache / sessão | Redis (go-redis) |
| Hash de senha | golang.org/x/crypto (argon2) |
| Tokens | o1ecc8o/paseto |
| Validação | go-playground/validator |
| Testes | testing padrão + testify |

**Notas:** GORM foi substituído por sqlc porque Go favorece queries explícitas e tipagem gerada — é mais seguro, mais rápido e mais idiomático. tRPC removido (JS/TS only).

---

## 7. Estrutura de Pastas

A estrutura abaixo é a referência canônica. Cada implementação segue essa organização com as adaptações mínimas que o idioma exige.

### Frontends (todos)

```
src/
  assets/        → CSS do design system, imagens, fontes
  components/    → Componentes reutilizáveis (botão, modal, tabela...)
  layouts/       → Layouts de página (público, autenticado, admin)
  pages/         → Telas da aplicação organizadas por grupo
  services/      → Cliente HTTP configurado e funções de chamada de API
  stores/        → Estado global (auth, preferências)
  schemas/       → Schemas Zod compartilhados
  utils/         → Funções auxiliares
  types/         → Tipos e interfaces TypeScript
```

### Backends (todos, adaptado ao idioma)

```
  config/        → Configurações e variáveis de ambiente
  routes/        → Definição de rotas (ou handlers/, controllers/)
  services/      → Lógica de negócio
  repositories/  → Acesso a dados (ou models/, queries/)
  middlewares/   → Auth, rate limit, logging, erros
  schemas/       → Validação de entrada e saída
  migrations/    → Migrations versionadas
  utils/         → Funções auxiliares
  tests/         → Testes organizados por módulo
```

---

## 8. Repositório e Organização do Projeto

O projeto viverá em um monorepo com a seguinte organização:

```
stacks-base/
  shared/
    openapi.yaml         → Contrato de API (fonte de verdade)
    schema.sql           → Schema do banco (referência)
    design-system/       → CSS, tokens, assets
    .env.example         → Variáveis de ambiente de referência
    docker-compose.yml   → Infraestrutura compartilhada

  frontends/
    react/
    solidjs/
    svelte/
    vue/
    angular/

  backends/
    node-fastify/
    bun-elysia/
    python-fastapi/
    go-chi/

  docs/
    → Documentação do projeto, decisões de arquitetura, guias
```

Cada pasta de frontend e backend é um projeto independente com seu próprio `package.json` (ou `go.mod`, `pyproject.toml`), seu próprio README e seus próprios scripts de desenvolvimento e build.

### Estado Implementado Agora

```text
stacks-base/
  shared/
    openapi.yaml
    schema.sql
    .env.example
    docker-compose.yml
    design-system/

  frontends/
    solidjs/

  backends/
    go-net-http/

  docs/
```

Neste momento SolidJS e Go net/http sao a referencia executavel inicial do projeto. As demais stacks continuam planejadas e devem perseguir paridade a partir desta base, sem alterar as fontes de verdade compartilhadas.

---

## 9. Roadmap Sugerido

1. **Fase 1 — Fundação:** Definir e documentar o contrato OpenAPI, o schema SQL e o design system CSS. Estes três artefatos são a fonte de verdade de todo o projeto.
2. **Fase 2 — Baseline de referência:** Implementar e testar completamente a baseline SolidJS + Go net/http sobre PostgreSQL local. Esta fase ja esta iniciada e validada localmente.
3. **Fase 3 — Consolidação da baseline:** Formalizar commit inicial, consolidar README, TODO, ADR e MANIFEST e estabilizar a baseline como referencia oficial.
4. **Fase 4 — Demais backends:** Implementar Node Fastify, Python FastAPI e Bun Elysia, validando cada um contra o contrato OpenAPI e o frontend de referencia.
5. **Fase 5 — Demais frontends:** Implementar React, Svelte, Vue e Angular, validando cada um contra o backend de referencia.
6. **Fase 6 — Matriz de compatibilidade:** Testar todas as combinações frontend × backend. Documentar e publicar.

---

## 10. O que Este Projeto Não É

- Não é um starter kit opinativo para um projeto específico
- Não é um gerador de código ou CLI
- Não é uma tentativa de criar "o melhor stack" — todos são equivalentes por design
- Não é um tutorial — é código de produção simplificado, mas não didático

A intenção é que qualquer desenvolvedor, independente do seu stack preferido, encontre aqui uma base honesta, sem gambiarras, sem atalhos e sem excesso de abstração — pronta para ser o ponto de partida de algo real.

---

## FUTURO

- Adicionar SvelteKit como variante de frontend
- Adicionar Qwik como variante de frontend
- Adicionar Rust (Axum) como variante de backend
- Adicionar GraphQL como variante de contrato (com Yoga, Apollo ou Strawberry)
- Adicionar testes de integração e2e com Playwright
