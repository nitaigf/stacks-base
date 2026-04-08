import type { JSX } from 'solid-js';
import { Show } from 'solid-js';
import type { User } from '../types/auth';

type PrivateLayoutProps = {
  children: JSX.Element;
  user: User | null;
  onLogout: () => void;
  onAdmin: () => void;
  onChangePassword: () => void;
  onFatalError: () => void;
};

export function PrivateLayout(props: PrivateLayoutProps) {
  return (
    <div class="layout-private">
      <nav class="layout-nav layout-nav-justified">
        <span class="layout-nav-user">{props.user?.name}</span>
        <div class="layout-nav-actions">
          <button class="button button-secondary button-sm" type="button" onClick={props.onChangePassword}>
            Alterar senha
          </button>
          <Show when={props.user?.role === 'admin'}>
            <button class="button button-secondary button-sm" type="button" onClick={props.onAdmin}>
              Area admin
            </button>
          </Show>
          <button class="button button-secondary button-sm" type="button" onClick={props.onLogout}>
            Sair
          </button>
        </div>
      </nav>
      {props.children}
    </div>
  );
}
