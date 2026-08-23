<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import type { Snippet } from "svelte";

  export interface ModalButton {
    label: string;
    variant?: "primary" | "secondary" | "danger" | "success" | "warning";
    onclick: () => void | Promise<void>;
    disabled?: boolean;
    type?: "button" | "submit";
  }

  let {
    show = $bindable(false),
    title = "",
    buttons = [],
    children,
  }: {
    show: boolean;
    title: string;
    buttons?: ModalButton[];
    children: Snippet;
  } = $props();

  function getButtonClass(variant: string = "primary") {
    switch (variant) {
      case "primary":
        return "px-4 py-2.5 rounded-xl border-none cursor-pointer text-xs font-bold text-white bg-emerald-600 hover:bg-emerald-700 transition-colors shadow-md shadow-emerald-500/20 disabled:bg-slate-400 dark:disabled:bg-slate-800 disabled:cursor-not-allowed disabled:text-slate-450";
      case "secondary":
        return "px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-800 cursor-pointer text-xs font-bold text-slate-700 dark:text-slate-350 hover:bg-slate-100 dark:hover:bg-slate-850 bg-white dark:bg-transparent transition-colors disabled:opacity-50 disabled:cursor-not-allowed";
      case "danger":
        return "px-4 py-2.5 rounded-xl border-none cursor-pointer text-xs font-bold text-white bg-red-500 hover:bg-red-600 transition-colors shadow-md shadow-red-500/20 disabled:bg-slate-450 disabled:cursor-not-allowed";
      case "success":
        return "px-4 py-2.5 rounded-xl border-none cursor-pointer text-xs font-bold text-white bg-emerald-600 hover:bg-emerald-700 transition-colors shadow-md shadow-emerald-500/20 disabled:bg-slate-450 disabled:cursor-not-allowed";
      case "warning":
        return "px-4 py-2.5 rounded-xl border-none cursor-pointer text-xs font-bold text-white bg-amber-500 hover:bg-amber-600 transition-colors shadow-md shadow-amber-500/20 disabled:bg-slate-450 disabled:cursor-not-allowed";
      default:
        return "px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-800 cursor-pointer text-xs font-bold text-slate-700 dark:text-slate-350 hover:bg-slate-100 dark:hover:bg-slate-850 bg-white dark:bg-transparent transition-colors";
    }
  }
</script>

{#if show}
  <div
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
  >
    <div
      class="bg-white dark:bg-[#0b0f19] border border-slate-200 dark:border-slate-800 rounded-2xl w-125 max-w-full max-h-[90vh] flex flex-col p-6 shadow-2xl animate-scaleIn text-slate-800 dark:text-slate-200"
    >
      <!-- Header -->
      <div
        class="flex justify-between items-center pb-4 border-b border-slate-200 dark:border-slate-800 shrink-0"
      >
        <h2 class="text-lg font-bold text-slate-900 dark:text-white">
          {title}
        </h2>
        <button
          type="button"
          class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 text-xl cursor-pointer bg-transparent border-none"
          onclick={() => (show = false)}
        >
          ✕
        </button>
      </div>

      <!-- Body -->
      <div
        class="py-4 overflow-y-auto flex-1 flex flex-col min-h-0 space-y-4 pr-1"
      >
        {@render children()}
      </div>

      <!-- Footer -->
      <div
        class="flex gap-3 justify-end pt-4 border-t border-slate-200 dark:border-slate-800 shrink-0"
      >
        <!-- Cancel Button (Standard for all modals) -->
        <button
          type="button"
          class="px-4 py-2.5 rounded-xl border-none cursor-pointer text-xs font-bold text-white bg-red-500 hover:bg-red-600 transition-colors shadow-md shadow-red-500/20"
          onclick={() => (show = false)}
        >
          {t("common.cancel")}
        </button>

        {#each buttons as btn}
          <button
            type={btn.type || "button"}
            class={getButtonClass(btn.variant)}
            disabled={btn.disabled}
            onclick={btn.onclick}
          >
            {btn.label}
          </button>
        {/each}
      </div>
    </div>
  </div>
{/if}
