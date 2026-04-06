# ARCHITECTURE

## Objetivo

Este documento define como o monorepo e organizado e quais fronteiras devem ser respeitadas para que a baseline inicial possa ser replicada nas demais stacks com o menor atrito possivel.

## Layout Canonico

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
  README.md
  TODO.md
  CONSTITUTION.md
  ARCHITECTURE.md
  MANIFEST.md
  AGENTS.md
  ADR.md
```

## Boundaries

- shared/ contem as fontes de verdade compartilhadas.
- shared/bruno/ contem a collection Bruno versionada com requests e assertions para validacao de contrato.
- backends/ contem implementacoes independentes, sem codigo compartilhado em runtime.
- frontends/ contem implementacoes independentes, sem dependencia em biblioteca de UI externa.
- docs/ contem materiais operacionais, guias e detalhamentos nao normativos.
- docs/flows/ contem diagramas Mermaid de navegacao e sequencias de API, mantidos como documentacao viva.

## Baseline Inicial

### Backend Go

- Entrada em cmd/server.
- Codigo de aplicacao em internal/.
- Roteamento com net/http.
- Persistencia em PostgreSQL.
- Migrations versionadas no proprio backend.

### Frontend SolidJS

- Aplicacao SPA com Vite.
- Cliente HTTP centralizado em services/.
- Estado de autenticacao em stores/.
- Schemas de validacao em schemas/.
- Consumo direto do design system compartilhado via src/assets.

## Regras de Evolucao

- Novas stacks devem reproduzir primeiro o corte vertical inicial antes de expandir funcionalidade.
- Qualquer nova camada compartilhada deve ser declarativa e nao pode introduzir acoplamento em runtime.
- O contrato OpenAPI e o schema SQL sao upstream; implementacoes nunca sao a fonte de verdade.

## Ordem Operacional

Ao construir ou expandir uma baseline, a ordem de trabalho esperada e:

1. Testes unitarios de frontend, frontend e e2e de frontend.
2. Banco de dados unico e compartilhado pelo projeto.
3. Testes unitarios de backend, backend e e2e de backend quando aplicavel.
4. Atualizacao da documentacao normativa e sincronizacao de shared/openapi.yaml e das demais especificacoes.

Essa ordem deve ser entendida como disciplina de entrega, nao como preferencia local da stack atual.