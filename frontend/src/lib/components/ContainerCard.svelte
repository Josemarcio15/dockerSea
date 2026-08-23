<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import type { ContainerInfo } from "$lib/server/docker";
  import { statsState } from "$lib/stores/stats.svelte";
  import { BrandButton, Button, DangerButton } from "$lib/components/buttons";

  let {
    container,
    checked = false,
    on_toggle = () => {},
    on_open_logs = (name: string) => {},
  }: {
    container: ContainerInfo;
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
    const s = container.status;
    if (s.includes("Up")) {
      return ["border-l-emerald-400", "text-emerald-600 dark:text-emerald-400"];
    } else if (s.includes("Exited") || s.includes("Paused")) {
      return ["border-l-amber-400", "text-amber-600 dark:text-amber-400"];
    } else {
      return ["border-l-red-400", "text-red-600 dark:text-red-400"];
    }
  });

  const statusDotColor = $derived(
    container.status.includes("Up")
      ? "bg-emerald-500"
      : container.status.includes("Exited") ||
          container.status.includes("Paused")
        ? "bg-amber-500"
        : "bg-red-500",
  );

  const pulseClass = $derived(
    container.status.includes("Up") ? "animate-pulse" : "",
  );

  const statusBg = $derived(
    container.status.includes("Up")
      ? "bg-emerald-50 dark:bg-emerald-950/40 text-emerald-600 dark:text-emerald-400 border-emerald-200 dark:border-emerald-900/60"
      : container.status.includes("Exited") ||
          container.status.includes("Paused")
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

  const shortCid = $derived(container.id.substring(0, 12));

  const restartPolicyDisplay = $derived.by(() => {
    const policy = container.restartPolicy || (container as any).restart_policy;
    switch (policy) {
      case "always":
        return t("containers.card_restart_always");
      case "unless-stopped":
        return t("containers.card_restart_unless_stopped");
      case "on-failure":
        return t("containers.card_restart_on_failure");
      case "":
        return "—";
      default:
        return t("containers.card_restart_no");
    }
  });

  // Networks list mapped to array (compatível com objetos ou arrays)
  const netItems = $derived(
    Object.entries(container.networks || {}).map(([name, netVal]: [string, any]) => {
      const type =
        name === "bridge"
          ? "bridge"
          : name === "host"
            ? "host"
            : name === "none"
              ? "none"
              : name === "overlay"
                ? "overlay"
                : "custom";
      const ip = Array.isArray(netVal) ? netVal[0] : (netVal?.ipAddress || netVal?.IPAddress || "");
      const gw = Array.isArray(netVal) ? netVal[1] : (netVal?.gateway || netVal?.Gateway || "");
      return { name, ip, gw, type };
    }),
  );

  // Mounts list (compatível com MountInfo objeto ou tupla antiga)
  const mountItems = $derived(
    (container.mounts || []).map((m: any) => {
      const typ = Array.isArray(m) ? m[0] : (m.type || "volume");
      const name = Array.isArray(m) ? m[1] : (m.name || "");
      const src = Array.isArray(m) ? m[2] : (m.source || "");
      const dest = Array.isArray(m) ? m[3] : (m.destination || "");
      const ro = Array.isArray(m) ? m[4] : (m.readOnly || false);

      const color =
        typ === "volume"
          ? "bg-blue-50 dark:bg-blue-950/20 border-blue-200 dark:border-blue-900/40 text-blue-700 dark:text-blue-400"
          : typ === "bind"
            ? "bg-amber-50 dark:bg-amber-950/20 border-amber-200 dark:border-amber-900/40 text-amber-700 dark:text-amber-400"
            : typ === "tmpfs"
              ? "bg-purple-50 dark:bg-purple-950/20 border-purple-200 dark:border-purple-900/40 text-purple-700 dark:text-purple-400"
              : "bg-slate-50 dark:bg-slate-850 border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-400";
      return { typ, name, src, dest, color, ro };
    }),
  );

  const cbClass = $derived(
    checked
      ? "bg-violet-600 border-violet-600 text-white animate-scaleIn"
      : "border-slate-300 dark:border-slate-700 bg-transparent hover:border-violet-500",
  );
</script>

<div
  class=" relative rounded-2xl bg-slate-100 dark:bg-[#0c1220] border border-slate-200 dark:border-slate-800/80 hover:border-slate-300 dark:hover:border-slate-700 transition-all duration-300 flex flex-col justify-between overflow-hidden self-start w-full text-slate-700 dark:text-slate-200 shadow-md dark:shadow-lg dark:shadow-black/40 p-4 gap-4"
>
  <!-- Header Row -->
  <div class=" flex items-center justify-between gap-3">
    <button
      type="button"
      class="w-6 h-6 rounded-lg border-2 flex items-center justify-center cursor-pointer transition-all duration-150 shrink-0 {checked
        ? 'bg-violet-600/20 border-violet-500 text-violet-500'
        : 'border-slate-300 dark:border-slate-700 bg-transparent hover:border-violet-500'}"
      onclick={on_toggle}
    >
      {#if checked}
        <span class="text-violet-500 text-xs font-bold leading-none">✓</span>
      {/if}
    </button>

    <button
      type="button"
      class="flex-1 font-mono font-bold text-sm tracking-tight truncate text-slate-750 dark:text-slate-100 px-4 py-2 rounded-2xl bg-slate-100/60 dark:bg-slate-800/30 border border-slate-200/50 dark:border-slate-700/30 flex items-center justify-between gap-2.5 cursor-pointer hover:bg-slate-200/50 dark:hover:bg-slate-800/50 transition-colors text-left"
      onclick={() => (expanded = !expanded)}
    >
      <div class="flex items-center gap-2 truncate">
        <span class="truncate font-semibold text-slate-855 dark:text-white grow"
          >{container.name}</span
        >
      </div>
      <span class="text-[10px] text-slate-400 dark:text-slate-500 shrink-0"
        >{expanded ? "▲" : "▼"}</span
      >
    </button>

    <button
      type="button"
      class="flex items-center gap-1.5 px-3 py-2 rounded-xl text-xs font-bold border cursor-pointer bg-transparent hover:opacity-90 transition-opacity shrink-0 border-slate-500/20 {container.status.includes(
        'Up',
      )
        ? 'text-emerald-500 border-emerald-500/30 bg-emerald-500/10'
        : 'text-amber-500 border-amber-500/30 bg-amber-500/10'}"
      onclick={() => (expanded = !expanded)}
    >
      <div class="w-2 h-2 rounded-full {statusDotColor} {pulseClass}"></div>
      <span class={statusColor}>{container.status}</span>
    </button>
  </div>

  <!-- Body (Visible only when expanded) -->
  {#if expanded}
    <div class="flex flex-col gap-3 text-xs pt-1">
      <!-- DOCKER IMAGE -->
      <div
        class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-violet-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
      >
        <div class="grow flex flex-col gap-0.5 min-w-0">
          <span class="text-[9px] text-violet-500 font-bold uppercase tracking-wider">
            {t("containers.card_image")}
          </span>
          <span class="font-semibold text-slate-700 dark:text-slate-200 truncate font-mono text-xs">
            {container.image}
          </span>
        </div>
        {#if shortCid}
          <div class="hidden md:flex items-center gap-1.5 text-right shrink-0">
            <span class="text-[10px] text-slate-400 dark:text-slate-500 font-bold uppercase">
              {t("containers.card_id")}:
            </span>
            <span class="font-mono text-[10px] text-blue-500 dark:text-blue-400 font-bold" title={container.id}>
              {shortCid}
            </span>
          </div>
        {/if}
      </div>

      <!-- Live Stats Row (CPU / RAM) -->
      {#if myStats}
        <div class="grid grid-cols-2 gap-3">
          <!-- CPU -->
          <div
            class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-emerald-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
          >
            <div class="flex flex-col gap-0.5">
              <span class="text-[9px] text-emerald-500 font-bold uppercase tracking-wider">CPU</span>
              <span class="font-mono text-emerald-600 dark:text-emerald-400 font-bold text-xs">
                {myStats.CPUPerc}
              </span>
            </div>
          </div>
          <!-- RAM -->
          <div
            class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-blue-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
          >
            <div class="flex flex-col gap-0.5">
              <span class="text-[9px] text-blue-500 font-bold uppercase tracking-wider">RAM</span>
              <span class="font-mono text-blue-600 dark:text-blue-400 font-bold text-xs">
                {myStats.MemUsage} ({myStats.MemPerc})
              </span>
            </div>
          </div>
        </div>
      {/if}

      <!-- Created & Ports Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        {#if createdStr}
          <div
            class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-rose-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
          >
            <div class="flex flex-col gap-0.5">
              <span class="text-[9px] text-rose-500 font-bold uppercase tracking-wider">
                {t("containers.card_created")}
              </span>
              <span class="font-semibold text-slate-700 dark:text-slate-200 text-xs">
                {createdStr}
              </span>
            </div>
          </div>
        {/if}

        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-purple-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <div class="flex flex-col gap-0.5 min-w-0">
            <span class="text-[9px] text-purple-500 font-bold uppercase tracking-wider">
              {t("containers.card_ports")}
            </span>
            {#if portItems.length > 0}
              <div class="flex flex-wrap gap-1 mt-0.5">
                {#each portItems as port}
                  <span
                    class="font-mono font-semibold text-blue-600 dark:text-blue-400 text-[10px] px-2 py-0.5 rounded border shadow-xs {port.class}"
                  >
                    {port.tag}
                    {port.formatted}
                  </span>
                {/each}
              </div>
            {:else}
              <span class="text-slate-400 dark:text-slate-500 italic text-xs">—</span>
            {/if}
          </div>
        </div>
      </div>

      <!-- Docker Networks -->
      {#if netItems.length > 0}
        <div
          class="flex flex-col gap-2 p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-red-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <span class="text-[9px] text-red-500 font-bold uppercase tracking-wider">
            {t("containers.card_networks")}
          </span>
          <div class="flex flex-col gap-1.5">
            {#each netItems as net}
              <div class="flex items-center justify-between text-xs pt-1">
                <div class="flex items-center gap-2">
                  <span
                    class="px-1.5 py-0.5 rounded text-[8px] font-bold uppercase tracking-wide bg-emerald-50 dark:bg-[#0c2920]/60 text-emerald-600 dark:text-emerald-400 border border-emerald-200 dark:border-[#174635]/60"
                  >
                    {net.type}
                  </span>
                  <span class="font-semibold text-slate-700 dark:text-slate-200">{net.name}</span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="font-mono font-bold text-emerald-600 dark:text-emerald-400">{net.ip}</span>
                  {#if net.gw}
                    <span class="text-[9px] text-slate-400 dark:text-slate-500 font-mono">
                      {t("containers.card_gateway")}: {net.gw}
                    </span>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Mounts / Volumes -->
      {#if mountItems.length > 0}
        <div
          class="flex flex-col gap-2 p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-amber-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <span class="text-[9px] text-amber-500 font-bold uppercase tracking-wider">
            {t("containers.card_mounts")}
          </span>
          <div class="flex flex-col gap-2.5 pt-1">
            {#each mountItems as mount}
              <div class="flex flex-col gap-1 text-xs">
                <div class="flex items-center gap-2">
                  <span class="px-1.5 py-0.5 rounded text-[8px] font-bold uppercase tracking-wide border {mount.color}">
                    {mount.typ}
                  </span>
                  {#if mount.name}
                    <span class="font-mono font-semibold text-slate-700 dark:text-slate-200">{mount.name}</span>
                  {/if}
                </div>
                <div class="flex flex-col gap-1 text-[11px] text-slate-550 dark:text-slate-400 font-mono">
                  <div class="flex items-center flex-wrap gap-1.5">
                    <span class="font-semibold text-slate-400 dark:text-slate-500">{t("containers.path_system")}:</span>
                    <span
                      class="text-[10px] px-1.5 py-0.5 rounded-lg border bg-rose-50 dark:bg-rose-950/40 text-rose-700 dark:text-rose-400 border-rose-200 dark:border-rose-900/60 truncate max-w-100 shadow-xs"
                      title={mount.src}>{mount.src}</span
                    >
                  </div>
                  <div class="flex items-center flex-wrap gap-1.5">
                    <span class="font-semibold text-slate-400 dark:text-slate-500">{t("containers.path_container")}:</span>
                    <span
                      class="text-[10px] px-1.5 py-0.5 rounded-lg border bg-rose-50 dark:bg-rose-950/40 text-rose-700 dark:text-rose-400 border-rose-200 dark:border-rose-900/60 truncate max-w-100 shadow-xs"
                      title={mount.dest}>{mount.dest}</span
                    >
                  </div>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Restart Policy -->
      <div
        class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-emerald-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
      >
        <div class="flex flex-col gap-0.5">
          <span class="text-[9px] text-emerald-500 font-bold uppercase tracking-wider">
            {t("containers.card_restart")}
          </span>
          <span class="font-semibold text-slate-700 dark:text-slate-200 text-xs">
            {restartPolicyDisplay}
          </span>
        </div>
      </div>
    </div>

    <!-- Footer / Action Buttons -->
    <div
      class="flex gap-2 justify-end items-center mt-2 pt-2 border-t border-slate-200/30 dark:border-slate-800/30"
    >
      <BrandButton
        size="sm"
        onclick={() => on_open_logs(container.name)}
      >
        {t("containers.card_view_logs")}
      </BrandButton>
      <Button
        variant="info"
        size="sm"
        onclick={() => (showEnv = true)}
      >
        {t("containers.card_view_env")}
      </Button>
    </div>
  {/if}
</div>

<!-- Environment Variables Modal -->
{#if showEnv}
  <div
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="bg-white dark:bg-[#0b0f19] border border-slate-200 dark:border-slate-800 rounded-2xl w-200 max-w-full max-h-[90vh] flex flex-col p-6 shadow-2xl animate-scaleIn text-slate-800 dark:text-slate-200"
    >
      <!-- Header -->
      <div
        class="flex justify-between items-center pb-4 border-b border-slate-200 dark:border-slate-800 shrink-0"
      >
        <h2
          class="text-lg font-bold text-slate-900 dark:text-white flex items-center gap-2"
        >
          <span>🌿</span>
          <span>{t("containers.card_env_title")}</span>
          <span
            class="text-sm font-mono font-normal text-slate-500 dark:text-slate-400"
            >({container.name})</span
          >
        </h2>
        <button
          type="button"
          class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 text-xl cursor-pointer bg-transparent border-none"
          onclick={() => (showEnv = false)}
        >
          ✕
        </button>
      </div>

      <!-- Body -->
      <div class="py-4 overflow-y-auto flex-1 min-h-0">
        {#if container.env && container.env.length > 0}
          <div class="grid grid-cols-1 gap-1.5">
            {#each container.env as envVar}
              {@const idx = envVar.indexOf("=")}
              {@const key = idx > 0 ? envVar.substring(0, idx) : envVar}
              {@const value = idx > 0 ? envVar.substring(idx + 1) : ""}
              <div
                class="flex items-start gap-2 p-2.5 rounded-xl bg-slate-50 dark:bg-[#0c101b] border border-slate-200 dark:border-slate-800/60 text-xs font-mono"
              >
                <span
                  class="font-bold text-cyan-600 dark:text-cyan-400 shrink-0 whitespace-nowrap"
                  >{key}</span
                >
                <span class="text-slate-600 dark:text-slate-300 break-all"
                  >=</span
                >
                <span class="text-slate-700 dark:text-slate-200 break-all"
                  >{value}</span
                >
              </div>
            {/each}
          </div>
        {:else}
          <div
            class="text-sm text-slate-400 dark:text-slate-500 py-12 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-2xl"
          >
            {t("containers.card_env_empty")}
          </div>
        {/if}
      </div>

      <!-- Footer -->
      <div
        class="flex gap-3 justify-end pt-4 border-t border-slate-200 dark:border-slate-800 shrink-0"
      >
        <DangerButton
          size="sm"
          onclick={() => (showEnv = false)}
        >
          {t("common.cancel")}
        </DangerButton>
      </div>
    </div>
  </div>
{/if}
