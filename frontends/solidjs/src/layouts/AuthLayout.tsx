import type { JSX } from 'solid-js';

type AuthLayoutProps = {
  children: JSX.Element;
  onHome: () => void;
};

export function AuthLayout(props: AuthLayoutProps) {
  return (
    <div class="marketing-shell">
      <section class="marketing-panel marketing-panel-auth">
        <div class="marketing-topbar">
          <span class="eyebrow">Acesso seguro</span>
          <span class="marketing-topbar-link">Fluxo inspirado no auth example</span>
        </div>
        <h1 class="hero-title">Entrar, recuperar ou redefinir sem sair de um layout mais editorial e enxuto.</h1>
        <p class="hero-copy">
          O fluxo de entrada agora conversa melhor com o workspace interno, com uma coluna visual forte e um card de
          autenticacao mais direto, no estilo das referências do shadcn.
        </p>
        <ul class="stack-list">
          <li>
            <span>Login</span>
            <strong>Sessao com token real</strong>
          </li>
          <li>
            <span>Conta</span>
            <strong>Registro sem mocks</strong>
          </li>
          <li>
            <span>Senha</span>
            <strong>Forgot, reset e change</strong>
          </li>
        </ul>
      </section>
      <div class="marketing-main">
        <nav class="marketing-nav">
          <button class="button-link" type="button" onClick={props.onHome}>
            Voltar para a home
          </button>
        </nav>
        {props.children}
      </div>
    </div>
  );
}
