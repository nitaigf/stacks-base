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