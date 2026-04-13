# CONSTITUTION

## 1. Escopo

Stacks Base existe para fornecer implementacoes reais e equivalentes da mesma aplicacao de referencia em stacks diferentes. Nenhum trabalho pode reduzir, desviar ou reinterpretar esse objetivo sem ADR aprovado.

## 2. Regras Inegociaveis

1. Paridade funcional maxima entre stacks.
2. Contrato de API unico definido em `specs/openapi.yaml`.
3. Schema SQL canonico definido em `shared/schema.sql`.
4. Design system compartilhado definido em `shared/design-system/`.
5. Estrutura de pastas padronizada, com adaptacoes minimas exigidas pelo ecossistema.
6. Qualquer divergencia intencional entre stacks exige ADR antes da implementacao.

## 3. Contrato de API

- Base path: /api/v1
- Recursos REST em kebab-case e plural.
- JSON em camelCase.
- Erro padrao: { error: { code, message, details? } }
- Paginacao padrao: { data, meta: { page, perPage, total, totalPages } }
- Access token via Bearer token.
- Refresh token via cookie httpOnly.

## 4. Qualidade Minima

- Toda stack deve expor /health.
- Toda stack deve cobrir autenticacao completa: registro, login, logout, users/me, forgot-password, reset-password e change-password.
- Toda stack deve cobrir gestao completa de usuarios: listar paginado, visualizar um, criar, editar, soft-delete, restore, inativar, reativar, hard-delete, exportar CSV, exportar XLSX e gerar PDF imprimivel.
- Toda stack deve cobrir auditoria real persistida em banco para autenticacao, leitura, escrita, exportacao e impressao.
- Testes, lint e formatacao nao sao opcionais.
- Segredos nao podem ser versionados.
- Operacoes de gravacao devem usar transacao quando envolverem mais de um efeito persistente ou quando a consistencia exigir rollback.
- Leituras e escritas devem retornar erros previsiveis e coerentes com o contrato publicado.
- Frontends devem contemplar paginas publicas, de autenticacao, privadas, administrativas e de erro.

## 5. Ordem de Construcao

Toda evolucao funcional deve respeitar a seguinte sequencia operacional:

1. Testes unitarios de frontend, depois frontend, depois e2e de frontend.
2. Banco de dados unico do projeto.
3. Testes unitarios de backend, depois backend, depois e2e de backend quando fizer sentido para a stack.
4. Atualizacao da documentacao, do OpenAPI e das especificacoes compartilhadas.

Essa ordem existe para impedir deriva de escopo, divergencia entre stacks e atraso na sincronizacao das fontes de verdade.

## 6. Artefatos de Teste Obrigatorios

- Collection Bruno em `specs/bruno/` validando todos os endpoints do contrato OpenAPI.
- Diagramas de fluxo em `docs/flows/` descrevendo navegacao frontend e sequencias de API.
- Testes unitarios de frontend e backend presentes e passando antes de promover qualquer fase.
- Testes E2E de frontend via Playwright cobrindo o fluxo autenticado completo.
- OpenAPI permanece hand-written em `specs/openapi.yaml`. Geradores stack-specific nao sao permitidos.

## 7. Rail Guards

- Nao introduzir framework, biblioteca de UI, protocolo ou padrao arquitetural fora do README sem ADR.
- Nao criar documentacao operacional fora da pasta docs/, nem especificacoes executaveis fora de `specs/`.
- Nao marcar fase como concluida no TODO sem validacao real do artefato correspondente.
- Nao commitar mudanca que altere estado validado sem atualizar o TODO no mesmo conjunto de mudancas.
- Nao tratar a baseline inicial como excecao: ela deve seguir as mesmas regras que serao exigidas das demais stacks.
- Nao inverter a ordem de construcao definida nesta constituicao sem ADR.

## 8. Baseline Inicial

- Frontend de referencia: SolidJS.
- Backend de referencia: Go com net/http.
- Banco de dados principal para desenvolvimento: PostgreSQL local em localhost:5432.
- Docker Compose e opcional para infraestrutura auxiliar e reproducao local.

## 9. Filosofia dos Ecossistemas

As escolhas tecnicas dentro deste monorepo obedecem as seguintes diretrizes absolutas de foco:

1. **Baseline de Referencia (SolidJS e Go net/http + sqlc):** Foco exclusivo em performance *bare metal*, controle absoluto e minimo de abstracoes de terceiros. Estas ferramentas devem extrair o maximo da biblioteca padrao (std lib) e primitives de baixo nivel.
2. **Demais Stacks (React, NestJS, etc.):** Foco maximo em produtividade guiada por bibliotecas estabelecidas e padroes idiomaticos de mercado. Deve-se aproveitar ao maximo o ferramental "pronto" que cada ecossistema oferece com o objetivo de compor arquiteturas rapidas e sem invencoes extraordinarias.
3. **Design System:** Fica restrito a Vanilla CSS com uma organizacao puramente Atomica e alta aderencia a componentes ARIA e semantica nativa HTML. E terminantemente proibida a injecao de engrenagens utilitarias (TailwindCSS, UnoCSS) ou primitivas robustas externas (Radix, Shadcn, Fluent2 pre-fabricados, Bootstrap) no design system. O visual deve se manter minimalista, garantindo que nao roubara o foco arquitetural do projeto.
