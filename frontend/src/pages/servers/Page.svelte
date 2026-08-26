<script lang="ts">
  import { onMount } from "svelte";
  import { t } from "$shared/stores/locale.svelte";
  import StatusBanner from "$shared/components/StatusBanner.svelte";
  import { notifySuccess } from "$shared/stores/notification.svelte";
  import { ButtonPurple } from "$shared/components/buttons";
  import ServerCard from "./components/ServerCard.svelte";
  import * as api from "./api";
  import type { Route } from "../../navigation/navigation.types";

  let { data, navigate }: { data: any; navigate?: (route: Route) => void } =
    $props();
  let systemUsageCache = $state<Record<string, any>>({});
  let loadingUsage = $state<Record<string, boolean>>({});

  async function fetchServerUsage(server: any, force = false) {
    if (!server?.id || (!force && systemUsageCache[server.id])) return;
    loadingUsage[server.id] = true;
    try {
      const usage = await api.getSystemUsage(server);
      if (usage) systemUsageCache[server.id] = usage;
    } catch (error) {
      console.warn("Não foi possível obter recursos do servidor:", error);
    } finally {
      loadingUsage[server.id] = false;
    }
  }

  async function activateServer(server: any) {
    try {
      await api.setActiveServer(server.id);
      data.activeVps = server;
      data.servers = (await api.listServers()) || data.servers;
      notifySuccess(`Servidor '${server.name}' ativado com sucesso!`);
    } catch {
      data.activeVps = server;
      notifySuccess(`Servidor '${server.name}' selecionado.`);
    }
    fetchServerUsage(server);
  }

  async function loadServers() {
    try {
      const servers = await api.listServers();
      if (!servers) return;
      data.servers = servers;
      const active = servers.find((server: any) => server.isActive);
      if (active) {
        data.activeVps = active;
        fetchServerUsage(active);
      }
    } catch (error) {
      console.warn("Erro ao carregar servidores do SQLite:", error);
    }
  }

  onMount(loadServers);
</script>

<div class="space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div
      class="inline-flex items-center px-5 py-2.5 rounded-2xl bg-linear-to-br from-violet-100 to-fuchsia-100 dark:from-violet-950/40 dark:to-fuchsia-950/40 border border-violet-200/50 dark:border-violet-800/50 self-start shadow-sm"
    >
      <h1
        class="text-2xl font-bold bg-linear-to-r from-violet-600 to-fuchsia-500 bg-clip-text text-transparent m-0"
      >
        {t("sidebar.devices")}
      </h1>
    </div>
    <ButtonPurple size="sm" onclick={() => navigate?.("config")}
      >{t("devices.manage_vps")}</ButtonPurple
    >
  </div>

  <StatusBanner />

  {#if data.servers.length === 0}
    <div
      class="flex flex-col items-center justify-center py-20 px-6 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-white dark:bg-[#0b0f19] shadow-sm max-w-2xl mx-auto"
    >
      <h2 class="text-xl font-bold text-slate-800 dark:text-white mb-2">
        {t("devices.empty_title")}
      </h2>
      <p class="text-sm text-slate-500 dark:text-slate-400 mb-6 max-w-md">
        {t("devices.empty_desc")}
      </p>
      <ButtonPurple onclick={() => navigate?.("config")}
        >{t("devices.add_first")}</ButtonPurple
      >
    </div>
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each data.servers as server (server.id)}
        <ServerCard
          {server}
          usage={systemUsageCache[server.id]}
          isActive={data.activeVps?.id === server.id}
          isLoading={loadingUsage[server.id]}
          onRefresh={() => fetchServerUsage(server, true)}
          onActivate={() => activateServer(server)}
          onViewContainers={() => navigate?.("containers")}
        />
      {/each}
    </div>
  {/if}
</div>
