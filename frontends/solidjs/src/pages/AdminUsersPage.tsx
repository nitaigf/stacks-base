import { For, createEffect, createSignal } from 'solid-js';
import {
  deactivateUser,
  exportUsersCsv,
  exportUsersXlsx,
  hardDeleteUser,
  listUsers,
  printUsers,
  reactivateUser,
  restoreUser,
  softDeleteUser,
} from '../services/users';
import type { User, UserListFilters, UsersEnvelope } from '../types/auth';

type AdminUsersPageProps = {
  onCreateUser: () => void;
  onViewUser: (userId: string) => void;
  onFatalError: () => void;
};

export function AdminUsersPage(props: AdminUsersPageProps) {
  const [users, setUsers] = createSignal<User[]>([]);
  const [page, setPage] = createSignal(1);
  const [perPage] = createSignal(10);
  const [query, setQuery] = createSignal('');
  const [role, setRole] = createSignal('');
  const [status, setStatus] = createSignal('');
  const [includeDeleted, setIncludeDeleted] = createSignal(false);
  const [meta, setMeta] = createSignal<UsersEnvelope['meta']>({
    page: 1,
    perPage: 10,
    total: 0,
    totalPages: 0,
  });
  const [loading, setLoading] = createSignal(true);
  const [feedback, setFeedback] = createSignal<string | null>(null);

  const currentFilters = (): UserListFilters => ({
    page: page(),
    perPage: perPage(),
    query: query(),
    role: role() as UserListFilters['role'],
    status: status() as UserListFilters['status'],
    includeDeleted: includeDeleted(),
  });

  const load = async () => {
    setLoading(true);
    try {
      const response = await listUsers(currentFilters());
      setUsers(response.data);
      setMeta(response.meta);
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao carregar usuarios.';
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

  const runAction = async (action: () => Promise<unknown>, successMessage: string) => {
    try {
      await action();
      setFeedback(successMessage);
      await load();
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao processar a acao.';
      setFeedback(message);
      if (message.toLowerCase().includes('internal')) {
        props.onFatalError();
      }
    }
  };

  const handleCsv = async () => {
    try {
      const blob = await exportUsersCsv(currentFilters());
      downloadBlob(blob, 'users.csv');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao exportar CSV.';
      setFeedback(message);
    }
  };

  const handleXlsx = async () => {
    try {
      const blob = await exportUsersXlsx(currentFilters());
      downloadBlob(blob, 'users.xlsx');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao exportar XLSX.';
      setFeedback(message);
    }
  };

  const handlePrint = async () => {
    try {
      const blob = await printUsers(currentFilters());
      downloadBlob(blob, 'users.pdf');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao gerar PDF.';
      setFeedback(message);
    }
  };

  return (
    <section class="surface-card dashboard-card">
      <div class="page-header">
        <span class="metric-pill metric-pill-tag">Usuarios reais do banco</span>
        <h2 class="form-title">Gestao de usuarios</h2>
        <p class="form-copy">Listagem paginada, exportacoes reais e acoes administrativas sem mock.</p>
      </div>

      <div class="field-grid">
        <label class="field-group">
          <span class="field-label">Busca</span>
          <input class="field-input" value={query()} onInput={(event) => setQuery(event.currentTarget.value)} />
        </label>
        <label class="field-group">
          <span class="field-label">Papel</span>
          <select class="field-input" value={role()} onInput={(event) => setRole(event.currentTarget.value)}>
            <option value="">Todos</option>
            <option value="admin">Admin</option>
            <option value="member">Member</option>
          </select>
        </label>
        <label class="field-group">
          <span class="field-label">Status</span>
          <select class="field-input" value={status()} onInput={(event) => setStatus(event.currentTarget.value)}>
            <option value="">Todos</option>
            <option value="active">Ativo</option>
            <option value="inactive">Inativo</option>
          </select>
        </label>
        <label class="field-group">
          <span class="field-label">Excluidos</span>
          <select
            class="field-input"
            value={includeDeleted() ? 'true' : 'false'}
            onInput={(event) => setIncludeDeleted(event.currentTarget.value === 'true')}
          >
            <option value="false">Ocultar</option>
            <option value="true">Mostrar</option>
          </select>
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
        <button class="button button-secondary" type="button" onClick={props.onCreateUser}>
          Novo usuario
        </button>
        <button class="button button-secondary" type="button" onClick={() => void handleCsv()}>
          Exportar CSV
        </button>
        <button class="button button-secondary" type="button" onClick={() => void handleXlsx()}>
          Exportar XLSX
        </button>
        <button class="button button-secondary" type="button" onClick={() => void handlePrint()}>
          Imprimir PDF
        </button>
      </div>

      {feedback() ? <p class="feedback feedback-success">{feedback()}</p> : null}

      {loading() ? (
        <p class="form-copy">Carregando usuarios...</p>
      ) : (
        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>Acoes</th>
                <th>Nome</th>
                <th>Email</th>
                <th>Papel</th>
                <th>Status</th>
                <th>Excluido em</th>
                <th>Ultimo login</th>
              </tr>
            </thead>
            <tbody>
              <For each={users()}>
                {(user) => (
                  <tr>
                    <td>
                      <div class="action-buttons">
                        <button class="button button-sm button-secondary" type="button" title="Visualizar" onClick={() => props.onViewUser(user.id)}>
                          Ver
                        </button>
                        {user.status === 'active' ? (
                          <button
                            class="button button-sm button-secondary"
                            type="button"
                            title="Inativar"
                            onClick={() => void runAction(() => deactivateUser(user.id), 'Usuario inativado com sucesso.')}
                          >
                            Inativar
                          </button>
                        ) : (
                          <button
                            class="button button-sm button-secondary"
                            type="button"
                            title="Reativar"
                            onClick={() => void runAction(() => reactivateUser(user.id), 'Usuario reativado com sucesso.')}
                          >
                            Reativar
                          </button>
                        )}
                        {!user.deletedAt ? (
                          <button
                            class="button button-sm button-secondary"
                            type="button"
                            title="Apagar"
                            onClick={() => void runAction(() => softDeleteUser(user.id), 'Usuario apagado com sucesso.')}
                          >
                            Apagar
                          </button>
                        ) : (
                          <button
                            class="button button-sm button-secondary"
                            type="button"
                            title="Restaurar"
                            onClick={() => void runAction(() => restoreUser(user.id), 'Usuario restaurado com sucesso.')}
                          >
                            Restaurar
                          </button>
                        )}
                        <button
                          class="button button-sm button-secondary"
                          type="button"
                          title="Excluir"
                          onClick={() => {
                            if (window.confirm(`Excluir definitivamente ${user.name}?`)) {
                              void runAction(() => hardDeleteUser(user.id), 'Usuario excluido definitivamente.');
                            }
                          }}
                        >
                          Excluir
                        </button>
                      </div>
                    </td>
                    <td>{user.name}</td>
                    <td>{user.email}</td>
                    <td>{user.role}</td>
                    <td>{user.status}</td>
                    <td>{formatDate(user.deletedAt)}</td>
                    <td>{formatDate(user.lastLoginAt)}</td>
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

function downloadBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement('a');
  anchor.href = url;
  anchor.download = filename;
  anchor.click();
  URL.revokeObjectURL(url);
}

function formatDate(value?: string | null) {
  if (!value) {
    return '—';
  }
  return new Date(value).toLocaleString('pt-BR');
}
