import { createSignal } from 'solid-js';
import { changePassword } from '../services/auth';
import { FormField } from '../components/FormField';

type ChangePasswordPageProps = {
  onCompleted: () => void;
  onBack: () => void;
  onFatalError: () => void;
};

export function ChangePasswordPage(props: ChangePasswordPageProps) {
  const [currentPassword, setCurrentPassword] = createSignal('');
  const [newPassword, setNewPassword] = createSignal('');
  const [feedback, setFeedback] = createSignal<string | null>(null);
  const [submitting, setSubmitting] = createSignal(false);

  const submit = async (event: SubmitEvent) => {
    event.preventDefault();
    setSubmitting(true);
    setFeedback(null);

    try {
      const response = await changePassword(currentPassword(), newPassword());
      setFeedback(response.data.message);
      props.onCompleted();
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao alterar senha.';
      setFeedback(message);
      if (message.toLowerCase().includes('internal')) {
        props.onFatalError();
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <section class="surface-card">
      <span class="metric-pill metric-pill-tag">Seguranca da conta</span>
      <h2 class="form-title">Alterar senha</h2>
      <p class="form-copy">Apos a troca, a sessao atual e encerrada e um novo login sera necessario.</p>

      <form onSubmit={submit}>
        <FormField
          label="Senha atual"
          name="currentPassword"
          type="password"
          value={currentPassword()}
          onInput={(event) => setCurrentPassword(event.currentTarget.value)}
        />
        <FormField
          label="Nova senha"
          name="newPassword"
          type="password"
          value={newPassword()}
          onInput={(event) => setNewPassword(event.currentTarget.value)}
        />

        <div class="button-row">
          <button class="button button-primary" type="submit" disabled={submitting()}>
            {submitting() ? 'Atualizando...' : 'Salvar nova senha'}
          </button>
          <button class="button button-secondary" type="button" onClick={props.onBack}>
            Voltar
          </button>
        </div>

        {feedback() ? <p class="feedback feedback-success">{feedback()}</p> : null}
      </form>
    </section>
  );
}
