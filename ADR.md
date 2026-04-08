# ADR

## ADR-001 - Baseline Inicial do Projeto

- Data: 2026-04-06
- Status: accepted

### Contexto

O projeto precisava sair do estado de proposta para uma fundacao executavel sem comprometer a regra de paridade entre stacks.

### Decisao

- A baseline inicial sera composta por SolidJS no frontend e Go com net/http no backend.
- O banco principal de desenvolvimento sera PostgreSQL local em localhost:5432.
- Docker Compose sera mantido como opcao secundaria para reproducao e apoio de infraestrutura.
- O TODO sera mantido por fases, com planejamento futuro permitido e progresso concluido restrito a trabalho validado.
- Os arquivos normativos da raiz serao README, TODO, CONSTITUTION, ARCHITECTURE, MANIFEST, AGENTS e ADR.

### Consequencias

- A primeira implementacao passa a ser a referencia para todas as demais combinacoes.
- Qualquer desvio da baseline ou do README exige novo ADR.

## ADR-002 - Ordem de Construcao da Baseline e das Futuras Stacks

- Data: 2026-04-06
- Status: accepted

### Contexto

O projeto precisava de uma disciplina de entrega explicita para impedir que frontend, backend, banco e especificacoes evoluissem fora de sincronia.

### Decisao

- A ordem de construcao obrigatoria passa a ser: testes unitarios de frontend, frontend, e2e de frontend, database unico do projeto, testes unitarios de backend, backend, e2e de backend quando aplicavel, e por fim atualizacao da documentacao, OpenAPI e especificacoes compartilhadas.

### Consequencias

- Contributors e agentes devem planejar entregas seguindo essa ordem.
- Inversoes dessa sequencia passam a exigir justificativa formal em ADR.

## ADR-003 - Previsibilidade de Erro, Transacoes e Entrega de E-mail na Baseline

- Data: 2026-04-06
- Status: accepted

### Contexto

O projeto tambem serve como referencia de boas praticas. A baseline precisava demonstrar tratamento previsivel de erros, transacoes explicitas em gravacoes e envio real de e-mail no ambiente local de desenvolvimento.

### Decisao

- Operacoes de gravacao sensiveis passam a usar transacoes no backend com rollback em caso de falha.
- Leituras e escritas devem retornar erros de aplicacao previsiveis, com codigo e mensagem coerentes.
- O ambiente local de desenvolvimento passa a usar Mailtrap SMTP configurado por `.env` na raiz.
- A baseline passa a incluir paginas publicas, de autenticacao, privadas, administrativas e de erro.

### Consequencias

- Novas stacks devem reproduzir nao so o comportamento funcional, mas tambem a disciplina de erro e transacao.
- Documentacao e especificacoes precisam carregar essa expectativa como parte da paridade.

## ADR-004 - Testabilidade por IA, Diagramas Mermaid, Bruno e Estrategia OpenAPI

- Data: 2026-04-06
- Status: accepted

### Contexto

A baseline precisa ser inteiramente testavel por agentes de IA. Ferramentas de documentacao e teste devem ser legiveis por maquina e mantidas como artefatos vivos no repositorio. A escolha de como especificar APIs impacta diretamente a paridade multi-stack.

### Decisao

- Todo artefato executavel do projeto deve ser legivel e acionavel por agentes automatizados.
- Diagramas minimalistas de fluxo ficam em `docs/flows/`.
- OpenAPI permanece hand-written e geradores stack-specific nao sao permitidos.
- Bruno versionado substitui o papel de colecoes ad hoc de API.

### Consequencias

- Novas stacks devem garantir que seus endpoints passem na mesma collection Bruno.
- Diagramas Mermaid devem ser atualizados quando rotas ou fluxos mudarem.
- Nenhum gerador de OpenAPI pode ser introduzido sem ADR que demonstre paridade preservada.

## ADR-005 - Padrao de Coluna de Acoes em Tabelas

- Data: 2026-04-06
- Status: accepted

### Contexto

A interface do usuario precisa de um padrao consistente para acoes em tabelas de dados.

### Decisao

- A coluna de acoes em tabelas deve ser sempre a primeira coluna.
- As acoes devem ser implementadas como botoes com icones ou rotulos curtos.
- Cada botao deve ter tooltip ou hint descritivo.
- A ordem de exibicao deve refletir a importancia da acao.

### Consequencias

- Novas tabelas devem seguir esse padrao.
- Tabelas existentes devem convergir para esse padrao.

## ADR-006 - Consolidacao de Especificacoes e Package Manager Oficial da Baseline

- Data: 2026-04-08
- Status: accepted

### Contexto

A baseline passou a incluir escopo administrativo real, auditoria real, exportacoes e seeds demonstrativos. Os caminhos antigos das especificacoes e a multiplicidade de lockfiles no frontend criavam ambiguidade sobre quais artefatos eram canonicos e como o ambiente devia ser reproduzido.

### Decisao

- OpenAPI hand-written passa a viver em `specs/openapi.yaml`.
- A collection Bruno versionada passa a viver em `specs/bruno/`.
- `shared/` permanece reservado para artefatos compartilhados de runtime e desenvolvimento, como schema SQL, design system, `.env.example` e `docker-compose`.
- O package manager oficial do frontend SolidJS passa a ser `npm`, com `package-lock.json` como lockfile canonico da baseline atual.
- Toda mudanca comitada que altere o estado validado do sistema deve atualizar `TODO.md` no mesmo conjunto de mudancas.

### Consequencias

- A raiz do repositorio passa a separar claramente runtime compartilhado de especificacoes executaveis.
- A baseline deixa de depender de interpretacao sobre caminhos antigos como `shared/openapi.yaml` e `shared/bruno/`.
- O frontend passa a ter reproducao univoca de dependencias enquanto Bun permanece como candidato de stack futura, nao como package manager normativo desta baseline.
