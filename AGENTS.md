# AGENTS

## Finalidade

Este arquivo define papeis e guard rails de execucao para contribuidores humanos e agentes automatizados dentro deste repositorio.

## Papeis

### Maintainer

- Aprova mudancas em README, CONSTITUTION, ARCHITECTURE, MANIFEST e ADR.
- Decide excecoes de escopo.
- Valida promocao de fase no TODO.

### Stack Owner

- Mantem uma implementacao especifica alinhada as fontes de verdade compartilhadas.
- Corrige desvios de paridade antes de adicionar novas features.

### Contributor

- Implementa trabalho aprovado dentro do escopo atual.
- Atualiza testes e documentacao correspondente.
- Respeita a ordem de construcao definida em CONSTITUTION e ARCHITECTURE.

### Coding Agent

- Deve tratar README, CONSTITUTION, ARCHITECTURE, AGENTS e ADR como prioridade sobre docs historicos.
- Nao pode expandir escopo por iniciativa propria.
- Deve atualizar TODO quando o estado validado mudar.
- Deve registrar em ADR qualquer decisao estrutural efetivamente aplicada.
- Deve seguir a ordem: testes unitarios de frontend -> frontend -> e2e -> database unico -> testes unitarios de backend -> backend -> e2e quando aplicavel -> documentacao, OpenAPI e especificacoes.

## Guard Rails Operacionais

1. Nao criar stacks novas antes de a baseline inicial estar funcional.
2. Nao introduzir dependencia compartilhada em runtime entre implementacoes.
3. Nao criar documentos fora da raiz quando forem normativos, nem fora de docs/ quando forem operacionais.
4. Nao marcar itens como concluidos sem validacao tecnica correspondente.
5. Nao substituir as fontes de verdade compartilhadas por implementacoes locais.
6. Nao quebrar a ordem de construcao sem justificativa formalizada em ADR.

## Fluxo de Escalacao

1. Divergencia de implementacao -> Stack Owner.
2. Divergencia de arquitetura ou escopo -> Maintainer.
3. Mudanca de regra ou excecao permanente -> ADR.