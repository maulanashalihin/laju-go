<script lang="ts">
  export type InputType = 'text' | 'email' | 'password' | 'number' | 'url' | 'tel';

  interface Props {
    id?: string;
    name?: string;
    type?: InputType;
    value?: string;
    placeholder?: string;
    label?: string;
    error?: string;
    required?: boolean;
    disabled?: boolean;
    class?: string;
  }

  let {
    id = '',
    name = '',
    type = 'text',
    value = $bindable(''),
    placeholder = '',
    label = '',
    error = '',
    required = false,
    disabled = false,
    class: className = '',
  }: Props = $props();

  const baseClasses = 'w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500';

  let classes = $derived(`${baseClasses} ${error ? 'border-red-300 focus:ring-red-500 focus:border-red-500' : 'border-gray-300'} ${className}`);
</script>

{#if label}
  <label for={id} class="form-label">{label}</label>
{/if}

<input
  {id}
  {name}
  {type}
  bind:value
  {placeholder}
  {required}
  {disabled}
  class={classes}
/>

{#if error}
  <p class="mt-1 text-sm text-red-600">{error}</p>
{/if}
