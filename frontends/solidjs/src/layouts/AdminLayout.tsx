import type { JSX } from 'solid-js';
import type { User } from '../types/auth';
import { AppShellLayout } from './AppShellLayout';

type AdminLayoutProps = {
  children: JSX.Element;
  user: User | null;
  section: 'users' | 'audit';
  onBackToApp: () => void;
  onUsers: () => void;
  onAuditLogs: () => void;
};

export function AdminLayout(props: AdminLayoutProps) {
  return (
    <AppShellLayout
      badge="Admin"
      title={props.section === 'users' ? 'Gestao de usuarios' : 'Auditoria do sistema'}
      description={
        props.section === 'users'
          ? 'Controle identidades, status de acesso e fluxos administrativos da baseline.'
          : 'Acompanhe trilhas de auditoria, eventos de acesso e operacoes criticas.'
      }
      user={props.user}
      navGroups={[
        {
          label: 'Administracao',
          items: [
            {
              label: 'Usuarios',
              hint: 'CRUD e exportacoes',
              active: props.section === 'users',
              onClick: props.onUsers,
            },
            {
              label: 'Auditoria',
              hint: 'Eventos e trilhas',
              active: props.section === 'audit',
              onClick: props.onAuditLogs,
            },
          ],
        },
        {
          label: 'Workspace',
          items: [
            {
              label: 'Voltar ao app',
              hint: 'Painel autenticado',
              onClick: props.onBackToApp,
            },
          ],
        },
      ]}
      actions={[
        { label: 'Usuarios', onClick: props.onUsers },
        { label: 'Auditoria', onClick: props.onAuditLogs },
        { label: 'Voltar ao app', variant: 'primary' as const, onClick: props.onBackToApp },
      ]}
    >
      {props.children}
    </AppShellLayout>
  );
}
