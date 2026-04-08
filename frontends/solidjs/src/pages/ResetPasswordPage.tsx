import { createSignal } from 'solid-js';
import { resetPassword } from '../services/auth';
import { FormField } from '../components/FormField';

type ResetPasswordPageProps = {
  token: string;
  onBackToLogin: () => void;
  onFatalError: () => void;
};

export function ResetPasswordPage(props: ResetPasswordPageProps) {
  const [token, setToken] = createSignal(props.token);
  const [newPassword, setNewPassword] = createSignal('');
  const [feedback, setFeedback] = createSignal<string | null>(null);
  const [submitting, setSubmitting] = createSignal(false);

  const submit = async (event: SubmitEvent) => {
    event.preventDefault();
    setSubmitting(true);
    setFeedback(null);

    try {
      const response = await resetPassword(token(), newPassword());
      setFeedback(response.data.message);
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao redefinir senha.';
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
      <span class="metric-pill metric-pill-tag">Nova senha</span>
      <h2 class="form-title">Redefinir senha</h2>
      <p class="form-copy">Use o token enviado por e-mail para definir uma nova senha.</p>

      <form onSubmit={submit}>
        <FormField
          label="Token"
          name="token"
          value={token()}
          onInput={(event) => setToken(event.currentTarget.value)}
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
            {submitting() ? 'Atualizando...' : 'Redefinir senha'}
          </button>
          <button class="button button-secondary" type="button" onClick={props.onBackToLogin}>
            Voltar ao login
          </button>
        </div>

        {feedback() ? <p class="feedback feedback-success">{feedback()}</p> : null}
      </form>
    </section>
  );
}
