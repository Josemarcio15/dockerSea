<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import type { DockerNetwork } from "$lib/server/docker";

  let {
    network,
    checked = false,
    on_toggle = () => {},
    on_disconnect = (containerName: string) => {},
    on_delete = () => {},
    on_connect = () => {},
  }: {
    network: DockerNetwork;
    checked?: boolean;
    on_toggle?: () => void;
    on_disconnect?: (containerName: string) => void | Promise<void>;
    on_delete?: () => void | Promise<void>;
    on_connect?: () => void;
  } = $props();

  let expanded = $state(false);

  const isDefaultNetwork = $derived(
    network.name === "bridge" ||
      network.name === "host" ||
      network.name === "none",
  );

  const cbClass = $derived(
    checked
      ? "bg-violet-600 border-violet-600 text-white animate-scaleIn"
      : isDefaultNetwork
        ? "border-slate-200 dark:border-slate-800 bg-slate-100 dark:bg-slate-900 cursor-not-allowed"
        : "border-slate-300 dark:border-slate-700 bg-transparent hover:border-violet-500 cursor-pointer",
  );

  // Curated accent border and badge based on driver
  const [accentClass, driverBadgeClass] = $derived.by(() => {
    switch (network.driver) {
      case "bridge":
        return [
          "border-l-blue-500",
          "bg-blue-50 dark:bg-blue-950/30 border-blue-200 dark:border-blue-900/40 text-blue-700 dark:text-blue-400",
        ];
      case "host":
        return [
          "border-l-amber-500",
          "bg-amber-50 dark:bg-amber-950/30 border-amber-200 dark:border-amber-900/40 text-amber-700 dark:text-amber-400",
        ];
      case "overlay":
        return [
          "border-l-purple-500",
          "bg-purple-50 dark:bg-purple-950/30 border-purple-200 dark:border-purple-900/40 text-purple-700 dark:text-purple-400",
        ];
      case "macvlan":
        return [
          "border-l-rose-500",
          "bg-rose-50 dark:bg-rose-950/30 border-rose-200 dark:border-rose-900/40 text-rose-700 dark:text-rose-400",
        ];
      default:
        return [
          "border-l-slate-400",
          "bg-slate-50 dark:bg-slate-800 border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-400",
        ];
    }
  });

  const shortId = $derived(network.id.substring(0, 12));
  const containersCount = $derived(network.containers?.length || 0);
</script>

<div
  class="relative rounded-2xl bg-[#f0f3f8] dark:bg-[#0c1220] border border-slate-300/80 dark:border-slate-800/80 hover:border-slate-400 dark:hover:border-slate-700 transition-all duration-300 flex flex-col justify-between overflow-hidden self-start w-full text-slate-700 dark:text-slate-200 shadow-md dark:shadow-lg dark:shadow-black/40 p-4 gap-3.5"
>
  <!-- Header -->
  <div class="flex items-center justify-between gap-3">
    <button
      type="button"
      class="w-5.5 h-5.5 rounded-lg border-2 flex items-center justify-center transition-all duration-150 shrink-0 {checked
        ? 'bg-violet-600 border-violet-500 text-white'
        : isDefaultNetwork
          ? 'border-slate-200 dark:border-slate-800 bg-slate-100 dark:bg-slate-900 cursor-not-allowed'
          : 'border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-900/60 hover:border-violet-500 cursor-pointer'}"
      disabled={isDefaultNetwork}
      onclick={on_toggle}
    >
      {#if checked}
        <span class="text-white text-xs font-bold leading-none">✓</span>
      {/if}
    </button>

    <button
      type="button"
      class="flex-1 font-mono font-bold text-sm tracking-tight truncate text-slate-850 dark:text-slate-100 px-4 py-2 rounded-2xl bg-white dark:bg-slate-800/30 border border-slate-200/80 dark:border-slate-700/30 flex items-center justify-between gap-2.5 cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors text-left shadow-xs"
      onclick={() => (expanded = !expanded)}
    >
      <div class="flex items-center gap-2 truncate">
        <span class="truncate font-semibold text-slate-855 dark:text-white grow"
          >{network.name}</span
        >
      </div>
      <span class="text-[10px] text-slate-400 dark:text-slate-500 shrink-0"
        >{expanded ? "▲" : "▼"}</span
      >
    </button>

    <div class="flex items-center gap-2 shrink-0">
      <span class="px-2.5 py-1.5 rounded-xl text-xs font-bold border {driverBadgeClass}">
        {network.driver}
      </span>
      <span
        class="bg-white dark:bg-slate-800/40 text-slate-600 dark:text-slate-400 text-xs px-2.5 py-1.5 rounded-xl font-semibold border border-slate-200/80 dark:border-slate-700/50 shadow-xs"
      >
        {t("networks.containers_count", { count: containersCount })}
      </span>
    </div>
  </div>

  <!-- Expanded Details -->
  {#if expanded}
    <div class="flex flex-col gap-3 text-xs pt-1">
      <!-- ID & Scope Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <!-- ID -->
        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-violet-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <div class="grow flex flex-col gap-0.5 min-w-0">
            <span class="text-[9px] text-violet-500 font-bold uppercase tracking-wider">
              {t("networks.card_id")}
            </span>
            <span class="font-mono text-xs font-bold text-blue-500 dark:text-blue-400 truncate">
              {network.id.substring(0, 12)}
            </span>
          </div>
        </div>

        <!-- Scope -->
        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-cyan-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <div class="flex flex-col gap-0.5">
            <span class="text-[9px] text-cyan-500 font-bold uppercase tracking-wider">
              {t("networks.card_scope")}
            </span>
            <span class="font-semibold text-slate-700 dark:text-slate-200 capitalize text-xs">
              {network.scope}
            </span>
          </div>
        </div>
      </div>

      <!-- Subnet & Gateway Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <!-- Subnet -->
        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-blue-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <div class="flex flex-col gap-0.5">
            <span class="text-[9px] text-blue-500 font-bold uppercase tracking-wider">
              {t("networks.card_subnet")}
            </span>
            <span class="font-mono font-semibold text-slate-700 dark:text-slate-200 text-xs">
              {network.subnet}
            </span>
          </div>
        </div>

        <!-- Gateway -->
        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-emerald-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <div class="flex flex-col gap-0.5">
            <span class="text-[9px] text-emerald-500 font-bold uppercase tracking-wider">
              {t("networks.card_gateway")}
            </span>
            <span class="font-mono font-semibold text-slate-700 dark:text-slate-200 text-xs">
              {network.gateway}
            </span>
          </div>
        </div>
      </div>

      <!-- Connected Containers -->
      <div
        class="flex flex-col gap-2 p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-indigo-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
      >
        <span class="text-[9px] text-indigo-500 font-bold uppercase tracking-wider">
          {t("networks.card_containers_connected")}
        </span>
        {#if network.containers && network.containers.length > 0}
          <div class="flex flex-col gap-1.5 pt-0.5">
            {#each network.containers as container}
              <div class="flex items-center justify-between text-xs pt-1">
                <div class="flex items-center gap-2">
                  <span class="font-semibold text-slate-700 dark:text-slate-200">{container.name}</span>
                  <span class="font-mono text-[10px] text-slate-400 dark:text-slate-500">({container.ip})</span>
                </div>
                <button
                  type="button"
                  class="px-2.5 py-1 text-[10px] font-bold text-red-600 hover:text-white bg-red-50 dark:bg-red-950/20 hover:bg-red-500 dark:hover:bg-red-600 border border-red-200 dark:border-red-900/50 rounded-lg cursor-pointer transition-colors"
                  onclick={() => on_disconnect(container.name)}
                >
                  {t("networks.card_disconnect_btn")}
                </button>
              </div>
            {/each}
          </div>
        {:else}
          <div class="py-2 text-left italic text-slate-400 dark:text-slate-500 text-xs">
            {t("networks.card_no_containers")}
          </div>
        {/if}
      </div>

      <!-- Footer Action Area -->
      <div class="flex justify-end pt-1">
        <button
          type="button"
          class="px-3.5 py-1.5 text-xs font-bold rounded-xl cursor-pointer transition-colors bg-blue-600 hover:bg-blue-700 text-white shadow-md shadow-blue-500/20"
          onclick={on_connect}
        >
          {t("networks.connect_title")}
        </button>
      </div>
    </div>
  {/if}
</div>
