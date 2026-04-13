# Navegacao Frontend

State machine da resolucao de rotas no frontend.
Fonte de verdade: `frontends/solidjs/src/router.tsx`.

```mermaid
stateDiagram-v2
  [*] --> Public : GET /

  Public --> AuthLogin : click "Entrar"
  Public --> AuthRegister : click "Criar conta"

  AuthLogin --> AuthRegister : click "Quero criar conta"
  AuthRegister --> AuthLogin : click "Ja tenho conta"
  AuthLogin --> ForgotPassword : click "Esqueci minha senha"
  ForgotPassword --> AuthLogin : click "Voltar ao login"
  ResetPassword --> AuthLogin : click "Voltar ao login"

  AuthLogin --> Private : login success
  AuthRegister --> Private : register success

  Private --> Admin : click "Area admin" (role=admin)
  Private --> ChangePassword : click "Alterar senha"
  Private --> Error403 : navigate /admin (role!=admin)
  Private --> AuthLogin : logout
  ChangePassword --> AuthLogin : change password success

  Admin --> UserDetails : click "Visualizar usuario"
  Admin --> UserDetails : click "Novo usuario"
  UserDetails --> Admin : click "Voltar"
  UserDetails --> UserDetails : save success (redirect to view)
  Admin --> AuditLogs : click "Auditoria"
  AuditLogs --> Admin : click "Usuarios"
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
  state "/auth/forgot-password (ForgotPasswordPage)" as ForgotPassword
  state "/auth/reset-password (ResetPasswordPage)" as ResetPassword
  state "/app (DashboardPage)" as Private
  state "/app/change-password (ChangePasswordPage)" as ChangePassword
  state "/admin/users (AdminUsersPage)" as Admin
  state "/admin/users/:userId (UserEditorPage)" as UserDetails
  state "/admin/audit-logs (AuditLogsPage)" as AuditLogs
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
| `/auth/forgot-password` | Auth | Auth | Auth |
| `/auth/reset-password` | Auth | Auth | Auth |
| `/app` | redirect → `/auth/login` | Private | Private |
| `/app/change-password` | redirect → `/auth/login` | Private | Private |
| `/admin/users` | redirect → `/auth/login` | Error 403 | Admin |
| `/admin/users/:userId` | redirect → `/auth/login` | Error 403 | Admin |
| `/admin/audit-logs` | redirect → `/auth/login` | Error 403 | Admin |
| desconhecida | Error 404 | Error 404 | Error 404 |
