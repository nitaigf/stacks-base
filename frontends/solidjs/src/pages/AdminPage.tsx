type AdminPageProps = {
  onBackToApp: () => void;
};

export function AdminPage(props: AdminPageProps) {
  return (
    <section class="surface-card dashboard-card">
      <span class="eyebrow">Pagina administrativa</span>
      <h2 class="form-title">Area administrativa da referencia</h2>
      <p class="form-copy">
        A baseline precisa ter espaco explicito para administracao, mesmo quando a primeira conta local ainda nao for admin.
      </p>
      <div class="metric-row">
        <span class="metric-pill">admin only</span>
        <span class="metric-pill">403 previsivel</span>
      </div>
      <div class="button-row">
        <button class="button button-primary" type="button" onClick={props.onBackToApp}>
          Voltar ao app
        </button>
      </div>
    </section>
  );
}