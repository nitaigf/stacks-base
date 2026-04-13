import { For, createEffect, createSignal } from 'solid-js';
import { listAuditLogs } from '../services/audit';
import { isApiClientError } from '../services/api';
import type { AuditLog } from '../types/auth';

type AuditLogsPageProps = {
  onFatalError: () => void;
};

export function AuditLogsPage(props: AuditLogsPageProps) {
  const [logs, setLogs] = createSignal<AuditLog[]>([]);
  const [page, setPage] = createSignal(1);
  const [query, setQuery] = createSignal('');
  const [action, setAction] = createSignal('');
  const [resource, setResource] = createSignal('');
  const [loading, setLoading] = createSignal(true);
  const [feedback, setFeedback] = createSignal<string | null>(null);
  const [meta, setMeta] = createSignal({
    page: 1,
    perPage: 20,
    total: 0,
    totalPages: 0,
  });
  const feedbackVariant = () =>
    feedback() && (feedback()!.toLowerCase().includes('falha') || feedback()!.toLowerCase().includes('invalid'))
      ? 'feedback-error'
      : 'feedback-success';
  const distinctActions = () => new Set(logs().map((entry) => entry.action)).size;
  const distinctActors = () => new Set(logs().map((entry) => entry.actorEmail ?? entry.actorName ?? 'anonimo')).size;

  const load = async () => {
    setLoading(true);
    try {
      const response = await listAuditLogs({
        page: page(),
        perPage: 20,
        query: query(),
        action: action(),
        resource: resource(),
      });
      setLogs(response.data);
      setMeta(response.meta);
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao carregar auditoria.';
      setFeedback(message);
      if (isApiClientError(error) && error.isServerError) {
        props.onFatalError();
      }
    } finally {
      setLoading(false);
    }
  };

  createEffect(() => {
    page();
    void load();
  });

  return (
    <div class="content-stack">
      <section class="surface-card dashboard-card">
        <div class="page-header">
          <div>
            <span class="metric-pill metric-pill-tag">Auditoria real</span>
            <h2 class="form-title">Trilha observavel do sistema</h2>
            <p class="form-copy">Eventos, atores e rotas em um painel de consulta mais proximo de um dashboard operacional.</p>
          </div>
        </div>

        <div class="summary-grid">
          <article class="stat-card">
            <span class="stat-label">Total encontrado</span>
            <p class="stat-value">{meta().total}</p>
            <p class="stat-copy">Eventos compatíveis com os filtros atuais.</p>
          </article>
          <article class="stat-card">
            <span class="stat-label">Linhas na pagina</span>
            <p class="stat-value">{logs().length}</p>
            <p class="stat-copy">Recorte imediato exibido no painel.</p>
          </article>
          <article class="stat-card">
            <span class="stat-label">Acoes distintas</span>
            <p class="stat-value">{distinctActions()}</p>
            <p class="stat-copy">Diversidade de operacoes observadas na pagina atual.</p>
          </article>
          <article class="stat-card">
            <span class="stat-label">Atores distintos</span>
            <p class="stat-value">{distinctActors()}</p>
            <p class="stat-copy">Usuarios ou agentes que aparecem no recorte carregado.</p>
          </article>
        </div>
      </section>

      <section class="surface-card dashboard-card">
        <div class="panel-heading">
          <div>
            <h3 class="section-title">Filtros de investigacao</h3>
            <p class="section-copy">Afine por acao, recurso e texto livre para localizar eventos com mais rapidez.</p>
          </div>
          <span class="metric-pill metric-pill-tag">
            Pagina {meta().page} de {meta().totalPages || 1}
          </span>
        </div>

        <div class="field-grid">
          <label class="field-group">
            <span class="field-label">Busca</span>
            <input class="field-input" value={query()} onInput={(event) => setQuery(event.currentTarget.value)} />
          </label>
          <label class="field-group">
            <span class="field-label">Acao</span>
            <input class="field-input" value={action()} onInput={(event) => setAction(event.currentTarget.value)} />
          </label>
          <label class="field-group">
            <span class="field-label">Recurso</span>
            <input class="field-input" value={resource()} onInput={(event) => setResource(event.currentTarget.value)} />
          </label>
        </div>

        <div class="button-row">
          <button
            class="button button-primary"
            type="button"
            onClick={() => {
              setPage(1);
              void load();
            }}
          >
            Aplicar filtros
          </button>
        </div>

        {feedback() ? <p class={`feedback ${feedbackVariant()}`}>{feedback()}</p> : null}
      </section>

      <section class="surface-card dashboard-card">
        <div class="panel-heading">
          <div>
            <h3 class="section-title">Eventos carregados</h3>
            <p class="section-copy">Consulta paginada com rota, metodo, IP, navegador, ator e recurso.</p>
          </div>
        </div>

        {loading() ? (
          <p class="form-copy">Carregando auditoria...</p>
        ) : (
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Quando</th>
                  <th>Quem</th>
                  <th>Acao</th>
                  <th>Recurso</th>
                  <th>Rota</th>
                  <th>IP</th>
                  <th>Navegador</th>
                </tr>
              </thead>
              <tbody>
                <For each={logs()}>
                  {(entry) => (
                    <tr>
                      <td>{formatDate(entry.createdAt)}</td>
                      <td>{entry.actorName || entry.actorEmail || 'anonimo'}</td>
                      <td>{entry.action}</td>
                      <td>{entry.resource}</td>
                      <td>{entry.method} {entry.route}</td>
                      <td>{entry.ipAddress || '—'}</td>
                      <td>{entry.userAgent || '—'}</td>
                    </tr>
                  )}
                </For>
              </tbody>
            </table>
          </div>
        )}

        <div class="button-row button-row-compact">
          <button class="button button-secondary" type="button" disabled={meta().page <= 1} onClick={() => setPage(page() - 1)}>
            Pagina anterior
          </button>
          <button
            class="button button-secondary"
            type="button"
            disabled={meta().page >= meta().totalPages}
            onClick={() => setPage(page() + 1)}
          >
            Proxima pagina
          </button>
        </div>
      </section>
    </div>
  );
}

function formatDate(value: string) {
  return new Date(value).toLocaleString('pt-BR');
}
