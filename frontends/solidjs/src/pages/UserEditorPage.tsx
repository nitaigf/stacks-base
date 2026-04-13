import { Show, createEffect, createSignal, on, type Setter } from 'solid-js';
import { createUser, getUser, updateUser } from '../services/users';
import { isApiClientError } from '../services/api';
import { FormField } from '../components/FormField';
import {
  userCreateSchema,
  userUpdateSchema,
} from '../schemas/auth';

type UserEditorPageProps = {
  mode: 'create' | 'edit';
  userId?: string;
  onBack: () => void;
  onSaved: (userId: string) => void;
  onFatalError: () => void;
};

export function UserEditorPage(props: UserEditorPageProps) {
  const [name, setName] = createSignal('');
  const [email, setEmail] = createSignal('');
  const [role, setRole] = createSignal<'admin' | 'member'>('member');
  const [status, setStatus] = createSignal<'active' | 'inactive'>('active');
  const [password, setPassword] = createSignal('');
  const [lastLoginAt, setLastLoginAt] = createSignal<string | null>(null);
  const [deletedAt, setDeletedAt] = createSignal<string | null>(null);
  const [errors, setErrors] = createSignal<Record<string, string[] | undefined>>({});
  const [feedback, setFeedback] = createSignal<string | null>(null);
  const [feedbackTone, setFeedbackTone] = createSignal<'success' | 'error' | null>(null);
  const [loading, setLoading] = createSignal(props.mode === 'edit');
  const [loadFailed, setLoadFailed] = createSignal(false);
  const [submitting, setSubmitting] = createSignal(false);

  const loadUser = async (userId: string) => {
    setLoading(true);
    setLoadFailed(false);
    setErrors({});
    setFeedback(null);
    setFeedbackTone(null);

    try {
      const response = await getUser(userId);
      setName(response.data.name);
      setEmail(response.data.email);
      setRole(response.data.role);
      setStatus(response.data.status);
      setLastLoginAt(response.data.lastLoginAt ?? null);
      setDeletedAt(response.data.deletedAt ?? null);
    } catch (error) {
      setLoadFailed(true);
      const message = error instanceof Error ? error.message : 'Falha ao carregar usuario.';
      setFeedback(message);
      setFeedbackTone('error');
      if (isApiClientError(error) && error.isServerError) {
        props.onFatalError();
      }
    } finally {
      setLoading(false);
    }
  };

  createEffect(
    on(
      () => [props.mode, props.userId] as const,
      ([mode, userId]) => {
        if (mode !== 'edit' || !userId) {
          setLoadFailed(false);
          setLoading(false);
          setErrors({});
          setFeedbackTone(null);
          return;
        }

        void loadUser(userId);
      },
    ),
  );

  const submit = async (event: SubmitEvent) => {
    event.preventDefault();
    setFeedback(null);
    setFeedbackTone(null);

    if (props.mode === 'create') {
      const result = userCreateSchema.safeParse({
        name: name(),
        email: email(),
        password: password(),
        role: role(),
        status: status(),
      });

      if (!result.success) {
        setErrors(result.error.flatten().fieldErrors);
        return;
      }

      setErrors({});
      setSubmitting(true);

      try {
        const response = await createUser({
          name: result.data.name,
          email: result.data.email,
          password: result.data.password,
          role: result.data.role,
          status: result.data.status,
        });
        setFeedback('Usuario criado com sucesso.');
        setFeedbackTone('success');
        props.onSaved(response.data.id);
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Falha ao salvar usuario.';
        setFeedback(message);
        setFeedbackTone('error');
        if (isApiClientError(error) && error.isServerError) {
          props.onFatalError();
        }
      } finally {
        setSubmitting(false);
      }

      return;
    }

    const result = userUpdateSchema.safeParse({
      name: name(),
      email: email(),
      role: role(),
    });

    if (!result.success) {
      setErrors(result.error.flatten().fieldErrors);
      return;
    }

    if (!props.userId) {
      return;
    }

    setErrors({});
    setSubmitting(true);

    try {
      const response = await updateUser(props.userId, {
        name: result.data.name,
        email: result.data.email,
        role: result.data.role,
      });
      setFeedback('Usuario atualizado com sucesso.');
      setFeedbackTone('success');
      props.onSaved(response.data.id);
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao salvar usuario.';
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
    <section class="surface-card dashboard-card form-card">
      <span class="metric-pill metric-pill-tag">{props.mode === 'create' ? 'Novo usuario' : 'Visualizar e editar'}</span>
      <h2 class="form-title">{props.mode === 'create' ? 'Criar usuario' : 'Detalhes do usuario'}</h2>
      <p class="form-copy">Tela ligada ao backend real para leitura individual, criacao e edicao.</p>

      {loading() ? (
        <p class="form-copy">Carregando usuario...</p>
      ) : loadFailed() ? (
        <>
          {feedback() ? <p class="feedback feedback-error">{feedback()}</p> : null}
          <div class="button-row">
            <button class="button button-secondary" type="button" onClick={props.onBack}>
              Voltar
            </button>
          </div>
        </>
      ) : (
        <form noValidate onSubmit={submit}>
          <FormField
            label="Nome"
            name="name"
            value={name()}
            error={errors().name?.[0]}
            onInput={(event) => {
              setName(event.currentTarget.value);
              clearFieldError(setErrors, 'name');
            }}
          />
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
          <label class="field-group">
            <span class="field-label">Papel</span>
            <select
              class="field-input"
              value={role()}
              onInput={(event) => {
                setRole(event.currentTarget.value as 'admin' | 'member');
                clearFieldError(setErrors, 'role');
              }}
            >
              <option value="member">Member</option>
              <option value="admin">Admin</option>
            </select>
            {errors().role?.[0] ? <span class="field-error">{errors().role?.[0]}</span> : null}
          </label>
          <Show when={props.mode === 'create'}>
            <label class="field-group">
              <span class="field-label">Status inicial</span>
              <select
                class="field-input"
                value={status()}
                onInput={(event) => {
                  setStatus(event.currentTarget.value as 'active' | 'inactive');
                  clearFieldError(setErrors, 'status');
                }}
              >
                <option value="active">Ativo</option>
                <option value="inactive">Inativo</option>
              </select>
              {errors().status?.[0] ? <span class="field-error">{errors().status?.[0]}</span> : null}
            </label>
            <FormField
              label="Senha inicial"
              name="password"
              type="password"
              value={password()}
              error={errors().password?.[0]}
              onInput={(event) => {
                setPassword(event.currentTarget.value);
                clearFieldError(setErrors, 'password');
              }}
            />
          </Show>

          <Show when={props.mode === 'edit'}>
            <div class="metric-row">
              <span class="metric-pill metric-pill-tag">Status: {status()}</span>
              <span class="metric-pill metric-pill-tag">Ultimo login: {formatDate(lastLoginAt())}</span>
              <span class="metric-pill metric-pill-tag">Excluido em: {formatDate(deletedAt())}</span>
            </div>
          </Show>

          <div class="button-row">
            <button class="button button-primary" type="submit" disabled={submitting()}>
              {submitting() ? 'Salvando...' : props.mode === 'create' ? 'Criar usuario' : 'Salvar alteracoes'}
            </button>
            <button class="button button-secondary" type="button" onClick={props.onBack}>
              Voltar
            </button>
          </div>

          {feedback() ? (
            <p class={`feedback ${feedbackTone() === 'success' ? 'feedback-success' : 'feedback-error'}`}>
              {feedback()}
            </p>
          ) : null}
        </form>
      )}
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

function formatDate(value: string | null) {
  if (!value) {
    return '—';
  }
  return new Date(value).toLocaleString('pt-BR');
}
