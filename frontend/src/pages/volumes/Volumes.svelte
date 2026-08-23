<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import { useRefreshKey, triggerRefresh } from "$lib/stores/refresh.svelte";
  import { notifySuccess, notifyError } from "$lib/stores/notification.svelte";
  import DockerseaLoading from "$lib/components/DockerseaLoading.svelte";
  import StatusBanner from "$lib/components/StatusBanner.svelte";
  import VolumeCard from "$lib/components/VolumeCard.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import VpsSelectWarning from "$lib/components/VpsSelectWarning.svelte";

  let { data } = $props();

  let selectedNames = $state<string[]>([]);
  let searchQuery = $state("");

  // Modal: Novo Volume
  let showCreateModal = $state(false);
  let modalName = $state("");
  let modalDriver = $state("");
  let labelEntries = $state<Array<{ key: string; value: string }>>([]);

  // Modal: Confirmar Limpeza
  let showPruneModal = $state(false);

  // Client-side fetched volumes
  let volumes = $state<any[]>([]);
  let loading = $state(true);
  let fetchError = $state("");

  // Reactive fetch: runs on mount and whenever triggerRefresh() is called
  async function fetchVolumes() {
    if (!data?.activeVps) {
      loading = false;
      return;
    }
    try {
      const res = await fetch("/api/volumes");
      const json = await res.json();
      volumes = json.volumes || [];
      fetchError = json.error || "";
    } catch (e: any) {
      fetchError = e.message;
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    useRefreshKey();
    fetchVolumes();
  });

  // Derived: Filtered volumes
  let filteredVolumes = $derived(
    (volumes || []).filter((v) =>
      v.name.toLowerCase().includes(searchQuery.toLowerCase()),
    ),
  );

  // Toggle single checkbox
  function toggleChecked(name: string) {
    if (selectedNames.includes(name)) {
      selectedNames = selectedNames.filter((n) => n !== name);
    } else {
      selectedNames = [...selectedNames, name];
    }
  }

  // Toggle all filtered volumes
  function toggleAll() {
    if (selectedNames.length === filteredVolumes.length) {
      selectedNames = [];
    } else {
      selectedNames = filteredVolumes.map((v) => v.name);
    }
  }

  // Submit Action helper
  async function runApiAction(payload: any) {
    try {
      const response = await fetch("/api/volumes", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      const result = await response.json();

      if (result.success) {
        notifySuccess(result.message);
        if (payload.action === "delete") {
          selectedNames = [];
        }
      } else {
        notifyError(result.message || "Erro ao processar volume");
      }

      await invalidateAll();
      triggerRefresh();
    } catch (e: any) {
      notifyError(e.message);
    }
  }

  // Execute prune (clean unused volumes)
  async function doPrune() {
    showPruneModal = false;
    notifySuccess("Limpando volumes não utilizados...");
    await runApiAction({ action: "prune" });
  }

  // Execute delete selected
  async function doDeleteSelected() {
    if (selectedNames.length === 0) return;
    notifySuccess(`Removendo ${selectedNames.length} volume(s)...`);
    await runApiAction({ action: "delete", names: selectedNames });
  }

  // Execute create volume
  async function doCreate() {
    if (!modalName.trim()) return;

    // Convert labelEntries list to record
    const labelsObj: Record<string, string> = {};
    for (const entry of labelEntries) {
      if (entry.key.trim()) {
        labelsObj[entry.key.trim()] = entry.value.trim();
      }
    }

    notifySuccess(`Criando volume '${modalName}'...`);
    await runApiAction({
      action: "create",
      name: modalName,
      driver: modalDriver,
      labels: labelsObj,
    });
    showCreateModal = false;
  }

  function addLabel() {
    labelEntries = [...labelEntries, { key: "", value: "" }];
  }

  function removeLabel(index: number) {
    labelEntries = labelEntries.filter((_, i) => i !== index);
  }
</script>

{#if !data.activeVps}
  <VpsSelectWarning />
{:else}
  <div class="space-y-6">
    <!-- Top Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div
        class="inline-flex items-center px-5 py-2.5 rounded-2xl bg-linear-to-br from-violet-100 to-fuchsia-100 dark:from-violet-950/40 dark:to-fuchsia-950/40 border border-violet-200/50 dark:border-violet-800/50 self-start shadow-sm"
      >
        <h1
          class="text-2xl font-bold bg-linear-to-r from-violet-600 to-fuchsia-500 bg-clip-text text-transparent m-0 flex items-center gap-2"
        >
          {t("volumes.title")}
        </h1>
      </div>

      <div class="flex items-center gap-2">
        <input
          type="text"
          placeholder={t("common.search")}
          class="px-4 py-2 text-xs rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-[#0b0f19] text-slate-850 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all w-60"
          bind:value={searchQuery}
        />
        <button
          type="button"
          class="px-4 py-2 text-xs rounded-xl border-none cursor-pointer font-bold text-white bg-blue-600 hover:bg-blue-700 transition-colors shadow-md shadow-blue-500/20 whitespace-nowrap"
          onclick={() => {
            modalName = "";
            modalDriver = "";
            labelEntries = [];
            showCreateModal = true;
          }}
        >
          + {t("volumes.new_volume")}
        </button>
      </div>
    </div>

    <!-- Status Alerts -->
    <StatusBanner />

    {#if loading}
      <DockerseaLoading message={t("common.loading")} />
    {:else if fetchError}
    <div
      class="p-6 rounded-2xl border border-red-200 dark:border-red-900/60 bg-red-50 dark:bg-red-950/20 text-red-700 dark:text-red-400"
    >
      <h3 class="font-bold text-sm mb-1">Erro de Conexão</h3>
      <p class="text-xs whitespace-pre-wrap">{fetchError}</p>
    </div>
  {:else}
    <!-- Action Toolbar -->
    <div
      class="flex flex-wrap items-center justify-between gap-3 bg-white dark:bg-[#0b0f19] border border-slate-200/80 dark:border-slate-800/80 p-3.5 rounded-2xl shadow-sm"
    >
      <div class="flex items-center gap-2">
        <button
          type="button"
          class="px-4 py-2.5 rounded-xl border-none cursor-pointer text-xs font-bold text-white bg-amber-500 hover:bg-amber-600 transition-colors shadow-md shadow-amber-500/20"
          onclick={toggleAll}
        >
          {selectedNames.length === filteredVolumes.length
            ? t("common.deselect_all")
            : t("common.select_all")}
        </button>

        <span class="text-xs text-slate-400 dark:text-slate-500 px-2">
          {filteredVolumes.length}
          {t("volumes.vol_count")}
          {#if selectedNames.length > 0}
            — <span class="font-semibold text-violet-600 dark:text-violet-400"
              >{selectedNames.length} {t("volumes.selected_count")}</span
            >
          {/if}
        </span>
      </div>

      <div class="flex items-center gap-2">
        <button
          type="button"
          class="px-4 py-2.5 rounded-xl border-none cursor-pointer text-xs font-bold text-white bg-emerald-600 hover:bg-emerald-700 transition-colors shadow-md shadow-emerald-500/20"
          onclick={() => (showPruneModal = true)}
        >
          {t("volumes.prune_btn")}
        </button>
        <button
          type="button"
          class="px-4 py-2.5 rounded-xl border-none cursor-pointer text-xs font-bold text-white transition-colors shadow-md shadow-red-500/10 {selectedNames.length >
          0
            ? 'bg-red-500 hover:bg-red-600'
            : 'bg-slate-400 cursor-not-allowed'}"
          disabled={selectedNames.length === 0}
          onclick={doDeleteSelected}
        >
          {t("volumes.delete_selected")}
        </button>
      </div>
    </div>

    <!-- Volumes Grid -->
    {#if filteredVolumes.length === 0}
      <div
        class="text-sm text-slate-400 dark:text-slate-500 py-16 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-slate-50/20"
      >
        {t("volumes.empty")}
      </div>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        {#each filteredVolumes as vol (vol.name)}
          <VolumeCard
            {vol}
            checked={selectedNames.includes(vol.name)}
            on_toggle={() => toggleChecked(vol.name)}
          />
        {/each}
      </div>
    {/if}
  {/if}
</div>
{/if}

<!-- Modal: Novo Volume -->
<Modal
  bind:show={showCreateModal}
  title={t("volumes.modal_title")}
  buttons={[
    {
      label: t("volumes.create_btn"),
      variant: "primary",
      onclick: doCreate,
      disabled: !modalName.trim(),
    },
  ]}
>
  <div class="flex flex-col gap-1.5">
    <label
      for="vol-name"
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >
      {t("volumes.field_name")}
    </label>
    <input
      id="vol-name"
      type="text"
      class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
      placeholder="meu-volume"
      bind:value={modalName}
    />
  </div>

  <div class="flex flex-col gap-1.5">
    <label
      for="vol-driver"
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >
      {t("volumes.field_driver")}
    </label>
    <input
      id="vol-driver"
      type="text"
      class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
      placeholder="local (padrão)"
      bind:value={modalDriver}
    />
  </div>

  <div class="flex flex-col gap-2">
    <div class="flex items-center justify-between">
      <span
        class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
      >
        {t("volumes.field_labels")}
      </span>
      <button
        type="button"
        class="text-xs text-violet-600 dark:text-violet-400 font-bold cursor-pointer bg-transparent border-none hover:underline"
        onclick={addLabel}
      >
        {t("volumes.add_label")}
      </button>
    </div>

    <div class="space-y-2 max-h-40 overflow-y-auto pr-1">
      {#each labelEntries as entry, idx}
        <div class="flex gap-2 items-center">
          <input
            type="text"
            class="flex-1 px-3 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
            placeholder={t("volumes.placeholder_key")}
            bind:value={entry.key}
          />
          <input
            type="text"
            class="flex-1 px-3 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
            placeholder={t("volumes.placeholder_value")}
            bind:value={entry.value}
          />
          <button
            type="button"
            class="p-2 rounded-xl text-red-500 hover:bg-red-50 dark:hover:bg-red-950/20 cursor-pointer transition-colors border-none bg-transparent"
            onclick={() => removeLabel(idx)}
          >
            ✕
          </button>
        </div>
      {/each}
    </div>
  </div>
</Modal>

<!-- Modal: Confirmar Limpeza -->
<Modal
  bind:show={showPruneModal}
  title={t("volumes.prune_confirm_title")}
  buttons={[
    {
      label: t("volumes.prune_confirm_btn"),
      variant: "success",
      onclick: doPrune,
    },
  ]}
>
  <p class="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
    {t("volumes.prune_confirm_msg")}
  </p>
</Modal>
