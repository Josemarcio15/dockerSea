<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { useRefreshKey } from "$shared/stores/refresh.svelte";
  import * as ExtraService from "../../../bindings/go-walis/internal/extras/extraservice.js";

  let { activeVps }: { activeVps: any } = $props();
  let rows = $state<any[]>([]);
  let loading = $state(true);
  let error = $state("");
  let query = $state("");
  let minPort = $state("");
  let maxPort = $state("");
  let protocol = $state("all");
  let sortOrder = $state("asc");
  let ipVersion = $state("all");

  async function refresh() {
    if (!activeVps) {
      loading = false;
      return;
    }
    loading = true;
    error = "";
    try {
      const list = await ExtraService.ListListeningPorts(activeVps);
      rows = list || [];
    } catch (e: any) {
      error = e.message || String(e);
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    useRefreshKey();
    if (activeVps) refresh();
    else loading = false;
  });

  let filtered = $derived(
    rows
      .filter((row) => {
        const isIpv6 = row.address.includes(":");
        return (
          (!query ||
            `${row.processName} ${row.address}`
              .toLowerCase()
              .includes(query.toLowerCase())) &&
          (!minPort || row.port >= Number(minPort)) &&
          (!maxPort || row.port <= Number(maxPort)) &&
          (protocol === "all" || row.protocol === protocol) &&
          (ipVersion === "all" || (ipVersion === "ipv6" ? isIpv6 : !isIpv6))
        );
      })
      .sort((a, b) =>
        sortOrder === "asc" ? a.port - b.port : b.port - a.port,
      ),
  );
</script>

<div class="space-y-5">
  <div class="flex items-center justify-between gap-3">
    <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100">
      {t("ports.title")}
    </h2>
  </div>

  <div
    class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-6 gap-3 bg-slate-50 dark:bg-slate-950/50 p-4 rounded-2xl border border-slate-200/60 dark:border-slate-800/60"
  >
    <select
      bind:value={sortOrder}
      class="rounded-xl text-xs px-3 py-2 bg-white dark:bg-slate-950 text-slate-800 dark:text-slate-200 border border-slate-200 dark:border-slate-800 focus:outline-none focus:border-violet-500"
    >
      <option value="asc">{t("ports.order_asc")}</option>
      <option value="desc">{t("ports.order_desc")}</option>
    </select>
    <input
      bind:value={query}
      placeholder={t("ports.filter_placeholder")}
      class="rounded-xl text-xs px-3 py-2 bg-white dark:bg-slate-950 text-slate-800 dark:text-slate-200 border border-slate-200 dark:border-slate-800 focus:outline-none focus:border-violet-500"
    />
    <input
      bind:value={minPort}
      type="number"
      placeholder={t("ports.min_port_placeholder")}
      class="rounded-xl text-xs px-3 py-2 bg-white dark:bg-slate-950 text-slate-800 dark:text-slate-200 border border-slate-200 dark:border-slate-800 focus:outline-none focus:border-violet-500"
    />
    <input
      bind:value={maxPort}
      type="number"
      placeholder={t("ports.max_port_placeholder")}
      class="rounded-xl text-xs px-3 py-2 bg-white dark:bg-slate-950 text-slate-800 dark:text-slate-200 border border-slate-200 dark:border-slate-800 focus:outline-none focus:border-violet-500"
    />
    <select
      bind:value={protocol}
      class="rounded-xl text-xs px-3 py-2 bg-white dark:bg-slate-950 text-slate-800 dark:text-slate-200 border border-slate-200 dark:border-slate-800 focus:outline-none focus:border-violet-500"
    >
      <option value="all">{t("ports.all_protocols")}</option>
      <option value="tcp">{t("ports.tcp_only")}</option>
      <option value="udp">{t("ports.udp_only")}</option>
    </select>
    <select
      bind:value={ipVersion}
      class="rounded-xl text-xs px-3 py-2 bg-white dark:bg-slate-950 text-slate-800 dark:text-slate-200 border border-slate-200 dark:border-slate-800 focus:outline-none focus:border-violet-500"
    >
      <option value="all">{t("ports.all_ip_versions")}</option>
      <option value="ipv4">{t("ports.ipv4_only")}</option>
      <option value="ipv6">{t("ports.ipv6_only")}</option>
    </select>
  </div>

  {#if loading}
    <div class="p-12 text-center text-sm text-slate-500 dark:text-slate-400">
      {t("ports.loading")}
    </div>
  {:else if error}
    <div
      class="p-5 rounded-2xl border border-red-200 dark:border-red-900/50 bg-red-50 dark:bg-red-950/30 text-red-700 dark:text-red-400 text-sm"
    >
      {error}
    </div>
  {:else}
    <div
      class="overflow-x-auto rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-[#0c1220] shadow-sm"
    >
      <table class="w-full text-left text-xs text-slate-700 dark:text-slate-300">
        <thead
          class="bg-slate-50 dark:bg-slate-950/70 text-slate-500 dark:text-slate-400 uppercase border-b border-slate-200 dark:border-slate-800"
        >
          <tr>
            <th class="p-4">{t("ports.col_port")}</th>
            <th class="p-4">{t("ports.col_protocol")}</th>
            <th class="p-4">{t("ports.col_process")}</th>
            <th class="p-4">{t("ports.col_pid")}</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as row}
            <tr
              class="border-t border-slate-100 dark:border-slate-800/80 hover:bg-slate-50/50 dark:hover:bg-slate-800/30 transition-colors"
            >
              <td class="p-4 font-bold text-violet-600 dark:text-violet-400 font-mono">
                {row.address}:{row.port}
              </td>
              <td class="p-4 uppercase font-semibold">{row.protocol}</td>
              <td class="p-4 font-semibold">{row.processName}</td>
              <td class="p-4 font-mono text-slate-500">{row.pid ?? "—"}</td>
            </tr>
          {:else}
            <tr>
              <td colspan="4" class="p-12 text-center text-slate-500">
                {t("ports.empty")}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>
