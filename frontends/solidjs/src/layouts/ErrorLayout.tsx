import type { JSX } from 'solid-js';

type ErrorLayoutProps = {
  children: JSX.Element;
};

export function ErrorLayout(props: ErrorLayoutProps) {
  return <>{props.children}</>;
}
