<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import { useRefreshKey } from "$lib/stores/refresh.svelte";
  import DockerseaLoading from "$lib/components/DockerseaLoading.svelte";
  import StatusBanner from "$lib/components/StatusBanner.svelte";
  import VpsSelectWarning from "$lib/components/VpsSelectWarning.svelte";

  import NetworkToolbar from "./components/NetworkToolbar.svelte";
  import NetworkCard from "./components/NetworkCard.svelte";
  import CreateNetworkModal from "./components/CreateNetworkModal.svelte";
  import ConnectContainerModal from "./components/ConnectContainerModal.svelte";
  import PruneNetworksModal from "./components/PruneNetworksModal.svelte";
  import { createNetworksState } from "./useNetworks.svelte.js";

  let { data } = $props();

  const netState = createNetworksState(() => data?.activeVps);

  $effect(() => {
    useRefreshKey();
    if (data?.activeVps?.id) {
      netState.fetchAll();
    }
  });
</script>

{#if !data.activeVps}
  <VpsSelectWarning />
{:else}
  <div class="space-y-6">
    <!-- Toolbar / Header -->
    <NetworkToolbar
      bind:searchQuery={netState.searchQuery}
      selectedCount={netState.selectedNetworkNames.length}
      totalFilteredCount={netState.filteredNetworks.length}
      allSelected={netState.selectedNetworkNames.length ===
        netState.filteredNetworks.length &&
        netState.filteredNetworks.length > 0}
      onRefresh={() => netState.fetchAll()}
      onNewNetwork={netState.openCreateModal}
      onToggleAll={netState.toggleAll}
      onPrune={() => (netState.showPruneModal = true)}
      onDeleteSelected={netState.doDeleteSelected}
    />

    <!-- Status Alerts -->
    <StatusBanner />

    {#if netState.loading}
      <DockerseaLoading message={t("common.loading")} />
    {:else if netState.fetchError}
      <div
        class="p-6 rounded-2xl border border-red-200 dark:border-red-900/60 bg-red-50 dark:bg-red-950/20 text-red-700 dark:text-red-400"
      >
        <h3 class="font-bold text-sm mb-1">{t("common.error")}</h3>
        <p class="text-xs whitespace-pre-wrap">{netState.fetchError}</p>
      </div>
    {:else}
      <!-- Networks Grid -->
      {#if netState.filteredNetworks.length === 0}
        <div
          class="text-sm text-slate-400 dark:text-slate-500 py-16 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-slate-50/20"
        >
          {t("networks.empty")}
        </div>
      {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          {#each netState.filteredNetworks as network (network.id)}
            <NetworkCard
              {network}
              checked={netState.selectedNetworkNames.includes(network.name)}
              on_toggle={() => netState.handleToggleSelect(network.name)}
              on_delete={() => netState.doDelete(network.name)}
              on_disconnect={(cName) =>
                netState.doDisconnect(network.name, cName)}
              on_connect={() => netState.openConnectModal(network.name)}
            />
          {/each}
        </div>
      {/if}
    {/if}
  </div>
{/if}

<!-- Modals -->
<CreateNetworkModal
  bind:show={netState.showCreateModal}
  bind:name={netState.modalName}
  bind:driver={netState.modalDriver}
  bind:subnet={netState.modalSubnet}
  bind:gateway={netState.modalGateway}
  onSubmit={netState.doCreate}
/>

<ConnectContainerModal
  bind:show={netState.showConnectModal}
  networkName={netState.selectedNetworkName}
  bind:selectedContainer={netState.selectedContainerName}
  containers={netState.containers}
  onSubmit={netState.doConnect}
/>

<PruneNetworksModal
  bind:show={netState.showPruneModal}
  onConfirm={netState.doPrune}
/>
