import type { JSX } from 'solid-js';

type AuthLayoutProps = {
  children: JSX.Element;
  onHome: () => void;
};

export function AuthLayout(props: AuthLayoutProps) {
  return (
    <div class="layout-auth">
      <nav class="layout-nav">
        <button class="button-link" type="button" onClick={props.onHome}>
          ← Voltar para a home
        </button>
      </nav>
      {props.children}
    </div>
  );
}
