<script lang="ts">
  import Button, { type ButtonSize } from "./Button.svelte";
  import type { Snippet } from "svelte";

  let {
    type = "button",
    size = "md",
    disabled = false,
    loading = false,
    title = "",
    class: customClass = "",
    onclick,
    children,
    icon,
  }: {
    type?: "button" | "submit" | "reset";
    size?: ButtonSize;
    disabled?: boolean;
    loading?: boolean;
    title?: string;
    class?: string;
    onclick?: (e: MouseEvent) => unknown;
    children?: Snippet;
    icon?: Snippet;
  } = $props();
</script>

<Button
  {type}
  {size}
  {disabled}
  {loading}
  {title}
  class="bg-violet-600 hover:bg-violet-700 text-white shadow-md shadow-violet-500/20 border-transparent active:bg-violet-800 {customClass}"
  {onclick}
  {icon}
>
  {#if children}
    {@render children()}
  {/if}
</Button>
