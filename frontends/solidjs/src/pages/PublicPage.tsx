type PublicPageProps = {
  onLogin: () => void;
  onRegister: () => void;
};

export function PublicPage(props: PublicPageProps) {
  return (
    <section class="surface-card dashboard-card public-card">
      <span class="eyebrow">Pagina publica</span>
      <h2 class="form-title">Entrada publica da baseline de referencia</h2>
      <p class="form-copy">
        Esta area existe para fixar a separacao entre pagina publica, paginas de autenticacao, paginas privadas, paginas administrativas e paginas de erro.
      </p>
      <div class="metric-row metric-row-tags">
        <span class="metric-pill metric-pill-tag">publico</span>
        <span class="metric-pill metric-pill-tag">auth</span>
        <span class="metric-pill metric-pill-tag">privado</span>
        <span class="metric-pill metric-pill-tag">admin</span>
        <span class="metric-pill metric-pill-tag">errors</span>
      </div>
      <div class="button-row button-row-compact">
        <button class="button button-primary" type="button" onClick={props.onRegister}>
          Criar conta
        </button>
        <button class="button button-secondary" type="button" onClick={props.onLogin}>
          Entrar
        </button>
      </div>
    </section>
  );
}