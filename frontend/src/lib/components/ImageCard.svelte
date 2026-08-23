<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import type { DockerImage } from "$lib/server/docker";

  let {
    img,
    checked = false,
    on_toggle = () => {},
    on_build = () => {},
    on_delete = () => {},
  }: {
    img: DockerImage;
    checked?: boolean;
    on_toggle?: () => void;
    on_build?: () => void;
    on_delete?: () => void;
  } = $props();

  let expanded = $state(false);

  const isInUse = $derived(img.containersUsing && img.containersUsing.length > 0);
  const containerCount = $derived(img.containersUsing ? img.containersUsing.length : 0);
  const repoLower = $derived(img.repo.toLowerCase());

  // Helper para determinar ícone/logo do serviço
  const iconUrl = $derived.by(() => {
    if (repoLower.includes("wordpress"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/wordpress/wordpress-plain.svg";
    if (repoLower.includes("postgres"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/postgresql/postgresql-original.svg";
    if (repoLower.includes("mariadb"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mariadb/mariadb-original.svg";
    if (repoLower.includes("mysql"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mysql/mysql-original.svg";
    if (repoLower.includes("debian"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/debian/debian-original.svg";
    if (repoLower.includes("alpine"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/alpinejs/alpinejs-original.svg";
    if (repoLower.includes("redis"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/redis/redis-original.svg";
    if (repoLower.includes("nginx"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nginx/nginx-original.svg";
    if (repoLower.includes("node"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nodejs/nodejs-original.svg";
    return null;
  });
</script>

<div
  class="relative rounded-2xl bg-[#f0f3f8] dark:bg-[#0c1220] border border-slate-300/80 dark:border-slate-800/80 hover:border-slate-400 dark:hover:border-slate-700 transition-all duration-300 flex flex-col justify-between overflow-hidden self-start w-full text-slate-700 dark:text-slate-200 shadow-md dark:shadow-lg dark:shadow-black/40 p-4 gap-3.5 group"
>
  <!-- Card Header (Checkbox + Logo + Name + Status + Expand Button) -->
  <div class="flex items-center justify-between gap-3">
    <div class="flex items-center gap-3 min-w-0 flex-1">
      <!-- Checkbox -->
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

      <!-- Icon / Logo -->
      <div
        class="w-10 h-10 rounded-xl bg-white dark:bg-slate-900/80 border border-slate-200 dark:border-slate-800 flex items-center justify-center p-2 shrink-0 shadow-xs"
      >
        {#if iconUrl}
          <img src={iconUrl} alt={img.repo} class="w-full h-full object-contain" />
        {:else}
          <span class="text-lg">📦</span>
        {/if}
      </div>

      <!-- Expandable Title Bar -->
      <button
        type="button"
        class="flex-1 font-mono font-bold text-sm tracking-tight text-slate-850 dark:text-slate-100 px-3 py-1.5 rounded-xl bg-white dark:bg-slate-800/30 border border-slate-200/80 dark:border-slate-700/30 flex items-center justify-between gap-2 cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors text-left shadow-xs min-w-0"
        onclick={() => (expanded = !expanded)}
      >
        <div class="flex flex-col items-start gap-1 min-w-0">
          <!-- Status Badge alinhada à esquerda no topo -->
          {#if isInUse}
            <span
              class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-[9px] font-semibold bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20 shrink-0"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 dark:bg-emerald-400 animate-pulse"
              ></span>
              {t("images.status_in_use")}
            </span>
          {:else if img.repo === "<none>"}
            <span
              class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-[9px] font-semibold bg-red-500/10 text-red-600 dark:text-red-400 border border-red-500/20 shrink-0"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-red-500 dark:bg-red-400"></span>
              {t("images.status_dangling")}
            </span>
          {:else}
            <span
              class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-[9px] font-semibold bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/20 shrink-0"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-amber-500 dark:bg-amber-400"></span>
              {t("images.status_unused")}
            </span>
          {/if}

          <!-- Nome da imagem abaixo da badge -->
          <span class="font-bold text-sm text-slate-850 dark:text-white break-all leading-tight" title={img.repo}>
            {img.repo}
          </span>
        </div>
        <span class="text-[10px] text-slate-400 dark:text-slate-500 shrink-0 ml-1">
          {expanded ? "▲" : "▼"}
        </span>
      </button>
    </div>
  </div>

  <!-- Body (Visible only when expanded) -->
  {#if expanded}
    <div class="flex flex-col gap-3.5 pt-1 text-xs">
      <!-- Item Elevado: TAG e TAMANHO Grid -->
      <div class="grid grid-cols-2 gap-3">
        <!-- Tag -->
        <div
          class="flex flex-col gap-0.5 p-3 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-purple-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all min-w-0"
        >
          <span class="text-[9px] text-purple-500 font-bold uppercase tracking-wider">
            {t("images.tag_label_lower")}
          </span>
          <span class="font-bold text-slate-700 dark:text-slate-200 truncate font-mono text-xs" title={img.tag}>
            {img.tag}
          </span>
        </div>

        <!-- Size -->
        <div
          class="flex flex-col gap-0.5 p-3 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-blue-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all min-w-0"
        >
          <span class="text-[9px] text-blue-500 font-bold uppercase tracking-wider">
            {t("images.size_label_lower")}
          </span>
          <span class="font-bold text-slate-700 dark:text-slate-200 truncate font-mono text-xs">
            {img.size}
          </span>
        </div>
      </div>

      <!-- Item Elevado: Containers Em Uso -->
      <div
        class="flex flex-col gap-2 p-3 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-emerald-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
      >
        <span class="text-[9px] text-emerald-500 font-bold uppercase tracking-wider">Containers</span>
        <div class="flex items-center gap-1.5 text-xs font-semibold text-slate-700 dark:text-slate-200">
          <span>📦</span>
          <span
            >{containerCount}
            {containerCount === 1
              ? t("images.containers_count_one")
              : t("images.containers_count_other")}</span
          >
        </div>
        {#if img.containersUsing && img.containersUsing.length > 0}
          <div class="flex flex-wrap gap-1.5 pt-1 border-t border-slate-100 dark:border-slate-800/60">
            {#each img.containersUsing as cName}
              <span
                class="px-2 py-0.5 rounded-lg border text-[11px] font-mono font-semibold bg-emerald-50 dark:bg-emerald-950/30 text-emerald-700 dark:text-emerald-400 border-emerald-200 dark:border-emerald-900/50 break-all"
              >
                {cName}
              </span>
            {/each}
          </div>
        {/if}
      </div>

      <!-- Item Elevado: ID da Imagem -->
      <div
        class="flex items-center justify-between p-3 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-violet-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all text-xs"
      >
        <span class="text-[9px] text-violet-500 font-bold uppercase tracking-wider">IMAGE ID</span>
        <span class="font-mono text-xs font-bold text-blue-500 dark:text-blue-400">
          {img.id.substring(0, 12)}
        </span>
      </div>

      <!-- Card Action Buttons (Build e Deletar) -->
      <div class="grid grid-cols-4 gap-1.5 pt-2 border-t border-slate-200/60 dark:border-slate-800/60">
        <button
          type="button"
          class="col-span-3 flex items-center justify-center py-2 px-2 text-xs font-bold text-white bg-blue-600 hover:bg-blue-500 rounded-xl transition-all shadow-md shadow-blue-600/20 cursor-pointer"
          onclick={on_build}
        >
          {t("images.build_btn_action")}
        </button>

        <button
          type="button"
          class="col-span-1 flex items-center justify-center p-2 text-xs font-bold text-red-400 hover:text-white bg-red-500/10 hover:bg-red-600 border border-red-500/20 rounded-xl transition-all cursor-pointer"
          title={t("images.delete_btn")}
          onclick={on_delete}
        >
          🗑️
        </button>
      </div>
    </div>
  {/if}
</div>
