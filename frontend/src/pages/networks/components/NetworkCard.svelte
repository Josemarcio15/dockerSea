<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { ButtonGreen, ButtonRed } from "$shared/components/buttons";
  import type { DockerNetwork } from "$lib/domains/networks";
  import { isDefaultNetwork as checkIsDefault } from "$lib/domains/networks";

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

  const isDefaultNetwork = $derived(checkIsDefault(network.name));

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

  const containersCount = $derived(network.containers?.length || 0);
</script>

<div
  class="relative rounded-2xl bg-white dark:bg-[#0b101d] border border-slate-200/80 dark:border-slate-800/80 hover:border-violet-500/40 dark:hover:border-violet-500/40 hover:shadow-lg dark:hover:shadow-violet-950/20 transition-all duration-200 flex flex-col justify-between overflow-hidden self-start w-full text-slate-700 dark:text-slate-200 shadow-sm p-3.5 gap-3 group"
>
  <!-- Card Header Compacto -->
  <div class="flex items-center gap-2.5 min-w-0">
    <!-- Checkbox -->
    <button
      type="button"
      class="w-5 h-5 rounded-lg border flex items-center justify-center transition-all duration-150 shrink-0 {checked
        ? 'bg-violet-600 border-violet-500 text-white shadow-xs'
        : isDefaultNetwork
          ? 'border-slate-200 dark:border-slate-800 bg-slate-100 dark:bg-slate-900 cursor-not-allowed opacity-40'
          : 'border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 hover:border-violet-400 cursor-pointer'}"
      disabled={isDefaultNetwork}
      onclick={on_toggle}
    >
      {#if checked}
        <span class="text-white text-[11px] font-bold leading-none">✓</span>
      {/if}
    </button>

    <!-- Ícone de Rede -->
    <div
      class="w-10 h-10 rounded-xl bg-slate-50 dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 p-2 flex items-center justify-center shrink-0 shadow-inner text-base"
    >
      🌐
    </div>

    <!-- Nome + Driver + Status -->
    <div class="flex flex-col min-w-0 flex-1">
      <span
        class="font-bold text-xs text-slate-900 dark:text-white truncate"
        title={network.name}
      >
        {network.name}
      </span>

      <div class="flex items-center gap-1.5 mt-0.5 min-w-0">
        <span
          class="text-[10px] font-mono px-1.5 py-0.2 rounded font-semibold border {driverBadgeClass}"
        >
          {network.driver}
        </span>

        <span class="text-[10px] font-medium flex items-center gap-1 text-slate-500 dark:text-slate-400">
          <span
            class="w-1.5 h-1.5 rounded-full {containersCount > 0
              ? 'bg-emerald-500'
              : 'bg-slate-300 dark:bg-slate-600'} shrink-0"
          ></span>
          {containersCount > 0 ? `${containersCount} container(s)` : 'Sem containers'}
        </span>
      </div>
    </div>

    <!-- Botão Expandir -->
    <button
      type="button"
      class="w-7 h-7 rounded-lg bg-slate-100 dark:bg-slate-800/60 hover:bg-slate-200 dark:hover:bg-slate-800 text-slate-500 dark:text-slate-400 flex items-center justify-center cursor-pointer transition-colors text-[10px] shrink-0 border border-slate-200/50 dark:border-slate-700/50"
      onclick={() => (expanded = !expanded)}
      title="Mais detalhes"
    >
      {expanded ? "▲" : "▼"}
    </button>
  </div>

  <!-- Expanded Details -->
  {#if expanded}
    <div class="flex flex-col gap-2.5 text-xs pt-1 border-t border-slate-100 dark:border-slate-800/60">
      <!-- ID & Subnet Grid -->
      <div class="grid grid-cols-2 gap-2 text-[11px]">
        <div class="flex flex-col p-2 rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200/50 dark:border-slate-800/50">
          <span class="text-[9px] text-slate-400 uppercase font-bold tracking-wider">ID</span>
          <span class="font-mono text-blue-600 dark:text-blue-400 font-semibold">{network.id.substring(0, 8)}</span>
        </div>

        <div class="flex flex-col p-2 rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200/50 dark:border-slate-800/50">
          <span class="text-[9px] text-slate-400 uppercase font-bold tracking-wider">Subnet</span>
          <span class="font-mono text-slate-700 dark:text-slate-200 font-semibold truncate">
            {network.subnet || "—"}
          </span>
        </div>
      </div>
      <!-- Escopo e Gateway -->
      <div class="flex items-center justify-between text-[11px] px-1">
        <span class="text-slate-400 text-[10px] uppercase font-bold tracking-wider">
          {t("networks.card_scope")}
        </span>
        <span class="font-medium text-slate-700 dark:text-slate-300">
          {network.scope}
        </span>
      </div>

      {#if network.gateway}
        <div class="flex items-center justify-between text-[11px] px-1">
          <span class="text-slate-400 text-[10px] uppercase font-bold tracking-wider">
            {t("networks.card_gateway")}
          </span>
          <span class="font-mono text-slate-700 dark:text-slate-300">
            {network.gateway}
          </span>
        </div>
      {/if}

      <!-- Connected Containers -->
      <div class="flex flex-col gap-1.5 p-2 rounded-xl bg-purple-50/40 dark:bg-purple-950/20 border border-purple-200/50 dark:border-purple-900/30">
        <span class="text-[9px] text-purple-600 dark:text-purple-400 font-bold uppercase tracking-wider">
          {t("networks.card_containers_connected")}
        </span>
        {#if network.containers && network.containers.length > 0}
          <div class="flex flex-col gap-1">
            {#each network.containers as container}
              <div class="flex items-center justify-between text-xs py-0.5">
                <div class="flex items-center gap-1.5 truncate">
                  <span class="font-semibold text-slate-800 dark:text-slate-200 text-[11px] truncate">
                    {container.name}
                  </span>
                  {#if container.ip}
                    <span class="font-mono text-[10px] text-slate-400 dark:text-slate-500">
                      ({container.ip})
                    </span>
                  {/if}
                </div>
                <ButtonRed
                  size="xs"
                  onclick={() => on_disconnect(container.name)}
                >
                  {t("networks.card_disconnect_btn")}
                </ButtonRed>
              </div>
            {/each}
          </div>
        {:else}
          <p class="text-[11px] italic text-slate-400 dark:text-slate-500 my-0.5">
            {t("networks.card_no_containers")}
          </p>
        {/if}
      </div>

      <!-- Actions -->
      <div class="flex gap-2 pt-1">
        <ButtonGreen
          size="xs"
          class="flex-1"
          onclick={on_connect}
        >
          {t("networks.connect_title")}
        </ButtonGreen>
        {#if !isDefaultNetwork}
          <ButtonRed
            size="xs"
            onclick={on_delete}
          >
            {t("common.delete")}
          </ButtonRed>
        {/if}
      </div>
    </div>
  {/if}
</div>
