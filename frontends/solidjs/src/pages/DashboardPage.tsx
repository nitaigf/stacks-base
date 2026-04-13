import { createSignal } from 'solid-js';
import { me } from '../services/auth';
import { isApiClientError } from '../services/api';
import { authStore } from '../stores/auth';

type DashboardPageProps = {
  onFatalError: () => void;
};

export function DashboardPage(props: DashboardPageProps) {
  const [feedback, setFeedback] = createSignal<string | null>(null);

  const currentUser = () => authStore.currentUser();
  const isSuccessfulFeedback = () => feedback() === 'Perfil sincronizado com sucesso.';

  const refreshUser = async () => {
    try {
      const response = await me();
      authStore.setCurrentUser(response.data);
      setFeedback('Perfil sincronizado com sucesso.');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao carregar perfil.';
      setFeedback(message);
      if (isApiClientError(error) && error.isServerError) {
        props.onFatalError();
      }
    }
  };

  return (
    <div class="content-stack">
      <section class="surface-card dashboard-card">
        <div class="page-header">
          <div>
            <span class="metric-pill metric-pill-tag">Painel protegido</span>
            <h2 class="form-title">Resumo do workspace</h2>
            <p class="form-copy">A sessao autenticada agora vive em um shell mais proximo de um dashboard de produto.</p>
          </div>
          <button class="button button-primary" type="button" onClick={refreshUser}>
            Recarregar perfil
          </button>
        </div>

        <div class="summary-grid">
          <article class="stat-card">
            <span class="stat-label">Usuario atual</span>
            <p class="stat-value">{currentUser()?.name ?? 'Nao carregado'}</p>
            <p class="stat-copy">Identidade usada para navegar nas areas privadas e administrativas.</p>
          </article>
          <article class="stat-card">
            <span class="stat-label">Papel</span>
            <p class="stat-value">{currentUser()?.role === 'admin' ? 'Admin' : 'Member'}</p>
            <p class="stat-copy">Permissao aplicada pelos guards do frontend e pelo backend.</p>
          </article>
          <article class="stat-card">
            <span class="stat-label">Status</span>
            <p class="stat-value">{currentUser()?.status === 'active' ? 'Ativo' : 'Inativo'}</p>
            <p class="stat-copy">Estado operacional da conta sincronizado com a API.</p>
          </article>
          <article class="stat-card">
            <span class="stat-label">Area admin</span>
            <p class="stat-value">{currentUser()?.role === 'admin' ? 'Liberada' : 'Restrita'}</p>
            <p class="stat-copy">Atalho direto para usuarios e auditoria quando o perfil permite.</p>
          </article>
        </div>

        {feedback() ? <p class={`feedback ${isSuccessfulFeedback() ? 'feedback-success' : 'feedback-error'}`}>{feedback()}</p> : null}
      </section>

      <section class="surface-card dashboard-card">
        <div class="panel-heading">
          <div>
            <h3 class="section-title">Sessao e contexto</h3>
            <p class="section-copy">Detalhes do usuario autenticado usados para compor o workspace atual.</p>
          </div>
        </div>

        <div class="details-grid">
          <dl class="detail-block">
            <dt>E-mail</dt>
            <dd>{currentUser()?.email ?? '—'}</dd>
          </dl>
          <dl class="detail-block">
            <dt>Perfil</dt>
            <dd>{currentUser()?.role ?? '—'}</dd>
          </dl>
          <dl class="detail-block">
            <dt>Status da conta</dt>
            <dd>{currentUser()?.status ?? '—'}</dd>
          </dl>
          <dl class="detail-block">
            <dt>Persistencia local</dt>
            <dd>{authStore.accessToken() ? 'Sessao reidratavel' : 'Sem token local'}</dd>
          </dl>
        </div>
      </section>
    </div>
  );
}
