<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { ButtonOrange, ButtonRed } from "$shared/components/buttons";

  let {
    show = $bindable(false),
    targetName = "",
    deleteVolumes = $bindable(false),
    loading = false,
    onConfirm = () => {},
  }: {
    show: boolean;
    targetName: string;
    deleteVolumes: boolean;
    loading: boolean;
    onConfirm: () => void;
  } = $props();
</script>

{#if show}
  <div
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
  >
    <div
      class="bg-white dark:bg-[#0b0f19] border border-slate-200 dark:border-slate-800 rounded-2xl w-140 max-w-full flex flex-col p-6 shadow-2xl animate-scaleIn text-slate-800 dark:text-slate-200 gap-4"
    >
      <div class="flex justify-between items-center pb-3 border-b border-slate-200 dark:border-slate-800">
        <h3 class="text-base font-bold text-slate-900 dark:text-white flex items-center gap-2">
          {t("stacks.remove_remote_btn")}: {targetName}
        </h3>
        <button
          type="button"
          class="text-slate-400 hover:text-slate-200 text-lg bg-transparent border-none cursor-pointer"
          onclick={() => (show = false)}
        >
          ✕
        </button>
      </div>

      <div class="space-y-4">
        <p class="text-xs text-slate-600 dark:text-slate-400 m-0">
          {t("stacks.remove_remote_desc")}
        </p>

        <!-- Checkbox de volumes persistentes (DESMARCADA por padrão) -->
        <div class="p-3.5 rounded-xl border border-rose-500/20 bg-rose-500/5 dark:bg-rose-950/20 flex flex-col gap-2">
          <label class="flex items-center gap-2 text-xs font-bold text-rose-600 dark:text-rose-400 cursor-pointer">
            <input
              type="checkbox"
              class="rounded border-rose-400 text-rose-600 focus:ring-rose-500 w-4 h-4 cursor-pointer"
              bind:checked={deleteVolumes}
            />
            {t("stacks.delete_volumes_checkbox")}
          </label>
          {#if deleteVolumes}
            <p class="text-[11px] text-rose-500 m-0 animate-fadeIn">
              {t("stacks.delete_volumes_warning")}
            </p>
          {/if}
        </div>
      </div>

      <div class="flex gap-2.5 justify-end items-center pt-3 border-t border-slate-200 dark:border-slate-800 shrink-0">
        <ButtonOrange size="sm" disabled={loading} onclick={() => (show = false)}>
          {t("common.cancel")}
        </ButtonOrange>
        <ButtonRed size="sm" {loading} disabled={loading} onclick={onConfirm}>
          {loading ? t("stacks.executing") : "Confirmar Remoção na VPS"}
        </ButtonRed>
      </div>
    </div>
  </div>
{/if}
