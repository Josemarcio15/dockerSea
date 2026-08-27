<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { ButtonPink } from "$shared/components/buttons";
  import type { ExtrasMainTab } from "../types";

  let {
    mainTab = $bindable(),
    refresh,
  }: { mainTab: ExtrasMainTab; refresh: () => void } = $props();
</script>

<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
  <div
    class="inline-flex items-center px-5 py-2.5 rounded-2xl bg-linear-to-br from-violet-100 to-fuchsia-100 dark:from-violet-950/40 dark:to-fuchsia-950/40 border border-violet-200/50 dark:border-violet-800/50 self-start shadow-sm"
  >
    <h1
      class="text-2xl font-bold bg-linear-to-r from-violet-600 to-fuchsia-500 bg-clip-text text-transparent m-0 flex items-center gap-2"
    >
      {t("extras.title")}
    </h1>
  </div>
  <div class="flex items-center gap-3">
    <div
      class="inline-flex rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-[#0b0f19] p-1 shadow-xs"
    >
      <button
        type="button"
        class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {mainTab ===
        'nginx'
          ? 'bg-violet-600 text-white shadow-xs'
          : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}"
        onclick={() => (mainTab = "nginx")}>{t("extras.tab_nginx")}</button
      >
      <button
        type="button"
        class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {mainTab ===
        'ports'
          ? 'bg-violet-600 text-white shadow-xs'
          : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}"
        onclick={() => (mainTab = "ports")}>{t("extras.tab_ports")}</button
      >
      <button type="button" class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {mainTab === 'deploy_temp' ? 'bg-violet-600 text-white shadow-xs' : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}" onclick={() => (mainTab = "deploy_temp")}>{t("extras.tab_deploy_temp")}</button>
    </div>
    <ButtonPink size="sm" title={t("common.refresh")} onclick={refresh}
      >{t("common.refresh")}</ButtonPink
    >
  </div>
</div>
