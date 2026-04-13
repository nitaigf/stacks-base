import type { JSX } from 'solid-js';

type ErrorLayoutProps = {
  children: JSX.Element;
};

export function ErrorLayout(props: ErrorLayoutProps) {
  return (
    <div class="error-shell">
      <div class="error-shell-content">{props.children}</div>
    </div>
  );
}
