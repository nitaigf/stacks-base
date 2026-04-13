# Fluxo de Autenticacao

Diagramas de sequencia dos endpoints de auth.
Fonte de verdade: `specs/openapi.yaml` e `backends/go-net-http/internal/services/auth_service.go`.

## Register

```mermaid
sequenceDiagram
  participant B as Browser
  participant F as Frontend
  participant API as Backend (Go)
  participant DB as PostgreSQL
  participant Mail as Mailtrap SMTP

  B->>F: preenche nome, email, senha
  F->>F: valida com Zod (registerSchema)
  F->>API: POST /api/v1/auth/register {name, email, password}
  API->>API: valida payload (schemas.RegisterRequest)
  API->>API: hash senha (Argon2id)

  rect rgb(240, 248, 255)
    note right of API: BEGIN TRANSACTION
    API->>DB: SELECT user por email
    alt email ja existe
      API-->>F: 409 {error: {code: "email_taken"}}
    else email disponivel
      API->>DB: INSERT user
      API->>API: gera access token (PASETO v4)
      API->>API: gera refresh token + hash (SHA-256)
      API->>DB: INSERT refresh_token
    end
    note right of API: COMMIT
  end

  API->>DB: INSERT audit_log (async)
  API-->>Mail: SMTP email de boas-vindas (async)
  API-->>F: 201 {data: {accessToken, user}} + Set-Cookie refresh_token
  F->>F: authStore.setSession(token, user)
  F->>B: navega para /app
```

## Login

```mermaid
sequenceDiagram
  participant B as Browser
  participant F as Frontend
  participant API as Backend (Go)
  participant DB as PostgreSQL

  B->>F: preenche email, senha
  F->>F: valida com Zod (loginSchema)
  F->>API: POST /api/v1/auth/login {email, password}
  API->>API: valida payload

  rect rgb(240, 248, 255)
    note right of API: BEGIN TRANSACTION
    API->>DB: SELECT user por email
    alt usuario nao encontrado
      API-->>F: 401 {error: {code: "invalid_credentials"}}
    else usuario encontrado
      API->>API: verifica senha (Argon2id)
      alt senha incorreta
        API-->>F: 401 {error: {code: "invalid_credentials"}}
      else senha correta
        API->>API: gera access token (PASETO v4)
        API->>API: gera refresh token + hash
        API->>DB: INSERT refresh_token
      end
    end
    note right of API: COMMIT
  end

  API->>DB: INSERT audit_log (async)
  API-->>F: 200 {data: {accessToken, user}} + Set-Cookie refresh_token
  F->>F: authStore.setSession(token, user)
  F->>B: navega para /app
```

## Users/Me

```mermaid
sequenceDiagram
  participant F as Frontend
  participant API as Backend (Go)
  participant DB as PostgreSQL

  F->>API: GET /api/v1/users/me (Authorization: Bearer PASETO)
  API->>API: middleware extrai e valida token
  alt token invalido ou ausente
    API-->>F: 401 {error: {code: "unauthorized"}}
  else token valido
    API->>DB: SELECT user por ID (do token)
    API-->>F: 200 {data: {id, name, email, role, status}}
  end
```

## Logout

```mermaid
sequenceDiagram
  participant F as Frontend
  participant API as Backend (Go)
  participant DB as PostgreSQL

  F->>API: POST /api/v1/auth/logout (Authorization: Bearer + Cookie refresh_token)
  API->>API: middleware valida access token

  rect rgb(240, 248, 255)
    note right of API: BEGIN TRANSACTION
    API->>API: hash do refresh token (SHA-256)
    API->>DB: UPDATE refresh_token SET revoked_at = now()
    note right of API: COMMIT
  end

  API->>DB: INSERT audit_log (async)
  API-->>F: 204 No Content
  F->>F: authStore.clearSession()
  F->>F: navega para /auth/login
```

## Forgot Password

```mermaid
sequenceDiagram
  participant B as Browser
  participant F as Frontend
  participant API as Backend (Go)
  participant DB as PostgreSQL
  participant Mail as Mailtrap SMTP

  B->>F: informa email
  F->>API: POST /api/v1/auth/forgot-password {email}
  
  rect rgb(240, 248, 255)
    note right of API: BEGIN TRANSACTION
    API->>DB: SELECT user por email
    alt usuario encontrado
      API->>DB: UPDATE password_reset_tokens (revoke active ones)
      API->>API: gera random reset token
      API->>API: hash do token (SHA-256)
      API->>DB: INSERT password_reset_token
      API-->>Mail: SMTP link de redefinicao (async)
    end
    note right of API: COMMIT
  end

  API-->>F: 200 {message: "If the account exists..."}
  F->>B: exibe confirmacao
```

## Reset Password

```mermaid
sequenceDiagram
  participant B as Browser
  participant F as Frontend
  participant API as Backend (Go)
  participant DB as PostgreSQL

  B->>F: informa nova senha + token da URL
  F->>API: POST /api/v1/auth/reset-password {token, password}
  
  rect rgb(240, 248, 255)
    note right of API: BEGIN TRANSACTION
    API->>API: hash do token recebido
    API->>DB: SELECT token hash (not used & not expired)
    alt token invalido ou expirado
      API-->>F: 400 {error: {code: "invalid_reset_token"}}
    else token valido
      API->>API: hash da nova senha (Argon2id)
      API->>DB: UPDATE user password
      API->>DB: UPDATE password_reset_token (set used_at)
      API->>DB: UPDATE refresh_tokens (revoke all for user)
    end
    note right of API: COMMIT
  end

  API->>DB: INSERT audit_log (async)
  API-->>F: 200 {message: "Password changed successfully"}
  F->>B: navega para /auth/login
```

## Change Password (Authenticated)

```mermaid
sequenceDiagram
  participant B as Browser
  participant F as Frontend
  participant API as Backend (Go)
  participant DB as PostgreSQL

  B->>F: informa senha atual e nova senha
  F->>API: POST /api/v1/auth/change-password {currentPassword, newPassword}
  
  rect rgb(240, 248, 255)
    note right of API: BEGIN TRANSACTION
    API->>DB: SELECT user (ID do token)
    API->>API: verifica senha atual (Argon2id)
    alt senha incorreta
      API-->>F: 400 {error: {code: "invalid_password"}}
    else senha correta
      API->>API: hash da nova senha
      API->>DB: UPDATE user password
      API->>DB: UPDATE refresh_tokens (revoke all for user)
    end
    note right of API: COMMIT
  end

  API->>DB: INSERT audit_log (async)
  API-->>F: 200 {message: "Password changed successfully"}
  F->>F: authStore.clearSession()
  F->>B: navega para /auth/login
```

## Observacao Atual

- A baseline atual persiste `accessToken` e `currentUser` no frontend.
- Nao existe endpoint dedicado de refresh de sessao proativa no frontend nesta baseline.
- Quando o `accessToken` expira, o frontend limpa a sessao e exige novo login.
- Toda emissao de `forgot-password` invalida automaticamente tokens de reset anteriores para o mesmo usuario.
- Toda troca de senha (reset ou change) revoga todos os `refresh_tokens` ativos do usuario.
