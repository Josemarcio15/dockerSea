<script lang="ts">
  import Button, { type ButtonSize } from "./Button.svelte";
  import type { Snippet } from "svelte";

  let {
    type = "button",
    size = "sm",
    disabled = false,
    loading = false,
    title = "Excluir",
    class: customClass = "",
    onclick,
    children,
  }: {
    type?: "button" | "submit" | "reset";
    size?: ButtonSize;
    disabled?: boolean;
    loading?: boolean;
    title?: string;
    class?: string;
    onclick?: (e: MouseEvent) => unknown;
    children?: Snippet;
  } = $props();
</script>

{#snippet defaultTrashIcon()}
  <svg
    xmlns="http://www.w3.org/2000/svg"
    class="w-3.5 h-3.5"
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
    stroke-width="1.75"
  >
    <path
      stroke-linecap="round"
      stroke-linejoin="round"
      d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0"
    />
  </svg>
{/snippet}

<Button
  {type}
  {size}
  {disabled}
  {loading}
  {title}
  class="bg-red-500 hover:bg-red-600 text-white shadow-md shadow-red-500/20 border-transparent active:bg-red-700 {customClass}"
  {onclick}
  icon={defaultTrashIcon}
>
  {#if children}
    {@render children()}
  {/if}
</Button>
