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
    onclick?: (e: MouseEvent) => void | Promise<void>;
    children?: Snippet;
    icon?: Snippet;
  } = $props();
</script>

<Button
  {type}
  variant="success"
  {size}
  {disabled}
  {loading}
  {title}
  class={customClass}
  {onclick}
  {children}
  {icon}
/>
