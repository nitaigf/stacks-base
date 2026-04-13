import { For, Show } from 'solid-js';
import type { JSX } from 'solid-js';
import type { User } from '../types/auth';

type NavigationItem = {
  label: string;
  hint?: string;
  active?: boolean;
  onClick: () => void;
};

type NavigationGroup = {
  label: string;
  items: NavigationItem[];
};

type TopbarAction = {
  label: string;
  variant?: 'primary' | 'secondary';
  onClick: () => void;
};

type AppShellLayoutProps = {
  children: JSX.Element;
  badge: string;
  title: string;
  description: string;
  user?: User | null;
  navGroups: NavigationGroup[];
  actions?: TopbarAction[];
};

export function AppShellLayout(props: AppShellLayoutProps) {
  return (
    <div class="workspace-shell">
      <aside class="workspace-sidebar">
        <div class="workspace-brand">
          <span class="workspace-brand-mark">SB</span>
          <div>
            <p class="workspace-brand-title">Stacks Base</p>
            <p class="workspace-brand-copy">Baseline operacional com SolidJS + Go.</p>
          </div>
        </div>

        <div class="workspace-nav">
          <For each={props.navGroups}>
            {(group) => (
              <section class="workspace-nav-group">
                <span class="workspace-nav-label">{group.label}</span>
                <div class="workspace-nav-list">
                  <For each={group.items}>
                    {(item) => (
                      <button
                        class={`workspace-nav-item${item.active ? ' workspace-nav-item-active' : ''}`}
                        type="button"
                        onClick={item.onClick}
                      >
                        <span>{item.label}</span>
                        <Show when={item.hint}>
                          <small>{item.hint}</small>
                        </Show>
                      </button>
                    )}
                  </For>
                </div>
              </section>
            )}
          </For>
        </div>

        <div class="workspace-sidebar-spacer" />

        <Show when={props.user}>
          <div class="workspace-user-card">
            <div class="workspace-user-avatar">{props.user?.name?.slice(0, 1) ?? 'U'}</div>
            <div>
              <p class="workspace-user-name">{props.user?.name}</p>
              <p class="workspace-user-meta">{props.user?.email}</p>
            </div>
          </div>
        </Show>
      </aside>

      <main class="workspace-main">
        <header class="workspace-topbar">
          <div class="workspace-topbar-copy">
            <span class="metric-pill metric-pill-tag">{props.badge}</span>
            <h1 class="workspace-title">{props.title}</h1>
            <p class="workspace-description">{props.description}</p>
          </div>
          <div class="workspace-topbar-actions">
            <For each={props.actions ?? []}>
              {(action) => (
                <button
                  class={`button ${action.variant === 'primary' ? 'button-primary' : 'button-secondary'} button-sm`}
                  type="button"
                  onClick={action.onClick}
                >
                  {action.label}
                </button>
              )}
            </For>
          </div>
        </header>

        <div class="workspace-content">{props.children}</div>
      </main>
    </div>
  );
}
