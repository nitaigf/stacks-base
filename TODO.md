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

# TODO

## Regras

- Este arquivo e organizado por fases.
- Itens em Progresso Validado ou Concluido so podem mudar de estado depois de acao validada no codigo e registrada no repositorio.
- Toda mudanca comitada que altere o estado validado do sistema deve atualizar este arquivo no mesmo conjunto de mudancas.
- Planejamento Futuro pode existir aqui, mas nao conta como progresso executado.
- Mudancas fora do escopo do README exigem registro previo em ADR.

## Validado no Codigo e em Checks Locais Atuais

- Backend Go validado com `go test ./...`.
- Frontend SolidJS validado com `npm test` e `npm run build`.
- Frontend SolidJS validado com `npx playwright test` cobrindo home publica, navegacao auth e fluxo admin critico com backend real local.
- Frontend SolidJS com rotas publicas, autenticacao, area privada, area administrativa, auditoria e erros ligado ao backend real.
- Backend Go com `health`, autenticacao completa, `users/me`, gestao administrativa de usuarios, exportacoes e auditoria real persistida.
- `specs/openapi.yaml` e `specs/bruno/` presentes e cobrindo o contrato publicado da baseline atual.
- Playwright configurado no frontend com suite E2E inicial versionada e runner validado com `npx playwright test --list`.
- Seeds reais de admin, usuarios demonstrativos e auditoria inicial ativos no backend.
- Migrations e schema canonico presentes para `users`, `refresh_tokens`, `password_reset_tokens` e `audit_logs`.

## Consolidacao Obrigatoria Antes da Proxima Stack

- [x] Implementar E2E de frontend via Playwright cobrindo pelo menos o fluxo autenticado completo e a navegacao administrativa critica.
- [x] Ampliar os testes automatizados do backend para os fluxos administrativos e de auditoria, nao apenas auth e exportacoes.
- [x] Remover ou arquivar o frontend legado com mocks em `frontends/solidjs/src/pages/AdminPage.tsx` e `frontends/solidjs/src/components/UserTable.tsx`.
- [x] Revalidar a sessao do frontend ao entrar em rotas privadas e administrativas, tratando token expirado com limpeza de sessao.
- [x] Padronizar tratamento de erro do frontend por status e codigo da API, sem depender de substring de mensagem.
- [x] Adicionar validacao client-side para forgot password, reset password, change password e editor administrativo de usuario.
- [x] Corrigir o fluxo de logout apos change password para encerrar a sessao tambem no backend.
- [x] Sincronizar a documentacao operacional em `docs/flows/` com os caminhos e arquivos reais da baseline atual.

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

- [x] Revisar a collection Bruno para explicitar cobertura e asserts do escopo administrativo detalhado, incluindo cenarios negativos relevantes.

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

- [x] Ampliar a cobertura automatizada de testes do backend para create, show, update, deactivate, reactivate, soft-delete, restore, hard-delete e listagem de auditoria.
- [x] Tratar violacoes de unicidade de email como erro de contrato previsivel em create, update e restore, inclusive em concorrencia.
- [x] Adicionar indice apropriado para `refresh_tokens.token_hash` e revisar a politica de unicidade e revogacao por hash.
- [x] Revisar a politica de `password_reset_tokens` para evitar multiplos tokens ativos simultaneamente para o mesmo usuario.
- [ ] Decidir formalmente a estrategia de refresh de sessao: implementar endpoint dedicado ou documentar explicitamente a exigencia de novo login apos expiracao do access token.
- [x] Completar a historia de rollback das migrations sem quebrar a baseline canonica atual.
- [x] Tornar `RunInTx` resiliente a panic.
- [ ] Avaliar estrategia de indices e retencao para `audit_logs` conforme o volume crescer.
- [ ] Revisar readiness operacional do backend, incluindo health check com dependencia de banco e documentacao fiel dos requisitos de SMTP e seeds.

## Fase 3 - Frontend de Referencia em SolidJS

### Progresso Validado

- [x] Criar aplicacao SolidJS com Vite.
- [x] Migrar o roteamento para TanStack Router.
- [x] Implementar login, registro, dashboard autenticado e fluxos de senha.
- [x] Integrar cliente HTTP com backend Go.
- [x] Implementar telas administrativas reais de usuarios e auditoria, sem mocks.
- [x] Adicionar testes e build de producao.

### Planejamento Futuro

- [x] Implementar Playwright cobrindo login, register, forgot password, reset password, dashboard, change password e logout.
- [x] Ampliar a cobertura E2E para listagem de usuarios, visualizacao individual, criacao, edicao, status, soft-delete, restore, hard-delete, exportacoes e auditoria.
- [x] Adicionar testes de componente ou integracao para router, guards, reidratacao de sessao e tratamento de 401, 403 e 500.
- [x] Remover o codigo legado de admin nao roteado e manter um unico fluxo administrativo real sem mocks.
- [x] Melhorar a UX do shell privado para nao expor entrada administrativa a usuarios `member`.

## Fase 4 - Integracao e Promocao da Baseline

### Progresso Validado

- [x] Validar localmente o fluxo autenticado e administrativo contra PostgreSQL real.
- [x] Atualizar MANIFEST com o status real da baseline.
- [x] Registrar as decisoes implementadas em ADR.
- [x] Sincronizar README, CONSTITUTION, ARCHITECTURE, TECHNOLOGIES e HOW-TO-USE com o estado real do sistema.

### Planejamento Futuro

- [x] Sincronizar `README.md`, `HOW-TO-USE.md`, `MANIFEST.md` e `docs/flows/` com o estado real do codigo antes de promover a baseline.
- [x] Revisar a documentacao operacional do backend para remover subdeclaracoes de escopo e explicitar dependencias obrigatorias de ambiente local.
- [ ] Preparar checklist de replicacao para a proxima stack somente apos consolidacao final desta baseline.
