import { For, createEffect, createSignal } from 'solid-js';
import { listAuditLogs } from '../services/audit';
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
      if (message.toLowerCase().includes('internal')) {
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
    <section class="surface-card dashboard-card">
      <div class="page-header">
        <span class="metric-pill metric-pill-tag">Auditoria real</span>
        <h2 class="form-title">Logs auditaveis do sistema</h2>
        <p class="form-copy">Consulta paginada com rota, metodo, IP, navegador, ator e metadata.</p>
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

      <div class="button-row button-row-compact">
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

      {feedback() ? <p class="feedback feedback-success">{feedback()}</p> : null}

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
        <span class="metric-pill metric-pill-tag">
          Pagina {meta().page} de {meta().totalPages || 1}
        </span>
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
  );
}

function formatDate(value: string) {
  return new Date(value).toLocaleString('pt-BR');
}
