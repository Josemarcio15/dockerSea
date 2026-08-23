<script lang="ts">
  import type { Snippet } from "svelte";

  export type ButtonVariant =
    | "brand"
    | "info"
    | "primary"
    | "pink"
    | "success"
    | "warning"
    | "danger"
    | "secondary"
    | "ghost";
  export type ButtonSize = "xs" | "sm" | "md" | "lg";

  let {
    type = "button",
    variant = "brand",
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
    variant?: ButtonVariant;
    size?: ButtonSize;
    disabled?: boolean;
    loading?: boolean;
    title?: string;
    class?: string;
    onclick?: (e: MouseEvent) => void | Promise<void>;
    children?: Snippet;
    icon?: Snippet;
  } = $props();

  const variantClasses: Record<ButtonVariant, string> = {
    brand:
      "bg-violet-600 hover:bg-violet-700 text-white shadow-md shadow-violet-500/20 border-transparent active:bg-violet-800",
    info:
      "bg-blue-600 hover:bg-blue-700 text-white shadow-md shadow-blue-500/20 border-transparent active:bg-blue-800",
    primary:
      "bg-blue-600 hover:bg-blue-700 text-white shadow-md shadow-blue-500/20 border-transparent active:bg-blue-800",
    pink:
      "bg-pink-600 hover:bg-pink-700 text-white shadow-md shadow-pink-500/20 border-transparent active:bg-pink-800",
    success:
      "bg-emerald-600 hover:bg-emerald-700 text-white shadow-md shadow-emerald-500/20 border-transparent active:bg-emerald-800",
    warning:
      "bg-amber-500 hover:bg-amber-600 text-white shadow-md shadow-amber-500/20 border-transparent active:bg-amber-700",
    danger:
      "bg-red-500 hover:bg-red-600 text-white shadow-md shadow-red-500/20 border-transparent active:bg-red-700",
    secondary:
      "bg-white dark:bg-slate-800 hover:bg-slate-100 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 border border-slate-200/80 dark:border-slate-700 active:bg-slate-200 dark:active:bg-slate-600 shadow-xs",
    ghost:
      "bg-transparent hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 border-transparent",
  };

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
  class="inline-flex items-center justify-center transition-all duration-150 select-none cursor-pointer border font-sans disabled:opacity-50 disabled:cursor-not-allowed disabled:shadow-none {variantClasses[
    variant
  ]} {sizeClasses[size]} {customClass}"
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
