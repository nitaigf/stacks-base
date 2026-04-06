# HOW TO USE - Guia para Novos Colaboradores

Este guia explica como configurar e rodar a stack Stacks Base localmente para desenvolvimento e testes.

## Pré-requisitos

### Software Necessário
- **Node.js** 20+ (para frontend)
- **Go** 1.22+ (para backend)
- **Docker** 20.10+ (para PostgreSQL)
- **Docker Compose** (para infraestrutura)
- **Git** (para controle de versão)

### Ferramentas Recomendadas
- **VS Code** com extensões para Go, TypeScript e Docker
- **Postico** ou **pgAdmin** (para visualizar o banco)
- **Bruno** ou **Postman** (para testar API)

## Passo 1: Clonar e Configurar o Projeto

```bash
# Clone o repositório
git clone <repository-url>
cd stacks-base

# Copie as variáveis de ambiente de exemplo
cp shared/.env.example .env
```

### Configurar Variáveis de Ambiente

Edite o arquivo `.env` com suas configurações:

```bash
# Configurações básicas
APP_ENV=development
APP_NAME=stacks-base

# Database (já configurado para Docker)
DATABASE_HOST=127.0.0.1
DATABASE_PORT=5432
DATABASE_NAME=stacks_base
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres

# Backend e Frontend
BACKEND_HOST=127.0.0.1
BACKEND_PORT=8080
FRONTEND_PORT=3000
VITE_API_BASE_URL=http://127.0.0.1:8080

# Segurança - GERE CHAVES NOVAS!
ACCESS_TOKEN_SECRET=sua-chave-secreta-long-aqui
REFRESH_TOKEN_SECRET=outra-chave-secreta-diferente-aqui

# E-mail (opcional para desenvolvimento)
SMTP_HOST=sandbox.smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USERNAME=seu-usuario-mailtrap
SMTP_PASSWORD=sua-senha-mailtrap

# Seed de admin (já configurado)
ADMIN_SEED_ENABLED=true
ADMIN_INITIAL_NAME=Admin
ADMIN_INITIAL_EMAIL=admin@stacks-base.local
ADMIN_INITIAL_PASSWORD=Admin@123456
```

## Passo 2: Iniciar Infraestrutura (PostgreSQL)

```bash
# Na raiz do projeto, inicie o PostgreSQL
docker-compose -f shared/docker-compose.yml up -d

# Verifique se o container está rodando
docker ps
```

O PostgreSQL estará disponível em `localhost:5432` com:
- Banco: `stacks_base`
- Usuário: `postgres`
- Senha: `postgres`

O schema será aplicado automaticamente na primeira inicialização.

## Passo 3: Configurar e Rodar o Backend

```bash
# Entre na pasta do backend
cd backends/go-net-http

# Instale dependências Go
go mod download

# Configure variáveis de ambiente (se ainda não fez)
export $(cat ../../.env | xargs)

# Rode o servidor
go run ./cmd/server
```

O backend estará rodando em `http://127.0.0.1:8080`

### Verificar Backend

```bash
# Health check
curl http://127.0.0.1:8080/health

# Deve retornar: {"status":"ok"}
```

### Seed do Admin

Com `ADMIN_SEED_ENABLED=true`, o backend criará automaticamente o usuário admin na primeira inicialização:

- **Email**: `admin@stacks-base.local`
- **Senha**: `Admin@123456`
- **Nome**: `Admin`

## Passo 4: Configurar e Rodar o Frontend

```bash
# Em outro terminal, entre na pasta do frontend
cd frontends/solidjs

# Copie variáveis de ambiente
cp .env.example .env

# Instale dependências
npm install

# Configure variáveis de ambiente (se ainda não fez)
export $(cat ../../.env | xargs)

# Rode o servidor de desenvolvimento
npm run dev
```

O frontend estará rodando em `http://127.0.0.1:3000`

## Passo 5: Acessar e Testar a Aplicação

### Acessar as Áreas

1. **Página Pública**: `http://127.0.0.1:3000/`
2. **Login**: `http://127.0.0.1:3000/login`
3. **Registro**: `http://127.0.0.1:3000/register`

### Login como Admin

1. Acesse `http://127.0.0.1:3000/login`
2. Use as credenciais:
   - Email: `admin@stacks-base.local`
   - Senha: `Admin@123456`

### Fluxo de Teste Completo

1. **Registro de Novo Usuário**:
   - Acesse `/register`
   - Preencha o formulário com email válido
   - Após registro, você será redirecionado para login

2. **Login e Acesso**:
   - Faça login com o novo usuário
   - Você será redirecionado para a área privada (`/private`)
   - Verifique se os dados do usuário são carregados

3. **Acesso Admin**:
   - Faça login como admin
   - Acesse a área administrativa (`/admin`)
   - Visualize a lista de usuários (se implementado)

4. **Logout**:
   - Clique em logout
   - Verifique se o token é invalidado
   - Tente acessar páginas privadas (deve redirecionar para login)

## Testes Automatizados

### Testes de Frontend

```bash
cd frontends/solidjs
npm test
```

### Testes de Backend

```bash
cd backends/go-net-http
go test ./...
```

### Testes de Contrato (API)

```bash
# Se tiver Bruno CLI instalado
cd shared/bruno
bru run
```

## Estrutura de Pastas Importantes

```
stacks-base/
├── shared/                 # Recursos compartilhados
│   ├── openapi.yaml       # Contrato da API
│   ├── schema.sql         # Schema do banco
│   ├── design-system/     # CSS compartilhado
│   └── docker-compose.yml # Infraestrutura
├── frontends/solidjs/     # Frontend de referência
│   ├── src/
│   │   ├── assets/        # CSS do design system
│   │   ├── pages/         # Páginas da aplicação
│   │   ├── services/      # Cliente HTTP
│   │   └── stores/        # Estado global
├── backends/go-net-http/  # Backend de referência
│   ├── cmd/server/        # Ponto de entrada
│   ├── internal/          # Código da aplicação
│   └── migrations/        # Migrations do banco
└── docs/                  # Documentação adicional
```

## Comandos Úteis

### Docker

```bash
# Ver logs do PostgreSQL
docker-compose -f shared/docker-compose.yml logs -f postgres

# Parar infraestrutura
docker-compose -f shared/docker-compose.yml down

# Reiniciar banco (perde dados)
docker-compose -f shared/docker-compose.yml down -v
docker-compose -f shared/docker-compose.yml up -d
```

### Backend

```bash
# Build para produção
cd backends/go-net-http
go build -o server ./cmd/server

# Rodar com rebuild automático
go run ./cmd/server
```

### Frontend

```bash
# Build para produção
cd frontends/solidjs
npm run build

# Preview do build
npm run preview
```

## Troubleshooting Comum

### Backend não inicia
- Verifique se PostgreSQL está rodando: `docker ps`
- Verifique variáveis de ambiente: `export $(cat .env | xargs)`
- Verifique se a porta 8080 está livre: `lsof -i :8080`

### Frontend não conecta no backend
- Verifique se `VITE_API_BASE_URL` está correto
- Verifique CORS no backend (`BACKEND_ALLOWED_ORIGIN`)
- Verifique se backend está rodando: `curl http://127.0.0.1:8080/health`

### PostgreSQL não conecta
- Reinicie o container: `docker-compose -f shared/docker-compose.yml restart postgres`
- Verifique se a porta 5432 está livre: `lsof -i :5432`
- Aguarde o health check: `docker-compose -f shared/docker-compose.yml logs postgres`

### Admin seed não funciona
- Verifique `ADMIN_SEED_ENABLED=true`
- Verifique se o email já existe no banco
- Delete o usuário manualmente se necessário

## Próximos Passos

1. Explore o código fonte em `frontends/solidjs/src/` e `backends/go-net-http/internal/`
2. Leia `README.md` para entender a arquitetura
3. Consulte `TECHNOLOGIES.md` para detalhes das tecnologias
4. Verifique `TODO.md` para tarefas pendentes
5. Leia `ARCHITECTURE.md` para regras do projeto

## Suporte

- **Documentação**: Verifique os arquivos `.md` na raiz
- **Contrato API**: `shared/openapi.yaml`
- **Schema DB**: `shared/schema.sql`
- **Issues**: Abra uma issue no repositório
