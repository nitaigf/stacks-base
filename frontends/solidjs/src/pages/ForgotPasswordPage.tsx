import { createSignal } from 'solid-js';
import { forgotPassword } from '../services/auth';
import { FormField } from '../components/FormField';

type ForgotPasswordPageProps = {
  onBackToLogin: () => void;
  onFatalError: () => void;
};

export function ForgotPasswordPage(props: ForgotPasswordPageProps) {
  const [email, setEmail] = createSignal('');
  const [feedback, setFeedback] = createSignal<string | null>(null);
  const [submitting, setSubmitting] = createSignal(false);

  const submit = async (event: SubmitEvent) => {
    event.preventDefault();
    setSubmitting(true);
    setFeedback(null);

    try {
      const response = await forgotPassword(email());
      setFeedback(response.data.message);
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao solicitar recuperacao de acesso.';
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
      <span class="metric-pill metric-pill-tag">Recuperacao de acesso</span>
      <h2 class="form-title">Solicitar redefinicao de senha</h2>
      <p class="form-copy">Informe o e-mail da conta para receber o link real de recuperacao.</p>

      <form onSubmit={submit}>
        <FormField
          label="E-mail"
          name="email"
          type="email"
          value={email()}
          onInput={(event) => setEmail(event.currentTarget.value)}
        />

        <div class="button-row">
          <button class="button button-primary" type="submit" disabled={submitting()}>
            {submitting() ? 'Enviando...' : 'Enviar link'}
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
