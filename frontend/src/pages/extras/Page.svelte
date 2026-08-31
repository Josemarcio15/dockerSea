<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { useRefreshKey, triggerRefresh } from "$shared/stores/refresh.svelte";
  import DockerseaLoading from "$shared/components/DockerseaLoading.svelte";
  import StatusBanner from "$shared/components/StatusBanner.svelte";
  import VpsSelectWarning from "$shared/components/VpsSelectWarning.svelte";
  import PortsPanel from "$shared/components/PortsPanel.svelte";
  import NginxLogsModal from "$shared/components/NginxLogsModal.svelte";
  import ConfirmDialog from "$shared/components/ConfirmDialog.svelte";
  import { Button, ButtonRed } from "$shared/components/buttons";

  import ExtrasHeader from "./components/ExtrasHeader.svelte";
  import NginxSitesPanel from "./components/NginxSitesPanel.svelte";
  import DeployTempPanel from "./components/DeployTempPanel.svelte";

  import { extrasStore as store } from "./store.svelte";
  import * as api from "./api";

  let { data } = $props();
  let mainTab = $state<"nginx" | "ports" | "deploy_temp">("nginx");

  // Estado Deploy Temp
  let deployFiles = $state<any[]>([]);
  let deployLoading = $state(false);
  let deployPath = $state("$HOME/.docksea");
  let selectedDeployPaths = $state<string[]>([]);
  let showDeleteDeployDialog = $state(false);

  // Modais
  let showDeleteDialog = $state(false);
  let showLogs = $state(false);

  // Derivações reativas da store
  const site = $derived(store.site);
  const busy = $derived(store.busy);
  const available = $derived(store.available);
  const enabled = $derived(store.enabled);
  const loadingSites = $derived(store.loading);
  const loadingFile = $derived(store.loadingFile);
  const fetchError = $derived(store.error);
  const editorKey = $derived(store.editorKey);

  // Ações de Deploy Temp
  async function loadDeployFiles() {
    if (!data?.activeVps) return;
    deployLoading = true;
    try {
      deployFiles = (await api.listDeployTempFilesAt(data.activeVps, deployPath)) || [];
    } finally {
      deployLoading = false;
    }
  }

  async function openDeployFolder(path: string) {
    deployPath = path;
    await loadDeployFiles();
  }

  async function goDeployParent() {
    const parts = deployPath.split("/");
    if (parts.length > 4) {
      parts.pop();
      deployPath = parts.join("/");
      await loadDeployFiles();
    }
  }

  async function deleteSelectedDeployPaths() {
    for (const path of selectedDeployPaths) {
      await api.deleteDeployTempPath(data.activeVps, path);
    }
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
    <ExtrasHeader
      bind:mainTab
      refresh={() =>
        mainTab === "nginx"
          ? store.load()
          : mainTab === "deploy_temp"
            ? loadDeployFiles()
            : triggerRefresh()}
    />

    <StatusBanner />

    {#if mainTab === "ports"}
      <section
        class="w-full bg-white dark:bg-[#0b0f19] border border-slate-200/80 dark:border-slate-800/80 rounded-2xl p-5 shadow-sm"
      >
        <PortsPanel activeVps={data.activeVps} />
      </section>
    {:else if mainTab === "deploy_temp"}
      {#if deployLoading}
        <DockerseaLoading message={t("common.loading")} />
      {:else}
        <DeployTempPanel
          bind:deployPath
          {deployFiles}
          bind:selectedDeployPaths
          onGoParent={goDeployParent}
          onRequestDelete={() => {
            if (selectedDeployPaths.length) showDeleteDeployDialog = true;
          }}
          onOpenFolder={openDeployFolder}
        />
      {/if}
    {:else if loadingSites}
      <DockerseaLoading message={t("common.loading")} />
    {:else if fetchError}
      <div
        class="p-6 rounded-2xl border border-red-200 dark:border-red-900/60 bg-red-50 dark:bg-red-950/20 text-red-700 dark:text-red-400"
      >
        <h3 class="font-bold text-sm mb-1">Erro de Conexão Nginx</h3>
        <p class="text-xs whitespace-pre-wrap">{fetchError}</p>
      </div>
    {:else}
      <NginxSitesPanel
        bind:site={store.site}
        bind:content={store.content}
        bind:activeTab={store.activeTab}
        {available}
        {enabled}
        {editorKey}
        {loadingFile}
        {busy}
        onOpenSite={(f) => store.open(f)}
        onNewSite={() => store.newSite()}
        onRun={(action) => store.run(action)}
        onRequestDelete={() => (showDeleteDialog = true)}
        onViewLogs={() => (showLogs = true)}
      />
    {/if}
  </div>
{/if}

<NginxLogsModal bind:show={showLogs} activeVps={data.activeVps} />

<ConfirmDialog
  bind:show={showDeleteDeployDialog}
  title="Remover arquivos temporários"
  message={`Tem certeza de que deseja remover ${selectedDeployPaths.length} item(ns) selecionado(s)?`}
  confirmText="Remover selecionados"
  type="danger"
  onConfirm={deleteSelectedDeployPaths}
/>

{#if showDeleteDialog}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4"
  >
    <div
      class="w-full max-w-md rounded-2xl bg-white dark:bg-[#0c1220] p-6 space-y-4"
    >
      <h2 class="text-base font-bold">{t("extras.delete_confirm_title")}</h2>
      <p class="text-xs">{t("extras.delete_confirm_msg", { site })}</p>
      <div class="flex justify-end gap-2.5">
        <Button size="sm" onclick={() => (showDeleteDialog = false)}>
          {t("common.cancel")}
        </Button>
        <ButtonRed size="sm" onclick={() => store.remove()}>
          {t("common.delete")}
        </ButtonRed>
      </div>
    </div>
  </div>
{/if}
