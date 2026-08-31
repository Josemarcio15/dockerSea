<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import type { Container } from "$lib/domains/containers";
  import { statsState } from "$shared/stores/stats.svelte";
  import { ButtonPurple, ButtonCyan } from "$shared/components/buttons";
  import FormModal from "$shared/components/FormModal.svelte";

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

  const statusLabel = $derived.by(() => {
    const rawStatus = (container.status || "").trim();
    if (rawStatus.includes("Up")) {
      const uptime = rawStatus.replace(/^Up\s*/i, "").trim();
      return uptime ? `Ativo (${uptime})` : "Ativo";
    }
    return rawStatus.includes("Exited") ? t("containers.status_stop") : rawStatus.includes("Paused") ? "Pausado" : "Indisponível";
  });

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

  const networkItems = $derived.by(() => {
    if (!container.networks) return [];
    return Object.entries(container.networks).map(([netName, netEndpoint]) => {
      return {
        name: netName,
        ip: netEndpoint?.ipAddress || "—",
        gateway: netEndpoint?.gateway || "",
      };
    });
  });

  const shortCid = $derived((container.id || "").substring(0, 12));

  const restartPolicyDisplay = $derived.by(() => {
    const policy =
      container.restartPolicy || (container as any).restart_policy;
    if (!policy) return "—";
    switch (policy.toLowerCase()) {
      case "always":
        return t("containers.card_restart_always");
      case "unless-stopped":
        return t("containers.card_restart_unless_stopped");
      case "on-failure":
        return t("containers.card_restart_on_failure");
      case "no":
        return t("containers.card_restart_no");
      default:
        return policy;
    }
  });

  const iconUrl = $derived.by(() => {
    const target = `${container.name} ${container.image || ""}`.toLowerCase();
    if (target.includes("postgres") || target.includes("postgre"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/postgresql/postgresql-original.svg";
    if (target.includes("mariadb"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mariadb/mariadb-original.svg";
    if (target.includes("mysql"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mysql/mysql-original.svg";
    if (target.includes("mongo"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mongodb/mongodb-original.svg";
    if (target.includes("redis"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/redis/redis-original.svg";
    if (target.includes("wordpress"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/wordpress/wordpress-plain.svg";
    if (target.includes("nginx"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nginx/nginx-original.svg";
    if (target.includes("node"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nodejs/nodejs-original.svg";
    if (target.includes("python"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/python/python-original.svg";
    if (target.includes("golang") || target.includes("go"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg";
    if (target.includes("debian"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/debian/debian-original.svg";
    if (target.includes("alpine"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/alpinejs/alpinejs-original.svg";
    if (target.includes("ubuntu"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/ubuntu/ubuntu-plain.svg";
    return null;
  });
</script>

<div
  class="relative rounded-2xl bg-white dark:bg-[#0b101d] border border-slate-200/80 dark:border-slate-800/80 hover:border-violet-500/40 dark:hover:border-violet-500/40 hover:shadow-lg dark:hover:shadow-violet-950/20 transition-all duration-200 flex flex-col justify-between overflow-hidden self-start w-full text-slate-700 dark:text-slate-200 shadow-sm p-3.5 gap-3 group"
>
  <!-- Card Header Compacto -->
  <div class="flex items-center gap-2.5 min-w-0">
    <!-- Checkbox -->
    <button
      type="button"
      class="w-5 h-5 rounded-lg border flex items-center justify-center cursor-pointer transition-all duration-150 shrink-0 {checked
        ? 'bg-violet-600 border-violet-500 text-white shadow-xs'
        : 'border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 hover:border-violet-400'}"
      onclick={on_toggle}
    >
      {#if checked}
        <span class="text-white text-[11px] font-bold leading-none">✓</span>
      {/if}
    </button>

    <!-- Logo / Ícone Tecnológico -->
    <div
      class="w-10 h-10 rounded-xl bg-slate-50 dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 p-2 flex items-center justify-center shrink-0 shadow-inner"
    >
      {#if iconUrl}
        <img
          src={iconUrl}
          alt={container.name}
          class="w-full h-full object-contain"
        />
      {:else}
        <span class="text-base">📦</span>
      {/if}
    </div>

    <!-- Nome + Imagem + Status -->
    <div class="flex flex-col min-w-0 flex-1">
      <span
        class="font-bold text-sm leading-5 text-slate-900 dark:text-white whitespace-normal break-words line-clamp-2"
        title={container.name}
      >
        {container.name}
      </span>

      <div class="flex items-center gap-1.5 mt-0.5 min-w-0">
        <span
          class="text-xs font-mono px-1.5 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 font-semibold border border-slate-200/60 dark:border-slate-700/50 truncate max-w-[130px]"
          title={container.image}
        >
          {container.image}
        </span>

        <span class="text-xs font-medium flex items-center gap-1 shrink-0 {statusColor}">
          <span class="w-1.5 h-1.5 rounded-full {statusDotColor} {pulseClass} shrink-0"></span>
          {statusLabel}
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

  <!-- Expanded Details (Mostra todos os dados apenas ao expandir) -->
  {#if expanded}
    <div class="flex flex-col gap-2.5 text-xs pt-1 border-t border-slate-100 dark:border-slate-800/60">
      <!-- ID & Criado Em Grid -->
      <div class="grid grid-cols-2 gap-2 text-[11px]">
        <div class="flex flex-col p-2 rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200/50 dark:border-slate-800/50">
          <span class="text-[9px] text-slate-400 uppercase font-bold tracking-wider">ID</span>
          <span class="font-mono text-blue-600 dark:text-blue-400 font-semibold">{shortCid}</span>
        </div>

        <div class="flex flex-col p-2 rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200/50 dark:border-slate-800/50">
          <span class="text-[9px] text-slate-400 uppercase font-bold tracking-wider">
            {t("containers.card_created")}
          </span>
          <span class="font-medium text-slate-700 dark:text-slate-300 truncate">
            {createdStr || "—"}
          </span>
        </div>
      </div>

      <!-- Live Stats se ativo -->
      {#if myStats}
        <div class="grid grid-cols-2 gap-2 text-[11px] p-2 rounded-xl bg-violet-50/30 dark:bg-violet-950/20 border border-violet-200/40 dark:border-violet-900/30 font-mono">
          <div>
            <span class="text-[9px] uppercase font-bold text-violet-400">CPU</span>
            <p class="font-semibold text-violet-700 dark:text-violet-300 m-0">{myStats.CPUPerc || "0%"}</p>
          </div>
          <div>
            <span class="text-[9px] uppercase font-bold text-blue-400">RAM</span>
            <p class="font-semibold text-blue-700 dark:text-blue-300 m-0">{myStats.MemUsage || "0B"}</p>
          </div>
        </div>
      {/if}

      <!-- Política de Reinício -->
      <div class="flex flex-col gap-1 p-2 rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200/50 dark:border-slate-800/50 text-[11px]">
        <span class="text-[9px] text-slate-400 uppercase font-bold tracking-wider">
          {t("containers.card_restart")}
        </span>
        <span class="font-medium text-slate-700 dark:text-slate-300">
          {restartPolicyDisplay}
        </span>
      </div>

      <!-- Redes Conectadas -->
      {#if networkItems.length > 0}
        <div class="flex flex-col gap-1.5 p-2 rounded-xl bg-purple-50/40 dark:bg-purple-950/20 border border-purple-200/50 dark:border-purple-900/30">
          <span class="text-[9px] text-purple-600 dark:text-purple-400 font-bold uppercase tracking-wider">
            {t("containers.card_networks")}
          </span>
          <div class="flex flex-wrap gap-1">
            {#each networkItems as net}
              <span class="px-2 py-0.5 rounded-md bg-white dark:bg-slate-900 border border-purple-200/60 dark:border-purple-900/50 text-[10px] text-purple-700 dark:text-purple-300 font-medium">
                {net.name} {net.ip !== '—' ? `(${net.ip})` : ''}
              </span>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Portas Mapeadas -->
      {#if portItems.length > 0}
        <div class="flex flex-col gap-1.5 p-2 rounded-xl bg-blue-50/40 dark:bg-blue-950/20 border border-blue-200/50 dark:border-blue-900/30">
          <span class="text-[9px] text-blue-600 dark:text-blue-400 font-bold uppercase tracking-wider">
            {t("containers.card_ports")}
          </span>
          <div class="flex flex-wrap gap-1">
            {#each portItems as port}
              {#if port}
                <span class="px-2 py-0.5 rounded-md bg-white dark:bg-slate-900 border font-mono text-[10px] font-semibold {port.class}">
                  {port.formatted}
                </span>
              {/if}
            {/each}
          </div>
        </div>
      {/if}

      <!-- Actions -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 pt-1">
        <ButtonCyan
          size="xs"
          class="w-full"
          onclick={() => (showEnv = true)}
        >
          {t("containers.card_view_env")}
        </ButtonCyan>
        <ButtonPurple
          size="xs"
          class="w-full"
          onclick={() => on_open_logs(container.name)}
        >
          {t("containers.card_view_logs")}
        </ButtonPurple>
      </div>
    </div>
  {/if}

  <FormModal bind:show={showEnv} cancelLabel={t("common.close")} title={`${t("containers.card_env_title")} — ${container.name}`}>
          {#if container.env?.length}
            <div class="grid grid-cols-[minmax(8rem,0.8fr)_minmax(0,2fr)] overflow-hidden rounded-xl border border-cyan-200 dark:border-cyan-900/50">
              <div class="bg-slate-100 px-3 py-2 text-[10px] font-bold uppercase tracking-wider text-slate-500 dark:bg-slate-900 dark:text-slate-400">{t("containers.card_env_key")}</div>
              <div class="bg-slate-100 px-3 py-2 text-[10px] font-bold uppercase tracking-wider text-slate-500 dark:bg-slate-900 dark:text-slate-400">{t("containers.card_env_value")}</div>
              {#each container.env as env}
                {@const separator = env.indexOf("=")}
                {@const key = separator >= 0 ? env.slice(0, separator) : env}
                {@const value = separator >= 0 ? env.slice(separator + 1) : ""}
                <code class="block break-all border-t border-slate-200 bg-slate-50 px-3 py-2 text-xs font-semibold text-orange-800 dark:border-slate-800 dark:bg-slate-950/60 dark:text-orange-400">{key}</code>
                <code class="block break-all border-t border-slate-200 bg-slate-50 px-3 py-2 text-xs text-slate-700 dark:border-slate-800 dark:bg-slate-950/60 dark:text-slate-300">{value}</code>
              {/each}
            </div>
    {:else}
      <p class="text-sm text-slate-500">{t("containers.card_env_empty")}</p>
    {/if}
  </FormModal>
</div>
