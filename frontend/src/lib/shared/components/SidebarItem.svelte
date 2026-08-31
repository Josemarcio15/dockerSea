<script lang="ts">
  import type { Snippet } from "svelte";

  let {
    icon,
    iconBg,
    label,
    active = false,
    disabled = false,
    href = "#",
    onclick,
  }: {
    icon: Snippet;
    iconBg: string;
    label: string;
    active?: boolean;
    disabled?: boolean;
    href?: string;
    onclick?: () => void;
  } = $props();

  let itemClass = $derived(
    disabled
      ? "text-white/20 cursor-not-allowed"
      : active
        ? "bg-white/15 text-white font-medium cursor-pointer"
        : "text-violet-300/70 hover:text-white hover:bg-white/10 cursor-pointer",
  );

  let dotClass = $derived(
    disabled
      ? "from-slate-700 to-slate-800 opacity-20"
      : active
        ? iconBg
        : `${iconBg} opacity-60`,
  );

  let hoverClass = $derived(!disabled ? "group-hover:scale-110" : "");
</script>

{#if disabled}
  <div
    class="{itemClass} px-4 py-2.5 text-sm rounded-xl transition-all duration-200 flex items-center gap-3 group"
  >
    <div
      class="w-8 h-8 rounded-lg bg-linear-to-br {dotClass} flex items-center justify-center text-sm text-white shadow-lg transition-transform duration-200 {hoverClass}"
    >
      {@render icon()}
    </div>
    <span class="tracking-wide">{label}</span>
    <span class="ml-auto text-xs opacity-40">🔒</span>
  </div>
{:else}
  <button
    type="button"
    {onclick}
    class="{itemClass} w-full text-left px-4 py-2.5 text-sm rounded-xl transition-all duration-200 flex items-center gap-3 group border-0 bg-transparent"
  >
    <div
      class="w-8 h-8 rounded-lg bg-linear-to-br {dotClass} flex items-center justify-center text-sm text-white shadow-lg transition-transform duration-200 {hoverClass}"
    >
      {@render icon()}
    </div>
    <span class="tracking-wide">{label}</span>
    {#if active}
      <div
        class="w-1 h-5 rounded-full bg-violet-400 ml-auto shadow-lg shadow-violet-400/50"
      ></div>
    {/if}
  </button>
{/if}
