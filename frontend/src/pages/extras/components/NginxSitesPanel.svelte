<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { notifyWarning } from "$shared/stores/notification.svelte";
  import {
    ButtonBlue,
    ButtonGreen,
    ButtonYellow,
    ButtonRed,
    ButtonIndigo,
    ButtonFuchsia,
    ButtonSky,
  } from "$shared/components/buttons";
  import NginxEditor from "./NginxEditor.svelte";

  let {
    site = $bindable(),
    content = $bindable(),
    activeTab = $bindable(),
    available = [],
    enabled = [],
    editorKey = 0,
    loadingFile = false,
    busy = null,
    onOpenSite,
    onNewSite,
    onRun,
    onRequestDelete,
    onViewLogs,
  }: {
    site: string;
    content: string;
    activeTab: "available" | "enabled";
    available: string[];
    enabled: string[];
    editorKey: number;
    loadingFile: boolean;
    busy: string | null;
    onOpenSite: (filename: string) => void;
    onNewSite: () => void;
    onRun: (action: "enable" | "test" | "restart" | "save") => void;
    onRequestDelete: () => void;
    onViewLogs: () => void;
  } = $props();
</script>

<section
  class="space-y-6 rounded-2xl border border-slate-200/80 dark:border-slate-800/80 bg-white dark:bg-[#0b0f19] p-6 shadow-sm"
>
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div
      class="inline-flex rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-100 dark:bg-slate-900 p-1"
    >
      <button
        type="button"
        class="px-4 py-1.5 rounded-lg text-xs font-bold {activeTab ===
        'available'
          ? 'bg-violet-600 text-white'
          : 'text-slate-600'}"
        onclick={() => (activeTab = "available")}
      >
        {t("extras.sites_available")} ({available.length})
      </button>
      <button
        type="button"
        class="px-4 py-1.5 rounded-lg text-xs font-bold {activeTab === 'enabled'
          ? 'bg-violet-600 text-white'
          : 'text-slate-600'}"
        onclick={() => (activeTab = "enabled")}
      >
        {t("extras.sites_enabled")} ({enabled.length})
      </button>
    </div>

    <div class="flex flex-wrap items-center gap-2">
      <ButtonSky
        size="sm"
        loading={busy === "test"}
        onclick={() => onRun("test")}
      >
        {t("extras.test_nginx")}
      </ButtonSky>
      <ButtonYellow
        size="sm"
        loading={busy === "restart"}
        onclick={() => onRun("restart")}
      >
        {t("extras.restart_nginx")}
      </ButtonYellow>
      <ButtonFuchsia size="sm" onclick={onViewLogs}>
        {t("extras.view_logs")}
      </ButtonFuchsia>
      <ButtonRed
        size="sm"
        disabled={!site.trim() || !!busy}
        onclick={() =>
          site.trim()
            ? onRequestDelete()
            : notifyWarning(t("extras.select_file_warn"))}
      >
        {t("extras.delete_file")}
      </ButtonRed>
      <ButtonBlue size="sm" onclick={onNewSite}>
        {t("extras.new_site")}
      </ButtonBlue>
    </div>
  </div>

  <div class="min-h-12 flex items-center">
    <div class="flex flex-wrap gap-2">
      {#each activeTab === "available" ? available : enabled as filename}
        <button
          type="button"
          class="px-3 py-1.5 rounded-xl text-xs font-mono border transition-colors {site ===
          filename
            ? 'border-violet-500 bg-violet-600 text-white shadow-md shadow-violet-500/25 dark:bg-violet-600 dark:text-white'
            : 'border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-400 hover:border-violet-400 hover:bg-violet-50 dark:hover:bg-violet-950/30'}"
          onclick={() => onOpenSite(filename)}
        >
          {filename}
        </button>
      {/each}
    </div>
  </div>

  <NginxEditor
    bind:site
    bind:content
    tab={activeTab}
    {editorKey}
    loading={loadingFile}
  />

  <div class="flex flex-wrap gap-2.5 pt-2">
    <ButtonIndigo
      size="md"
      loading={busy === "enable"}
      onclick={() => onRun("enable")}
    >
      {t("extras.btn_enable")}
    </ButtonIndigo>
    <ButtonGreen
      size="md"
      loading={busy === "save"}
      onclick={() => onRun("save")}
    >
      {t("extras.btn_save")}
    </ButtonGreen>
  </div>
</section>
