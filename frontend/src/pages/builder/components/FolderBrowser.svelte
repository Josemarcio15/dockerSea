<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import {
    notifySuccess,
    notifyError,
  } from "$shared/stores/notification.svelte";
  import {
    ButtonGreen,
    ButtonYellow,
    ButtonPurple,
  } from "$shared/components/buttons";
  import { folderNameFromPath } from "../service";
  import type { BuilderFolder, BuilderStore } from "../types";

  let { store }: { store: BuilderStore } = $props();
  let editablePath = $state("");
  let lastSyncedPath = $state("");

  $effect(() => {
    if (store.currentPath && store.currentPath !== lastSyncedPath) {
      editablePath = store.currentPath;
      lastSyncedPath = store.currentPath;
    }
  });

  async function navigateToPath() {
    const path = editablePath.trim();
    if (path && path !== store.currentPath) await store.browse(path);
  }
  async function save() {
    try {
      await store.saveCurrentPath();
      notifySuccess("Pasta salva com sucesso!");
    } catch (error: any) {
      notifyError(`Erro ao salvar pasta: ${error?.message || error}`);
    }
  }
  async function remove(path: string, event: Event) {
    event.stopPropagation();
    try {
      await store.removeSavedPath(path);
      notifySuccess("Pasta removida dos favoritos!");
    } catch (error: any) {
      notifyError(`Erro ao remover pasta: ${error?.message || error}`);
    }
  }
</script>

<div class="space-y-4">
  <div
    class="bg-white dark:bg-[#0b0f19] border border-slate-200/70 dark:border-slate-800/80 p-5 rounded-2xl shadow-sm space-y-4"
  >
    <h3
      class="text-sm font-bold text-slate-800 dark:text-slate-200 uppercase tracking-wider"
    >
      {t("builder.nav_title")}
    </h3>
    <div class="flex items-center gap-2 flex-wrap">
      <ButtonPurple size="xs" onclick={() => store.browse()}
        >{t("builder.nav_home")}</ButtonPurple
      >
      {#if store.parentPath}<ButtonYellow
          size="xs"
          onclick={() => store.browse(store.parentPath ?? "")}
          >{t("builder.nav_up")}</ButtonYellow
        >{/if}
    </div>
    {#if store.savedPaths.length > 0}
      <div class="flex items-center gap-2 flex-wrap pt-1 border-t border-slate-100 dark:border-slate-800/60">
      {#each store.savedPaths as path}
        <div class="relative inline-block group">
          <ButtonGreen size="xs" onclick={() => store.browse(path)} title={path}
            ><span class="mr-1">📌</span>{folderNameFromPath(path)}</ButtonGreen
          ><span
            class="absolute -top-1.5 -right-1.5 w-3.5 h-3.5 rounded-full bg-red-500 text-white text-[8px] flex items-center justify-center opacity-0 group-hover:opacity-100 cursor-pointer z-10"
            role="button"
            tabindex="0"
            onclick={(event) => remove(path, event)}
            onkeydown={(event) => {
              if (event.key === "Enter" || event.key === " ")
                remove(path, event);
            }}
            title="Remover atalho">✕</span
          >
        </div>
      {/each}
      </div>
    {/if}
    {#if store.currentPath}<input
        type="text"
        aria-label="Caminho da pasta do projeto"
        class="w-full px-3 py-2 rounded-xl bg-slate-100 dark:bg-slate-900/60 border border-slate-200 dark:border-slate-800 text-xs font-mono text-slate-600 dark:text-slate-400 focus:outline-none focus:border-violet-500"
        bind:value={editablePath}
        title={store.currentPath}
        onkeydown={(event) => { if (event.key === "Enter") void navigateToPath(); }}
        onblur={() => void navigateToPath()}
      />{/if}
    <div
      class="border border-slate-200 dark:border-slate-800 rounded-xl max-h-64 overflow-y-auto bg-slate-50 dark:bg-slate-900/20"
    >
      {#if store.loading}<div class="flex items-center justify-center py-8">
          <span
            class="animate-spin inline-block w-5 h-5 border-2 border-violet-500 border-t-transparent rounded-full"
          ></span>
        </div>{:else if store.folders.length === 0 && !store.hasDockerfile}<div
          class="text-center py-8 text-xs text-slate-400 italic"
        >
          Nenhuma subpasta encontrada.
        </div>{:else}{#each store.folders as folder: BuilderFolder}<button
            type="button"
            class="w-full flex items-center gap-2.5 px-4 py-2.5 text-xs text-left hover:bg-slate-100 dark:hover:bg-slate-800/60 border-b border-slate-100 dark:border-slate-800/40 last:border-none cursor-pointer transition-colors"
            onclick={() => store.browse(folder.path)}
            ><span class="text-base shrink-0">📂</span><span
              class="font-mono font-semibold text-slate-700 dark:text-slate-300 truncate"
              >{folder.name}</span
            ></button
          >{/each}{/if}
    </div>
    {#if store.hasDockerfile}<div class="flex flex-wrap items-center gap-2">
        <span
          class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-bold bg-emerald-500/15 text-emerald-600 border border-emerald-500/20"
          >✓ {t("builder.detect_dockerfile")}</span
        >{#if store.hasDockerignore}<span
            class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-bold bg-blue-500/15 text-blue-600 border border-blue-500/20"
            >✓ {t("builder.detect_dockerignore")}</span
          >{:else}<span
            class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-medium bg-amber-500/15 text-amber-600 border border-amber-500/20"
            >⚠️ {t("builder.warn_no_dockerignore")}</span
          >{/if}
      </div>{/if}
    {#if store.hasDockerignore && store.ignoredFiles.length > 0}
      <details class="rounded-xl border border-blue-200 dark:border-blue-900/50 bg-blue-50/50 dark:bg-blue-950/20 p-3">
        <summary class="cursor-pointer text-xs font-bold text-blue-700 dark:text-blue-300">Arquivos ignorados ({store.ignoredFiles.length})</summary>
        <ul class="mt-2 max-h-32 overflow-y-auto space-y-1 text-[11px] font-mono text-slate-600 dark:text-slate-400">{#each store.ignoredFiles as file}<li class="break-all">{file}</li>{/each}</ul>
      </details>
    {/if}
  </div>
</div>
