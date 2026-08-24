<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import type { DockerVolume } from "$lib/domains/volumes";
  import type { VpsServer } from "../../../../bindings/go-walis/internal/core/db/models.js";
  import { getVolumeSize } from "$lib/domains/volumes";

  let {
    vol,
    server,
    checked = false,
    on_toggle = () => {},
  }: {
    vol: DockerVolume;
    server?: VpsServer;
    checked?: boolean;
    on_toggle?: () => void;
  } = $props();

  let expanded = $state(false);
  let fetchedSize = $state<string | null>(null);
  let loadingSize = $state(false);
  let displaySize = $derived(fetchedSize || vol.size || "—");

  // Fetch size when expanded
  $effect(() => {
    if (expanded && (displaySize === "—" || !displaySize) && !loadingSize && server) {
      loadingSize = true;
      getVolumeSize(server, vol.name)
        .then((s) => {
          if (s) {
            fetchedSize = s;
          }
        })
        .catch((err) => {
          console.error("Erro ao buscar tamanho do volume:", err);
        })
        .finally(() => {
          loadingSize = false;
        });
    }
  });

  const borderAccent = $derived(
    vol.inUse ? "border-l-emerald-400" : "border-l-red-400",
  );
  const badgeClass = $derived(
    vol.inUse
      ? "text-emerald-700 dark:text-emerald-400 bg-emerald-100 dark:bg-emerald-950/30 border-emerald-250 dark:border-emerald-900/50"
      : "text-red-700 dark:text-red-400 bg-red-100 dark:bg-red-950/30 border-red-250 dark:border-red-900/50",
  );
  const tagText = $derived(
    vol.inUse ? t("volumes.in_use") : t("volumes.free"),
  );

  function formatDate(iso: string) {
    if (!iso) return "";
    const date = new Date(iso);
    if (isNaN(date.getTime())) return iso;
    const months = [
      "jan",
      "fev",
      "mar",
      "abr",
      "mai",
      "jun",
      "jul",
      "ago",
      "set",
      "out",
      "nov",
      "dez",
    ];
    const day = date.getDate();
    const month = months[date.getMonth()];
    const year = date.getFullYear();
    return `${day} ${month} ${year}`;
  }

  const formattedDate = $derived(formatDate(vol.createdAt));
  const labels = $derived(Object.entries(vol.labels || {}));
  const containers = $derived(
    ((vol.containers || []).filter(Boolean) as string[][]).filter(
      (c) => c && c.length >= 2,
    ),
  );
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

    <!-- Ícone de Volume / Disco -->
    <div
      class="w-10 h-10 rounded-xl bg-slate-50 dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 p-2 flex items-center justify-center shrink-0 shadow-inner text-base"
    >
      💾
    </div>

    <!-- Nome + Driver + Status -->
    <div class="flex flex-col min-w-0 flex-1">
      <span
        class="font-bold text-xs text-slate-900 dark:text-white truncate"
        title={vol.name}
      >
        {vol.name}
      </span>

      <div class="flex items-center gap-1.5 mt-0.5 min-w-0">
        <span
          class="text-[10px] font-mono px-1.5 py-0.2 rounded bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 font-semibold border border-slate-200/60 dark:border-slate-700/50 truncate"
        >
          {vol.driver || "local"}
        </span>

        <span
          class="text-[10px] font-medium flex items-center gap-1 shrink-0 {vol.inUse
            ? 'text-emerald-600 dark:text-emerald-400'
            : 'text-slate-400 dark:text-slate-500'}"
        >
          <span
            class="w-1.5 h-1.5 rounded-full {vol.inUse
              ? 'bg-emerald-500'
              : 'bg-slate-300 dark:bg-slate-600'} shrink-0"
          ></span>
          {tagText}
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
      <!-- Escopo e Tamanho Grid -->
      <div class="grid grid-cols-2 gap-2 text-[11px]">
        <div class="flex flex-col p-2 rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200/50 dark:border-slate-800/50">
          <span class="text-[9px] text-slate-400 uppercase font-bold tracking-wider">Escopo</span>
          <span class="font-mono text-blue-600 dark:text-blue-400 font-semibold">{vol.scope || "local"}</span>
        </div>

        <div class="flex flex-col p-2 rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200/50 dark:border-slate-800/50">
          <span class="text-[9px] text-slate-400 uppercase font-bold tracking-wider">
            {t("volumes.label_size")}
          </span>
          <span class="font-semibold text-slate-700 dark:text-slate-200 truncate">
            {displaySize}
          </span>
        </div>
      </div>

      <!-- Data de Criação -->
      <div class="flex items-center justify-between text-[11px] px-1">
        <span class="text-slate-400 text-[10px] uppercase font-bold tracking-wider">
          {t("volumes.label_created")}
        </span>
        <span class="font-medium text-slate-700 dark:text-slate-300">
          {formattedDate || "—"}
        </span>
      </div>

      <!-- Ponto de Montagem -->
      <div class="flex flex-col gap-1 p-2 rounded-xl bg-slate-50 dark:bg-slate-900/40 border border-slate-200/60 dark:border-slate-800/60">
        <span class="text-[9px] text-slate-400 font-bold uppercase tracking-wider">
          {t("volumes.label_mount")}
        </span>
        <span class="font-mono text-[10px] text-slate-700 dark:text-slate-300 break-all select-all">
          {vol.mountpoint}
        </span>
      </div>

      <!-- Containers Vinculados -->
      {#if containers.length > 0}
        <div class="flex flex-col gap-1.5 p-2 rounded-xl bg-purple-50/40 dark:bg-purple-950/20 border border-purple-200/50 dark:border-purple-900/30">
          <span class="text-[9px] text-purple-600 dark:text-purple-400 font-bold uppercase tracking-wider">
            {t("volumes.connected_containers")}
          </span>
          <div class="flex flex-wrap gap-1">
            {#each containers as c}
              <span class="px-2 py-0.5 rounded-md bg-white dark:bg-slate-900 border border-purple-200/60 dark:border-purple-900/50 text-[10px] text-purple-700 dark:text-purple-300 font-medium">
                {c[0]} <span class="opacity-70 text-[9px]">({c[1]})</span>
              </span>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Labels -->
      {#if labels.length > 0}
        <div class="flex flex-col gap-1.5 p-2 rounded-xl bg-pink-50/40 dark:bg-pink-950/20 border border-pink-200/50 dark:border-pink-900/30">
          <span class="text-[9px] text-pink-600 dark:text-pink-400 font-bold uppercase tracking-wider">
            {t("volumes.label_labels")}
          </span>
          <div class="flex flex-wrap gap-1">
            {#each labels as [k, v]}
              <span class="font-mono text-[9px] px-1.5 py-0.5 rounded bg-white dark:bg-slate-900 border border-pink-200/50 text-slate-600 dark:text-slate-300 truncate max-w-full">
                {k}: {v}
              </span>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  {/if}
</div>
