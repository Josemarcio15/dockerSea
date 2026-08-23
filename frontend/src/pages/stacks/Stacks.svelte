<script lang="ts">
  import StatusBanner from "$lib/components/StatusBanner.svelte";
  import { notifySuccess, notifyError } from "$lib/stores/notification.svelte";
  import StackRow from "$lib/components/StackRow.svelte";
  import PullProgressModal from "$lib/components/PullProgressModal.svelte";
  import TerminalModal from "$lib/components/TerminalModal.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import CodeEditor from "$lib/components/CodeEditor.svelte";
  import VpsSelectWarning from "$lib/components/VpsSelectWarning.svelte";
  import { t, getLocale } from "$lib/stores/locale.svelte";

  let { data } = $props();

  let searchQuery = $state("");

  // Editor state
  let showEditor = $state(false);
  let editorId = $state("");
  let editorName = $state("");
  let editorProjectName = $state("");
  let editorYaml = $state("");

  // Logs Terminal state
  let showLogsTerminal = $state(false);
  let terminalTitle = $state("");
  let terminalLogs = $state<string[]>([]);
  let terminalLoading = $state(false);

  // Deploy modal state (SSE)
  let showDeployModal = $state(false);
  let deploySseUrl = $state("");
  let deployTitle = $state("");

  // Delete confirmation modal
  let showDeleteModal = $state(false);
  let deleteTargetId = $state("");
  let deleteTargetName = $state("");

  // Derived: Filtered stacks
  let filteredStacks = $derived(
    (data.stacks || []).filter((s) =>
      s.name.toLowerCase().includes(searchQuery.toLowerCase()),
    ),
  );

  const defaultYamlTemplate = `name: my-project
services:
  web:
    image: nginx:alpine
    ports:
      - "8080:80"
    restart: always
`;

  // Submit Action helper
  async function runAction(
    action: "save" | "delete" | "stop" | "logs",
    formData: FormData,
  ) {
    try {
      const res = await fetch(`?/${action}`, {
        method: "POST",
        body: formData,
      });
      const result = deserialize(await res.text()) as any;

      let parsedData: any = null;
      if (result && typeof result === "object") {
        if (result.type === "success" && result.data) {
          parsedData = result.data;
        } else if (result.data) {
          parsedData = result.data;
        } else {
          parsedData = result;
        }
      }

      if (parsedData && parsedData.success) {
        notifySuccess(parsedData.message);

        if (action === "logs" || parsedData.logs) {
          terminalTitle = t("stacks.operation_logs_title");
          terminalLogs = (parsedData.logs || parsedData.message || "").split(
            "\n",
          );
          showLogsTerminal = true;
        }
      } else {
        notifyError(
          parsedData?.message ||
            parsedData?.error ||
            result?.error?.message ||
            result?.error ||
            result?.message ||
            (result ? JSON.stringify(result).slice(0, 300) : "Resposta vazia"),
        );
      }

      await invalidateAll();
    } catch (e: any) {
      notifyError(e.message);
    }
  }

  async function doSave() {
    const projectName = editorProjectName.trim();
    if (!editorName.trim() || !projectName || !editorYaml.trim()) return;
    notifySuccess(t("stacks.saving_status").replace("{name}", editorName));

    // Inject/replace `name:` at the top of the YAML
    let yaml = editorYaml.trimStart();
    const nameLine = `name: ${projectName}`;
    if (/^name:\s*/.test(yaml)) {
      yaml = yaml.replace(/^name:.*(\r?\n)?/, `${nameLine}\n`);
    } else {
      yaml = `${nameLine}\n${yaml}`;
    }

    const formData = new FormData();
    if (editorId) formData.append("id", editorId);
    formData.append("name", editorName);
    formData.append("yamlContent", yaml);

    showEditor = false;
    await runAction("save", formData);
  }

  async function doDelete(id: string, name: string) {
    deleteTargetId = id;
    deleteTargetName = name;
    showDeleteModal = true;
  }

  async function confirmDelete() {
    showDeleteModal = false;
    notifySuccess(
      t("stacks.deleting_status").replace("{name}", deleteTargetName),
    );

    const formData = new FormData();
    formData.append("id", deleteTargetId);

    await runAction("delete", formData);
  }

  async function doDeploy(id: string, name: string) {
    if (!confirm(t("stacks.confirm_deploy_recreate"))) return;
    deploySseUrl = `/api/stacks/deploy-stream?id=${id}&locale=${getLocale()}`;
    deployTitle = `${t("stacks.deploy_btn")}: ${name}`;
    showDeployModal = true;
  }

  function onDeployComplete() {
    invalidateAll();
  }

  async function doStop(id: string, name: string) {
    notifySuccess(t("stacks.stopping_status").replace("{name}", name));

    const formData = new FormData();
    formData.append("id", id);

    await runAction("stop", formData);
  }

  async function doLogs(id: string, name: string) {
    notifySuccess(t("stacks.fetching_logs_status").replace("{name}", name));

    const formData = new FormData();
    formData.append("id", id);

    terminalLoading = true;
    showLogsTerminal = true;

    // Execute action
    try {
      const res = await fetch(`?/logs`, {
        method: "POST",
        body: formData,
      });
      const result = deserialize(await res.text()) as any;
      let parsedData: any = result?.data || result;
      if (parsedData?.success) {
        terminalTitle = t("stacks.stack_logs_title").replace("{name}", name);
        terminalLogs = (parsedData.logs || t("stacks.no_logs_output")).split(
          "\n",
        );
        notifySuccess(t("stacks.logs_loaded_status").replace("{name}", name));
      } else {
        notifyError(parsedData?.message || t("common.error"));
        showLogsTerminal = false;
      }
    } catch (e: any) {
      notifyError(e.message);
      showLogsTerminal = false;
    } finally {
      terminalLoading = false;
    }
  }

  function openCreateEditor() {
    editorId = "";
    editorName = "";
    editorProjectName = "";
    editorYaml = defaultYamlTemplate;
    showEditor = true;
  }

  function openEditEditor(stack: any) {
    editorId = stack.id;
    editorName = stack.name;
    const match = (stack.yamlContent || "").match(/^name:\s*(.+)/m);
    editorProjectName = match ? match[1].trim() : "";
    editorYaml = stack.yamlContent;
    showEditor = true;
  }
</script>

{#if !data.activeVps}
  <VpsSelectWarning />
{:else}
  <div class="space-y-6">
    <!-- Top Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
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
          class="px-4 py-2 text-xs rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-[#0b0f19] text-slate-850 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all w-60"
          bind:value={searchQuery}
        />
        <button
          type="button"
          class="px-4 py-2 text-xs rounded-xl border-none cursor-pointer font-bold text-white bg-blue-600 hover:bg-blue-700 transition-colors shadow-md shadow-blue-500/20 whitespace-nowrap"
          onclick={openCreateEditor}
        >
          {t("stacks.new_stack")}
        </button>
      </div>
    </div>

    <!-- Status Alerts -->
    <StatusBanner />

    <!-- Stacks List -->
    {#if filteredStacks.length === 0}
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
            on_deploy={() => doDeploy(stack.id, stack.name)}
            on_stop={() => doStop(stack.id, stack.name)}
            on_logs={() => doLogs(stack.id, stack.name)}
            on_edit={() => openEditEditor(stack)}
            on_delete={() => doDelete(stack.id, stack.name)}
          />
        {/each}
      </div>
    {/if}
  </div>
{/if}

<!-- Modal: Editor de YAML -->
{#if showEditor}
  <div
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
  >
    <div
      class="bg-white dark:bg-[#0b0f19] border border-slate-200 dark:border-slate-800 rounded-2xl w-175 max-w-full h-[85vh] flex flex-col p-6 shadow-2xl animate-scaleIn text-slate-800 dark:text-slate-200"
    >
      <div
        class="flex justify-between items-center pb-4 border-b border-slate-200 dark:border-slate-800"
      >
        <h2 class="text-lg font-bold text-slate-900 dark:text-white">
          {editorId
            ? t("stacks.modal_edit_title")
            : t("stacks.modal_new_title")}
        </h2>
        <button
          type="button"
          class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 text-xl cursor-pointer bg-transparent border-none"
          onclick={() => (showEditor = false)}
        >
          ✕
        </button>
      </div>

      <div class="py-4 grow flex flex-col min-h-0 space-y-4">
        <div class="flex flex-col gap-1.5 shrink-0">
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
            disabled={!!editorId}
          />
        </div>

        <div class="flex flex-col gap-1.5 shrink-0">
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

        <div class="flex flex-col gap-1.5 grow min-h-0">
          <label
            for="stack-yaml"
            class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
          >
            {t("stacks.field_yaml")}
          </label>
          {#key showEditor}
            <CodeEditor
              mode="yaml"
              value={editorYaml}
              minHeight="280px"
              maxHeight="420px"
              onchange={(v) => (editorYaml = v)}
            />
          {/key}
        </div>
      </div>

      <div
        class="flex gap-3 justify-end pt-4 border-t border-slate-200 dark:border-slate-800 shrink-0"
      >
        <button
          type="button"
          class="px-4 py-2.5 rounded-xl border-none cursor-pointer text-xs font-bold text-white bg-red-500 hover:bg-red-600 transition-colors shadow-md shadow-red-500/20"
          onclick={() => (showEditor = false)}
        >
          {t("common.cancel")}
        </button>
        <button
          type="button"
          class="px-4 py-2.5 rounded-xl border-none cursor-pointer text-xs font-bold text-white transition-colors shadow-md shadow-emerald-500/20 {editorName.trim() &&
          editorProjectName.trim() &&
          editorYaml.trim()
            ? 'bg-emerald-600 hover:bg-emerald-700'
            : 'bg-slate-400 dark:bg-slate-800 cursor-not-allowed text-slate-400'}"
          disabled={!editorName.trim() ||
            !editorProjectName.trim() ||
            !editorYaml.trim()}
          onclick={doSave}
        >
          {t("common.save")}
        </button>
      </div>
    </div>
  </div>
{/if}

<TerminalModal
  bind:show={showLogsTerminal}
  title={terminalTitle}
  loading={terminalLoading}
  logs={terminalLogs}
  console_id="stack-logs-terminal"
/>

<Modal
  bind:show={showDeleteModal}
  title={t("stacks.delete_confirm").replace("{name}", deleteTargetName)}
  buttons={[
    { label: t("common.delete"), variant: "warning", onclick: confirmDelete },
  ]}
>
  <p class="text-sm text-slate-600 dark:text-slate-400">
    {t("stacks.delete_confirm").replace("{name}", deleteTargetName)}
  </p>
</Modal>

<PullProgressModal
  bind:show={showDeployModal}
  sseUrl={deploySseUrl}
  title={deployTitle}
  oncomplete={onDeployComplete}
/>
