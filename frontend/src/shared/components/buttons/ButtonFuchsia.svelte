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
  class="bg-[#86088d] hover:bg-[#700775] active:bg-[#5d0562] text-white shadow-md shadow-[#86088d]/25 border-transparent {customClass}"
  {onclick}
  {icon}
>
  {#if children}
    {@render children()}
  {/if}
</Button>
