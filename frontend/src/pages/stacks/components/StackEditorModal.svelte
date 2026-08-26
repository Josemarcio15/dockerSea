<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { ButtonOrange, ButtonGreen, ButtonBlue } from "$shared/components/buttons";
  import CodeEditor from "$shared/components/CodeEditor.svelte";

  let {
    show = $bindable(false),
    editorId = "",
    editorName = $bindable(""),
    editorProjectName = $bindable(""),
    editorSourceType = $bindable<"editor" | "folder">("editor"),
    editorFolderPath = $bindable(""),
    editorYaml = $bindable(""),
    onOpenFolderBrowser = () => {},
    onSave = () => {},
  }: {
    show: boolean;
    editorId: string;
    editorName: string;
    editorProjectName: string;
    editorSourceType: "editor" | "folder";
    editorFolderPath: string;
    editorYaml: string;
    onOpenFolderBrowser: () => void;
    onSave: () => void;
  } = $props();
</script>

{#if show}
  <div
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
  >
    <div
      class="bg-white dark:bg-[#0b0f19] border border-slate-200 dark:border-slate-800 rounded-2xl w-190 max-w-full h-[88vh] flex flex-col p-6 shadow-2xl animate-scaleIn text-slate-800 dark:text-slate-200"
    >
      <div
        class="flex justify-between items-center pb-4 border-b border-slate-200 dark:border-slate-800"
      >
        <div class="flex items-center gap-3">
          <h2 class="text-lg font-bold text-slate-900 dark:text-white">
            {editorId
              ? t("stacks.modal_edit_title")
              : t("stacks.modal_new_title")}
          </h2>
        </div>
        <button
          type="button"
          class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 text-xl cursor-pointer bg-transparent border-none"
          onclick={() => (show = false)}
        >
          ✕
        </button>
      </div>

      <div class="py-4 grow flex flex-col min-h-0 space-y-4 overflow-y-auto pr-1">
        <!-- Alternador de Tipo de Origem (Tabs) -->
        <div class="flex items-center gap-2 p-1 bg-slate-100 dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800">
          <button
            type="button"
            class="flex-1 py-2 text-xs font-bold rounded-lg transition-all flex items-center justify-center gap-2 cursor-pointer border-none {editorSourceType === 'editor'
              ? 'bg-white dark:bg-slate-800 text-violet-600 dark:text-violet-400 shadow-sm'
              : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300 bg-transparent'}"
            onclick={() => (editorSourceType = "editor")}
          >
            {t("stacks.source_type_editor")}
          </button>
          <button
            type="button"
            class="flex-1 py-2 text-xs font-bold rounded-lg transition-all flex items-center justify-center gap-2 cursor-pointer border-none {editorSourceType === 'folder'
              ? 'bg-white dark:bg-slate-800 text-violet-600 dark:text-violet-400 shadow-sm'
              : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300 bg-transparent'}"
            onclick={() => (editorSourceType = "folder")}
          >
            {t("stacks.source_type_folder")} (Dockerfile + Compose)
          </button>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 shrink-0">
          <div class="flex flex-col gap-1.5">
            <label
              for="stack-name"
              class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
            >
              {t("stacks.field_name")}
            </label>
            <input
              id="stack-name"
              type="text"
              class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all font-semibold"
              placeholder={t("stacks.placeholder_name_wp")}
              bind:value={editorName}
            />
          </div>

          <div class="flex flex-col gap-1.5">
            <label
              for="stack-project"
              class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
            >
              {t("stacks.field_project")}
            </label>
            <input
              id="stack-project"
              type="text"
              class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all font-semibold"
              placeholder={t("stacks.placeholder_project_ex")}
              bind:value={editorProjectName}
            />
          </div>
        </div>

        {#if editorSourceType === "folder"}
          <!-- Modo Pasta Local -->
          <div class="flex flex-col gap-3 p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/20">
            <label
              for="stack-folder"
              class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
            >
              {t("stacks.field_folder_path")}
            </label>
            <div class="flex gap-2 items-center">
              <input
                id="stack-folder"
                type="text"
                class="flex-1 px-3.5 py-2.5 text-xs font-mono border border-slate-200 dark:border-slate-800 rounded-xl bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500"
                placeholder={t("stacks.placeholder_folder_path")}
                bind:value={editorFolderPath}
              />
              <ButtonBlue size="sm" onclick={onOpenFolderBrowser}>
                {t("stacks.browse_folder_btn")}
              </ButtonBlue>
            </div>

            <div class="text-[11px] text-slate-500 dark:text-slate-400 flex items-center gap-2 mt-1">
              <span>Dica: Todos os arquivos da pasta (Dockerfile, .dockerignore, docker-compose.yml e código-fonte) serão enviados para a VPS para build nativo.</span>
            </div>
          </div>
        {:else}
          <!-- Modo Editor YAML -->
          <div class="flex flex-col gap-1.5 grow min-h-0">
            <label
              for="stack-yaml"
              class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
            >
              {t("stacks.field_yaml")}
            </label>
            {#key show}
              <CodeEditor
                mode="yaml"
                value={editorYaml}
                minHeight="280px"
                maxHeight="420px"
                onchange={(v) => (editorYaml = v)}
              />
            {/key}
          </div>
        {/if}
      </div>

      <div
        class="flex gap-3 justify-end pt-4 border-t border-slate-200 dark:border-slate-800 shrink-0"
      >
        <ButtonOrange onclick={() => (show = false)}>
          {t("common.cancel")}
        </ButtonOrange>
        <ButtonGreen
          disabled={!editorName.trim() ||
            !editorProjectName.trim() ||
            (editorSourceType === "editor" && !editorYaml.trim()) ||
            (editorSourceType === "folder" && !editorFolderPath.trim())}
          onclick={onSave}
        >
          {t("common.save")}
        </ButtonGreen>
      </div>
    </div>
  </div>
{/if}
