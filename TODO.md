# TODO

## Regras

- Este arquivo e organizado por fases.
- Itens em Progresso Validado ou Concluido so podem mudar de estado depois de acao validada no codigo e registrada no repositorio.
- Toda mudanca comitada que altere o estado validado do sistema deve atualizar este arquivo no mesmo conjunto de mudancas.
- Planejamento Futuro pode existir aqui, mas nao conta como progresso executado.
- Mudancas fora do escopo do README exigem registro previo em ADR.

## Validado Localmente Nesta Sessao

- Backend Go validado com `go test ./...`.
- Frontend SolidJS validado com `npm test` e `npm run build`.
- Fluxo real validado contra PostgreSQL local em `127.0.0.1:5432` para:
  - health
  - login admin
  - listagem administrativa de usuarios
  - criacao, edicao, soft-delete, restore e hard-delete
  - exportacoes CSV, XLSX e PDF
  - change password
  - forgot password
  - reset password
  - logout
  - listagem de auditoria
  - persistencia dos eventos de auditoria no banco
- Frontend SolidJS migrado para TanStack Router.
- Especificacoes movidas para `specs/openapi.yaml` e `specs/bruno/`.
- Seeds reais de admin, usuarios demonstrativos e auditoria inicial ativos no backend.
- Baseline sem mocks em telas administrativas e de auditoria.
- Documentacao normativa e operacional sincronizada com o estado real do sistema.

## Pronto Para Primeiro Commit

- chore: baseline alinhada ao escopo prometido, validada localmente e documentada com fidelidade.

## Fase 0 - Governanca e Guard Rails

### Progresso Validado

- [x] Criar os documentos normativos da raiz.
- [x] Definir baseline inicial SolidJS + Go net/http + PostgreSQL local 5432.
- [x] Formalizar regra de paridade maxima e fontes de verdade compartilhadas.

### Planejamento Futuro

- [ ] Formalizar a politica de promocao de stack de `validated-local` para referencia estabilizada.

## Fase 1 - Espinha Dorsal Compartilhada

### Progresso Validado

- [x] Criar `specs/openapi.yaml` com o contrato canonico da baseline.
- [x] Criar `shared/schema.sql` para users, refresh_tokens, password_reset_tokens e audit_logs.
- [x] Criar design system CSS compartilhado.
- [x] Criar `.env.example` e `docker-compose` opcional.
- [x] Criar `specs/bruno/` como collection versionada da API.

### Planejamento Futuro

- [ ] Expandir a collection Bruno para cobrir todo o escopo administrativo detalhado.

## Fase 2 - Backend de Referencia em Go

### Progresso Validado

- [x] Subir servidor net/http com health check.
- [x] Implementar registro, login, logout e users/me.
- [x] Implementar forgot password, reset password e change password.
- [x] Implementar gestao completa de usuarios com CRUD administrativo, soft-delete, restore, inativacao, reativacao, hard-delete e exportacoes.
- [x] Persistir dados no PostgreSQL local 5432.
- [x] Adicionar migrations, seeds reais e testes do backend.
- [x] Implementar auditoria real persistida com rota, metodo, IP, user agent e metadata.

### Planejamento Futuro

- [ ] Ampliar a cobertura automatizada de testes do backend para os fluxos administrativos novos.

## Fase 3 - Frontend de Referencia em SolidJS

### Progresso Validado

- [x] Criar aplicacao SolidJS com Vite.
- [x] Migrar o roteamento para TanStack Router.
- [x] Implementar login, registro, dashboard autenticado e fluxos de senha.
- [x] Integrar cliente HTTP com backend Go.
- [x] Implementar telas administrativas reais de usuarios e auditoria, sem mocks.
- [x] Adicionar testes e build de producao.

### Planejamento Futuro

- [ ] Ampliar a cobertura E2E para o escopo administrativo completo.

## Fase 4 - Integracao e Promocao da Baseline

### Progresso Validado

- [x] Validar localmente o fluxo autenticado e administrativo contra PostgreSQL real.
- [x] Atualizar MANIFEST com o status real da baseline.
- [x] Registrar as decisoes implementadas em ADR.
- [x] Sincronizar README, CONSTITUTION, ARCHITECTURE, TECHNOLOGIES e HOW-TO-USE com o estado real do sistema.

### Planejamento Futuro

- [ ] Preparar checklist de replicacao para a proxima stack somente apos consolidacao final desta baseline.
