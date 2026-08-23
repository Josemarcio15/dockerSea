<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import type { DockerVolume } from "$lib/domains/volumes";
  import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
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
  let size = $state(vol.size || "—");
  let loadingSize = $state(false);

  // Fetch size when expanded
  $effect(() => {
    if (expanded && (size === "—" || !size) && !loadingSize && server) {
      loadingSize = true;
      getVolumeSize(server, vol.name)
        .then((s) => {
          if (s) {
            size = s;
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
  class="relative rounded-2xl bg-[#f0f3f8] dark:bg-[#0c1220] border border-slate-300/80 dark:border-slate-800/80 hover:border-slate-400 dark:hover:border-slate-700 transition-all duration-300 flex flex-col justify-between overflow-hidden self-start w-full text-slate-700 dark:text-slate-200 shadow-md dark:shadow-lg dark:shadow-black/40 p-4 gap-3.5"
>
  <!-- Header Row -->
  <div class="flex items-center justify-between gap-3">
    <button
      type="button"
      class="w-5.5 h-5.5 rounded-lg border-2 flex items-center justify-center cursor-pointer transition-all duration-150 shrink-0 {checked
        ? 'bg-violet-600 border-violet-500 text-white'
        : 'border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-900/60 hover:border-violet-500'}"
      onclick={on_toggle}
    >
      {#if checked}
        <span class="text-white text-xs font-bold leading-none">✓</span>
      {/if}
    </button>

    <button
      type="button"
      class="flex-1 font-mono font-bold text-sm tracking-tight truncate text-slate-855 dark:text-slate-100 px-4 py-2 rounded-2xl bg-white dark:bg-slate-800/30 border border-slate-200/80 dark:border-slate-700/30 flex items-center justify-between gap-2.5 cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors text-left shadow-xs"
      onclick={() => (expanded = !expanded)}
    >
      <div class="flex items-center gap-2 truncate">
        <span class="truncate font-semibold text-slate-855 dark:text-white grow">
          {vol.name}
        </span>
      </div>
      <span class="text-[10px] text-slate-400 dark:text-slate-500 shrink-0">
        {expanded ? "▲" : "▼"}
      </span>
    </button>

    <div
      class="px-2.5 py-1 rounded-xl text-xs font-semibold border flex items-center gap-1.5 shrink-0 {badgeClass}"
    >
      <span
        class="w-1.5 h-1.5 rounded-full {vol.inUse
          ? 'bg-emerald-500'
          : 'bg-red-500'} shrink-0"
      ></span>
      {tagText}
    </div>
  </div>

  <!-- Expanded Details -->
  {#if expanded}
    <div class="flex flex-col gap-3 text-xs pt-1">
      <!-- Driver & Size Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <!-- Driver -->
        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-violet-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <div class="grow flex flex-col gap-0.5 min-w-0">
            <span
              class="text-[9px] text-violet-500 font-bold uppercase tracking-wider"
            >
              {t("volumes.card_driver")}
            </span>
            <span
              class="font-mono text-xs font-bold text-blue-500 dark:text-blue-400 truncate"
            >
              {vol.driver}
            </span>
          </div>
        </div>

        <!-- Size -->
        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-cyan-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <div class="flex flex-col gap-0.5">
            <span
              class="text-[9px] text-cyan-500 font-bold uppercase tracking-wider"
            >
              {t("volumes.card_size")}
            </span>
            <span
              class="font-semibold text-slate-700 dark:text-slate-200 capitalize text-xs"
            >
              {#if loadingSize}
                <span class="animate-pulse">{t("common.loading")}</span>
              {:else}
                {size}
              {/if}
            </span>
          </div>
        </div>
      </div>

      <!-- Created At & Scope Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <!-- Created At -->
        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-blue-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <div class="flex flex-col gap-0.5">
            <span
              class="text-[9px] text-blue-500 font-bold uppercase tracking-wider"
            >
              {t("volumes.card_created_at")}
            </span>
            <span
              class="font-semibold text-slate-700 dark:text-slate-200 text-xs"
            >
              {formattedDate || "—"}
            </span>
          </div>
        </div>

        <!-- Scope -->
        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-emerald-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <div class="flex flex-col gap-0.5">
            <span
              class="text-[9px] text-emerald-500 font-bold uppercase tracking-wider"
            >
              {t("volumes.card_scope")}
            </span>
            <span
              class="font-mono font-semibold text-slate-700 dark:text-slate-200 text-xs"
            >
              {vol.scope || "local"}
            </span>
          </div>
        </div>
      </div>

      <!-- Mountpoint -->
      <div
        class="flex flex-col gap-1 p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-indigo-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
      >
        <span
          class="text-[9px] text-indigo-500 font-bold uppercase tracking-wider"
        >
          {t("volumes.card_mountpoint")}
        </span>
        <span
          class="font-mono text-slate-700 dark:text-slate-300 break-all select-all text-xs"
        >
          {vol.mountpoint}
        </span>
      </div>

      <!-- Attached Containers -->
      {#if containers.length > 0}
        <div
          class="flex flex-col gap-2 p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-purple-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <span
            class="text-[9px] text-purple-500 font-bold uppercase tracking-wider"
          >
            {t("volumes.card_containers_using")}
          </span>
          <div class="flex flex-wrap gap-2">
            {#each containers as c}
              <div
                class="flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-slate-50 dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800"
              >
                <span class="font-medium text-slate-800 dark:text-slate-200">
                  {c[0]}
                </span>
                <span
                  class="text-[9px] font-mono uppercase px-1 py-0.5 rounded {c[1] ===
                  'ro'
                    ? 'bg-amber-500/10 text-amber-500'
                    : 'bg-emerald-500/10 text-emerald-500'}"
                >
                  {c[1]}
                </span>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Labels -->
      {#if labels.length > 0}
        <div
          class="flex flex-col gap-2 p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-pink-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <span
            class="text-[9px] text-pink-500 font-bold uppercase tracking-wider"
          >
            {t("volumes.card_labels")}
          </span>
          <div class="flex flex-wrap gap-1.5">
            {#each labels as [k, v]}
              <span
                class="font-mono text-[10px] px-2 py-0.5 rounded-md bg-slate-50 dark:bg-slate-900 border border-slate-200/60 dark:border-slate-800 text-slate-600 dark:text-slate-400"
              >
                {k}: <span class="text-slate-900 dark:text-slate-200">{v}</span>
              </span>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  {/if}
</div>
