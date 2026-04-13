# MANIFEST

## Estado Atual

| Grupo | Stack | Status | Papel |
| --- | --- | --- | --- |
| Frontend | SolidJS | validated-local | Baseline inicial de referencia |
| Backend | Go net/http | validated-local | Baseline inicial de referencia |

## Planejado

| Grupo | Stack | Status Previsto |
| --- | --- | --- |
| Frontend | React | planned |
| Frontend | Svelte | planned |
| Frontend | Vue | planned |
| Frontend | Angular | planned |
| Backend | Node Fastify | planned |
| Backend | Bun Elysia | planned |
| Backend | Python FastAPI | planned |

## Criterio de Promocao

Uma stack so deixa o estado bootstrap quando o corte vertical inicial estiver funcional, testado e alinhado a `specs/openapi.yaml`, `shared/schema.sql` e `shared/design-system/`.

### 1. SolidJS + Go (net/http)
- **Estado**: Consolidada (Gabarito Oficial)
- **Diretorio**: `frontends/solidjs` e `backends/go-net-http`
- **Observacao de Estado**: Cobertura E2E completa via Playwright, testes unitários de serviço no backend robustos, e segurança de tokens reforçada.

## Observacao de Estado

`validated-local` significa que a stack passou por validacao local tecnica e de runtime, mas ainda pode depender de consolidacao de commit e ampliacao de cobertura antes de ser tratada como referencia estabilizada.

Na baseline atual, essa validacao local inclui autenticacao completa, gestao real de usuarios, auditoria real persistida, exportacoes CSV/XLSX/PDF, seeds reais, tratamento transacional nas operacoes de escrita, layouts por zona no frontend, diagramas Mermaid em `docs/flows/` e collection Bruno em `specs/bruno/`.
