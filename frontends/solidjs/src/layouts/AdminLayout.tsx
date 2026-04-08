import type { JSX } from 'solid-js';

type AdminLayoutProps = {
  children: JSX.Element;
  onBackToApp: () => void;
  onUsers: () => void;
  onAuditLogs: () => void;
};

export function AdminLayout(props: AdminLayoutProps) {
  return (
    <div class="layout-admin">
      <nav class="layout-nav layout-nav-justified">
        <div class="layout-nav-actions">
          <button class="button button-secondary button-sm" type="button" onClick={props.onUsers}>
            Usuarios
          </button>
          <button class="button button-secondary button-sm" type="button" onClick={props.onAuditLogs}>
            Auditoria
          </button>
        </div>
        <div class="layout-nav-actions">
          <span class="layout-nav-breadcrumb">Admin</span>
          <button class="button button-secondary button-sm" type="button" onClick={props.onBackToApp}>
            Voltar ao app
          </button>
        </div>
      </nav>
      {props.children}
    </div>
  );
}
