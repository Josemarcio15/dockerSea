<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { ButtonPurple, ButtonRed } from "$shared/components/buttons";
  let { profile, active, canDelete, onSelect, onEdit, onDelete } = $props();
</script>

<div
  class="flex flex-col justify-between p-5 rounded-2xl bg-[radial-gradient(ellipse_at_center,_var(--tw-gradient-stops))] from-white via-slate-50 to-slate-100 dark:from-[#151c2e] dark:via-[#0e1526] dark:to-[#080d1a] border border-slate-200/80 dark:border-slate-800 text-slate-800 dark:text-slate-200 shadow-[0_10px_25px_-5px_rgba(0,0,0,0.05),0_8px_10px_-6px_rgba(0,0,0,0.03),inset_0_1px_0_0_rgba(255,255,255,0.9),inset_0_-1px_0_0_rgba(0,0,0,0.03)] dark:shadow-[0_10px_25px_-5px_rgba(0,0,0,0.5),inset_0_1px_0_0_rgba(255,255,255,0.05)] hover:-translate-y-1 hover:shadow-[0_20px_30px_-10px_rgba(0,0,0,0.08)] transition-all duration-200 active:translate-y-0"
>
  <div class="space-y-2">
    <div class="flex items-center justify-between">
      {#if active}
        <span
          class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-800/50 shadow-2xs"
        >
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
          {t("profiles.active_badge")}
        </span>
      {:else}
        <span
          class="px-2.5 py-0.5 rounded-full text-[11px] font-medium bg-slate-100 dark:bg-slate-800/70 text-slate-500 dark:text-slate-400 border border-slate-200/70 dark:border-slate-700/60"
        >
          {t("profiles.inactive_badge")}
        </span>
      {/if}

      <!-- Ícone decorativo 3D em card branco -->
      <div
        class="w-9 h-9 rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 shadow-sm shadow-[inset_0_-2px_0_0_rgba(0,0,0,0.05)] flex items-center justify-center text-slate-700 dark:text-slate-300"
      >
        <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
          <path d="M12 2L2 22l10-4 10 4L12 2z"/>
        </svg>
      </div>
    </div>

    <h3 class="font-extrabold text-slate-900 dark:text-white text-xl tracking-tight truncate">
      {profile.name}
    </h3>
  </div>

  <div
    class="flex items-center gap-2 mt-6 pt-4 border-t border-slate-100 dark:border-slate-800"
  >
    {#if !active}
      <button
        type="button"
        class="flex-1 px-4 py-2.5 rounded-xl font-bold text-xs bg-linear-to-b from-violet-500 to-violet-600 border border-violet-400/40 border-b-violet-700 border-b-2 text-white shadow-md shadow-violet-500/25 shadow-[inset_0_1px_0_0_rgba(255,255,255,0.35)] hover:from-violet-400 hover:to-violet-500 active:translate-y-0.5 active:border-b-0 transition-all cursor-pointer"
        onclick={onSelect}
      >
        {t("profiles.select_btn")}
      </button>
    {/if}

    <button
      type="button"
      class="px-3.5 py-2.5 rounded-xl font-bold text-xs bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-200 border border-slate-300 dark:border-slate-700 shadow-sm shadow-[inset_0_-2px_0_0_rgba(0,0,0,0.06)] hover:bg-slate-50 active:translate-y-0.5 transition-all cursor-pointer"
      onclick={onEdit}
    >
      {t("common.edit")}
    </button>

    {#if canDelete}
      <button
        type="button"
        class="px-3.5 py-2.5 rounded-xl font-bold text-xs bg-linear-to-b from-rose-500 to-rose-600 text-white border border-rose-400/40 border-b-rose-700 border-b-2 shadow-md shadow-rose-500/20 shadow-[inset_0_1px_0_0_rgba(255,255,255,0.35)] hover:from-rose-400 hover:to-rose-500 active:translate-y-0.5 active:border-b-0 transition-all cursor-pointer"
        onclick={onDelete}
      >
        {t("common.delete")}
      </button>
    {/if}
  </div>
</div>
