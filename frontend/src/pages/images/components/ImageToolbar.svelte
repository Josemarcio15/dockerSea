<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import ColumnLayoutSwitcher from "$shared/components/ColumnLayoutSwitcher.svelte";
  import { ButtonGreen, ButtonYellow, ButtonRed } from "$shared/components/buttons";

  let {
    diskUsage = "0 B",
    totalCount = 0,
    selectedCount = 0,
    allSelected = false,
    imageFilter = $bindable("all"),
    searchQuery = $bindable(""),
    countInUse = 0,
    countUnused = 0,
    countDangling = 0,
    onToggleAll = () => {},
    onPrune = () => {},
    onDeleteSelected = () => {},
  }: {
    diskUsage?: string;
    totalCount?: number;
    selectedCount?: number;
    allSelected?: boolean;
    imageFilter?: "all" | "in_use" | "unused" | "dangling";
    searchQuery?: string;
    countInUse?: number;
    countUnused?: number;
    countDangling?: number;
    onToggleAll?: () => void;
    onPrune?: () => void;
    onDeleteSelected?: () => void;
  } = $props();
</script>

<div
  class="flex flex-col gap-3 bg-white dark:bg-[#0b0f19] border border-slate-200/80 dark:border-slate-800/80 p-3.5 rounded-2xl shadow-sm"
>
  <!-- Top Row: Seleção (Esquerda como em Containers) & Ações/Disk stats (Direita) -->
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div class="flex items-center gap-2">
      <!-- Botão Marcar Todos (Verde) -->
      <ButtonGreen
        size="sm"
        onclick={onToggleAll}
      >
        {allSelected ? t("common.deselect_all") : t("common.select_all")}
      </ButtonGreen>

      {#if selectedCount > 0}
        <span
          class="text-xs font-semibold text-violet-600 dark:text-violet-400 animate-pulse px-2"
        >
          {selectedCount}
          {t("common.selected")}
        </span>
      {/if}

      <span
        class="text-xs text-slate-400 dark:text-slate-500 px-2 font-semibold"
      >
        {t("images.registered_images", { count: String(totalCount) })}
      </span>
    </div>

    <div class="flex items-center gap-2">
      <div
        class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-violet-500/10 border border-violet-500/20 text-violet-600 dark:text-violet-400 text-xs font-bold"
      >
        <span>📊 {t("images.disk_space")}</span>
        <span class="font-mono text-sm">{diskUsage}</span>
      </div>

      <ButtonYellow
        size="sm"
        onclick={onPrune}
      >
        {t("volumes.prune_btn")}
      </ButtonYellow>

      <ButtonRed
        size="sm"
        disabled={selectedCount === 0}
        onclick={onDeleteSelected}
      >
        {t("images.delete_selected")}
      </ButtonRed>
    </div>
  </div>

  <!-- Bottom Row: Filter chips & Search input -->
  <div
    class="flex flex-wrap items-center justify-between gap-3 pt-2 border-t border-slate-100 dark:border-slate-800/50"
  >
    <div class="flex items-center gap-1.5 flex-wrap">
      <button
        type="button"
        class="px-3 py-1.5 text-xs font-bold rounded-xl border transition-all cursor-pointer {imageFilter ===
        'all'
          ? 'bg-violet-600 text-white border-violet-600 shadow-md shadow-violet-500/20'
          : 'bg-slate-100 dark:bg-slate-900 text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-800 hover:border-violet-500'}"
        onclick={() => (imageFilter = "all")}
      >
        {t("images.filter_all", { count: String(totalCount) })}
      </button>
      <button
        type="button"
        class="px-3 py-1.5 text-xs font-bold rounded-xl border transition-all cursor-pointer {imageFilter ===
        'in_use'
          ? 'bg-emerald-600 text-white border-emerald-600 shadow-md shadow-emerald-500/20'
          : 'bg-slate-100 dark:bg-slate-900 text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-800 hover:border-emerald-500'}"
        onclick={() => (imageFilter = "in_use")}
      >
        {t("images.filter_in_use", { count: String(countInUse) })}
      </button>
      <button
        type="button"
        class="px-3 py-1.5 text-xs font-bold rounded-xl border transition-all cursor-pointer {imageFilter ===
        'unused'
          ? 'bg-blue-600 text-white border-blue-600 shadow-md shadow-blue-500/20'
          : 'bg-slate-100 dark:bg-slate-900 text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-800 hover:border-blue-500'}"
        onclick={() => (imageFilter = "unused")}
      >
        {t("images.filter_unused", { count: String(countUnused) })}
      </button>
      {#if countDangling > 0}
        <button
          type="button"
          class="px-3 py-1.5 text-xs font-bold rounded-xl border transition-all cursor-pointer {imageFilter ===
          'dangling'
            ? 'bg-amber-600 text-white border-amber-600 shadow-md shadow-amber-500/20'
            : 'bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/30 hover:border-amber-500'}"
          onclick={() => (imageFilter = "dangling")}
        >
          {t("images.filter_dangling", { count: String(countDangling) })}
        </button>
      {/if}
    </div>

    <div class="flex items-center gap-2">
      <!-- 1, 2, 3 Columns Switcher -->
      <ColumnLayoutSwitcher />

      <input
        type="text"
        placeholder={t("common.search")}
        class="px-4 py-2 text-xs rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900 text-slate-855 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:border-violet-500 transition-all w-60"
        bind:value={searchQuery}
      />
    </div>
  </div>
</div>
