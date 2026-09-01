<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import ColumnLayoutSwitcher from "$shared/components/ColumnLayoutSwitcher.svelte";
  import {
    ButtonPurple,
    ButtonGreen,
    ButtonBlue,
    ButtonYellow,
    ButtonRed,
    ButtonPink,
  } from "$shared/components/buttons";

  let {
    searchQuery = $bindable(""),
    selectedCount = 0,
    totalCount = 0,
    allSelected = false,
    onRefresh = () => {},
    onToggleAll = () => {},
    onStart = () => {},
    onStop = () => {},
    onRestart = () => {},
    onRemove = () => {},
  }: {
    searchQuery?: string;
    selectedCount?: number;
    totalCount?: number;
    allSelected?: boolean;
    onRefresh?: () => void;
    onToggleAll?: () => void;
    onStart?: () => void;
    onStop?: () => void;
    onRestart?: () => void;
    onRemove?: () => void;
  } = $props();
</script>

<div class="space-y-6">
  <!-- Top Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div
      class="inline-flex items-center px-5 py-2.5 rounded-2xl bg-linear-to-br from-violet-100 to-fuchsia-100 dark:from-violet-950/40 dark:to-fuchsia-950/40 border border-violet-200/50 dark:border-violet-800/50 self-start shadow-sm"
    >
      <h1
        class="text-2xl font-bold bg-linear-to-r from-violet-600 to-fuchsia-500 bg-clip-text text-transparent m-0 flex items-center gap-2"
      >
        {t("containers.title")}
      </h1>
    </div>

    <div class="flex items-center gap-2">
      <!-- 1, 2, 3 Colunas Switcher -->
      <ColumnLayoutSwitcher />

      <input
        type="text"
        placeholder={t("common.search")}
        class="px-4 py-2 text-xs rounded-xl border border-slate-300 dark:border-slate-800 bg-white dark:bg-[#0b0f19] text-slate-800 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 shadow-2xs transition-all w-60"
        bind:value={searchQuery}
      />
      <ButtonPink
        size="sm"
        title={t("common.refresh")}
        onclick={onRefresh}
      >
        {#snippet icon()}
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="w-3.5 h-3.5"
          >
            <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.19" />
          </svg>
        {/snippet}
        {t("common.refresh")}
      </ButtonPink>
    </div>
  </div>

  <!-- Action Toolbar -->
  <div
    class="flex flex-wrap items-center justify-between gap-3 bg-white dark:bg-[#0b0f19] border border-slate-200 dark:border-slate-800/80 p-3.5 rounded-2xl shadow-sm hover:shadow-md transition-shadow"
  >
    <div class="flex items-center gap-2">
      <ButtonGreen
        size="sm"
        onclick={onToggleAll}
      >
        {#snippet icon()}
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="w-3.5 h-3.5"
          >
            <polyline points="9 11 12 14 22 4" />
            <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11" />
          </svg>
        {/snippet}
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
        {totalCount} {t("containers.title")}
      </span>
    </div>

    <div class="flex items-center gap-2">
      <!-- Iniciar (Verde) -->
      <ButtonGreen
        size="sm"
        disabled={selectedCount === 0}
        onclick={onStart}
      >
        {#snippet icon()}
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="currentColor"
            class="w-3.5 h-3.5"
          >
            <path
              fill-rule="evenodd"
              d="M4.5 5.653c0-1.427 1.529-2.33 2.779-1.643l11.54 6.347c1.295.712 1.295 2.573 0 3.286L7.28 19.99c-1.25.687-2.779-.217-2.779-1.643V5.653Z"
              clip-rule="evenodd"
            />
          </svg>
        {/snippet}
        {t("containers.start")}
      </ButtonGreen>

      <!-- Reiniciar (Azul) -->
      <ButtonBlue
        size="sm"
        disabled={selectedCount === 0}
        onclick={onRestart}
      >
        {#snippet icon()}
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="w-3.5 h-3.5"
          >
            <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8" />
            <path d="M21 3v5h-5" />
            <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16" />
            <path d="M8 16H3v5" />
          </svg>
        {/snippet}
        {t("containers.restart")}
      </ButtonBlue>

      <!-- Parar (Amarelo) -->
      <ButtonYellow
        size="sm"
        disabled={selectedCount === 0}
        onclick={onStop}
      >
        {#snippet icon()}
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="currentColor"
            class="w-3.5 h-3.5"
          >
            <rect width="12" height="12" x="6" y="6" rx="2" />
          </svg>
        {/snippet}
        {t("containers.stop")}
      </ButtonYellow>

      <!-- Remover (Vermelho) -->
      <ButtonRed
        size="sm"
        disabled={selectedCount === 0}
        onclick={onRemove}
      >
        {#snippet icon()}
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="w-3.5 h-3.5"
          >
            <path d="M3 6h18" />
            <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
            <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
          </svg>
        {/snippet}
        {t("containers.delete")}
      </ButtonRed>
    </div>
  </div>
</div>
