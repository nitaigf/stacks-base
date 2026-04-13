import type { JSX } from 'solid-js';

type PublicLayoutProps = {
  children: JSX.Element;
};

export function PublicLayout(props: PublicLayoutProps) {
  return (
    <div class="marketing-shell">
      <section class="marketing-panel">
        <div class="marketing-topbar">
          <span class="eyebrow">Stacks Base</span>
          <span class="marketing-topbar-link">Design system orientado a produto</span>
        </div>
        <h1 class="hero-title">Uma baseline executavel com visual mais proximo de uma plataforma real.</h1>
        <p class="hero-copy">
          Interface publica inspirada pela homepage do shadcn, com hero forte, superfícies limpas e foco em componentes,
          contrato e consistencia entre as areas do sistema.
        </p>
        <ul class="stack-list">
          <li>
            <span>Frontend</span>
            <strong>SolidJS + TanStack Router</strong>
          </li>
          <li>
            <span>Backend</span>
            <strong>Go net/http</strong>
          </li>
          <li>
            <span>Contrato</span>
            <strong>OpenAPI + Bruno</strong>
          </li>
          <li>
            <span>Dados</span>
            <strong>PostgreSQL local</strong>
          </li>
        </ul>
      </section>
      <div class="marketing-main">{props.children}</div>
    </div>
  );
}
