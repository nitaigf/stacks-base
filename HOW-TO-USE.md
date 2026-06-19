<!-- KHKB GOVERNANCE PORTAL START -->
---

# 🏛️ Diretrizes de Governança (Database-Centric Docs)

Para manter este projeto limpo e compatível com as regras de desenvolvimento da Kovesh-Halutz, toda a documentação de arquitetura (ADRs), especificações de produto (Specs), instruções de IA (Agents/Skills) e Constituições residem **diretamente no banco de dados central PostgreSQL** (`kh_knowledge`).

Não crie nem edite arquivos markdown locais em pastas como `docs/`, `specs/` ou na raiz do projeto. 

### 🛠️ Comandos de Consulta e Registro (Uso por Humanos e IAs)

*   **Pesquisar decisões/especificações (Busca Semântica):**
    ```bash
    ./kh-ask "Sua dúvida, conceito ou termo técnico"
    ```
*   **Buscar e Sintetizar Respostas (RAG com IA Local):**
    ```bash
    ./kh-ask "Como implementamos o isolamento RLS neste projeto?" --rag
    ```
*   **Registrar novas decisões (ADRs) ou requisitos (Specs) direto no banco:**
    Executa a CLI interativa para criar novos registros indexados vetorialmente:
    ```bash
    ./kh-write
    ```
    Ou via argumentos (ideal para agentes de IA atualizarem o banco programaticamente):
    ```bash
    ./kh-write --type decision --project stacks-base --title "ADR 008" --content "Conteúdo..."
    ```

---
<!-- KHKB GOVERNANCE PORTAL END -->

# HOW TO USE - Guia para Novos Colaboradores

Este guia explica como configurar e rodar a stack Stacks Base localmente para desenvolvimento e testes.

## Pre-requisitos

### Software Necessario

- **Node.js** 20+ (frontend)
- **npm** 11+ (frontend atual)
- **Go** 1.22+ (backend)
- **PostgreSQL** 16+ rodando localmente em `127.0.0.1:5432`
- **Git**

### Ferramentas Recomendadas

- **VS Code** com extensoes para Go e TypeScript
- **Postico** ou **pgAdmin**
- **Bruno**
- **Playwright**
- **Docker Compose** como alternativa opcional para infraestrutura

## Passo 1: Clonar e Configurar o Projeto

```bash
git clone <repository-url>
cd stacks-base
cp shared/.env.example .env
```

Edite o `.env` com suas configuracoes locais quando necessario.

## Passo 2: Garantir o PostgreSQL Local

O fluxo padrao de desenvolvimento usa PostgreSQL local, nao Docker.

```bash
createdb -h 127.0.0.1 -p 5432 stacks_base
```

Ou, se preferir:

```bash
psql -h 127.0.0.1 -p 5432 -d postgres -c 'create database stacks_base;'
```

### Alternativa Opcional com Docker

```bash
docker compose -f shared/docker-compose.yml up -d
```

## Passo 3: Rodar o Backend

```bash
cd backends/go-net-http
go mod download
export $(cat ../../.env | xargs)
go run ./cmd/server
```

O backend fica em `http://127.0.0.1:8080`.

### Seed do Admin

Com `ADMIN_SEED_ENABLED=true`, o backend cria automaticamente:

- **Email**: `admin@stacks-base.local`
- **Senha**: `Admin@123456`
- **Nome**: `Admin`

### Seed Demonstrativo

Com `DEMO_SEED_ENABLED=true`, o backend popula o banco com usuarios reais de exemplo e logs iniciais de auditoria.

## Passo 4: Rodar o Frontend

```bash
cd frontends/solidjs
cp .env.example .env
npm install
export $(cat ../../.env | xargs)
npm run dev
```

O frontend fica em `http://127.0.0.1:3000`.

## Passo 5: Exercitar a Aplicacao

### Areas Principais

1. Pagina publica: `http://127.0.0.1:3000/`
2. Login: `http://127.0.0.1:3000/auth/login`
3. Registro: `http://127.0.0.1:3000/auth/register`
4. Recuperacao de senha: `http://127.0.0.1:3000/auth/forgot-password`
5. Dashboard autenticado: `http://127.0.0.1:3000/app`
6. Usuarios admin: `http://127.0.0.1:3000/admin/users`
7. Auditoria admin: `http://127.0.0.1:3000/admin/audit-logs`

### Fluxo Recomendo de Validacao Manual

1. Registrar um novo usuario real
2. Fazer login e validar `users/me`
3. Alterar a senha
4. Entrar como admin
5. Validar gestao de usuarios:
   - listar
   - visualizar
   - criar
   - editar
   - soft-delete
   - restore
   - inativar
   - reativar
   - hard-delete
   - exportar CSV
   - exportar XLSX
   - gerar PDF
6. Validar auditoria administrativa

## Testes Automatizados

### Frontend

```bash
cd frontends/solidjs
npm test
npm run build
```

### Frontend E2E

Com backend, banco e frontend local disponiveis:

```bash
cd frontends/solidjs
npm run test:e2e
```

A suite agora cobre fluxos completos de registro, login, esqueci senha, gestao de usuarios, auditoria e paginas de erro.

Se quiser sobrescrever credenciais admin da seed:

```bash
cd frontends/solidjs
PW_ADMIN_EMAIL=admin@stacks-base.local PW_ADMIN_PASSWORD='Admin@123456' npm run test:e2e
```

### Novas Migrations

Sempre que houver mudancas estruturais no backend Go, aplique as migrations:

- `001_initial_schema.up.sql`: Schema basico.
- `002_user_seeds.up.sql`: Dados iniciais (opcional).
- `003_token_indexes.up.sql`: Indices de performance para tokens e auditoria.

### Backend

```bash
cd backends/go-net-http
go test ./...
```

### Contrato de API

```bash
cd specs/bruno
bru run
```

## Estrutura Importante

```text
stacks-base/
  shared/
    schema.sql
    design-system/
    docker-compose.yml
  specs/
    openapi.yaml
    bruno/
  frontends/solidjs/
  backends/go-net-http/
  docs/
```

## Comandos Uteis

### PostgreSQL Local

```bash
psql -h 127.0.0.1 -p 5432 -d stacks_base -c 'select now();'
```

### Docker Opcional

```bash
docker compose -f shared/docker-compose.yml up -d
docker compose -f shared/docker-compose.yml logs -f postgres
```

### Backend

```bash
cd backends/go-net-http
go build -o server ./cmd/server
go run ./cmd/server
```

### Frontend

```bash
cd frontends/solidjs
npm run build
npm run preview
npm run test:e2e
```

## Troubleshooting

### Backend nao inicia

- Verifique se PostgreSQL local esta rodando
- Verifique variaveis de ambiente com `export $(cat .env | xargs)`
- Verifique se a porta 8080 esta livre com `lsof -i :8080`

### Frontend nao conecta no backend

- Verifique `VITE_API_BASE_URL`
- Verifique CORS no backend com `BACKEND_ALLOWED_ORIGIN`
- Verifique o health check em `http://127.0.0.1:8080/health`
- Verifique se a sessao seed/admin esperada pelo Playwright corresponde ao `.env`

### PostgreSQL nao conecta

- Verifique a instancia local em `127.0.0.1:5432`
- Verifique se a porta 5432 esta livre com `lsof -i :5432`
- Se preferir, suba a alternativa Docker com `docker compose -f shared/docker-compose.yml up -d`

## Suporte

- **Documentacao**: arquivos `.md` da raiz
- **Contrato API**: `specs/openapi.yaml`
- **Schema DB**: `shared/schema.sql`
