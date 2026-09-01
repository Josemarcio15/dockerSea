<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { viewModeStore } from "$shared/stores/viewMode.svelte";
  import type { Container } from "../types";
  import ContainerCard from "./ContainerCard.svelte";

  let {
    containers = [],
    selectedNames = [],
    onToggle = () => {},
    onOpenLogs = () => {},
  }: {
    containers?: Container[];
    selectedNames?: string[];
    onToggle?: (name: string) => void;
    onOpenLogs?: (name: string) => void;
  } = $props();
</script>

{#if containers.length === 0}
  <div class="text-sm text-slate-400 dark:text-slate-500 py-16 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-slate-50/20">
    {t("containers.empty")}
  </div>
{:else}
  <div class={viewModeStore.getGridClass()}>
    {#each containers as container (container.id)}
      <ContainerCard
        {container}
        checked={selectedNames.includes(container.name)}
        on_toggle={() => onToggle(container.name)}
        on_open_logs={onOpenLogs}
      />
    {/each}
  </div>
{/if}
