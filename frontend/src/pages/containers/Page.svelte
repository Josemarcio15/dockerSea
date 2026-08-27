<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { useRefreshKey } from "$shared/stores/refresh.svelte";
  import DockerseaLoading from "$shared/components/DockerseaLoading.svelte";
  import StatusBanner from "$shared/components/StatusBanner.svelte";
  import VpsSelectWarning from "$shared/components/VpsSelectWarning.svelte";
  import TerminalModal from "$shared/components/TerminalModal.svelte";
  import ConfirmDialog from "$shared/components/ConfirmDialog.svelte";
  import ContainerToolbar from "./components/ContainerToolbar.svelte";
  import ContainerGrid from "./components/ContainerGrid.svelte";
  import { createContainersStore } from "./store.svelte.js";

  let { data } = $props();
  const store = createContainersStore(() => data?.activeVps);

  $effect(() => {
    useRefreshKey();
    if (data?.activeVps?.id) {
      store.fetchAll();
      return store.setupEventsStream();
    }
  });

  let showRemoveConfirm = $state(false);
  function requestRemove() {
    if (store.selectedNames.length > 0) showRemoveConfirm = true;
  }
</script>

{#if !data?.activeVps}
  <VpsSelectWarning />
{:else}
  <div class="space-y-6">
    <ContainerToolbar
      bind:searchQuery={store.searchQuery}
      selectedCount={store.selectedNames.length}
      totalCount={store.filteredContainers.length}
      allSelected={store.selectedNames.length ===
        store.filteredContainers.length && store.filteredContainers.length > 0}
      onRefresh={() => store.fetchAll()}
      onToggleAll={store.toggleAll}
      onStart={() => store.doActionSelected("start")}
      onStop={() => store.doActionSelected("stop")}
      onRestart={() => store.doActionSelected("restart")}
      onRemove={requestRemove}
    />
    <StatusBanner />
    {#if store.loading}
      <DockerseaLoading message={t("common.loading")} />
    {:else if store.fetchError}
      <div
        class="p-6 rounded-2xl border border-red-200 dark:border-red-900/60 bg-red-50 dark:bg-red-950/20 text-red-700 dark:text-red-400"
      >
        <h3 class="font-bold text-sm mb-1">{t("common.error")}</h3>
        <p class="text-xs whitespace-pre-wrap">{store.fetchError}</p>
      </div>
    {:else}
      <ContainerGrid
        containers={store.filteredContainers}
        selectedNames={store.selectedNames}
        onToggle={(name) => store.handleToggleSelect(name)}
        onOpenLogs={store.openLogs}
      />
    {/if}
  </div>
{/if}

<TerminalModal
  bind:show={store.showLogs}
  title={store.logsTitle}
  loading={store.logsLoading}
  logs={store.logsContent}
/>
<ConfirmDialog
  bind:show={showRemoveConfirm}
  title="Remover Contêiner(es)"
  message={`Tem certeza de que deseja remover ${store.selectedNames.length} contêiner(es) selecionado(s)?\n[${store.selectedNames.join(", ")}]\nEssa ação apagará os contêineres e seus dados não persistidos.`}
  confirmText="Remover Contêiner(es)"
  type="danger"
  onConfirm={() => store.doActionSelected("rm")}
/>
