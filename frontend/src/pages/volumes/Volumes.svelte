<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { useRefreshKey } from "$shared/stores/refresh.svelte";
  import DockerseaLoading from "$shared/components/DockerseaLoading.svelte";
  import StatusBanner from "$shared/components/StatusBanner.svelte";
  import VpsSelectWarning from "$shared/components/VpsSelectWarning.svelte";

  import VolumeToolbar from "./components/VolumeToolbar.svelte";
  import VolumeCard from "./components/VolumeCard.svelte";
  import CreateVolumeModal from "./components/CreateVolumeModal.svelte";
  import PruneVolumesModal from "./components/PruneVolumesModal.svelte";
  import { createVolumesState } from "./useVolumes.svelte.js";

  let { data } = $props();

  const vState = createVolumesState(() => data?.activeVps);

  $effect(() => {
    useRefreshKey();
    if (data?.activeVps?.id) {
      vState.fetchAll();
    }
  });
</script>

{#if !data.activeVps}
  <VpsSelectWarning />
{:else}
  <div class="space-y-6">
    <!-- Toolbar / Header -->
    <VolumeToolbar
      bind:searchQuery={vState.searchQuery}
      selectedCount={vState.selectedNames.length}
      totalCount={vState.filteredVolumes.length}
      allSelected={vState.selectedNames.length ===
        vState.filteredVolumes.length && vState.filteredVolumes.length > 0}
      onRefresh={() => vState.fetchAll()}
      onNewVolume={vState.openCreateModal}
      onToggleAll={vState.toggleAll}
      onPrune={() => (vState.showPruneModal = true)}
      onDeleteSelected={vState.doDeleteSelected}
    />

    <!-- Status Alerts -->
    <StatusBanner />

    {#if vState.loading}
      <DockerseaLoading message={t("common.loading")} />
    {:else if vState.fetchError}
      <div
        class="p-6 rounded-2xl border border-red-200 dark:border-red-900/60 bg-red-50 dark:bg-red-950/20 text-red-700 dark:text-red-400"
      >
        <h3 class="font-bold text-sm mb-1">{t("common.error")}</h3>
        <p class="text-xs whitespace-pre-wrap">{vState.fetchError}</p>
      </div>
    {:else}
      <!-- Volumes Grid -->
      {#if vState.filteredVolumes.length === 0}
        <div
          class="text-sm text-slate-400 dark:text-slate-500 py-16 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-slate-50/20"
        >
          {t("volumes.empty")}
        </div>
      {:else}
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-3.5">
          {#each vState.filteredVolumes as vol (vol.name)}
            <VolumeCard
              {vol}
              server={data.activeVps}
              checked={vState.selectedNames.includes(vol.name)}
              on_toggle={() => vState.toggleChecked(vol.name)}
            />
          {/each}
        </div>
      {/if}
    {/if}
  </div>
{/if}

<!-- Modals -->
<CreateVolumeModal
  bind:show={vState.showCreateModal}
  bind:name={vState.modalName}
  bind:driver={vState.modalDriver}
  bind:labelEntries={vState.labelEntries}
  onAddLabel={vState.addLabel}
  onRemoveLabel={vState.removeLabel}
  onSubmit={vState.doCreate}
/>

<PruneVolumesModal
  bind:show={vState.showPruneModal}
  onConfirm={vState.doPrune}
/>
