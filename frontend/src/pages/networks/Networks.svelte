<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import { useRefreshKey, triggerRefresh } from "$lib/stores/refresh.svelte";
  import { notifySuccess, notifyError } from "$lib/stores/notification.svelte";
  import DockerseaLoading from "$lib/components/DockerseaLoading.svelte";
  import StatusBanner from "$lib/components/StatusBanner.svelte";
  import NetworkCard from "$lib/components/NetworkCard.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import VpsSelectWarning from "$lib/components/VpsSelectWarning.svelte";

  let { data } = $props();

  let searchQuery = $state("");

  // Modal: Nova Rede
  let showCreateModal = $state(false);
  let modalName = $state("");
  let modalDriver = $state("bridge");

  // Modal: Conectar Container
  let showConnectModal = $state(false);
  let selectedNetworkName = $state("");
  let selectedContainerName = $state("");

  // Modal: Confirmar Limpeza
  let showPruneModal = $state(false);

  // Client-side selected networks
  let selectedNetworkNames = $state<string[]>([]);

  // Client-side fetched data
  let networks = $state<any[]>([]);
  let containers = $state<any[]>([]);
  let loading = $state(true);
  let fetchError = $state("");
  let statusMsg = $state("");
  let statusTipo = $state("");

  // Reactive fetch: runs on mount and whenever triggerRefresh() is called
  async function fetchNetworks() {
    if (!data?.activeVps) {
      loading = false;
      return;
    }
    try {
      const res = await fetch("/api/networks");
      const json = await res.json();
      networks = json.networks || [];
      containers = json.containers || [];
      fetchError = json.error || "";
    } catch (e: any) {
      fetchError = e.message;
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    useRefreshKey();
    fetchNetworks();
  });

  // Derived: Filtered networks
  let filteredNetworks = $derived(
    (networks || [])
      .filter(
        (n) => n.name !== "bridge" && n.name !== "host" && n.name !== "none",
      )
      .filter((n) => n.name.toLowerCase().includes(searchQuery.toLowerCase())),
  );

  let selectableFilteredNetworks = $derived(filteredNetworks);

  function toggleAll() {
    if (selectedNetworkNames.length === selectableFilteredNetworks.length) {
      selectedNetworkNames = [];
    } else {
      selectedNetworkNames = selectableFilteredNetworks.map((n) => n.name);
    }
  }

  function handleToggleSelect(name: string) {
    if (selectedNetworkNames.includes(name)) {
      selectedNetworkNames = selectedNetworkNames.filter((n) => n !== name);
    } else {
      selectedNetworkNames = [...selectedNetworkNames, name];
    }
  }

  // Submit Action helper
  async function runAction(
    action: "delete" | "prune" | "create" | "connect" | "disconnect",
    formData: FormData,
  ) {
    try {
      const res = await fetch(`?/${action}`, {
        method: "POST",
        body: formData,
      });
      const result = deserialize(await res.text()) as any;

      let parsedData: any = null;
      if (result && typeof result === "object") {
        if (result.type === "success" && result.data) {
          parsedData = result.data;
        } else if (result.data) {
          parsedData = result.data;
        } else {
          parsedData = result;
        }
      }

      if (parsedData && parsedData.success) {
        notifySuccess(parsedData.message);
      } else {
        notifyError(
          parsedData?.message ||
            parsedData?.error ||
            result?.error?.message ||
            result?.error ||
            result?.message ||
            (result ? JSON.stringify(result).slice(0, 300) : "Resposta vazia"),
        );
      }

      await invalidateAll();
      triggerRefresh();
    } catch (e: any) {
      notifyError(e.message);
    }
  }

  async function doCreate() {
    if (!modalName.trim()) return;
    statusMsg = t("networks.creating_status").replace("{name}", modalName);
    statusTipo = "Success";

    const formData = new FormData();
    formData.append("name", modalName);
    formData.append("driver", modalDriver);

    showCreateModal = false;
    await runAction("create", formData);
  }

  async function doDelete(name: string) {
    statusMsg = t("networks.deleting_status").replace("{name}", name);
    statusTipo = "Success";

    const formData = new FormData();
    formData.append("name", name);

    await runAction("delete", formData);
    selectedNetworkNames = selectedNetworkNames.filter(
      (n) => !name.split(",").includes(n),
    );
  }

  async function doDeleteSelected() {
    if (selectedNetworkNames.length === 0) return;
    await doDelete(selectedNetworkNames.join(","));
  }

  async function doPrune() {
    showPruneModal = false;
    statusMsg = t("networks.pruning_status");
    statusTipo = "Success";
    const formData = new FormData();
    await runAction("prune", formData);
  }

  async function doConnect() {
    if (!selectedNetworkName || !selectedContainerName) return;
    statusMsg = t("networks.connecting_status")
      .replace("{container}", selectedContainerName)
      .replace("{network}", selectedNetworkName);
    statusTipo = "Success";

    const formData = new FormData();
    formData.append("networkName", selectedNetworkName);
    formData.append("containerName", selectedContainerName);

    showConnectModal = false;
    await runAction("connect", formData);
  }

  async function doDisconnect(networkName: string, containerName: string) {
    statusMsg = t("networks.disconnecting_status")
      .replace("{container}", containerName)
      .replace("{name}", networkName);
    statusTipo = "Success";

    const formData = new FormData();
    formData.append("networkName", networkName);
    formData.append("containerName", containerName);

    await runAction("disconnect", formData);
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
          {t("networks.title")}
        </h1>
      </div>

      <div class="flex items-center gap-2">
        <input
          type="text"
          placeholder={t("networks.search_placeholder")}
          class="px-4 py-2 text-xs rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-[#0b0f19] text-slate-850 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all w-60"
          bind:value={searchQuery}
        />
        <button
          type="button"
          class="px-4 py-2 text-xs rounded-xl border-none cursor-pointer font-bold text-white bg-blue-600 hover:bg-blue-700 transition-colors shadow-md shadow-blue-500/20 whitespace-nowrap"
          onclick={() => {
            modalName = "";
            modalDriver = "bridge";
            showCreateModal = true;
          }}
        >
          {t("networks.new_network")}
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
      <h3 class="font-bold text-sm mb-1">{t("common.error")}</h3>
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
          {selectedNetworkNames.length === selectableFilteredNetworks.length
            ? t("common.deselect_all")
            : t("common.select_all")}
        </button>

        {#if selectedNetworkNames.length > 0}
          <span
            class="text-xs font-semibold text-violet-600 dark:text-violet-400 animate-pulse px-2"
          >
            {selectedNetworkNames.length}
            {t("networks.selected_count")}
          </span>
        {/if}

        <span
          class="text-xs text-slate-400 dark:text-slate-500 px-2 font-semibold"
        >
          {t("networks.net_count_found").replace(
            "{count}",
            String(filteredNetworks.length),
          )}
        </span>
      </div>

      <div class="flex items-center gap-2">
        <button
          type="button"
          class="px-4 py-2.5 rounded-xl border-none cursor-pointer text-xs font-bold text-white bg-emerald-600 hover:bg-emerald-700 transition-colors shadow-md shadow-emerald-500/20"
          onclick={() => (showPruneModal = true)}
        >
          {t("networks.prune_btn")}
        </button>
        <button
          type="button"
          class="px-4 py-2.5 rounded-xl border-none cursor-pointer text-xs font-bold text-white transition-colors shadow-md shadow-red-500/10 {selectedNetworkNames.length >
          0
            ? 'bg-red-500 hover:bg-red-600'
            : 'bg-slate-400 dark:bg-slate-800 cursor-not-allowed text-slate-400'}"
          onclick={doDeleteSelected}
          disabled={selectedNetworkNames.length === 0}
        >
          {t("networks.card_remove_btn")}
        </button>
      </div>
    </div>

    <!-- Networks Grid -->
    {#if filteredNetworks.length === 0}
      <div
        class="text-sm text-slate-400 dark:text-slate-500 py-16 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-slate-50/20"
      >
        {t("networks.empty")}
      </div>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        {#each filteredNetworks as network (network.id)}
          <NetworkCard
            {network}
            checked={selectedNetworkNames.includes(network.name)}
            on_toggle={() => handleToggleSelect(network.name)}
            on_delete={() => doDelete(network.name)}
            on_disconnect={(cName) => doDisconnect(network.name, cName)}
            on_connect={() => {
              selectedNetworkName = network.name;
              selectedContainerName = containers[0]?.name || "";
              showConnectModal = true;
            }}
          />
        {/each}
      </div>
    {/if}
  {/if}
</div>
{/if}

<!-- Modal: Nova Rede -->
<Modal
  bind:show={showCreateModal}
  title={t("networks.create_title")}
  buttons={[
    {
      label: t("networks.create_btn"),
      variant: "primary",
      onclick: doCreate,
      disabled: !modalName.trim(),
    },
  ]}
>
  <div class="flex flex-col gap-1.5">
    <label
      for="net-name"
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >
      {t("networks.field_name")}
    </label>
    <input
      id="net-name"
      type="text"
      class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
      placeholder={t("networks.placeholder_name")}
      bind:value={modalName}
    />
  </div>

  <div class="flex flex-col gap-1.5">
    <label
      for="net-driver"
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >
      {t("networks.field_driver")}
    </label>
    <select
      id="net-driver"
      class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
      bind:value={modalDriver}
    >
      <option value="bridge">{t("networks.driver_bridge")}</option>
      <option value="host">host</option>
      <option value="overlay">overlay</option>
      <option value="macvlan">macvlan</option>
    </select>
  </div>
</Modal>

<!-- Modal: Conectar Container -->
<Modal
  bind:show={showConnectModal}
  title={t("networks.connect_modal_title")}
  buttons={[
    {
      label: t("networks.connect_btn"),
      variant: "primary",
      onclick: doConnect,
      disabled: !selectedContainerName,
    },
  ]}
>
  <!-- Network info label (inherited from the card button) -->
  <div
    class="flex items-center gap-2 px-4 py-3 rounded-xl bg-slate-50 dark:bg-slate-900/30 border border-slate-200/60 dark:border-slate-800/60"
  >
    <span
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider shrink-0"
    >
      {t("networks.select_network_label")}
    </span>
    <span class="text-sm font-semibold text-violet-600 dark:text-violet-400">
      {selectedNetworkName}
    </span>
  </div>

  <div class="flex flex-col gap-1.5">
    <label
      for="select-cont"
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >
      {t("networks.select_container_label")}
    </label>
    <select
      id="select-cont"
      class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
      bind:value={selectedContainerName}
    >
      {#each containers as container}
        <option value={container.name}
          >{container.name} ({container.image})</option
        >
      {:else}
        <option value="" disabled>{t("networks.no_active_containers")}</option>
      {/each}
    </select>
  </div>
</Modal>

<!-- Modal: Confirmar Limpeza (Prune) -->
<Modal
  bind:show={showPruneModal}
  title={t("networks.prune_confirm_title")}
  buttons={[
    {
      label: t("common.confirm"),
      variant: "success",
      onclick: doPrune,
    },
  ]}
>
  <p class="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
    {t("networks.prune_confirm_msg")}
  </p>
</Modal>
