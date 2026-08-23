<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import { useRefreshKey, triggerRefresh } from "$lib/stores/refresh.svelte";
  import {
    notifySuccess,
    notifyWarning,
    notifyError,
  } from "$lib/stores/notification.svelte";
  import DockerseaLoading from "$lib/components/DockerseaLoading.svelte";
  import StatusBanner from "$lib/components/StatusBanner.svelte";
  import ContainerCard from "$lib/components/ContainerCard.svelte";
  import VpsSelectWarning from "$lib/components/VpsSelectWarning.svelte";
  import TerminalModal from "$lib/components/TerminalModal.svelte";
  import {
    BrandButton,
    PinkButton,
    SuccessButton,
    WarningButton,
    DangerButton,
    PrimaryButton,
    SecondaryButton,
  } from "$lib/components/buttons";
  import { Events } from "@wailsio/runtime";
  import * as ContainerService from "../../../bindings/go-walis/internal/containers/containerservice.js";

  let { data } = $props();

  let selectedNames = $state<string[]>([]);
  let searchQuery = $state("");

  // Logs Modal state
  let showLogs = $state(false);
  let logsTitle = $state("");
  let logsLoading = $state(false);
  let logsContent = $state<string[]>([]);

  // Client-side fetched containers
  let containers = $state<any[]>([]);
  let loading = $state(true);
  let fetchError = $state("");

  // Reactive fetch: runs on mount and whenever triggerRefresh() is called
  async function fetchContainers(silent = false) {
    if (!data?.activeVps) {
      loading = false;
      return;
    }
    if (!silent) {
      loading = true;
    }
    fetchError = "";
    try {
      const list = await ContainerService.ListContainers(data.activeVps, true);
      containers = list || [];
    } catch (e: any) {
      fetchError = e.message || String(e);
      if (!silent) {
        containers = [];
      }
    } finally {
      loading = false;
    }
  }

  // Inicializa a stream de eventos no backend e inscreve o listener reativo no frontend
  $effect(() => {
    useRefreshKey();
    if (data?.activeVps?.id) {
      fetchContainers();
      // Inicia a stream de eventos no backend para a VPS ativa
      ContainerService.StartEventsStream(data.activeVps).catch((err) => {
        console.warn("Não foi possível iniciar stream de eventos Docker:", err);
      });

      // Escuta eventos em tempo real do Docker Daemon (sem polling!)
      const unsubs = Events.On("docker:container:event", (event: any) => {
        console.log("[DockerEvents] Evento recebido no frontend:", event);
        const evData = event?.data?.[0] || event?.data || event;
        if (!evData?.serverId || evData.serverId === data.activeVps?.id) {
          // Atualiza a lista de containers imediatamente em background sem travar a UI
          fetchContainers(true);
        }
      });

      return () => {
        unsubs();
        if (data?.activeVps?.id) {
          ContainerService.StopEventsStream(data.activeVps.id);
        }
      };
    }
  });

  // Filter containers
  let filteredContainers = $derived(
    (containers || []).filter(
      (c) =>
        (c.name && c.name.toLowerCase().includes(searchQuery.toLowerCase())) ||
        (c.image && c.image.toLowerCase().includes(searchQuery.toLowerCase())),
    ),
  );

  // Toggle checked state
  function toggleChecked(name: string) {
    if (selectedNames.includes(name)) {
      selectedNames = selectedNames.filter((n) => n !== name);
    } else {
      selectedNames = [...selectedNames, name];
    }
  }

  // Toggle all
  function toggleAll() {
    if (selectedNames.length === filteredContainers.length) {
      selectedNames = [];
    } else {
      selectedNames = filteredContainers.map((c) => c.name);
    }
  }

  // Execute batch action
  async function doAction(actionType: "start" | "stop" | "restart" | "rm") {
    if (selectedNames.length === 0) {
      notifyWarning(t("containers.action_warning"));
      return;
    }

    if (!data?.activeVps) {
      notifyError("Nenhuma VPS ativa selecionada.");
      return;
    }

    const actionLabels = {
      start: t("containers.action_started"),
      stop: t("containers.action_stopped"),
      restart: t("containers.action_restarted"),
      rm: t("containers.action_deleted"),
    };

    notifySuccess(
      t("containers.action_progress", {
        action: actionLabels[actionType],
        count: String(selectedNames.length),
      }),
    );

    // Save names locally before clearing (API call is async)
    const namesToAct = [...selectedNames];

    // Clear selected names if they are being removed
    if (actionType === "rm") {
      selectedNames = [];
    }

    try {
      const result = await ContainerService.ExecuteAction(
        data.activeVps,
        actionType,
        namesToAct,
      );

      if (result && result.success) {
        notifySuccess(
          t("containers.action_success", {
            count: String(namesToAct.length),
          }),
        );
      } else {
        notifyError(
          t("containers.action_error", {
            count: String(namesToAct.length),
            errors: result?.message || "Falha ao executar ação",
          }),
        );
      }

      // Refresh containers list without page reload
      triggerRefresh();
    } catch (e: any) {
      notifyError(t("containers.connection_error", { error: e.message || String(e) }));
    }
  }

  // Open logs modal
  async function openLogs(containerName: string) {
    logsTitle = `📋 ${t("containers.card_logs_title")}${containerName}`;
    logsContent = [];
    logsLoading = true;
    showLogs = true;

    if (!data?.activeVps) {
      logsContent = ["Nenhuma VPS ativa selecionada."];
      logsLoading = false;
      return;
    }

    try {
      const logs = await ContainerService.GetLogs(data.activeVps, containerName, 200);
      if (logs) {
        logsContent = logs.trim().split("\n").filter(Boolean);
        if (logsContent.length === 0) {
          logsContent = [t("containers.card_no_logs")];
        }
      } else {
        logsContent = [t("containers.card_no_logs")];
      }
    } catch (e: any) {
      logsContent = [`Erro ao consultar logs: ${e.message || String(e)}`];
    } finally {
      logsLoading = false;
    }
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
          {t("containers.title")}
        </h1>
      </div>

      <div class="flex items-center gap-2">
        <input
          type="text"
          placeholder={t("common.search")}
          class="px-4 py-2 text-xs rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-[#0b0f19] text-slate-850 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all w-60"
          bind:value={searchQuery}
        />
        <PinkButton
          size="sm"
          title={t("common.refresh")}
          onclick={() => fetchContainers()}
        >
          {t("common.refresh")}
        </PinkButton>
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
        <BrandButton
          size="sm"
          onclick={toggleAll}
        >
          {selectedNames.length === filteredContainers.length
            ? t("common.deselect_all")
            : t("common.select_all")}
        </BrandButton>

        {#if selectedNames.length > 0}
          <span
            class="text-xs font-semibold text-violet-600 dark:text-violet-400 animate-pulse px-2"
          >
            {selectedNames.length}
            {t("images.selected_count")}
          </span>
        {/if}
      </div>

      <div class="flex items-center gap-2">
        <SuccessButton
          size="sm"
          onclick={() => doAction("start")}
        >
          {t("containers.start")}
        </SuccessButton>
        <WarningButton
          size="sm"
          onclick={() => doAction("stop")}
        >
          {t("containers.stop")}
        </WarningButton>
        <PrimaryButton
          size="sm"
          onclick={() => doAction("restart")}
        >
          {t("containers.restart")}
        </PrimaryButton>
        <DangerButton
          size="sm"
          onclick={() => doAction("rm")}
        >
          {t("containers.delete")}
        </DangerButton>
      </div>
    </div>

    <!-- Containers Grid -->
    {#if filteredContainers.length === 0}
      <div
        class="text-sm text-slate-400 dark:text-slate-500 py-16 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-slate-50/20"
      >
        {t("containers.empty")}
      </div>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        {#each filteredContainers as container (container.id)}
          <ContainerCard
            {container}
            checked={selectedNames.includes(container.name)}
            on_toggle={() => toggleChecked(container.name)}
            on_open_logs={openLogs}
          />
        {/each}
      </div>
    {/if}
  {/if}
</div>
{/if}

<!-- Logs terminal modal -->
<TerminalModal
  bind:show={showLogs}
  title={logsTitle}
  loading={logsLoading}
  logs={logsContent}
  console_id="container-logs-console"
/>
