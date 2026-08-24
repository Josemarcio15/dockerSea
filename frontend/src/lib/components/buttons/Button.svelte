<script lang="ts">
  import type { Snippet } from "svelte";

  export type ButtonSize = "xs" | "sm" | "md" | "lg";

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

  const sizeClasses: Record<ButtonSize, string> = {
    xs: "px-2.5 py-1 text-[11px] rounded-lg gap-1.5",
    sm: "px-3.5 py-1.5 text-xs rounded-xl gap-1.5 font-semibold",
    md: "px-4 py-2 text-xs rounded-xl gap-2 font-bold",
    lg: "px-5 py-2.5 text-sm rounded-xl gap-2.5 font-bold",
  };
</script>

<button
  {type}
  {title}
  disabled={disabled || loading}
  {onclick}
  class="inline-flex items-center justify-center transition-all duration-150 select-none cursor-pointer border font-sans disabled:opacity-50 disabled:cursor-not-allowed disabled:shadow-none {sizeClasses[
    size
  ]} {customClass}"
>
  {#if loading}
    <div
      class="w-3.5 h-3.5 border-2 border-current border-t-transparent rounded-full animate-spin shrink-0"
    ></div>
  {:else if icon}
    {@render icon()}
  {/if}

  {#if children}
    {@render children()}
  {/if}
</button>
