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

Uma stack so deixa o estado bootstrap quando o corte vertical inicial estiver funcional, testado e alinhado a shared/openapi.yaml, shared/schema.sql e shared/design-system/.

## Observacao de Estado

validated-local significa que a stack passou por validacao local tecnica e E2E, mas ainda pode depender de commit e consolidacao documental antes de ser tratada como referencia estabilizada.

Na baseline atual, essa validacao local inclui SMTP Mailtrap, paginas de erro e tratamento transacional nas operacoes de escrita ja implementadas.