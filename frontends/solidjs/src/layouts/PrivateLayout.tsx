import type { JSX } from 'solid-js';
import type { User } from '../types/auth';
import { AppShellLayout } from './AppShellLayout';

type PrivateLayoutProps = {
  children: JSX.Element;
  user: User | null;
  section: 'dashboard' | 'security';
  onDashboard: () => void;
  onLogout: () => void;
  onAdmin: () => void;
  onChangePassword: () => void;
};

export function PrivateLayout(props: PrivateLayoutProps) {
  const canOpenAdmin = () => props.user?.role === 'admin';

  return (
    <AppShellLayout
      badge="Workspace"
      title={props.section === 'dashboard' ? 'Dashboard' : 'Seguranca da conta'}
      description={
        props.section === 'dashboard'
          ? 'Visao geral do ambiente autenticado com dados de sessao e atalhos principais.'
          : 'Ajuste credenciais e mantenha a conta protegida sem sair do fluxo principal.'
      }
      user={props.user}
      navGroups={[
        {
          label: 'Aplicacao',
          items: [
            {
              label: 'Visao geral',
              hint: 'Painel principal',
              active: props.section === 'dashboard',
              onClick: props.onDashboard,
            },
            {
              label: 'Seguranca',
              hint: 'Senha e sessao',
              active: props.section === 'security',
              onClick: props.onChangePassword,
            },
          ],
        },
        ...(canOpenAdmin()
          ? [
              {
                label: 'Acesso',
                items: [
                  {
                    label: 'Area admin',
                    hint: 'Usuarios e auditoria',
                    onClick: props.onAdmin,
                  },
                ],
              },
            ]
          : []),
      ]}
      actions={[
        ...(canOpenAdmin() ? [{ label: 'Area admin', onClick: props.onAdmin }] : []),
        { label: 'Alterar senha', onClick: props.onChangePassword },
        { label: 'Sair', variant: 'primary' as const, onClick: props.onLogout },
      ]}
    >
      {props.children}
    </AppShellLayout>
  );
}
