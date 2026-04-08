import { Show, createSignal, onMount } from 'solid-js';
import { createUser, getUser, updateUser } from '../services/users';
import { FormField } from '../components/FormField';

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
  const [feedback, setFeedback] = createSignal<string | null>(null);
  const [loading, setLoading] = createSignal(props.mode === 'edit');
  const [submitting, setSubmitting] = createSignal(false);

  onMount(async () => {
    if (props.mode !== 'edit' || !props.userId) {
      return;
    }

    try {
      const response = await getUser(props.userId);
      setName(response.data.name);
      setEmail(response.data.email);
      setRole(response.data.role);
      setStatus(response.data.status);
      setLastLoginAt(response.data.lastLoginAt ?? null);
      setDeletedAt(response.data.deletedAt ?? null);
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao carregar usuario.';
      setFeedback(message);
      if (message.toLowerCase().includes('internal')) {
        props.onFatalError();
      }
    } finally {
      setLoading(false);
    }
  });

  const submit = async (event: SubmitEvent) => {
    event.preventDefault();
    setSubmitting(true);
    setFeedback(null);

    try {
      if (props.mode === 'create') {
        const response = await createUser({
          name: name(),
          email: email(),
          password: password(),
          role: role(),
          status: status(),
        });
        setFeedback('Usuario criado com sucesso.');
        props.onSaved(response.data.id);
      } else if (props.userId) {
        const response = await updateUser(props.userId, {
          name: name(),
          email: email(),
          role: role(),
        });
        setFeedback('Usuario atualizado com sucesso.');
        props.onSaved(response.data.id);
      }
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Falha ao salvar usuario.';
      setFeedback(message);
      if (message.toLowerCase().includes('internal')) {
        props.onFatalError();
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <section class="surface-card dashboard-card">
      <span class="metric-pill metric-pill-tag">{props.mode === 'create' ? 'Novo usuario' : 'Visualizar e editar'}</span>
      <h2 class="form-title">{props.mode === 'create' ? 'Criar usuario' : 'Detalhes do usuario'}</h2>
      <p class="form-copy">Tela ligada ao backend real para leitura individual, criacao e edicao.</p>

      {loading() ? (
        <p class="form-copy">Carregando usuario...</p>
      ) : (
        <form onSubmit={submit}>
          <FormField label="Nome" name="name" value={name()} onInput={(event) => setName(event.currentTarget.value)} />
          <FormField
            label="E-mail"
            name="email"
            type="email"
            value={email()}
            onInput={(event) => setEmail(event.currentTarget.value)}
          />
          <label class="field-group">
            <span class="field-label">Papel</span>
            <select class="field-input" value={role()} onInput={(event) => setRole(event.currentTarget.value as 'admin' | 'member')}>
              <option value="member">Member</option>
              <option value="admin">Admin</option>
            </select>
          </label>
          <Show when={props.mode === 'create'}>
            <label class="field-group">
              <span class="field-label">Status inicial</span>
              <select
                class="field-input"
                value={status()}
                onInput={(event) => setStatus(event.currentTarget.value as 'active' | 'inactive')}
              >
                <option value="active">Ativo</option>
                <option value="inactive">Inativo</option>
              </select>
            </label>
            <FormField
              label="Senha inicial"
              name="password"
              type="password"
              value={password()}
              onInput={(event) => setPassword(event.currentTarget.value)}
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

          {feedback() ? <p class="feedback feedback-success">{feedback()}</p> : null}
        </form>
      )}
    </section>
  );
}

function formatDate(value: string | null) {
  if (!value) {
    return '—';
  }
  return new Date(value).toLocaleString('pt-BR');
}
