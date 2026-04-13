import { createEffect, createSignal, on, type Setter } from 'solid-js';
import { resetPassword } from '../services/auth';
import { isApiClientError } from '../services/api';
import { FormField } from '../components/FormField';
import { resetPasswordSchema } from '../schemas/auth';

type ResetPasswordPageProps = {
  token: string;
  onBackToLogin: () => void;
  onFatalError: () => void;
};

export function ResetPasswordPage(props: ResetPasswordPageProps) {
  const [token, setToken] = createSignal(props.token);
  const [newPassword, setNewPassword] = createSignal('');
  const [confirmPassword, setConfirmPassword] = createSignal('');
  const [errors, setErrors] = createSignal<Record<string, string[] | undefined>>({});
  const [feedback, setFeedback] = createSignal<string | null>(null);
  const [feedbackTone, setFeedbackTone] = createSignal<'success' | 'error' | null>(null);
  const [submitting, setSubmitting] = createSignal(false);

  createEffect(
    on(
      () => props.token,
      (value) => {
        setToken(value);
      },
    ),
  );

  const submit = async (event: SubmitEvent) => {
    event.preventDefault();
    setFeedback(null);
    setFeedbackTone(null);

    const result = resetPasswordSchema.safeParse({
      token: token(),
      newPassword: newPassword(),
      confirmPassword: confirmPassword(),
    });

    if (!result.success) {
      setErrors(result.error.flatten().fieldErrors);
      return;
    }

    setErrors({});
    setSubmitting(true);

    try {
      const response = await resetPassword(result.data.token, result.data.newPassword);
      setFeedback(response.data.message);
      setFeedbackTone('success');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao redefinir senha.';
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
        <span class="metric-pill metric-pill-tag">Nova senha</span>
        <h2 class="form-title">Redefinir senha</h2>
        <p class="form-copy">Use o token enviado por e-mail para definir uma nova senha.</p>
      </div>

      <form noValidate onSubmit={submit}>
        <FormField
          label="Token"
          name="token"
          value={token()}
          error={errors().token?.[0]}
          onInput={(event) => {
            setToken(event.currentTarget.value);
            clearFieldError(setErrors, 'token');
          }}
        />
        <FormField
          label="Nova senha"
          name="newPassword"
          type="password"
          value={newPassword()}
          error={errors().newPassword?.[0]}
          onInput={(event) => {
            setNewPassword(event.currentTarget.value);
            clearFieldError(setErrors, 'newPassword');
          }}
        />
        <FormField
          label="Confirmar nova senha"
          name="confirmPassword"
          type="password"
          value={confirmPassword()}
          error={errors().confirmPassword?.[0]}
          onInput={(event) => {
            setConfirmPassword(event.currentTarget.value);
            clearFieldError(setErrors, 'confirmPassword');
          }}
        />

        <div class="button-row">
          <button class="button button-primary" type="submit" disabled={submitting()}>
            {submitting() ? 'Atualizando...' : 'Redefinir senha'}
          </button>
          <button class="button button-secondary" type="button" onClick={props.onBackToLogin}>
            Voltar ao login
          </button>
        </div>

        {feedback() ? <p class={`feedback ${feedbackTone() === 'success' ? 'feedback-success' : 'feedback-error'}`}>{feedback()}</p> : null}
      </form>
      <p class="auth-card-footer">A nova senha passa a valer imediatamente no backend.</p>
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
