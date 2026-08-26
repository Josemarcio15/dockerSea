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
  class="group relative rounded-2xl bg-white dark:bg-[#0c1220] border transition-all duration-300 flex flex-col justify-between p-5 gap-5 shadow-md dark:shadow-lg dark:shadow-black/40 hover:-translate-y-1 hover:border-slate-300 dark:hover:border-slate-700 {isActive
    ? 'border-violet-500 ring-2 ring-violet-500/10'
    : 'border-slate-200 dark:border-slate-800/80'}"
>
  <div class="space-y-4">
    <div class="flex items-start justify-between">
      {#if isActive}
        <span
          class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20 shadow-sm"
          >{t("profiles.active_badge")}</span
        >
      {:else}
        <span
          class="px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-slate-100 dark:bg-slate-900 text-slate-500 border border-slate-200/50 dark:border-slate-800/50"
          >{t("profiles.inactive_badge")}</span
        >
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

    {#if isLoading}
      <div
        class="p-3 rounded-xl bg-slate-50 dark:bg-slate-900/60 border border-slate-100 dark:border-slate-800/50 space-y-2 animate-pulse"
      >
        <div
          class="h-2.5 bg-slate-200 dark:bg-slate-700/50 rounded-full w-3/4"
        ></div>
        <div
          class="h-2 bg-slate-200 dark:bg-slate-700/40 rounded-full w-1/2"
        ></div>
      </div>
    {:else if usage}
      <div
        class="p-3.5 rounded-xl bg-slate-50 dark:bg-[#080d1a] border border-slate-200/60 dark:border-slate-800/60 space-y-3"
      >
        {#each [{ label: "🧠 RAM", used: usage.memUsed, total: usage.memTotal, percent: usage.memUsagePerc, color: usage.memUsagePerc > 85 ? "bg-rose-500" : usage.memUsagePerc > 65 ? "bg-amber-500" : "bg-emerald-500" }, { label: "💾 Disco (/)", used: usage.diskUsed, total: usage.diskTotal, percent: usage.diskUsagePerc, color: usage.diskUsagePerc > 85 ? "bg-rose-500" : usage.diskUsagePerc > 70 ? "bg-amber-500" : "bg-cyan-500" }] as metric}
          <div class="space-y-1">
            <div class="flex justify-between text-[11px] font-semibold">
              <span class="text-slate-500 dark:text-slate-400"
                >{metric.label}</span
              >
              <span class="text-slate-700 dark:text-slate-300 font-mono"
                >{formatBytes(metric.used)} / {formatBytes(metric.total)} ({metric.percent.toFixed(
                  1,
                )}%)</span
              >
            </div>
            <div
              class="h-1.5 w-full bg-slate-200 dark:bg-slate-800 rounded-full overflow-hidden"
            >
              <div
                class="h-full rounded-full transition-all duration-500 {metric.color}"
                style="width: {Math.min(100, Math.max(0, metric.percent))}%"
              ></div>
            </div>
          </div>
        {/each}
        {#if usage.swapTotal > 0}
          <div class="space-y-1">
            <div class="flex justify-between text-[11px] font-semibold">
              <span class="text-slate-500 dark:text-slate-400">🔄 Swap</span
              ><span class="text-slate-700 dark:text-slate-300 font-mono"
                >{formatBytes(usage.swapUsed)} / {formatBytes(usage.swapTotal)} ({usage.swapUsagePerc.toFixed(
                  1,
                )}%)</span
              >
            </div>
            <div
              class="h-1.5 w-full bg-slate-200 dark:bg-slate-800 rounded-full overflow-hidden"
            >
              <div
                class="h-full rounded-full transition-all duration-500 {usage.swapUsagePerc >
                70
                  ? 'bg-rose-500'
                  : 'bg-indigo-500'}"
                style="width: {Math.min(
                  100,
                  Math.max(0, usage.swapUsagePerc),
                )}%"
              ></div>
            </div>
          </div>
        {/if}
        {#if usage.uptime}<div
            class="pt-1 border-t border-slate-200/50 dark:border-slate-800/50 flex items-center justify-between text-[10px] text-slate-400"
          >
            <span class="truncate max-w-35" title={usage.uptime}
              >⏱️ {usage.uptime}</span
            >{#if usage.cpuCount}<span>⚡ {usage.cpuCount} CPU(s)</span>{/if}
          </div>{/if}
      </div>
    {/if}
  </div>
  <div
    class="flex items-center gap-2.5 border-t border-slate-100 dark:border-slate-900 pt-4"
  >
    {#if isActive}<ButtonBlue size="sm" onclick={onViewContainers}
        >{t("devices.view_containers")}</ButtonBlue
      >{:else}<ButtonPurple class="flex-1" onclick={onActivate}
        >{t("devices.activate")}</ButtonPurple
      >{/if}
  </div>
</div>
