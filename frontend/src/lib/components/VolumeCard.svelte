<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import type { DockerVolume } from "$lib/server/docker";

  let {
    vol,
    checked = false,
    on_toggle = () => {},
  }: {
    vol: DockerVolume;
    checked?: boolean;
    on_toggle?: () => void;
  } = $props();

  let expanded = $state(false);
  let size = $state("—");
  let loadingSize = $state(false);

  // Fetch size when expanded
  $effect(() => {
    if (expanded && size === "—" && !loadingSize) {
      loadingSize = true;
      fetch(`/volumes/size?name=${encodeURIComponent(vol.name)}`)
        .then((res) => res.json())
        .then((data) => {
          if (data && data.size) {
            size = data.size;
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
    vol.in_use ? "border-l-emerald-400" : "border-l-red-400",
  );
  const badgeClass = $derived(
    vol.in_use
      ? "text-emerald-700 dark:text-emerald-400 bg-emerald-100 dark:bg-emerald-950/30 border-emerald-250 dark:border-emerald-900/50"
      : "text-red-700 dark:text-red-400 bg-red-100 dark:bg-red-950/30 border-red-250 dark:border-red-900/50",
  );
  const tagText = $derived(
    vol.in_use ? t("volumes.in_use") : t("volumes.free"),
  );

  const cbClass = $derived(
    checked
      ? "bg-violet-600 border-violet-600 text-white"
      : "border-slate-300 dark:border-slate-700 bg-transparent hover:border-violet-500",
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

  const formattedDate = $derived(formatDate(vol.created_at));
  const labels = $derived(Object.entries(vol.labels || {}));
  const containers = $derived(vol.containers || []);
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
      class="flex-1 font-mono font-bold text-sm tracking-tight truncate text-slate-850 dark:text-slate-100 px-4 py-2 rounded-2xl bg-white dark:bg-slate-800/30 border border-slate-200/80 dark:border-slate-700/30 flex items-center justify-between gap-2.5 cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors text-left shadow-xs"
      onclick={() => (expanded = !expanded)}
    >
      <div class="flex items-center gap-2 truncate">
        <span class="truncate font-semibold text-slate-855 dark:text-white grow"
          >{vol.name}</span
        >
      </div>
      <span class="text-[10px] text-slate-400 dark:text-slate-500 shrink-0"
        >{expanded ? "▲" : "▼"}</span
      >
    </button>

    <button
      type="button"
      class="text-[11px] px-3 py-1.5 rounded-xl font-bold border cursor-pointer hover:opacity-90 transition-opacity shrink-0 {vol.in_use
        ? 'text-emerald-600 dark:text-emerald-400 border-emerald-500/30 bg-emerald-500/10'
        : 'text-slate-500 border-slate-400/30 bg-slate-400/10'}"
      onclick={() => (expanded = !expanded)}
    >
      {tagText}
    </button>
  </div>

  <!-- Body (Visible only when expanded) -->
  {#if expanded}
    <div class="flex flex-col gap-3 text-xs pt-1">
      <!-- Mount Point -->
      <div
        class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-violet-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
      >
        <div class="flex-1 min-w-0 flex flex-col gap-0.5">
          <span class="text-[9px] text-violet-500 font-bold uppercase tracking-wider">
            {t("volumes.label_mount")}
          </span>
          <span class="font-mono font-semibold text-slate-700 dark:text-slate-200 break-all leading-normal text-[11px] mt-0.5">
            {vol.mountpoint}
          </span>
        </div>
      </div>

      <!-- Size & Driver Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <!-- Size -->
        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-blue-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <div class="flex flex-col gap-0.5">
            <span class="text-[9px] text-blue-500 font-bold uppercase tracking-wider">
              {t("volumes.label_size")}
            </span>
            <span class="font-semibold text-slate-700 dark:text-slate-200 text-xs font-mono">
              {loadingSize ? t("common.loading") : size}
            </span>
          </div>
        </div>

        <!-- Driver -->
        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-fuchsia-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <div class="flex flex-col gap-0.5">
            <span class="text-[9px] text-fuchsia-500 font-bold uppercase tracking-wider">
              {t("volumes.label_driver")}
            </span>
            <span class="font-mono font-semibold text-slate-700 dark:text-slate-200 text-xs">
              {vol.driver}
            </span>
          </div>
        </div>
      </div>

      <!-- Creation date -->
      {#if formattedDate}
        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-emerald-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <div class="flex flex-col gap-0.5">
            <span class="text-[9px] text-emerald-500 font-bold uppercase tracking-wider">
              {t("volumes.label_created")}
            </span>
            <span class="font-semibold text-slate-700 dark:text-slate-200 text-xs">
              {formattedDate}
            </span>
          </div>
        </div>
      {/if}

      <!-- Connected Containers -->
      {#if containers.length > 0}
        <div
          class="flex flex-col gap-2 p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-teal-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <span class="text-[9px] text-teal-500 font-bold uppercase tracking-wider">
            {t("volumes.connected_containers")}
          </span>
          <div class="flex flex-wrap gap-1.5 pt-0.5">
            {#each containers as [cname, ro]}
              <div
                class="flex items-center gap-1.5 px-2 py-1 rounded-lg border text-[11px] bg-slate-50 dark:bg-[#0c101b] border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-mono font-semibold shadow-xs"
              >
                <span>{cname}</span>
                <span
                  class="text-[9px] px-1.5 py-0.5 rounded border bg-blue-50 dark:bg-blue-950/40 text-blue-700 dark:text-blue-400 border-blue-200 dark:border-blue-900/60 font-sans font-semibold"
                >
                  {ro ? t("containers.read_only") : t("containers.read_write")}
                </span>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Labels -->
      {#if labels.length > 0}
        <div
          class="flex flex-col gap-2 p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-amber-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <span class="text-[9px] text-amber-500 font-bold uppercase tracking-wider">
            {t("volumes.label_labels")}
          </span>
          <div class="flex flex-wrap gap-1.5 pt-0.5">
            {#each labels as [k, v]}
              <div
                class="flex items-center gap-1.5 px-2 py-1 rounded-lg border text-[11px] bg-slate-50 dark:bg-[#0c101b] border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-mono shadow-xs"
              >
                <span>{k}={v}</span>
              </div>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  {/if}
</div>
