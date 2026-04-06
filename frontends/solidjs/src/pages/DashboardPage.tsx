import { createSignal } from 'solid-js';
import { logout, me } from '../services/auth';
import { authStore } from '../stores/auth';

type DashboardPageProps = {
  onLoggedOut: () => void;
  onAdmin: () => void;
  onFatalError: () => void;
};

export function DashboardPage(props: DashboardPageProps) {
  const [feedback, setFeedback] = createSignal<string | null>(null);

  const currentUser = () => authStore.currentUser();

  const refreshUser = async () => {
    try {
      const response = await me();
      authStore.setSession(authStore.accessToken() ?? '', response.data);
      setFeedback('Perfil sincronizado com sucesso.');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao carregar perfil.';
      setFeedback(message);
      if (message.toLowerCase().includes('internal')) {
        props.onFatalError();
      }
    }
  };

  const endSession = async () => {
    try {
      await logout();
      authStore.clearSession();
      setFeedback('Sessao encerrada.');
      props.onLoggedOut();
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao encerrar sessao.';
      setFeedback(message);
      if (message.toLowerCase().includes('internal')) {
        props.onFatalError();
      }
    }
  };

  return (
    <section class="surface-card dashboard-card">
      <span class="eyebrow">Painel protegido</span>
      <h2 class="form-title">{currentUser()?.name}</h2>
      <p class="form-copy">Esta tela existe para validar o corte vertical inicial e a protecao do fluxo autenticado.</p>
      <div class="metric-row">
        <span class="metric-pill">{currentUser()?.email}</span>
        <span class="metric-pill">role: {currentUser()?.role}</span>
        <span class="metric-pill">status: {currentUser()?.status}</span>
      </div>
      <div class="button-row">
        <button class="button button-primary" type="button" onClick={refreshUser}>
          Recarregar perfil
        </button>
        <button class="button button-secondary" type="button" onClick={props.onAdmin}>
          Area admin
        </button>
        <button class="button button-secondary" type="button" onClick={endSession}>
          Logout
        </button>
      </div>
      {feedback() ? <p class="feedback feedback-success">{feedback()}</p> : null}
    </section>
  );
}