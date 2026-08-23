<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import type { Container } from "$lib/domains/containers";
  import { statsState } from "$lib/stores/stats.svelte";
  import { BrandButton, Button, DangerButton } from "$lib/components/buttons";

  let {
    container,
    checked = false,
    on_toggle = () => {},
    on_open_logs = (name: string) => {},
  }: {
    container: Container;
    checked?: boolean;
    on_toggle?: () => void;
    on_open_logs?: (name: string) => void;
  } = $props();

  let expanded = $state(false);
  let showEnv = $state(false);

  const myStats = $derived(
    statsState?.stats
      ? statsState.stats.find(
          (s) =>
            s &&
            s.ID &&
            (s.ID === container.id ||
              s.ID.startsWith(container.id.substring(0, 12))),
        )
      : undefined,
  );

  // Helpers
  const [accentBorder, statusColor] = $derived.by(() => {
    const s = container.status || "";
    if (s.includes("Up")) {
      return [
        "border-l-emerald-400",
        "text-emerald-600 dark:text-emerald-400",
      ];
    } else if (s.includes("Exited") || s.includes("Paused")) {
      return ["border-l-amber-400", "text-amber-600 dark:text-amber-400"];
    } else {
      return ["border-l-red-400", "text-red-600 dark:text-red-400"];
    }
  });

  const statusDotColor = $derived(
    (container.status || "").includes("Up")
      ? "bg-emerald-500"
      : (container.status || "").includes("Exited") ||
          (container.status || "").includes("Paused")
        ? "bg-amber-500"
        : "bg-red-500",
  );

  const pulseClass = $derived(
    (container.status || "").includes("Up") ? "animate-pulse" : "",
  );

  const statusBg = $derived(
    (container.status || "").includes("Up")
      ? "bg-emerald-50 dark:bg-emerald-950/40 text-emerald-600 dark:text-emerald-400 border-emerald-200 dark:border-emerald-900/60"
      : (container.status || "").includes("Exited") ||
          (container.status || "").includes("Paused")
        ? "bg-amber-50 dark:bg-amber-950/40 text-amber-600 dark:text-amber-400 border-amber-200 dark:border-amber-900/60"
        : "bg-red-50 dark:bg-red-950/40 text-red-600 dark:text-red-400 border-red-200 dark:border-red-900/60",
  );

  const portItems = $derived.by(() => {
    if (!container.ports) return [];
    return container.ports
      .split(", ")
      .map((entry) => {
        const e = entry.trim();
        if (!e) return null;
        const formatted = e.replace("->", " → ");
        const isIpv6 = e.includes("::");
        return {
          formatted,
          tag: isIpv6 ? "IPv6" : "IPv4",
          class: isIpv6
            ? "bg-purple-50 dark:bg-purple-950/30 border-purple-200 dark:border-purple-900/40 text-purple-700 dark:text-purple-400"
            : "bg-blue-50 dark:bg-blue-950/30 border-blue-200 dark:border-blue-900/40 text-blue-700 dark:text-blue-400",
        };
      })
      .filter((x) => x !== null);
  });

  const createdStr = $derived.by(() => {
    if (!container.created) return "";
    const date = new Date(container.created * 1000);
    const day = String(date.getDate()).padStart(2, "0");
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
  });

  const shortCid = $derived((container.id || "").substring(0, 12));

  const restartPolicyDisplay = $derived.by(() => {
    const policy =
      container.restartPolicy || (container as any).restart_policy;
    if (!policy) return "—";
    switch (policy.toLowerCase()) {
      case "always":
        return t("containers.restart_always");
      case "unless-stopped":
        return t("containers.restart_unless_stopped");
      case "on-failure":
        return t("containers.restart_on_failure");
      case "no":
        return t("containers.restart_no");
      default:
        return policy;
    }
  });
</script>

<div
  class="relative rounded-2xl bg-[#f0f3f8] dark:bg-[#0c1220] border border-slate-300/80 dark:border-slate-800/80 hover:border-slate-400 dark:hover:border-slate-700 transition-all duration-300 flex flex-col justify-between overflow-hidden self-start w-full text-slate-700 dark:text-slate-200 shadow-md dark:shadow-lg dark:shadow-black/40 p-4 gap-3.5"
>
  <!-- Card Header -->
  <div class="flex items-center justify-between gap-3">
    <!-- Checkbox -->
    <button
      type="button"
      class="w-5.5 h-5.5 rounded-lg border-2 flex items-center justify-center transition-all duration-150 shrink-0 cursor-pointer {checked
        ? 'bg-violet-600 border-violet-500 text-white'
        : 'border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-900/60 hover:border-violet-500'}"
      onclick={on_toggle}
    >
      {#if checked}
        <span class="text-white text-xs font-bold leading-none">✓</span>
      {/if}
    </button>

    <!-- Name button -->
    <button
      type="button"
      class="flex-1 font-mono font-bold text-sm tracking-tight truncate text-slate-855 dark:text-slate-100 px-4 py-2 rounded-2xl bg-white dark:bg-slate-800/30 border border-slate-200/80 dark:border-slate-700/30 flex items-center justify-between gap-2.5 cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors text-left shadow-xs"
      onclick={() => (expanded = !expanded)}
    >
      <div class="flex items-center gap-2 truncate">
        <span class="truncate font-semibold text-slate-855 dark:text-white grow">
          {container.name}
        </span>
      </div>
      <span class="text-[10px] text-slate-400 dark:text-slate-500 shrink-0">
        {expanded ? "▲" : "▼"}
      </span>
    </button>

    <!-- Status Badge -->
    <div
      class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl border text-xs font-semibold shrink-0 {statusBg}"
    >
      <span
        class="w-2 h-2 rounded-full {statusDotColor} {pulseClass} shrink-0"
      ></span>
      <span class="capitalize">{container.status}</span>
    </div>
  </div>

  <!-- CPU/Memory Realtime Mini-Stats -->
  {#if myStats}
    <div
      class="grid grid-cols-2 gap-2 text-xs bg-white dark:bg-[#070a12] p-2.5 rounded-xl border border-slate-200 dark:border-slate-800/80 shadow-xs"
    >
      <div class="flex flex-col">
        <span
          class="text-[9px] uppercase font-bold text-slate-400 dark:text-slate-500"
        >
          CPU
        </span>
        <span class="font-mono font-semibold text-violet-600 dark:text-violet-400">
          {myStats.CPUPerc || "0%"}
        </span>
      </div>
      <div class="flex flex-col">
        <span
          class="text-[9px] uppercase font-bold text-slate-400 dark:text-slate-500"
        >
          {t("containers.card_memory")}
        </span>
        <span class="font-mono font-semibold text-blue-600 dark:text-blue-400">
          {myStats.MemUsage || "0B"} / {myStats.MemPerc || "0%"}
        </span>
      </div>
    </div>
  {/if}

  <!-- Expanded details -->
  {#if expanded}
    <div class="flex flex-col gap-3 text-xs pt-1">
      <!-- ID & Image -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-violet-500 shadow-md dark:shadow-black/40"
        >
          <div class="grow flex flex-col gap-0.5 min-w-0">
            <span
              class="text-[9px] text-violet-500 font-bold uppercase tracking-wider"
            >
              {t("containers.card_container_id")}
            </span>
            <span
              class="font-mono text-xs font-bold text-blue-500 dark:text-blue-400 truncate"
            >
              {shortCid}
            </span>
          </div>
        </div>

        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-fuchsia-500 shadow-md dark:shadow-black/40"
        >
          <div class="grow flex flex-col gap-0.5 min-w-0">
            <span
              class="text-[9px] text-fuchsia-500 font-bold uppercase tracking-wider"
            >
              {t("containers.card_image")}
            </span>
            <span
              class="font-mono text-xs font-semibold text-slate-700 dark:text-slate-200 truncate"
            >
              {container.image}
            </span>
          </div>
        </div>
      </div>

      <!-- Created & Restart Policy -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-cyan-500 shadow-md dark:shadow-black/40"
        >
          <div class="flex flex-col gap-0.5">
            <span
              class="text-[9px] text-cyan-500 font-bold uppercase tracking-wider"
            >
              {t("containers.card_created")}
            </span>
            <span class="font-semibold text-slate-700 dark:text-slate-200 text-xs">
              {createdStr || "—"}
            </span>
          </div>
        </div>

        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-emerald-500 shadow-md dark:shadow-black/40"
        >
          <div class="flex flex-col gap-0.5">
            <span
              class="text-[9px] text-emerald-500 font-bold uppercase tracking-wider"
            >
              {t("containers.card_restart_policy")}
            </span>
            <span class="font-semibold text-slate-700 dark:text-slate-200 text-xs">
              {restartPolicyDisplay}
            </span>
          </div>
        </div>
      </div>

      <!-- Ports Panel -->
      {#if portItems.length > 0}
        <div
          class="p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-blue-500 shadow-md dark:shadow-black/40 flex flex-col gap-2"
        >
          <span
            class="text-[9px] text-blue-500 font-bold uppercase tracking-wider"
          >
            {t("containers.card_ports")}
          </span>
          <div class="flex flex-wrap gap-2">
            {#each portItems as port}
              {#if port}
                <span
                  class="px-2 py-1 rounded-lg border font-mono text-[11px] font-semibold {port.class}"
                >
                  {port.formatted}
                </span>
              {/if}
            {/each}
          </div>
        </div>
      {/if}

      <!-- Action Area -->
      <div class="flex justify-end pt-1">
        <BrandButton
          size="sm"
          onclick={() => on_open_logs(container.name)}
        >
          {t("containers.view_logs")}
        </BrandButton>
      </div>
    </div>
  {/if}
</div>
