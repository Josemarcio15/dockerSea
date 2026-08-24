<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import type { Snippet } from "svelte";
  import {
    ButtonOrange,
    ButtonGreen,
    ButtonRed,
    ButtonYellow,
    ButtonBlue,
    Button,
  } from "$lib/components/buttons";

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
        <!-- Cancel/Close Button -->
        <ButtonOrange onclick={() => (show = false)}>
          {t("common.cancel")}
        </ButtonOrange>

        {#each buttons as btn}
          {#if btn.variant === "danger"}
            <ButtonRed
              type={btn.type || "button"}
              disabled={btn.disabled}
              onclick={btn.onclick}
            >
              {btn.label}
            </ButtonRed>
          {:else if btn.variant === "warning"}
            <ButtonYellow
              type={btn.type || "button"}
              disabled={btn.disabled}
              onclick={btn.onclick}
            >
              {btn.label}
            </ButtonYellow>
          {:else if btn.variant === "success" || btn.variant === "primary"}
            <ButtonGreen
              type={btn.type || "button"}
              disabled={btn.disabled}
              onclick={btn.onclick}
            >
              {btn.label}
            </ButtonGreen>
          {:else}
            <ButtonBlue
              type={btn.type || "button"}
              disabled={btn.disabled}
              onclick={btn.onclick}
            >
              {btn.label}
            </ButtonBlue>
          {/if}
        {/each}
      </div>
    </div>
  </div>
{/if}
