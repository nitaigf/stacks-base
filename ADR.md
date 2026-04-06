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

A baseline precisa ser inteiramente testavel por agentes de IA — fronte, backend, banco e contratos. Ferramentas de documentacao e teste devem ser legíveis por maquina e mantidas como artefatos vivos no repositorio. A escolha de como especificar APIs impacta diretamente a paridade multi-stack.

### Decisao

- **Testabilidade por IA**: todo artefato executavel do projeto deve ser legivel e acionavel por agentes automatizados. Vitest e Go test cobrem unitarios, Playwright cobre E2E de frontend, Bruno CLI cobre contratos de API.
- **Mermaid**: diagramas minimalistas de fluxo ficam em `docs/flows/`. Funcionam como mapa rapido para agentes navegarem frontend e entenderem sequencias de API. Nao substituem o codigo como fonte de verdade.
- **OpenAPI hand-written**: `shared/openapi.yaml` permanece como especificacao unica, escrita manualmente. Geradores stack-specific como swaggo, tsoa ou similares nao sao permitidos porque violam a regra de paridade — cada stack geraria contratos ligeiramente diferentes.
- **Bruno versionado**: uma collection com requests e assertions vive em `shared/bruno/`. Pode ser executada via CLI para validacao de contrato. Substitui o papel de Postman como ferramenta exploratoria e de teste.

### Consequencias

- Novas stacks devem garantir que seus endpoints passem na mesma collection Bruno.
- Diagramas Mermaid devem ser atualizados quando rotas ou fluxos mudarem.
- Nenhum gerador de OpenAPI pode ser introduzido sem ADR que demonstre paridade preservada.

## ADR-005 - Padrão de Coluna de Ações em Tabelas

- Data: 2026-04-06
- Status: accepted

### Contexto

A interface do usuário precisa de um padrão consistente para ações em tabelas de dados. A posição e o formato da coluna de ações impactam diretamente a usabilidade e a experiência do usuário em telas administrativas e de listagem.

### Decisão

- A coluna de ações em tabelas deve ser sempre a primeira coluna (à esquerda)
- As ações devem ser implementadas como botões com ícones, não com texto
- Cada botão deve ter um tooltip/hint que descreve a ação ao passar o mouse
- Ícones devem seguir o design system compartilhado e ser semanticamente claros
- Ações devem seguir ordem de importância: primária (editar), secundária (excluir/desativar), etc.

### Consequências

- Todas as novas tabelas implementadas no projeto devem seguir este padrão
- Tabelas existentes devem ser migradas para conformidade com este ADR
- O design system deve incluir ícones padronizados para ações comuns
- A experiência do usuário se torna mais consistente em toda a aplicação