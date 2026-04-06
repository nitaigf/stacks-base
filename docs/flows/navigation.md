# Navegacao Frontend

State machine da resolucao de rotas no frontend.
Fonte de verdade: `frontends/solidjs/src/utils/router.ts` — funcao `resolveRoute()`.

```mermaid
stateDiagram-v2
  [*] --> Public : GET /

  Public --> AuthLogin : click "Entrar"
  Public --> AuthRegister : click "Criar conta"

  AuthLogin --> AuthRegister : click "Quero criar conta"
  AuthRegister --> AuthLogin : click "Ja tenho conta"

  AuthLogin --> Private : login success
  AuthRegister --> Private : register success

  Private --> Admin : click "Area admin" (role=admin)
  Private --> Error403 : navigate /admin (role!=admin)
  Private --> AuthLogin : logout

  Admin --> Private : click "Voltar ao app"

  Error403 --> Public : click "Voltar para a home"
  Error403 --> AuthLogin : click "Ir para login"
  Error404 --> Public : click "Voltar para a home"
  Error404 --> AuthLogin : click "Ir para login"
  Error500 --> Public : click "Voltar para a home"
  Error500 --> AuthLogin : click "Ir para login"

  state "/ (PublicPage)" as Public
  state "/auth/login (AuthPage)" as AuthLogin
  state "/auth/register (AuthPage)" as AuthRegister
  state "/app (DashboardPage)" as Private
  state "/admin (AdminPage)" as Admin
  state "/errors/403 (ErrorPage)" as Error403
  state "/errors/404 (ErrorPage)" as Error404
  state "/errors/500 (ErrorPage)" as Error500
```

## Guards

| Rota | Anonimo | role=member | role=admin |
|------|---------|-------------|------------|
| `/` | Public | Public | Public |
| `/auth/login` | Auth | Auth | Auth |
| `/auth/register` | Auth | Auth | Auth |
| `/app` | redirect → `/auth/login` | Private | Private |
| `/admin` | redirect → `/auth/login` | Error 403 | Admin |
| desconhecida | Error 404 | Error 404 | Error 404 |
