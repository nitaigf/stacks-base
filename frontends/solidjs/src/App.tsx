import { Outlet } from '@tanstack/solid-router';

export default function App() {
  return (
    <main class="app-shell">
      <div class="page-grid">
        <section class="hero-panel">
          <span class="eyebrow">Stacks Base</span>
          <h1 class="hero-title">SolidJS e Go seguem como a referencia executavel do projeto.</h1>
          <p class="hero-copy">
            Esta baseline agora concentra autenticacao, gestao real de usuarios, auditoria e exportacoes sobre
            PostgreSQL local com contrato unico do sistema.
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
              <span>Banco</span>
              <strong>PostgreSQL local 5432</strong>
            </li>
            <li>
              <span>Seeds</span>
              <strong>Admin + dados de exemplo</strong>
            </li>
          </ul>
        </section>
        <Outlet />
      </div>
    </main>
  );
}
