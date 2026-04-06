type PublicPageProps = {
  onLogin: () => void;
  onRegister: () => void;
};

export function PublicPage(props: PublicPageProps) {
  return (
    <section class="surface-card dashboard-card">
      <span class="eyebrow">Pagina publica</span>
      <h2 class="form-title">Entrada publica da baseline de referencia</h2>
      <p class="form-copy">
        Esta area existe para fixar a separacao entre pagina publica, paginas de autenticacao, paginas privadas, paginas administrativas e paginas de erro.
      </p>
      <div class="metric-row">
        <span class="metric-pill">publico</span>
        <span class="metric-pill">auth</span>
        <span class="metric-pill">privado</span>
        <span class="metric-pill">admin</span>
        <span class="metric-pill">errors</span>
      </div>
      <div class="button-row">
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