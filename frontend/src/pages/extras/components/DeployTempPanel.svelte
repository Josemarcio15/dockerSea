<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { Button, ButtonRed } from "$shared/components/buttons";

  let {
    deployPath = $bindable(),
    deployFiles = [],
    selectedDeployPaths = $bindable(),
    onGoParent,
    onRequestDelete,
    onOpenFolder,
  }: {
    deployPath: string;
    deployFiles: any[];
    selectedDeployPaths: string[];
    onGoParent: () => void;
    onRequestDelete: () => void;
    onOpenFolder: (path: string) => void;
  } = $props();

  function toggleDeployPath(path: string) {
    selectedDeployPaths = selectedDeployPaths.includes(path)
      ? selectedDeployPaths.filter((item) => item !== path)
      : [...selectedDeployPaths, path];
  }
</script>

<section
  class="space-y-4 rounded-2xl border border-slate-200/80 dark:border-slate-800/80 bg-white dark:bg-[#0b0f19] p-6 shadow-sm"
>
  <div class="flex items-center justify-between gap-3">
    <div>
      <h2 class="text-base font-bold">{t("extras.deploy_temp_title")}</h2>
      <p class="text-xs text-slate-500 font-mono break-all">
        {deployPath}
      </p>
    </div>
    <div class="flex gap-2">
      <Button
        size="sm"
        disabled={deployPath === "$HOME/.docksea"}
        onclick={onGoParent}>Voltar</Button
      >
      <ButtonRed
        size="sm"
        disabled={!selectedDeployPaths.length}
        onclick={onRequestDelete}
      >
        Remover selecionados
      </ButtonRed>
    </div>
  </div>

  {#if deployFiles.length === 0}
    <p class="text-sm text-slate-400 py-8 text-center">
      {t("extras.deploy_temp_empty")}
    </p>
  {:else}
    <div class="space-y-2">
      {#each deployFiles as file}
        <div
          class="w-full flex items-center gap-3 rounded-xl border border-slate-200 dark:border-slate-800 p-3 text-xs hover:bg-slate-50 dark:hover:bg-slate-900/50"
        >
          <input
            type="checkbox"
            checked={selectedDeployPaths.includes(file.path)}
            onchange={() => toggleDeployPath(file.path)}
          />
          <button
            type="button"
            class="flex flex-1 justify-between gap-3 text-left"
            onclick={() => file.isDir && onOpenFolder(file.path)}
          >
            <span class="font-mono break-all">
              {file.isDir ? "📁 " : "📄 "}{file.path.split("/").pop()}
            </span>
            <span class="shrink-0 text-slate-500">
              {file.isDir ? "Pasta" : `${file.size} bytes`}
            </span>
          </button>
        </div>
      {/each}
    </div>
  {/if}
</section>
