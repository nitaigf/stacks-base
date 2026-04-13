import { createSignal, type Setter } from 'solid-js';
import { forgotPassword } from '../services/auth';
import { isApiClientError } from '../services/api';
import { FormField } from '../components/FormField';
import { forgotPasswordSchema } from '../schemas/auth';

type ForgotPasswordPageProps = {
  onBackToLogin: () => void;
  onFatalError: () => void;
};

export function ForgotPasswordPage(props: ForgotPasswordPageProps) {
  const [email, setEmail] = createSignal('');
  const [errors, setErrors] = createSignal<Record<string, string[] | undefined>>({});
  const [feedback, setFeedback] = createSignal<string | null>(null);
  const [feedbackTone, setFeedbackTone] = createSignal<'success' | 'error' | null>(null);
  const [submitting, setSubmitting] = createSignal(false);

  const submit = async (event: SubmitEvent) => {
    event.preventDefault();
    setFeedback(null);
    setFeedbackTone(null);

    const result = forgotPasswordSchema.safeParse({
      email: email(),
    });

    if (!result.success) {
      setErrors(result.error.flatten().fieldErrors);
      return;
    }

    setErrors({});
    setSubmitting(true);

    try {
      const response = await forgotPassword(result.data.email);
      setFeedback(response.data.message);
      setFeedbackTone('success');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao solicitar recuperacao de acesso.';
      setFeedback(message);
      setFeedbackTone('error');
      if (isApiClientError(error) && error.isServerError) {
        props.onFatalError();
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <section class="surface-card form-card auth-card">
      <div class="auth-card-copy">
        <span class="metric-pill metric-pill-tag">Recuperacao de acesso</span>
        <h2 class="form-title">Solicitar redefinicao de senha</h2>
        <p class="form-copy">Informe o e-mail da conta para receber o link real de recuperacao.</p>
      </div>

      <form noValidate onSubmit={submit}>
        <FormField
          label="E-mail"
          name="email"
          type="email"
          value={email()}
          error={errors().email?.[0]}
          onInput={(event) => {
            setEmail(event.currentTarget.value);
            clearFieldError(setErrors, 'email');
          }}
        />

        <div class="button-row">
          <button class="button button-primary" type="submit" disabled={submitting()}>
            {submitting() ? 'Enviando...' : 'Enviar link'}
          </button>
          <button class="button button-secondary" type="button" onClick={props.onBackToLogin}>
            Voltar ao login
          </button>
        </div>

        {feedback() ? <p class={`feedback ${feedbackTone() === 'success' ? 'feedback-success' : 'feedback-error'}`}>{feedback()}</p> : null}
      </form>
      <p class="auth-card-footer">O envio usa o fluxo real do backend, sem etapa simulada.</p>
    </section>
  );
}

function clearFieldError(
  setErrors: Setter<Record<string, string[] | undefined>>,
  fieldName: string,
) {
  setErrors((current) => {
    if (!current[fieldName]) {
      return current;
    }

    const next = { ...current };
    delete next[fieldName];
    return next;
  });
}
