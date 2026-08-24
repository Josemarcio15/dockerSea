<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import { ButtonRed, ButtonYellow, ButtonGreen, ButtonBlue, ButtonOrange } from "$lib/components/buttons";

  let {
    show = $bindable(false),
    title = "Confirmar Ação",
    message = "Tem certeza de que deseja continuar?",
    confirmText = "Confirmar",
    cancelText = "Cancelar",
    type = "danger", // "danger" | "warning" | "info" | "success"
    onConfirm = () => {},
    onCancel = () => {},
  }: {
    show: boolean;
    title?: string;
    message?: string;
    confirmText?: string;
    cancelText?: string;
    type?: "danger" | "warning" | "info" | "success";
    onConfirm: () => void | Promise<void>;
    onCancel?: () => void;
  } = $props();

  let loading = $state(false);

  async function handleConfirm() {
    loading = true;
    try {
      await onConfirm();
      show = false;
    } finally {
      loading = false;
    }
  }

  function handleCancel() {
    if (onCancel) onCancel();
    show = false;
  }
</script>

{#if show}
  <div
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
  >
    <div
      class="bg-white dark:bg-[#0b0f19] border border-slate-200 dark:border-slate-800 rounded-2xl w-105 max-w-full flex flex-col p-6 shadow-2xl animate-scaleIn text-slate-800 dark:text-slate-200 gap-4"
    >
      <!-- Header com Ícone -->
      <div class="flex items-start gap-3.5">
        <div
          class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0 text-lg {type ===
          'danger'
            ? 'bg-red-500/10 text-red-500 border border-red-500/20'
            : type === 'success'
              ? 'bg-emerald-500/10 text-emerald-500 border border-emerald-500/20'
              : type === 'warning'
                ? 'bg-amber-500/10 text-amber-500 border border-amber-500/20'
                : 'bg-blue-500/10 text-blue-500 border border-blue-500/20'}"
        >
          {#if type === "danger"}
            ⚠️
          {:else if type === "success"}
            🚀
          {:else if type === "warning"}
            ⚡
          {:else}
            ℹ️
          {/if}
        </div>

        <div class="flex flex-col min-w-0 flex-1 pt-0.5">
          <h3 class="text-base font-bold text-slate-900 dark:text-white m-0">
            {title}
          </h3>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-1 leading-relaxed whitespace-pre-wrap">
            {message}
          </p>
        </div>
      </div>

      <!-- Footer Buttons -->
      <div class="flex gap-2.5 justify-end pt-3 border-t border-slate-100 dark:border-slate-800/80 items-center">
        <ButtonOrange
          size="sm"
          disabled={loading}
          onclick={handleCancel}
        >
          {cancelText || t("common.cancel")}
        </ButtonOrange>

        {#if type === "danger"}
          <ButtonRed size="sm" {loading} onclick={handleConfirm}>
            {confirmText}
          </ButtonRed>
        {:else if type === "success"}
          <ButtonGreen size="sm" {loading} onclick={handleConfirm}>
            {confirmText}
          </ButtonGreen>
        {:else if type === "warning"}
          <ButtonYellow size="sm" {loading} onclick={handleConfirm}>
            {confirmText}
          </ButtonYellow>
        {:else}
          <ButtonBlue size="sm" {loading} onclick={handleConfirm}>
            {confirmText}
          </ButtonBlue>
        {/if}
      </div>
    </div>
  </div>
{/if}
