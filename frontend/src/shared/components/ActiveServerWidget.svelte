<script lang="ts">
  import { statsState, subscribeToStats, reconnectStats } from "$shared/stores/stats.svelte";
  import { t } from "$shared/stores/locale.svelte";
  import { ButtonRed } from "$shared/components/buttons";

  let { vps = null }: { vps?: { name: string; id: string } | null } = $props();

  $effect(() => {
    if (vps?.id) {
      return subscribeToStats(vps.id);
    }
  });

  let isOnline = $derived(statsState.online);
  let isChecking = $derived(statsState.checking);

  let totalCpu = $derived(
    statsState.stats
      .reduce((acc, curr) => {
        const val = parseFloat(curr.CPUPerc.replace("%", ""));
        return acc + (isNaN(val) ? 0 : val);
      }, 0)
      .toFixed(2),
  );

  let totalMem = $derived(
    statsState.stats
      .reduce((acc, curr) => {
        const val = parseFloat(curr.MemPerc.replace("%", ""));
        return acc + (isNaN(val) ? 0 : val);
      }, 0)
      .toFixed(2),
  );

  let statusBadgeClass = $derived(
    isOnline
      ? "bg-emerald-500/15 text-emerald-300 border-emerald-500/30"
      : isChecking
        ? "bg-amber-500/15 text-amber-300 border-amber-500/30 animate-pulse"
        : "bg-rose-500/15 text-rose-300 border-rose-500/30",
  );

  let dotClass = $derived(
    isOnline
      ? "bg-emerald-400 shadow-lg shadow-emerald-400/50"
      : isChecking
        ? "bg-amber-400 shadow-lg shadow-amber-400/50"
        : "bg-rose-400 shadow-lg shadow-rose-400/50",
  );

  let statusText = $derived(
    isOnline
      ? t("app.online")
      : isChecking
        ? t("app.connecting")
        : t("app.offline"),
  );
</script>

{#if vps}
  <div
    class="mx-3 mt-4 p-3.5 rounded-xl bg-white/4 border border-white/10 backdrop-blur-md shadow-lg flex flex-col gap-2.5"
  >
    <div class="flex items-center justify-between">
      <span
        class="text-[10px] text-violet-300 font-bold uppercase tracking-wider"
        >{t("app.active_server")}</span
      >
      <span
        class="text-[9px] font-bold px-2 py-0.5 rounded-full border flex items-center gap-1.5 {statusBadgeClass}"
      >
        <div class="w-1.5 h-1.5 rounded-full {dotClass}"></div>
        {statusText}
      </span>
    </div>
    <div class="flex items-center gap-2.5">
      <div
        class="w-8 h-8 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center text-violet-300 shrink-0 shadow-inner"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-4 h-4"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M5.25 5.25h13.5m-13.5 0A2.25 2.25 0 0 0 3 7.5v9a2.25 2.25 0 0 0 2.25 2.25h13.5A2.25 2.25 0 0 0 21 16.5v-9a2.25 2.25 0 0 0-2.25-2.25M6.75 9.75h.008v.008H6.75V9.75Zm0 4.5h.008v.008H6.75v-.008Zm10.5-4.5h.008v.008h-.008V9.75Zm0 4.5h.008v.008h-.008v-.008Z"
          />
        </svg>
      </div>
      <div class="flex flex-col min-w-0 flex-1">
        <span class="text-sm text-violet-200 font-semibold truncate"
          >{vps.name}</span
        >
      </div>
    </div>
    {#if !isOnline && !isChecking}
      <ButtonRed
        size="xs"
        class="w-full mt-0.5"
        onclick={reconnectStats}
      >
        🔌 {t("app.reconnect")}
      </ButtonRed>
    {/if}
  </div>
{/if}
