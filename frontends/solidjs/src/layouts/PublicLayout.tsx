import type { JSX } from 'solid-js';

type PublicLayoutProps = {
  children: JSX.Element;
};

export function PublicLayout(props: PublicLayoutProps) {
  return <>{props.children}</>;
}
