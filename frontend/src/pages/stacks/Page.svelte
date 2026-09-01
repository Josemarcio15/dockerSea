<script lang="ts">
  import StackRow from "$shared/components/StackRow.svelte";
  import VpsSelectWarning from "$shared/components/VpsSelectWarning.svelte";
  import { t } from "$shared/stores/locale.svelte";
  import { ButtonBlue } from "$shared/components/buttons";
  import type { StackItem } from "./types";
  import { defaultYaml, folderName } from "./service";
  import StackModals from "./components/StackModals.svelte";
  import { createStacksStore } from "./store.svelte";

  let { data } = $props();

  const stackStore = createStacksStore(() => data);
  const stacksList = $derived(stackStore.stacks);
  const loadingStacks = $derived(stackStore.loading);
  const capabilities = $derived(stackStore.capabilities);

  // Editor & Folder state
  let showEditor = $state(false);
  let editorId = $state("");
  let editorName = $state("");
  let editorProjectName = $state("");
  let editorSourceType = $state<"editor" | "folder">("editor");
  let editorFolderPath = $state("");
  let editorYaml = $state("");

  let showFolderModal = $state(false);
  let browseCurrentPath = $state("");
  let browseParentPath = $state<string | null>(null);
  let browseFolders = $state<{ name: string; path: string }[]>([]);
  let browseHasDockerfile = $state(false);
  let browseHasDockerignore = $state(false);
  let browseLoading = $state(false);

  // Logs & Deploy state
  let showLogsTerminal = $state(false);
  let terminalTitle = $state("");
  let terminalLogs = $state<string[]>([]);
  let terminalLoading = $state(false);

  let showDeployModal = $state(false);
  let deployTitle = $state("");
  let showDeployConfirm = $state(false);
  let stackToDeploy = $state<{ id: string; name: string } | null>(null);

  // Delete & Remove state
  let showDeleteLocalModal = $state(false);
  let deleteTargetId = $state("");
  let deleteTargetName = $state("");

  let showRemoveRemoteModal = $state(false);
  let removeRemoteTargetId = $state("");
  let removeRemoteTargetName = $state("");
  let deleteVolumes = $state(false);
  let removingRemote = $state(false);

  let filteredStacks = $derived(stackStore.filteredStacks);

  async function browseFolder(targetPath: string = "") {
    browseLoading = true;
    const l = await stackStore.browseFolder(targetPath);
    if (l) {
      browseCurrentPath = l.currentPath;
      browseParentPath = l.parentPath || null;
      browseFolders = l.folders || [];
      browseHasDockerfile = l.hasDockerfile;
      browseHasDockerignore = l.hasDockerignore;
    }
    browseLoading = false;
  }

  function selectBrowseFolder() {
    editorFolderPath = browseCurrentPath;
    const name = folderName(browseCurrentPath);
    if (!editorProjectName)
      editorProjectName = name.toLowerCase().replace(/[^a-z0-9_-]/g, "-");
    if (!editorName) editorName = name;
    showFolderModal = false;
  }

  async function doSave() {
    const proj = editorProjectName.trim();
    if (!editorName.trim() || !proj) return;
    if (editorSourceType === "editor" && !editorYaml.trim()) return;
    if (editorSourceType === "folder" && !editorFolderPath.trim()) return;

    if (
      await stackStore.save({
        id: editorId,
        name: editorName.trim(),
        projectName: proj,
        sourceType: editorSourceType,
        folderPath: editorFolderPath.trim(),
        yamlContent: editorYaml.trim(),
      })
    )
      showEditor = false;
  }

  async function confirmDeleteLocal() {
    showDeleteLocalModal = false;
    await stackStore.deleteLocal(deleteTargetId);
  }

  async function confirmRemoveRemote() {
    if (!removeRemoteTargetId) return;
    removingRemote = true;
    try {
      await stackStore.removeRemote(removeRemoteTargetId, deleteVolumes);
      showRemoveRemoteModal = false;
    } finally {
      removingRemote = false;
    }
  }

  async function confirmDeploy() {
    if (!stackToDeploy) return;
    const target = { ...stackToDeploy };
    showDeployConfirm = false;
    deployTitle = `${t("stacks.deploy_btn")}: ${target.name}`;
    showDeployModal = true;

    // Dispara a chamada Wails em background sem bloquear o fechamento do diálogo
    void stackStore.deploy(target.id).catch(() => undefined);
  }

  async function doStop(id: string, name: string) {
    await stackStore.stop(id);
  }

  async function doLogs(id: string, name: string) {
    terminalLoading = true;
    showLogsTerminal = true;
    try {
      const logs = await stackStore.logs(id);
      terminalTitle = t("stacks.stack_logs_title").replace("{name}", name);
      terminalLogs = (logs || t("stacks.no_logs_output")).split("\n");
    } catch {
      showLogsTerminal = false;
    } finally {
      terminalLoading = false;
    }
  }

  function openEdit(stack: StackItem) {
    editorId = stack.id;
    editorName = stack.name;
    editorProjectName = stack.projectName;
    editorSourceType = (stack.sourceType as "editor" | "folder") || "editor";
    editorFolderPath = stack.folderPath || "";
    editorYaml = stack.yamlContent || "";
    showEditor = true;
  }

  function openCreate() {
    editorId = "";
    editorName = "";
    editorProjectName = "";
    editorSourceType = "editor";
    editorFolderPath = "";
    editorYaml = defaultYaml;
    showEditor = true;
  }

  $effect(() => {
    void stackStore.load();
  });
</script>

{#if !data.activeVps}
  <VpsSelectWarning />
{:else}
  <div class="space-y-6">
    <div
      class="flex flex-col sm:flex-row sm:items-center justify-between gap-4"
    >
      <div
        class="inline-flex items-center px-5 py-2.5 rounded-2xl bg-linear-to-br from-violet-100 to-fuchsia-100 dark:from-violet-950/40 dark:to-fuchsia-950/40 border border-violet-200/50 dark:border-violet-800/50 self-start shadow-sm"
      >
        <h1
          class="text-2xl font-bold bg-linear-to-r from-violet-600 to-fuchsia-500 bg-clip-text text-transparent m-0 flex items-center gap-2"
        >
          {t("stacks.title")}
        </h1>
      </div>

      <div class="flex items-center gap-2">
        <input
          type="text"
          placeholder={t("stacks.search_placeholder")}
          class="px-4 py-2 text-xs rounded-xl border border-slate-300 dark:border-slate-800 bg-white dark:bg-[#0b0f19] text-slate-800 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 shadow-2xs transition-all w-60"
          bind:value={stackStore.searchQuery}
        />
        <ButtonBlue size="sm" onclick={openCreate}
          >+ {t("stacks.new_stack")}</ButtonBlue
        >
      </div>
    </div>

    {#if capabilities && !capabilities.composeAvailable}
      <div
        class="p-4 rounded-2xl border border-amber-500/30 bg-amber-500/10 text-amber-600 dark:text-amber-400 text-xs flex items-center gap-3"
      >
        <p class="font-bold m-0 flex-1">
          {t("stacks.compose_unavailable_warning")}
        </p>
      </div>
    {/if}

    {#if loadingStacks && stacksList.length === 0}
      <div class="py-16 text-center text-xs text-slate-400">
        {t("common.loading")}
      </div>
    {:else if filteredStacks.length === 0}
      <div
        class="text-sm text-slate-400 dark:text-slate-500 py-16 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-slate-50/20"
      >
        {t("stacks.empty_msg")}
      </div>
    {:else}
      <div class="grid grid-cols-1 gap-4">
        {#each filteredStacks as stack (stack.id)}
          <StackRow
            {stack}
            on_deploy={() => {
              stackToDeploy = { id: stack.id, name: stack.name };
              showDeployConfirm = true;
            }}
            on_stop={() => doStop(stack.id, stack.name)}
            on_logs={() => doLogs(stack.id, stack.name)}
            on_edit={() => openEdit(stack)}
            on_delete_local={() => {
              deleteTargetId = stack.id;
              deleteTargetName = stack.name;
              showDeleteLocalModal = true;
            }}
            on_remove_remote={() => {
              removeRemoteTargetId = stack.id;
              removeRemoteTargetName = stack.name;
              deleteVolumes = false;
              showRemoveRemoteModal = true;
            }}
          />
        {/each}
      </div>
    {/if}
  </div>
{/if}

<StackModals
  bind:showEditor
  {editorId}
  bind:editorName
  bind:editorProjectName
  bind:editorSourceType
  bind:editorFolderPath
  bind:editorYaml
  onOpenFolderBrowser={() => {
    browseFolder(editorFolderPath || "");
    showFolderModal = true;
  }}
  onSave={doSave}
  bind:showFolderModal
  {browseCurrentPath}
  {browseParentPath}
  {browseFolders}
  {browseHasDockerfile}
  {browseHasDockerignore}
  {browseLoading}
  onBrowseNavigate={browseFolder}
  onBrowseSelect={selectBrowseFolder}
  bind:showRemoveRemoteModal
  {removeRemoteTargetName}
  bind:deleteVolumes
  {removingRemote}
  onConfirmRemoveRemote={confirmRemoveRemote}
  bind:showDeleteLocalModal
  {deleteTargetName}
  onConfirmDeleteLocal={confirmDeleteLocal}
  bind:showLogsTerminal
  {terminalTitle}
  {terminalLogs}
  {terminalLoading}
  bind:showDeployModal
  {deployTitle}
  onDeployComplete={() => {
    stackStore.load();
  }}
  bind:showDeployConfirm
  onConfirmDeployStack={confirmDeploy}
/>
