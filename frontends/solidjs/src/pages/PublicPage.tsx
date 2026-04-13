type PublicPageProps = {
  onLogin: () => void;
  onRegister: () => void;
};

export function PublicPage(props: PublicPageProps) {
  return (
    <section class="surface-card dashboard-card public-card public-home">
      <div class="page-header">
        <div>
          <span class="metric-pill metric-pill-tag">Pagina publica</span>
          <h2 class="form-title">Fundacao visual e funcional da baseline</h2>
          <p class="form-copy">
            A home agora funciona como vitrine do sistema, aproximando a linguagem da landing principal do shadcn.
          </p>
        </div>
      </div>
      <div class="summary-grid">
        <article class="stat-card">
          <span class="stat-label">Componentes vivos</span>
          <p class="stat-value">UI</p>
          <p class="stat-copy">Design system compartilhado e aplicado do publico ao admin.</p>
        </article>
        <article class="stat-card">
          <span class="stat-label">Contrato unico</span>
          <p class="stat-value">API</p>
          <p class="stat-copy">OpenAPI e Bruno versionados como fonte executavel de verdade.</p>
        </article>
        <article class="stat-card">
          <span class="stat-label">Dados reais</span>
          <p class="stat-value">DB</p>
          <p class="stat-copy">PostgreSQL local, sem mocks nas areas administrativas.</p>
        </article>
      </div>
      <div class="button-row button-row-compact public-actions">
        <button class="button button-primary" type="button" onClick={props.onRegister}>
          Criar conta
        </button>
        <button class="button button-secondary" type="button" onClick={props.onLogin}>
          Entrar
        </button>
      </div>
      <div class="metric-row metric-row-tags">
        <span class="metric-pill metric-pill-tag">publico</span>
        <span class="metric-pill metric-pill-tag">auth</span>
        <span class="metric-pill metric-pill-tag">privado</span>
        <span class="metric-pill metric-pill-tag">admin</span>
        <span class="metric-pill metric-pill-tag">errors</span>
      </div>
    </section>
  );
}
