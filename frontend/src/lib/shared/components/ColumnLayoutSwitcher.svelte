<script lang="ts">
  import { viewModeStore } from "$shared/stores/viewMode.svelte";
</script>

<div
  class="inline-flex rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-[#0b0f19] p-1 shadow-2xs gap-0.5"
>
  <!-- 1 per row (List) -->
  <button
    type="button"
    class="p-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {viewModeStore.columns ===
    1
      ? 'bg-violet-600 text-white shadow-xs'
      : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}"
    onclick={() => viewModeStore.setColumns(1)}
    title="1 per row (List)"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      class="w-4 h-4"
    >
      <line x1="4" x2="20" y1="6" y2="6" />
      <line x1="4" x2="20" y1="12" y2="12" />
      <line x1="4" x2="20" y1="18" y2="18" />
    </svg>
  </button>

  <!-- 2 per row (Double Grid) -->
  <button
    type="button"
    class="p-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {viewModeStore.columns ===
    2
      ? 'bg-violet-600 text-white shadow-xs'
      : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}"
    onclick={() => viewModeStore.setColumns(2)}
    title="2 por linha"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      class="w-4 h-4"
    >
      <rect width="7" height="16" x="4" y="4" rx="1.5" />
      <rect width="7" height="16" x="13" y="4" rx="1.5" />
    </svg>
  </button>

  <!-- 3 per row (Triple Grid) -->
  <button
    type="button"
    class="p-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {viewModeStore.columns ===
    3
      ? 'bg-violet-600 text-white shadow-xs'
      : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}"
    onclick={() => viewModeStore.setColumns(3)}
    title="3 por linha (Padrão)"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      class="w-4 h-4"
    >
      <rect width="4.5" height="16" x="3" y="4" rx="1" />
      <rect width="4.5" height="16" x="9.75" y="4" rx="1" />
      <rect width="4.5" height="16" x="16.5" y="4" rx="1" />
    </svg>
  </button>
</div>
