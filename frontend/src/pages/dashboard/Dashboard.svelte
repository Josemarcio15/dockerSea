<script lang="ts">
  import { onMount } from "svelte";
  import { t } from "$lib/stores/locale.svelte";
  import StatusBanner from "$lib/components/StatusBanner.svelte";
  import { notifySuccess, notifyError } from "$lib/stores/notification.svelte";
  import { ButtonPurple, ButtonBlue, ButtonPink } from "$lib/components/buttons";
  import * as ConfigService from "../../../bindings/go-walis/internal/config/configservice.js";

  let {
    data,
    navigate,
  }: {
    data: any;
    navigate?: (route: string) => void;
  } = $props();

  // Cache local em memória do dashboard para armazenar os dados capturados de cada VPS selecionada (1x)
  let systemUsageCache = $state<Record<string, any>>({});
  let loadingUsage = $state<Record<string, boolean>>({});

  function formatBytes(bytes: number, decimals = 1) {
    if (!bytes || bytes === 0) return "0 B";
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["B", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  }

  async function fetchServerUsage(server: any, force = false) {
    if (!server?.id) return;
    if (!force && systemUsageCache[server.id]) {
      return; // Já capturado 1x, não precisa buscar novamente
    }

    loadingUsage[server.id] = true;
    try {
      const usage = await ConfigService.GetSystemUsage(server);
      if (usage) {
        systemUsageCache[server.id] = usage;
      }
    } catch (e: any) {
      console.warn("Não foi possível obter recursos do servidor:", e);
    } finally {
      loadingUsage[server.id] = false;
    }
  }

  async function doSelectVps(server: any) {
    try {
      await ConfigService.SetActiveServer(server.id);
      data.activeVps = server;
      const servers = await ConfigService.ListServers();
      if (servers) {
        data.servers = servers;
      }
      notifySuccess(`Servidor '${server.name}' ativado com sucesso!`);
      // Busca 1x as métricas de hardware da VPS no momento em que ela for selecionada
      fetchServerUsage(server);
    } catch (e: any) {
      data.activeVps = server;
      notifySuccess(`Servidor '${server.name}' selecionado.`);
      fetchServerUsage(server);
    }
  }

  // Carregar servidores salvos no SQLite
  async function loadServers() {
    try {
      const servers = await ConfigService.ListServers();
      if (servers) {
        data.servers = servers;
        const active = servers.find((s: any) => s.isActive);
        if (active) {
          data.activeVps = active;
          fetchServerUsage(active);
        }
      }
    } catch (e: any) {
      console.warn("Erro ao carregar servidores do SQLite no Dashboard:", e);
    }
  }

  onMount(() => {
    // Busca sempre os servidores cadastrados no SQLite ao entrar na tela de Servidores
    loadServers();
  });
</script>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div
      class="inline-flex items-center px-5 py-2.5 rounded-2xl bg-linear-to-br from-violet-100 to-fuchsia-100 dark:from-violet-950/40 dark:to-fuchsia-950/40 border border-violet-200/50 dark:border-violet-800/50 self-start shadow-sm"
    >
      <h1
        class="text-2xl font-bold bg-linear-to-r from-violet-600 to-fuchsia-500 bg-clip-text text-transparent m-0 flex items-center gap-2"
      >
        {t("sidebar.devices")}
      </h1>
    </div>

    <ButtonPurple
      size="sm"
      onclick={() => navigate && navigate("config")}
    >
      {t("devices.manage_vps")}
    </ButtonPurple>
  </div>

  <!-- Status Banner -->
  <StatusBanner />

  {#if data.servers.length === 0}
    <!-- Empty state -->
    <div
      class="flex flex-col items-center justify-center py-20 px-6 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-white dark:bg-[#0b0f19] shadow-sm max-w-2xl mx-auto"
    >
      <h2 class="text-xl font-bold text-slate-800 dark:text-white mb-2">
        {t("devices.empty_title")}
      </h2>
      <p class="text-sm text-slate-500 dark:text-slate-400 mb-6 max-w-md">
        {t("devices.empty_desc")}
      </p>
      <ButtonPurple
        onclick={() => navigate && navigate("config")}
      >
        {t("devices.add_first")}
      </ButtonPurple>
    </div>
  {:else}
    <!-- Servers Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each data.servers as server (server.id)}
        {@const isActive = data.activeVps?.id === server.id}
        {@const usage = systemUsageCache[server.id]}
        {@const isLoading = loadingUsage[server.id]}
        <div
          class="group relative rounded-2xl bg-white dark:bg-[#0c1220] border transition-all duration-300 flex flex-col justify-between p-5 gap-5 shadow-md dark:shadow-lg dark:shadow-black/40 hover:-translate-y-1 hover:border-slate-300 dark:hover:border-slate-700 {isActive
            ? 'border-violet-500 ring-2 ring-violet-500/10'
            : 'border-slate-200 dark:border-slate-800/80'}"
        >
          <div class="space-y-4">
            <div class="flex items-start justify-between">
              {#if isActive}
                <span
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20 shadow-sm"
                >
                  {t("profiles.active_badge")}
                </span>
              {:else}
                <span
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-slate-100 dark:bg-slate-900 text-slate-500 border border-slate-200/50 dark:border-slate-800/50"
                >
                  {t("profiles.inactive_badge")}
                </span>
              {/if}

              {#if isActive}
                <ButtonPink
                  size="xs"
                  title="Atualizar estatísticas de hardware"
                  loading={isLoading}
                  onclick={() => fetchServerUsage(server, true)}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                    class="w-3.5 h-3.5"
                  >
                    <path
                      fill-rule="evenodd"
                      d="M15.312 11.424a5.5 5.5 0 01-9.201 2.466l-.312-.311h2.433a.75.75 0 000-1.5H3.75a.75.75 0 00-.75.75v4.482a.75.75 0 001.5 0v-2.122l.463.463a7 7 0 0011.712-3.138.75.75 0 00-1.363-.59zm-10.624-2.85a5.5 5.5 0 019.201-2.465l.312.311H11.77a.75.75 0 000 1.5h4.482a.75.75 0 00.75-.75V2.688a.75.75 0 00-1.5 0v2.122l-.463-.463A7 7 0 003.325 7.486a.75.75 0 001.363.588z"
                      clip-rule="evenodd"
                    />
                  </svg>
                </ButtonPink>
              {/if}
            </div>

            <div>
              <h3 class="font-bold text-slate-850 dark:text-white text-lg">
                {server.name}
              </h3>
              <p
                class="text-xs text-slate-400 dark:text-slate-500 font-mono mt-1 truncate"
              >
                {server.connectionType === "ssh"
                  ? `${server.username}@${server.host}:${server.port}`
                  : "Local Docker Engine"}
              </p>
            </div>

            <!-- Hardware & System Stats (Buscado apenas 1x na conexão) -->
            {#if isLoading}
              <div class="p-3 rounded-xl bg-slate-50 dark:bg-slate-900/60 border border-slate-100 dark:border-slate-800/50 space-y-2 animate-pulse">
                <div class="h-2.5 bg-slate-200 dark:bg-slate-700/50 rounded-full w-3/4"></div>
                <div class="h-2 bg-slate-200 dark:bg-slate-700/40 rounded-full w-1/2"></div>
              </div>
            {:else if usage}
              <div class="p-3.5 rounded-xl bg-slate-50 dark:bg-[#080d1a] border border-slate-200/60 dark:border-slate-800/60 space-y-3">
                <!-- Memória RAM -->
                <div class="space-y-1">
                  <div class="flex justify-between text-[11px] font-semibold">
                    <span class="text-slate-500 dark:text-slate-400 flex items-center gap-1">
                      🧠 RAM
                    </span>
                    <span class="text-slate-700 dark:text-slate-300 font-mono">
                      {formatBytes(usage.memUsed)} / {formatBytes(usage.memTotal)} ({usage.memUsagePerc.toFixed(1)}%)
                    </span>
                  </div>
                  <div class="h-1.5 w-full bg-slate-200 dark:bg-slate-800 rounded-full overflow-hidden">
                    <div
                      class="h-full rounded-full transition-all duration-500 {usage.memUsagePerc > 85 ? 'bg-rose-500' : usage.memUsagePerc > 65 ? 'bg-amber-500' : 'bg-emerald-500'}"
                      style="width: {Math.min(100, Math.max(0, usage.memUsagePerc))}%"
                    ></div>
                  </div>
                </div>

                <!-- Memória Swap (se existir) -->
                {#if usage.swapTotal > 0}
                  <div class="space-y-1">
                    <div class="flex justify-between text-[11px] font-semibold">
                      <span class="text-slate-500 dark:text-slate-400 flex items-center gap-1">
                        🔄 Swap
                      </span>
                      <span class="text-slate-700 dark:text-slate-300 font-mono">
                        {formatBytes(usage.swapUsed)} / {formatBytes(usage.swapTotal)} ({usage.swapUsagePerc.toFixed(1)}%)
                      </span>
                    </div>
                    <div class="h-1.5 w-full bg-slate-200 dark:bg-slate-800 rounded-full overflow-hidden">
                      <div
                        class="h-full rounded-full transition-all duration-500 {usage.swapUsagePerc > 70 ? 'bg-rose-500' : 'bg-indigo-500'}"
                        style="width: {Math.min(100, Math.max(0, usage.swapUsagePerc))}%"
                      ></div>
                    </div>
                  </div>
                {/if}

                <!-- Disco Root / -->
                <div class="space-y-1">
                  <div class="flex justify-between text-[11px] font-semibold">
                    <span class="text-slate-500 dark:text-slate-400 flex items-center gap-1">
                      💾 Disco (/)
                    </span>
                    <span class="text-slate-700 dark:text-slate-300 font-mono">
                      {formatBytes(usage.diskUsed)} / {formatBytes(usage.diskTotal)} ({usage.diskUsagePerc.toFixed(1)}%)
                    </span>
                  </div>
                  <div class="h-1.5 w-full bg-slate-200 dark:bg-slate-800 rounded-full overflow-hidden">
                    <div
                      class="h-full rounded-full transition-all duration-500 {usage.diskUsagePerc > 85 ? 'bg-rose-500' : usage.diskUsagePerc > 70 ? 'bg-amber-500' : 'bg-cyan-500'}"
                      style="width: {Math.min(100, Math.max(0, usage.diskUsagePerc))}%"
                    ></div>
                  </div>
                </div>

                <!-- Detalhes leves (Uptime / Kernel) -->
                {#if usage.uptime}
                  <div class="pt-1 border-t border-slate-200/50 dark:border-slate-800/50 flex items-center justify-between text-[10px] text-slate-400">
                    <span class="truncate max-w-[140px]" title={usage.uptime}>⏱️ {usage.uptime}</span>
                    {#if usage.cpuCount}
                      <span>⚡ {usage.cpuCount} CPU(s)</span>
                    {/if}
                  </div>
                {/if}
              </div>
            {/if}
          </div>

          <div
            class="flex items-center gap-2.5 border-t border-slate-100 dark:border-slate-900 pt-4"
          >
            {#if isActive}
              <ButtonBlue size="sm" onclick={() => navigate && navigate("containers")}>
                {t("devices.view_containers")}
              </ButtonBlue>
            {:else}
              <ButtonPurple
                class="flex-1"
                onclick={() => doSelectVps(server)}
              >
                {t("devices.activate")}
              </ButtonPurple>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

