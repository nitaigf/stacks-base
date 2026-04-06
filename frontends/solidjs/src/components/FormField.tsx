import type { JSX } from 'solid-js';

type FormFieldProps = {
  label: string;
  name: string;
  type?: string;
  value: string;
  error?: string;
  onInput: JSX.EventHandler<HTMLInputElement, InputEvent>;
};

export function FormField(props: FormFieldProps) {
  return (
    <div class="field-group">
      <label class="field-label" for={props.name}>
        {props.label}
      </label>
      <input
        id={props.name}
        name={props.name}
        class="field-input"
        type={props.type ?? 'text'}
        value={props.value}
        onInput={props.onInput}
      />
      {props.error ? <span class="field-error">{props.error}</span> : null}
    </div>
  );
}