const contentByStatus = {
  403: {
    title: 'Acesso negado',
    copy: 'Voce chegou a uma rota protegida sem permissao suficiente.',
  },
  404: {
    title: 'Pagina nao encontrada',
    copy: 'A rota solicitada nao existe na baseline atual.',
  },
  500: {
    title: 'Erro interno da interface',
    copy: 'A aplicacao encontrou um estado inesperado e redirecionou para uma tela de erro previsivel.',
  },
} as const;

type ErrorPageProps = {
  statusCode: 403 | 404 | 500;
  onHome: () => void;
  onLogin: () => void;
};

export function ErrorPage(props: ErrorPageProps) {
  const content = contentByStatus[props.statusCode];

  return (
    <section class="surface-card error-card form-card error-preview-card">
      <span class="metric-pill metric-pill-tag">Erro {props.statusCode}</span>
      <h2 class="error-code">{props.statusCode}</h2>
      <h3 class="form-title">{content.title}</h3>
      <p class="form-copy">{content.copy}</p>
      <div class="button-row">
        <button class="button button-primary" type="button" onClick={props.onHome}>
          Voltar para a home
        </button>
        <button class="button button-secondary" type="button" onClick={props.onLogin}>
          Ir para login
        </button>
      </div>
    </section>
  );
}
