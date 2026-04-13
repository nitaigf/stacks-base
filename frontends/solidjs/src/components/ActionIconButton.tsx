import type { JSX } from 'solid-js';

type ActionIconButtonProps = {
  label: string;
  title?: string;
  onClick: JSX.EventHandlerUnion<HTMLButtonElement, MouseEvent>;
  variant?: 'secondary' | 'danger';
  children: JSX.Element;
  disabled?: boolean;
};

export function ActionIconButton(props: ActionIconButtonProps) {
  return (
    <button
      class={`icon-button${props.variant === 'danger' ? ' icon-button-danger' : ''}`}
      type="button"
      title={props.title ?? props.label}
      aria-label={props.label}
      onClick={props.onClick}
      disabled={props.disabled}
    >
      {props.children}
    </button>
  );
}
