import { createSignal, Show } from 'solid-js';
import { loginSchema, registerSchema, type LoginInput, type RegisterInput } from '../schemas/auth';
import { login, register } from '../services/auth';
import { authStore } from '../stores/auth';
import { FormField } from '../components/FormField';

type Mode = 'login' | 'register';

type AuthPageProps = {
  mode: Mode;
  onModeChange: (mode: Mode) => void;
  onAuthenticated: () => void;
  onFatalError: () => void;
};

export function AuthPage(props: AuthPageProps) {
  const [name, setName] = createSignal('');
  const [email, setEmail] = createSignal('');
  const [password, setPassword] = createSignal('');
  const [errors, setErrors] = createSignal<Record<string, string>>({});
  const [feedback, setFeedback] = createSignal<string | null>(null);
  const [submitting, setSubmitting] = createSignal(false);

  const submit = async (event: SubmitEvent) => {
    event.preventDefault();
    setFeedback(null);

    const payload = {
      name: name(),
      email: email(),
      password: password(),
    };

    if (props.mode === 'register') {
      const result = registerSchema.safeParse(payload);
      if (!result.success) {
        setErrors(result.error.flatten().fieldErrors as Record<string, string>);
        return;
      }

      setErrors({});
      setSubmitting(true);
      try {
        const response = await register(result.data as RegisterInput);
        authStore.setSession(response.data.accessToken, response.data.user);
        setFeedback('Conta criada e sessao iniciada com sucesso.');
        props.onAuthenticated();
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Falha ao registrar usuario.';
        setFeedback(message);
        if (message.toLowerCase().includes('internal')) {
          props.onFatalError();
        }
      } finally {
        setSubmitting(false);
      }

      return;
    }

    const result = loginSchema.safeParse({ email: payload.email, password: payload.password });
    if (!result.success) {
      setErrors(result.error.flatten().fieldErrors as Record<string, string>);
      return;
    }

    setErrors({});
    setSubmitting(true);
    try {
      const response = await login(result.data as LoginInput);
      authStore.setSession(response.data.accessToken, response.data.user);
      setFeedback('Sessao iniciada com sucesso.');
      props.onAuthenticated();
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao autenticar.';
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
      <span class="eyebrow">Baseline inicial</span>
      <h2 class="form-title">{props.mode === 'login' ? 'Entrar na referencia' : 'Criar conta de referencia'}</h2>
      <p class="form-copy">
        SolidJS consome o contrato comum do projeto e autentica contra o backend Go em PostgreSQL local.
      </p>

      <form onSubmit={submit}>
        <Show when={props.mode === 'register'}>
          <FormField
            label="Nome"
            name="name"
            value={name()}
            error={errors().name?.[0] ?? errors().name}
            onInput={(event) => setName(event.currentTarget.value)}
          />
        </Show>

        <FormField
          label="E-mail"
          name="email"
          type="email"
          value={email()}
          error={errors().email?.[0] ?? errors().email}
          onInput={(event) => setEmail(event.currentTarget.value)}
        />

        <FormField
          label="Senha"
          name="password"
          type="password"
          value={password()}
          error={errors().password?.[0] ?? errors().password}
          onInput={(event) => setPassword(event.currentTarget.value)}
        />

        <div class="button-row">
          <button class="button button-primary" type="submit" disabled={submitting()}>
            {submitting() ? 'Processando...' : props.mode === 'login' ? 'Entrar' : 'Criar conta'}
          </button>
          <button
            class="button button-secondary"
            type="button"
            onClick={() => props.onModeChange(props.mode === 'login' ? 'register' : 'login')}
          >
            {props.mode === 'login' ? 'Quero criar conta' : 'Ja tenho conta'}
          </button>
        </div>

        <Show when={feedback()}>
          <p class={`feedback ${authStore.currentUser() ? 'feedback-success' : 'feedback-error'}`}>{feedback()}</p>
        </Show>
      </form>
    </section>
  );
}