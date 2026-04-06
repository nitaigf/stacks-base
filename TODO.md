# TODO

## Regras

- Este arquivo e organizado por fases.
- Itens em Progresso Validado ou Concluido so podem mudar de estado depois de acao validada no codigo e registrada no repositorio.
- Planejamento Futuro pode existir aqui, mas nao conta como progresso executado.
- Mudancas fora do escopo do README exigem registro previo em ADR.

## Validado Localmente Nesta Sessao

- Documentos normativos da raiz criados.
- Baseline inicial definida e materializada em shared/, backends/go-net-http/ e frontends/solidjs/.
- Backend Go validado com go test ./... e com fluxo real health -> register -> users/me -> logout contra PostgreSQL local em 127.0.0.1:5432.
- Frontend SolidJS validado com npm run build e npm run test.
- Integracao E2E da baseline validada em navegador real com fluxo register -> users/me -> logout.
- Fluxo E2E validado tambem para pagina publica, area admin negada e tela de erro 403.
- Backend endurecido com transacoes em gravacoes sensiveis e resposta previsivel de erros.
- SMTP Mailtrap configurado localmente via .env da raiz.
- ADR-004 registrado: testabilidade por IA, diagramas Mermaid, Bruno versionado, OpenAPI hand-written.
- Diagramas Mermaid criados em docs/flows/ para navegacao frontend e sequencia de auth.
- Collection Bruno criada em shared/bruno/ com requests e assertions para todo o corte vertical.
- Layouts por zona implementados no frontend: PublicLayout, AuthLayout, PrivateLayout, AdminLayout, ErrorLayout.
- Paginas refatoradas para delegar navegacao aos layouts.
- Testes unitarios do frontend passando (11 testes, 3 arquivos).
- E2E com layouts validado em navegador real (12 checkpoints).
- CONSTITUTION e ARCHITECTURE atualizados com artefatos de teste obrigatorios.

## Pronto Para Primeiro Commit

- docs: governanca da raiz, baseline validada localmente e README sincronizado com o estado atual.

## Fase 0 - Governanca e Guard Rails

### Progresso Validado

- [ ] Criar os documentos normativos da raiz.
- [ ] Definir baseline inicial SolidJS + Go net/http + PostgreSQL local 5432.
- [ ] Formalizar regra de paridade maxima e fontes de verdade compartilhadas.

### Planejamento Futuro

- [ ] Consolidar README apos o primeiro commit da baseline.

## Fase 1 - Espinha Dorsal Compartilhada

### Progresso Validado

- [ ] Criar shared/openapi.yaml com o corte vertical inicial.
- [ ] Criar shared/schema.sql para users, refresh_tokens e audit_logs.
- [ ] Criar design system CSS compartilhado.
- [ ] Criar .env.example e docker-compose opcional.

### Planejamento Futuro

- [ ] Definir estrategia de reutilizacao do design system entre stacks sem fork.

## Fase 2 - Backend de Referencia em Go

### Progresso Validado

- [ ] Subir servidor net/http com health check.
- [ ] Implementar registro, login, logout e users/me.
- [ ] Persistir dados no PostgreSQL local 5432.
- [ ] Adicionar migrations e testes do backend.

### Planejamento Futuro

- [ ] Adicionar rotacao completa de refresh token.
- [ ] Adicionar auditoria assincrona com fila simples.
- [ ] Expandir cobertura de transacoes para futuras operacoes de CRUD e administracao.

## Fase 3 - Frontend de Referencia em SolidJS

### Progresso Validado

- [ ] Criar aplicacao SolidJS com Vite.
- [ ] Implementar login, registro e dashboard autenticado.
- [ ] Integrar cliente HTTP com backend Go.
- [ ] Adicionar testes e build de producao.

### Planejamento Futuro

- [ ] Consolidar componentes reutilizaveis adicionais do design system.
- [ ] Evoluir a navegacao simples atual para um roteador idiomatico mantendo a mesma estrutura de paginas.

## Fase 4 - Integracao e Promocao da Baseline

### Progresso Validado

- [ ] Validar o fluxo register -> login -> users/me -> logout.
- [ ] Atualizar MANIFEST com o status real da baseline.
- [ ] Registrar as decisoes implementadas em ADR.

### Planejamento Futuro

- [ ] Preparar checklist de replicacao para as demais stacks.