<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { ButtonPurple, ButtonRed } from "$shared/components/buttons";
  let { profile, active, canDelete, onSelect, onEdit, onDelete } = $props();
</script>

<div
  class="flex flex-col justify-between p-5 rounded-2xl bg-white dark:bg-[#0c1220] border transition-all duration-300 shadow-md dark:shadow-lg dark:shadow-black/40 hover:-translate-y-1 hover:border-slate-300 dark:hover:border-slate-700 {active
    ? 'border-violet-500 ring-2 ring-violet-500/10'
    : 'border-slate-200 dark:border-slate-800/80'}"
>
  <div class="space-y-2">
    {#if active}<span
        class="px-2 py-0.5 rounded-full text-[9px] font-bold bg-emerald-500/15 text-emerald-500 border border-emerald-500/20"
        >{t("profiles.active_badge")}</span
      >{/if}
    <h3 class="font-bold text-slate-800 dark:text-white text-base truncate">
      {profile.name}
    </h3>
  </div>
  <div
    class="flex items-center gap-2 mt-6 pt-4 border-t border-slate-100 dark:border-slate-800/60"
  >
    {#if !active}<ButtonPurple size="xs" class="flex-1" onclick={onSelect}
        >{t("profiles.select_btn")}</ButtonPurple
      >{/if}<ButtonPurple size="xs" onclick={onEdit}
      >{t("common.edit")}</ButtonPurple
    >{#if canDelete}<ButtonRed size="xs" onclick={onDelete}
        >{t("common.delete")}</ButtonRed
      >{/if}
  </div>
</div>
