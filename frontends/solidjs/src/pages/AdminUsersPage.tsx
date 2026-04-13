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
import { isApiClientError } from '../services/api';
import { ActionIconButton } from '../components/ActionIconButton';
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

  const activeUsers = () => users().filter((user) => user.status === 'active' && !user.deletedAt).length;
  const adminUsers = () => users().filter((user) => user.role === 'admin').length;
  const deletedUsers = () => users().filter((user) => Boolean(user.deletedAt)).length;
  const feedbackVariant = () =>
    feedback() && (feedback()!.toLowerCase().includes('falha') || feedback()!.toLowerCase().includes('invalid'))
      ? 'feedback-error'
      : 'feedback-success';

  const load = async () => {
    setLoading(true);
    try {
      const response = await listUsers(currentFilters());
      setUsers(response.data);
      setMeta(response.meta);
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao carregar usuarios.';
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

  const runAction = async (action: () => Promise<unknown>, successMessage: string) => {
    try {
      await action();
      setFeedback(successMessage);
      await load();
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao processar a acao.';
      setFeedback(message);
      if (isApiClientError(error) && error.isServerError) {
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
    <div class="content-stack">
      <section class="surface-card dashboard-card">
        <div class="page-header">
          <div>
            <span class="metric-pill metric-pill-tag">Usuarios reais do banco</span>
            <h2 class="form-title">Controle operacional de usuarios</h2>
            <p class="form-copy">Listagem, leitura individual, status administrativos e exportacoes em um painel unico.</p>
          </div>
        </div>

        <div class="summary-grid">
          <article class="stat-card">
            <span class="stat-label">Total encontrado</span>
            <p class="stat-value">{meta().total}</p>
            <p class="stat-copy">Resultado total considerando os filtros atuais.</p>
          </article>
          <article class="stat-card">
            <span class="stat-label">Ativos na pagina</span>
            <p class="stat-value">{activeUsers()}</p>
            <p class="stat-copy">Usuarios ativos visiveis no recorte atual.</p>
          </article>
          <article class="stat-card">
            <span class="stat-label">Admins na pagina</span>
            <p class="stat-value">{adminUsers()}</p>
            <p class="stat-copy">Perfis com acesso a operacoes administrativas.</p>
          </article>
          <article class="stat-card">
            <span class="stat-label">Excluidos na pagina</span>
            <p class="stat-value">{deletedUsers()}</p>
            <p class="stat-copy">Contas ocultas que ainda podem ser restauradas.</p>
          </article>
        </div>
      </section>

      <section class="surface-card dashboard-card">
        <div class="panel-heading">
          <div>
            <h3 class="section-title">Filtros e operacoes</h3>
            <p class="section-copy">Combine busca, estado e exportacoes sem sair da mesma tela.</p>
          </div>
          <div class="button-row button-row-compact">
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
            <h3 class="section-title">Tabela operacional</h3>
            <p class="section-copy">Acoes rapidas na primeira coluna, mantendo o padrao administrativo da baseline.</p>
          </div>
          <span class="metric-pill metric-pill-tag">
            Pagina {meta().page} de {meta().totalPages || 1}
          </span>
        </div>

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
                          <ActionIconButton label="Visualizar usuario" title="Visualizar" onClick={() => props.onViewUser(user.id)}>
                            <ViewIcon />
                          </ActionIconButton>
                          {user.status === 'active' ? (
                            <ActionIconButton
                              label="Inativar usuario"
                              title="Inativar"
                              onClick={() => void runAction(() => deactivateUser(user.id), 'Usuario inativado com sucesso.')}
                            >
                              <PauseIcon />
                            </ActionIconButton>
                          ) : (
                            <ActionIconButton
                              label="Reativar usuario"
                              title="Reativar"
                              onClick={() => void runAction(() => reactivateUser(user.id), 'Usuario reativado com sucesso.')}
                            >
                              <PlayIcon />
                            </ActionIconButton>
                          )}
                          {!user.deletedAt ? (
                            <ActionIconButton
                              label="Apagar usuario"
                              title="Apagar"
                              onClick={() => void runAction(() => softDeleteUser(user.id), 'Usuario apagado com sucesso.')}
                            >
                              <ArchiveIcon />
                            </ActionIconButton>
                          ) : (
                            <ActionIconButton
                              label="Restaurar usuario"
                              title="Restaurar"
                              onClick={() => void runAction(() => restoreUser(user.id), 'Usuario restaurado com sucesso.')}
                            >
                              <RestoreIcon />
                            </ActionIconButton>
                          )}
                          <ActionIconButton
                            label="Excluir usuario definitivamente"
                            title="Excluir"
                            variant="danger"
                            onClick={() => {
                              if (window.confirm(`Excluir definitivamente ${user.name}?`)) {
                                void runAction(() => hardDeleteUser(user.id), 'Usuario excluido definitivamente.');
                              }
                            }}
                          >
                            <TrashIcon />
                          </ActionIconButton>
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

function ViewIcon() {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M2 12s3.5-6 10-6 10 6 10 6-3.5 6-10 6-10-6-10-6Z" fill="none" stroke="currentColor" stroke-width="1.8" />
      <circle cx="12" cy="12" r="3" fill="none" stroke="currentColor" stroke-width="1.8" />
    </svg>
  );
}

function PauseIcon() {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M8 5v14M16 5v14" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
    </svg>
  );
}

function PlayIcon() {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M8 6.5v11l8-5.5-8-5.5Z" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round" />
    </svg>
  );
}

function ArchiveIcon() {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M4 7h16v3H4zM6 10h12v9H6zM10 14h4" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
    </svg>
  );
}

function RestoreIcon() {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M8 8H4v4M4 8a8 8 0 1 1-1 6" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
    </svg>
  );
}

function TrashIcon() {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M4 7h16M9 7V4h6v3M7 7l1 12h8l1-12M10 11v5M14 11v5" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
    </svg>
  );
}
