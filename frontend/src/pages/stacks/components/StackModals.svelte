<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import TerminalModal from "$lib/components/TerminalModal.svelte";
  import PullProgressModal from "$lib/components/PullProgressModal.svelte";
  import ConfirmDialog from "$lib/components/ConfirmDialog.svelte";
  import StackEditorModal from "./StackEditorModal.svelte";
  import StackFolderBrowserModal from "./StackFolderBrowserModal.svelte";
  import StackRemoveRemoteModal from "./StackRemoveRemoteModal.svelte";

  let {
    // Editor modal
    showEditor = $bindable(false),
    editorId = "",
    editorName = $bindable(""),
    editorProjectName = $bindable(""),
    editorSourceType = $bindable<"editor" | "folder">("editor"),
    editorFolderPath = $bindable(""),
    editorYaml = $bindable(""),
    onOpenFolderBrowser = () => {},
    onSave = () => {},

    // Folder browser modal
    showFolderModal = $bindable(false),
    browseCurrentPath = "",
    browseParentPath = null,
    browseFolders = [],
    browseHasDockerfile = false,
    browseHasDockerignore = false,
    browseLoading = false,
    onBrowseNavigate = () => {},
    onBrowseSelect = () => {},

    // Remove remote modal
    showRemoveRemoteModal = $bindable(false),
    removeRemoteTargetName = "",
    deleteVolumes = $bindable(false),
    removingRemote = false,
    onConfirmRemoveRemote = () => {},

    // Delete local modal
    showDeleteLocalModal = $bindable(false),
    deleteTargetName = "",
    onConfirmDeleteLocal = () => {},

    // Terminal logs modal
    showLogsTerminal = $bindable(false),
    terminalTitle = "",
    terminalLogs = [],
    terminalLoading = false,

    // Deploy progress modal
    showDeployModal = $bindable(false),
    deployTitle = "",
    onDeployComplete = () => {},

    // Deploy confirm dialog
    showDeployConfirm = $bindable(false),
    onConfirmDeployStack = () => {},
  }: {
    showEditor: boolean;
    editorId: string;
    editorName: string;
    editorProjectName: string;
    editorSourceType: "editor" | "folder";
    editorFolderPath: string;
    editorYaml: string;
    onOpenFolderBrowser: () => void;
    onSave: () => void;

    showFolderModal: boolean;
    browseCurrentPath: string;
    browseParentPath: string | null;
    browseFolders: { name: string; path: string }[];
    browseHasDockerfile: boolean;
    browseHasDockerignore: boolean;
    browseLoading: boolean;
    onBrowseNavigate: (path: string) => void;
    onBrowseSelect: () => void;

    showRemoveRemoteModal: boolean;
    removeRemoteTargetName: string;
    deleteVolumes: boolean;
    removingRemote: boolean;
    onConfirmRemoveRemote: () => void;

    showDeleteLocalModal: boolean;
    deleteTargetName: string;
    onConfirmDeleteLocal: () => void;

    showLogsTerminal: boolean;
    terminalTitle: string;
    terminalLogs: string[];
    terminalLoading: boolean;

    showDeployModal: boolean;
    deployTitle: string;
    onDeployComplete: () => void;

    showDeployConfirm: boolean;
    onConfirmDeployStack: () => void;
  } = $props();
</script>

<!-- Modais da Stack -->
<StackEditorModal
  bind:show={showEditor}
  {editorId}
  bind:editorName
  bind:editorProjectName
  bind:editorSourceType
  bind:editorFolderPath
  bind:editorYaml
  {onOpenFolderBrowser}
  {onSave}
/>

<StackFolderBrowserModal
  bind:show={showFolderModal}
  currentPath={browseCurrentPath}
  parentPath={browseParentPath}
  folders={browseFolders}
  hasDockerfile={browseHasDockerfile}
  hasDockerignore={browseHasDockerignore}
  loading={browseLoading}
  onNavigate={onBrowseNavigate}
  onSelect={onBrowseSelect}
/>

<StackRemoveRemoteModal
  bind:show={showRemoveRemoteModal}
  targetName={removeRemoteTargetName}
  bind:deleteVolumes
  loading={removingRemote}
  onConfirm={onConfirmRemoveRemote}
/>

<Modal
  bind:show={showDeleteLocalModal}
  title={t("stacks.delete_local_btn")}
  buttons={[
    { label: t("common.delete"), variant: "warning", onclick: onConfirmDeleteLocal },
  ]}
>
  <div class="space-y-3">
    <p class="text-sm text-slate-800 dark:text-slate-200 font-semibold m-0">
      {t("stacks.delete_confirm").replace("{name}", deleteTargetName)}
    </p>
    <p class="text-xs text-slate-500 dark:text-slate-400 m-0">
      {t("stacks.delete_local_desc")}
    </p>
  </div>
</Modal>

<TerminalModal
  bind:show={showLogsTerminal}
  title={terminalTitle}
  loading={terminalLoading}
  logs={terminalLogs}
  console_id="stack-logs-terminal"
/>

<PullProgressModal
  bind:show={showDeployModal}
  eventPrefix="stacks:deploy"
  title={deployTitle}
  oncomplete={onDeployComplete}
/>

<ConfirmDialog
  bind:show={showDeployConfirm}
  title="Fazer Deploy de Stack"
  message={t("stacks.confirm_deploy_recreate")}
  confirmText="Confirmar Deploy"
  type="success"
  onConfirm={onConfirmDeployStack}
/>
