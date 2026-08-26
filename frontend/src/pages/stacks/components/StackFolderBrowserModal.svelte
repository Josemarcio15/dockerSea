<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { ButtonOrange, ButtonGreen } from "$shared/components/buttons";

  let {
    show = $bindable(false),
    currentPath = "",
    parentPath = null,
    folders = [],
    hasDockerfile = false,
    hasDockerignore = false,
    loading = false,
    onNavigate = () => {},
    onSelect = () => {},
  }: {
    show: boolean;
    currentPath: string;
    parentPath: string | null;
    folders: { name: string; path: string }[];
    hasDockerfile: boolean;
    hasDockerignore: boolean;
    loading: boolean;
    onNavigate: (path: string) => void;
    onSelect: () => void;
  } = $props();
</script>

{#if show}
  <div
    class="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-50 p-4"
  >
    <div
      class="bg-white dark:bg-[#0b0f19] border border-slate-200 dark:border-slate-800 rounded-2xl w-160 max-w-full max-h-[85vh] flex flex-col p-6 shadow-2xl animate-scaleIn text-slate-800 dark:text-slate-200 gap-4"
    >
      <div class="flex justify-between items-center pb-3 border-b border-slate-200 dark:border-slate-800">
        <h3 class="text-base font-bold text-slate-900 dark:text-white flex items-center gap-2">
          {t("stacks.browse_folder_modal_title")}
        </h3>
        <button
          type="button"
          class="text-slate-400 hover:text-slate-200 text-lg bg-transparent border-none cursor-pointer"
          onclick={() => (show = false)}
        >
          ✕
        </button>
      </div>

      <!-- Path bar -->
      <div class="flex items-center gap-2 px-3 py-2 bg-slate-100 dark:bg-slate-900/80 rounded-xl text-xs font-mono text-slate-700 dark:text-slate-300 border border-slate-200 dark:border-slate-800">
        <span class="text-slate-400">Path:</span>
        <span class="truncate flex-1 font-bold">{currentPath || "/"}</span>
        {#if parentPath}
          <ButtonOrange size="xs" onclick={() => onNavigate(parentPath || "")}>
            {t("stacks.folder_up_btn")}
          </ButtonOrange>
        {/if}
      </div>

      <!-- Detection badges -->
      <div class="flex items-center gap-2 flex-wrap">
        {#if hasDockerfile}
          <span class="px-2 py-0.5 rounded-md bg-emerald-500/10 text-emerald-500 font-mono text-[11px] border border-emerald-500/20">
            {t("stacks.folder_has_dockerfile")}
          </span>
        {/if}
        {#if hasDockerignore}
          <span class="px-2 py-0.5 rounded-md bg-blue-500/10 text-blue-500 font-mono text-[11px] border border-blue-500/20">
            {t("stacks.folder_has_dockerignore")}
          </span>
        {/if}
        {#if !hasDockerfile && !hasDockerignore}
          <span class="text-slate-400 text-xs italic">Nenhum Dockerfile ou .dockerignore nesta raiz</span>
        {/if}
      </div>

      <!-- Folders list -->
      <div class="flex-1 min-h-48 max-h-72 overflow-y-auto rounded-xl border border-slate-200 dark:border-slate-800 p-2 bg-slate-50/50 dark:bg-slate-900/30 space-y-1">
        {#if loading}
          <div class="py-8 text-center text-xs text-slate-400 animate-pulse">Carregando diretórios...</div>
        {:else if folders.length === 0}
          <div class="py-8 text-center text-xs text-slate-400">Nenhuma subpasta encontrada</div>
        {:else}
          {#each folders as folder (folder.path)}
            <button
              type="button"
              class="w-full text-left px-3 py-2 rounded-lg text-xs font-medium text-slate-700 dark:text-slate-300 hover:bg-violet-500/10 hover:text-violet-600 dark:hover:text-violet-400 flex items-center gap-2.5 transition-colors cursor-pointer border-none bg-transparent"
              onclick={() => onNavigate(folder.path)}
            >
              <span class="truncate">{folder.name}</span>
            </button>
          {/each}
        {/if}
      </div>

      <!-- Footer Action -->
      <div class="flex gap-2.5 justify-end items-center pt-3 border-t border-slate-200 dark:border-slate-800 shrink-0">
        <ButtonOrange size="sm" onclick={() => (show = false)}>
          {t("common.cancel")}
        </ButtonOrange>
        <ButtonGreen size="sm" onclick={onSelect}>
          Selecionar Esta Pasta
        </ButtonGreen>
      </div>
    </div>
  </div>
{/if}
