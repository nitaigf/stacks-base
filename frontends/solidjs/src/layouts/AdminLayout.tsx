import type { JSX } from 'solid-js';

type AdminLayoutProps = {
  children: JSX.Element;
  onBackToApp: () => void;
};

export function AdminLayout(props: AdminLayoutProps) {
  return (
    <div class="layout-admin">
      <nav class="layout-nav layout-nav-justified">
        <span class="layout-nav-breadcrumb">Admin</span>
        <button class="button button-secondary button-sm" type="button" onClick={props.onBackToApp}>
          Voltar ao app
        </button>
      </nav>
      {props.children}
    </div>
  );
}
