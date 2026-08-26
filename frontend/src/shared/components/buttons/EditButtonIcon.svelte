<script lang="ts">
  import Button, { type ButtonSize } from "./Button.svelte";
  import type { Snippet } from "svelte";

  let {
    type = "button",
    size = "sm",
    disabled = false,
    loading = false,
    title = "Editar",
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

{#snippet defaultEditIcon()}
  <svg
    xmlns="http://www.w3.org/2000/svg"
    class="w-3.5 h-3.5"
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
    stroke-width="2"
  >
    <path
      stroke-linecap="round"
      stroke-linejoin="round"
      d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L6.832 19.82a4.5 4.5 0 0 1-1.897 1.13l-2.685.8.8-2.685a4.5 4.5 0 0 1 1.13-1.897L16.863 4.487Zm0 0L19.5 7.125"
    />
  </svg>
{/snippet}

<Button
  {type}
  {size}
  {disabled}
  {loading}
  {title}
  class="bg-teal-600 hover:bg-teal-700 text-white shadow-md shadow-teal-500/20 border-transparent active:bg-teal-800 {customClass}"
  {onclick}
  icon={defaultEditIcon}
>
  {#if children}
    {@render children()}
  {/if}
</Button>
