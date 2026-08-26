<script lang="ts">
  import { onMount } from "svelte";
  import { t } from "$shared/stores/locale.svelte";
  import StatusBanner from "$shared/components/StatusBanner.svelte";
  import VpsSelectWarning from "$shared/components/VpsSelectWarning.svelte";
  import { Events } from "@wailsio/runtime";
  import { builderStore as store } from "./store.svelte";
  import FolderBrowser from "./components/FolderBrowser.svelte";
  import BuildControls from "./components/BuildControls.svelte";
  import BuildLogs from "./components/BuildLogs.svelte";

  let { data } = $props();

  function goToImages() {
    window.location.href = `/images?highlight=${encodeURIComponent(store.builtImage)}`;
  }

  onMount(() => {
    const unsubscribeProgress = Events.On("builder:progress", (event) => {
      if (event.data?.line) store.appendLog(event.data.line);
    });
    const unsubscribeComplete = Events.On("builder:complete", (event) => {
      store.completeBuild(
        event.data as { success?: boolean; image?: string; message?: string },
      );
    });

    store.loadSavedPaths();
    store.browse();

    return () => {
      unsubscribeProgress();
      unsubscribeComplete();
    };
  });
</script>

{#if !data.activeVps}
  <VpsSelectWarning />
{:else}
  <div class="space-y-6">
    <div
      class="flex flex-col sm:flex-row sm:items-center justify-between gap-4"
    >
      <div
        class="inline-flex items-center px-5 py-2.5 rounded-2xl bg-linear-to-br from-violet-100 to-fuchsia-100 dark:from-violet-950/40 dark:to-fuchsia-950/40 border border-violet-200/50 dark:border-violet-800/50 self-start shadow-sm"
      >
        <h1
          class="text-2xl font-bold bg-linear-to-r from-violet-600 to-fuchsia-500 bg-clip-text text-transparent m-0 flex items-center gap-2"
        >
          {t("sidebar.builder")}
        </h1>
      </div>
    </div>

    <StatusBanner />

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="space-y-4">
        <FolderBrowser {store} />
        <BuildControls {store} {goToImages} />
      </div>
      <BuildLogs logs={store.logs} status={store.status} />
    </div>
  </div>
{/if}
