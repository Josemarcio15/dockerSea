<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import {
    ButtonPurple,
    ButtonBlue,
    ButtonPink,
  } from "$shared/components/buttons";

  let {
    server,
    usage,
    isActive,
    isLoading,
    onRefresh,
    onActivate,
    onViewContainers,
  }: {
    server: any;
    usage?: any;
    isActive: boolean;
    isLoading: boolean;
    onRefresh: () => void;
    onActivate: () => void;
    onViewContainers: () => void;
  } = $props();

  function formatBytes(bytes: number, decimals = 1) {
    if (!bytes || bytes === 0) return "0 B";
    const k = 1024;
    const sizes = ["B", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return (
      parseFloat((bytes / Math.pow(k, i)).toFixed(Math.max(0, decimals))) +
      " " +
      sizes[i]
    );
  }
</script>

<div
  class="group relative rounded-2xl bg-white/90 dark:bg-[#0c1220]/90 backdrop-blur-sm border border-slate-200/80 dark:border-slate-800/80 transition-all duration-300 flex flex-col justify-between p-5 gap-5 shadow-sm hover:shadow-md dark:shadow-black/40 hover:-translate-y-0.5"
>
  <div class="space-y-4">
    <div class="flex items-start justify-between">
      {#if isActive}
        <span
          class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-800/50 shadow-2xs"
        >
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"
          ></span>
          {t("profiles.active_badge")}
        </span>
      {:else}
        <span
          class="px-2.5 py-0.5 rounded-full text-[11px] font-medium bg-slate-100 dark:bg-slate-800/70 text-slate-500 dark:text-slate-400 border border-slate-200/70 dark:border-slate-700/60"
        >
          {t("profiles.inactive_badge")}
        </span>
      {/if}
      {#if isActive}
        <ButtonPink
          size="xs"
          title="Atualizar estatísticas de hardware"
          loading={isLoading}
          onclick={onRefresh}
        >
          <span aria-hidden="true">↻</span>
        </ButtonPink>
      {/if}
    </div>

    <div>
      <h3
        class="font-bold text-slate-900 dark:text-white text-lg tracking-tight"
      >
        {server.name}
      </h3>
      <p
        class="text-xs text-indigo-600 dark:text-indigo-400 font-mono mt-1 truncate font-medium"
      >
        {server.connectionType === "ssh"
          ? `${server.username}@${server.host}:${server.port}`
          : "Local Docker Engine"}
      </p>
    </div>

    {#if isLoading}
      <div
        class="p-3.5 rounded-xl bg-slate-50/80 dark:bg-slate-900/60 border border-slate-200/60 dark:border-slate-800/50 space-y-2 animate-pulse"
      >
        <div
          class="h-2.5 bg-slate-200 dark:bg-slate-700/50 rounded-full w-3/4"
        ></div>
        <div
          class="h-2 bg-slate-200 dark:bg-slate-700/40 rounded-full w-1/2"
        ></div>
      </div>
    {:else if usage}
      <!-- Inner elevated ice-toned container, with RAM, Disk, Swap and Uptime -->
      <div
        class="p-4 rounded-2xl bg-[#f8fafc] dark:bg-[#0c1220] border border-slate-200/90 dark:border-slate-800 shadow-md divide-y divide-slate-200/60 dark:divide-slate-800/80 space-y-3.5"
      >
        {#each [{ label: "RAM", used: usage.memUsed, total: usage.memTotal, percent: usage.memUsagePerc, color: usage.memUsagePerc > 85 ? "bg-rose-500" : usage.memUsagePerc > 65 ? "bg-amber-500" : "bg-emerald-500", badge: "bg-purple-100 text-purple-700 dark:bg-purple-950/60 dark:text-purple-300" }, { label: "Disco (/)", used: usage.diskUsed, total: usage.diskTotal, percent: usage.diskUsagePerc, color: usage.diskUsagePerc > 85 ? "bg-rose-500" : usage.diskUsagePerc > 70 ? "bg-amber-500" : "bg-sky-500", badge: "bg-sky-100 text-sky-700 dark:bg-sky-950/60 dark:text-sky-300" }] as metric, i}
          <div class="space-y-2 {i > 0 ? 'pt-3.5' : ''}">
            <div class="flex justify-between items-center text-xs font-semibold">
              <span
                class="px-2 py-0.5 rounded-md font-bold text-[10px] uppercase tracking-wider {metric.badge}"
                >{metric.label}</span
              >
              <span
                class="text-slate-900 dark:text-slate-100 font-mono font-bold"
                >{formatBytes(metric.used)} / {formatBytes(metric.total)}
                <span class="text-indigo-600 dark:text-indigo-400"
                  >({metric.percent.toFixed(1)}%)</span
                ></span
              >
            </div>
            <div
              class="h-2.5 w-full bg-slate-200/80 dark:bg-slate-800 rounded-full overflow-hidden p-0.5 border border-slate-300/80 dark:border-slate-700 shadow-inner"
            >
              <div
                class="h-full rounded-full transition-all duration-500 {metric.color} shadow-xs"
                style="width: {Math.min(100, Math.max(0, metric.percent))}%"
              ></div>
            </div>
          </div>
        {/each}

        {#if usage.swapTotal > 0}
          <div class="pt-3.5 space-y-2">
            <div class="flex justify-between items-center text-xs font-semibold">
              <span
                class="px-2 py-0.5 rounded-md font-bold text-[10px] uppercase tracking-wider bg-indigo-100 text-indigo-700 dark:bg-indigo-950/60 dark:text-indigo-300"
                >Swap</span
              >
              <span
                class="text-slate-900 dark:text-slate-100 font-mono font-bold"
                >{formatBytes(usage.swapUsed)} / {formatBytes(usage.swapTotal)}
                <span class="text-indigo-600 dark:text-indigo-400"
                  >({usage.swapUsagePerc.toFixed(1)}%)</span
                ></span
              >
            </div>
            <div
              class="h-2.5 w-full bg-slate-200/80 dark:bg-slate-800 rounded-full overflow-hidden p-0.5 border border-slate-300/80 dark:border-slate-700 shadow-inner"
            >
              <div
                class="h-full rounded-full transition-all duration-500 {usage.swapUsagePerc >
                70
                  ? 'bg-rose-500'
                  : 'bg-indigo-500'} shadow-xs"
                style="width: {Math.min(
                  100,
                  Math.max(0, usage.swapUsagePerc),
                )}%"
              ></div>
            </div>
          </div>
        {/if}

        {#if usage.uptime}
          <div
            class="pt-3 flex items-center justify-between text-[11px] text-slate-600 dark:text-slate-400 font-medium"
          >
            <span class="truncate max-w-35 font-medium" title={usage.uptime}
              >{usage.uptime}</span
            >{#if usage.cpuCount}<span
                class="font-mono font-bold text-indigo-600 dark:text-indigo-400"
                >{usage.cpuCount} CPU(s)</span
              >{/if}
          </div>
        {/if}
      </div>
    {/if}
  </div>
  <div
    class="flex items-center gap-2.5 border-t border-slate-100 dark:border-slate-800/70 pt-4"
  >
    {#if isActive}
      <ButtonBlue size="sm" class="flex items-center justify-center gap-2" onclick={onViewContainers}>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="w-4 h-4"
        >
          <path
            d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"
          />
          <polyline points="3.27 6.96 12 12.01 20.73 6.96" />
          <line x1="12" y1="22.08" x2="12" y2="12" />
        </svg>
        <span>{t("devices.view_containers")}</span>
      </ButtonBlue>
    {:else}
      <ButtonPurple
        class="flex-1 flex items-center justify-center gap-2"
        onclick={onActivate}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2.2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="w-4 h-4"
        >
          <path d="M18.36 6.64a9 9 0 1 1-12.73 0" />
          <line x1="12" y1="2" x2="12" y2="12" />
        </svg>
        <span>{t("devices.activate")}</span>
      </ButtonPurple>
    {/if}
  </div>
</div>
