<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import { useRefreshKey } from "$lib/stores/refresh.svelte";
  import DockerseaLoading from "$lib/components/DockerseaLoading.svelte";
  import StatusBanner from "$lib/components/StatusBanner.svelte";
  import VpsSelectWarning from "$lib/components/VpsSelectWarning.svelte";
  import TerminalModal from "$lib/components/TerminalModal.svelte";

  import ContainerToolbar from "./components/ContainerToolbar.svelte";
  import ContainerCard from "./components/ContainerCard.svelte";
  import { createContainersState } from "./useContainers.svelte.js";

  let { data } = $props();

  const cState = createContainersState(() => data?.activeVps);

  $effect(() => {
    useRefreshKey();
    if (data?.activeVps?.id) {
      cState.fetchAll();
      return cState.setupEventsStream();
    }
  });
</script>

{#if !data.activeVps}
  <VpsSelectWarning />
{:else}
  <div class="space-y-6">
    <!-- Header / Toolbar -->
    <ContainerToolbar
      bind:searchQuery={cState.searchQuery}
      selectedCount={cState.selectedNames.length}
      totalCount={cState.filteredContainers.length}
      allSelected={cState.selectedNames.length ===
        cState.filteredContainers.length &&
        cState.filteredContainers.length > 0}
      onRefresh={() => cState.fetchAll()}
      onToggleAll={cState.toggleAll}
      onStart={() => cState.doActionSelected("start")}
      onStop={() => cState.doActionSelected("stop")}
      onRestart={() => cState.doActionSelected("restart")}
      onRemove={() => cState.doActionSelected("rm")}
    />

    <!-- Status Alerts -->
    <StatusBanner />

    {#if cState.loading}
      <DockerseaLoading message={t("common.loading")} />
    {:else if cState.fetchError}
      <div
        class="p-6 rounded-2xl border border-red-200 dark:border-red-900/60 bg-red-50 dark:bg-red-950/20 text-red-700 dark:text-red-400"
      >
        <h3 class="font-bold text-sm mb-1">{t("common.error")}</h3>
        <p class="text-xs whitespace-pre-wrap">{cState.fetchError}</p>
      </div>
    {:else}
      <!-- Containers Grid -->
      {#if cState.filteredContainers.length === 0}
        <div
          class="text-sm text-slate-400 dark:text-slate-500 py-16 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-slate-50/20"
        >
          {t("containers.empty")}
        </div>
      {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          {#each cState.filteredContainers as container (container.id)}
            <ContainerCard
              {container}
              checked={cState.selectedNames.includes(container.name)}
              on_toggle={() => cState.handleToggleSelect(container.name)}
              on_open_logs={(name) => cState.openLogs(name)}
            />
          {/each}
        </div>
      {/if}
    {/if}
  </div>
{/if}

<!-- Logs Modal -->
<TerminalModal
  bind:show={cState.showLogs}
  title={cState.logsTitle}
  loading={cState.logsLoading}
  lines={cState.logsContent}
/>
