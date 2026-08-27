<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { notifyWarning } from "$shared/stores/notification.svelte";
  import { useRefreshKey, triggerRefresh } from "$shared/stores/refresh.svelte";
  import DockerseaLoading from "$shared/components/DockerseaLoading.svelte";
  import StatusBanner from "$shared/components/StatusBanner.svelte";
  import VpsSelectWarning from "$shared/components/VpsSelectWarning.svelte";
  import PortsPanel from "$shared/components/PortsPanel.svelte";
  import NginxLogsModal from "$shared/components/NginxLogsModal.svelte";
  import ExtrasHeader from "./components/ExtrasHeader.svelte";
  import NginxEditor from "./components/NginxEditor.svelte";
  import ConfirmDialog from "$shared/components/ConfirmDialog.svelte";
  import { extrasStore as store } from "./store.svelte";
  import * as api from "./api";
  import { ButtonPurple, ButtonBlue, ButtonGreen, ButtonYellow, ButtonRed, ButtonIndigo, ButtonFuchsia, ButtonSky, Button } from "$shared/components/buttons";

  let { data } = $props();
  let mainTab = $state<"nginx" | "ports" | "deploy_temp">("nginx");
  let deployFiles = $state<any[]>([]);
  let deployLoading = $state(false);
  let deployPath = $state("$HOME/.docksea");
  let selectedDeployPaths = $state<string[]>([]);
  let showDeleteDeployDialog = $state(false);
  let showDeleteDialog = $state(false);
  let showLogs = $state(false);
  const site = $derived(store.site);
  const busy = $derived(store.busy);
  const activeTab = $derived(store.activeTab);
  const available = $derived(store.available);
  const enabled = $derived(store.enabled);
  const loadingSites = $derived(store.loading);
  const fetchError = $derived(store.error);
  const editorKey = $derived(store.editorKey);
  const loadSites = (silent = false) => store.load(silent);
  const openSite = (filename: string) => store.open(filename);
  const newSite = () => store.newSite();
  const run = (action: "enable" | "test" | "restart" | "save") => store.run(action);
  const deleteSite = () => store.remove();
  async function loadDeployFiles() {
    if (!data?.activeVps) return;
    deployLoading = true;
    try { deployFiles = (await api.listDeployTempFilesAt(data.activeVps, deployPath)) || []; } finally { deployLoading = false; }
  }
  async function openDeployFolder(path: string) { deployPath = path; await loadDeployFiles(); }
  async function goDeployParent() { const parts = deployPath.split("/"); if (parts.length > 4) { parts.pop(); deployPath = parts.join("/"); await loadDeployFiles(); } }
  function toggleDeployPath(path: string) { selectedDeployPaths = selectedDeployPaths.includes(path) ? selectedDeployPaths.filter((item) => item !== path) : [...selectedDeployPaths, path]; }
  function requestDeleteSelectedDeployPaths() { if (selectedDeployPaths.length) showDeleteDeployDialog = true; }
  async function deleteSelectedDeployPaths() {
    for (const path of selectedDeployPaths) await api.deleteDeployTempPath(data.activeVps, path);
    selectedDeployPaths = [];
    showDeleteDeployDialog = false;
    await loadDeployFiles();
  }

  $effect(() => {
    useRefreshKey();
    store.setVps(data?.activeVps);
    store.load();
  });

  $effect(() => {
    if (mainTab === "deploy_temp") void loadDeployFiles();
  });
</script>

{#if !data.activeVps}
  <VpsSelectWarning />
{:else}
  <div class="space-y-6">
    <ExtrasHeader bind:mainTab refresh={() => (mainTab === "nginx" ? loadSites() : mainTab === "deploy_temp" ? loadDeployFiles() : triggerRefresh())} />
    <StatusBanner />
    {#if mainTab === "ports"}
      <section class="w-full bg-white dark:bg-[#0b0f19] border border-slate-200/80 dark:border-slate-800/80 rounded-2xl p-5 shadow-sm"><PortsPanel activeVps={data.activeVps} /></section>
    {:else if mainTab === "deploy_temp"}
      <section class="space-y-4 rounded-2xl border border-slate-200/80 dark:border-slate-800/80 bg-white dark:bg-[#0b0f19] p-6 shadow-sm">
        <div class="flex items-center justify-between gap-3"><div><h2 class="text-base font-bold">{t("extras.deploy_temp_title")}</h2><p class="text-xs text-slate-500 font-mono break-all">{deployPath}</p></div><div class="flex gap-2"><Button size="sm" disabled={deployPath === "$HOME/.docksea"} onclick={goDeployParent}>Voltar</Button><ButtonRed size="sm" disabled={!selectedDeployPaths.length} onclick={requestDeleteSelectedDeployPaths}>Remover selecionados</ButtonRed></div></div>
        {#if deployFiles.length === 0}<p class="text-sm text-slate-400 py-8 text-center">{t("extras.deploy_temp_empty")}</p>{:else}<div class="space-y-2">{#each deployFiles as file}<div class="w-full flex items-center gap-3 rounded-xl border border-slate-200 dark:border-slate-800 p-3 text-xs hover:bg-slate-50 dark:hover:bg-slate-900/50"><input type="checkbox" checked={selectedDeployPaths.includes(file.path)} onchange={() => toggleDeployPath(file.path)} /><button type="button" class="flex flex-1 justify-between gap-3 text-left" onclick={() => file.isDir && openDeployFolder(file.path)}><span class="font-mono break-all">{file.isDir ? "📁 " : "📄 "}{file.path.split('/').pop()}</span><span class="shrink-0 text-slate-500">{file.isDir ? "Pasta" : `${file.size} bytes`}</span></button></div>{/each}</div>{/if}
      </section>
    {:else if loadingSites}
      <DockerseaLoading message={t("common.loading")} />
    {:else if fetchError}
      <div class="p-6 rounded-2xl border border-red-200 dark:border-red-900/60 bg-red-50 dark:bg-red-950/20 text-red-700 dark:text-red-400"><h3 class="font-bold text-sm mb-1">Erro de Conexão Nginx</h3><p class="text-xs whitespace-pre-wrap">{fetchError}</p></div>
    {:else}
      <section class="space-y-6 rounded-2xl border border-slate-200/80 dark:border-slate-800/80 bg-white dark:bg-[#0b0f19] p-6 shadow-sm">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div class="inline-flex rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-100 dark:bg-slate-900 p-1">
            <button type="button" class="px-4 py-1.5 rounded-lg text-xs font-bold {activeTab === 'available' ? 'bg-violet-600 text-white' : 'text-slate-600'}" onclick={() => (store.activeTab = "available")}>{t("extras.sites_available")} ({available.length})</button>
            <button type="button" class="px-4 py-1.5 rounded-lg text-xs font-bold {activeTab === 'enabled' ? 'bg-violet-600 text-white' : 'text-slate-600'}" onclick={() => (store.activeTab = "enabled")}>{t("extras.sites_enabled")} ({enabled.length})</button>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <ButtonSky size="sm" loading={busy === "test"} onclick={() => run("test")}>{t("extras.test_nginx")}</ButtonSky>
            <ButtonYellow size="sm" loading={busy === "restart"} onclick={() => run("restart")}>{t("extras.restart_nginx")}</ButtonYellow>
            <ButtonFuchsia size="sm" onclick={() => (showLogs = true)}>{t("extras.view_logs")}</ButtonFuchsia>
            <ButtonRed size="sm" disabled={!site.trim() || !!busy} onclick={() => site.trim() ? (showDeleteDialog = true) : notifyWarning(t("extras.select_file_warn"))}>{t("extras.delete_file")}</ButtonRed>
            <ButtonBlue size="sm" onclick={newSite}>{t("extras.new_site")}</ButtonBlue>
          </div>
        </div>
        <div class="min-h-12 flex items-center"><div class="flex flex-wrap gap-2">{#each (activeTab === "available" ? available : enabled) as filename}<button type="button" class="px-3 py-1.5 rounded-xl text-xs font-mono border transition-colors {site === filename ? 'border-violet-500 bg-violet-600 text-white shadow-md shadow-violet-500/25 dark:bg-violet-600 dark:text-white' : 'border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-400 hover:border-violet-400 hover:bg-violet-50 dark:hover:bg-violet-950/30'}" onclick={() => openSite(filename)}>{filename}</button>{/each}</div></div>
        <NginxEditor bind:site={store.site} bind:content={store.content} tab={activeTab} {editorKey} />
        <div class="flex flex-wrap gap-2.5 pt-2"><ButtonIndigo size="md" loading={busy === "enable"} onclick={() => run("enable")}>{t("extras.btn_enable")}</ButtonIndigo><ButtonGreen size="md" loading={busy === "save"} onclick={() => run("save")}>{t("extras.btn_save")}</ButtonGreen></div>
      </section>
    {/if}
  </div>
{/if}

<NginxLogsModal bind:show={showLogs} activeVps={data.activeVps} />
<ConfirmDialog bind:show={showDeleteDeployDialog} title="Remover arquivos temporários" message={`Tem certeza de que deseja remover ${selectedDeployPaths.length} item(ns) selecionado(s)?`} confirmText="Remover selecionados" type="danger" onConfirm={deleteSelectedDeployPaths} />
{#if showDeleteDialog}<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4"><div class="w-full max-w-md rounded-2xl bg-white dark:bg-[#0c1220] p-6 space-y-4"><h2 class="text-base font-bold">{t("extras.delete_confirm_title")}</h2><p class="text-xs">{t("extras.delete_confirm_msg", { site })}</p><div class="flex justify-end gap-2.5"><Button size="sm" onclick={() => (showDeleteDialog = false)}>{t("common.cancel")}</Button><ButtonRed size="sm" onclick={deleteSite}>{t("common.delete")}</ButtonRed></div></div></div>{/if}
