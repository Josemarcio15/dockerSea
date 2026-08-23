<script lang="ts">
  import { onMount } from "svelte";
  import { t, getLocale } from "$lib/stores/locale.svelte";
  import StatusBanner from "$lib/components/StatusBanner.svelte";
  import { notifySuccess, notifyError } from "$lib/stores/notification.svelte";

  import VpsCard from "./VpsCard.svelte";
  import VpsModal, { type VpsFormData } from "./VpsModal.svelte";
  import DiagnosticModal from "./DiagnosticModal.svelte";
  import { BrandButton } from "$lib/components/buttons";
  import * as ConfigService from "../../../bindings/go-walis/internal/config/configservice.js";

  let { data } = $props();

  // Initial Form State
  const defaultFormData: VpsFormData = {
    id: "",
    name: "",
    connectionType: "local",
    host: "",
    port: "22",
    username: "root",
    authType: "key",
    sshKeyPath: "",
    sshKeyPassphrase: "",
    sshPassword: "",
    sudoPassword: "",
    dockerSocketPath: "/var/run/docker.sock",
    dockerPath: "/usr/bin/docker",
    dockerComposePath: "",
  };

  let showVpsModal = $state(false);
  let formData = $state<VpsFormData>({ ...defaultFormData });

  // Carregar servidores do SQLite ao abrir a tela
  async function loadServers() {
    try {
      const servers = await ConfigService.ListServers();
      if (servers && servers.length > 0) {
        data.servers = servers;
        const active = servers.find((s: any) => s.isActive);
        if (active) {
          data.activeVps = active;
        }
      }
    } catch (e: any) {
      console.warn("Wails runtime ainda não disponível no modo dev/browser puro:", e);
    }
  }

  onMount(() => {
    loadServers();
  });

  // Diagnostic state
  let showDiagnosticModal = $state(false);
  let diagnosticTitle = $state("");
  let diagnosticLoading = $state(false);
  let diagnosticResult = $state<{
    success: boolean;
    message: string;
    steps: Array<{
      name: string;
      status: "success" | "error" | "warning";
      message: string;
    }>;
  } | null>(null);

  // Derived locale
  let activeLocale = $derived(data.activeProfile?.locale || getLocale());

  function openCreateModal() {
    formData = { ...defaultFormData };
    showVpsModal = true;
  }

  function openEditModal(vps: any) {
    formData = {
      id: vps.id,
      name: vps.name,
      connectionType: vps.connectionType,
      host: vps.host || "",
      port: String(vps.port || "22"),
      username: vps.username || "root",
      authType: vps.authType || "key",
      sshKeyPath: vps.sshKeyPath || "",
      sshKeyPassphrase: vps.sshKeyPassphrase || "",
      sshPassword: vps.sshPassword || "",
      sudoPassword: vps.sudoPassword || "",
      dockerSocketPath: vps.dockerSocketPath || "/var/run/docker.sock",
      dockerPath: vps.dockerPath || "/usr/bin/docker",
      dockerComposePath: vps.dockerComposePath || "",
    };
    showVpsModal = true;
  }

  async function runDiagnostic(vpsData?: VpsFormData) {
    diagnosticLoading = true;
    showDiagnosticModal = true;
    diagnosticResult = null;

    const target = vpsData || formData;
    diagnosticTitle = `🔍 Diagnóstico de Conexão: ${target.name || target.host || "Servidor"}`;

    try {
      const serverPayload: any = {
        id: target.id,
        name: target.name,
        connectionType: target.connectionType,
        host: target.host,
        port: parseInt(target.port) || 22,
        username: target.username,
        authType: target.authType,
        sshKeyPath: target.sshKeyPath,
        sshKeyPassphrase: target.sshKeyPassphrase,
        sshPassword: target.sshPassword,
        sudoPassword: target.sudoPassword,
        dockerSocketPath: target.dockerSocketPath,
        dockerPath: target.dockerPath,
        dockerComposePath: target.dockerComposePath,
        isActive: false,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      const result = await ConfigService.TestConnection(serverPayload);
      diagnosticResult = result as any;
    } catch (e: any) {
      diagnosticResult = {
        success: false,
        message: `Falha ao executar teste de conexão: ${e.message || e}`,
        steps: [
          {
            name: "Execução do Diagnóstico",
            status: "error",
            message: e.message || String(e),
          },
        ],
      };
    } finally {
      diagnosticLoading = false;
    }
  }

  async function handleSaveVps(savedForm: VpsFormData) {
    if (!savedForm.name.trim()) return;

    try {
      const serverPayload: any = {
        id: savedForm.id,
        name: savedForm.name,
        connectionType: savedForm.connectionType,
        host: savedForm.host,
        port: parseInt(savedForm.port) || 22,
        username: savedForm.username,
        authType: savedForm.authType,
        sshKeyPath: savedForm.sshKeyPath,
        sshKeyPassphrase: savedForm.sshKeyPassphrase,
        sshPassword: savedForm.sshPassword,
        sudoPassword: savedForm.sudoPassword,
        dockerSocketPath: savedForm.dockerSocketPath,
        dockerPath: savedForm.dockerPath,
        dockerComposePath: savedForm.dockerComposePath,
        isActive: false,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      await ConfigService.SaveServer(serverPayload);
      await loadServers();
      notifySuccess("Servidor salvo no SQLite com sucesso!");
    } catch (e: any) {
      notifyError(`Erro ao salvar no banco: ${e.message || e}`);
    }

    showVpsModal = false;
  }

  async function handleDeleteVps(server: any) {
    if (!confirm(`Tem certeza de que deseja remover a VPS '${server.name}'?`)) return;

    try {
      await ConfigService.DeleteServer(server.id);
      await loadServers();
      notifySuccess(`Servidor '${server.name}' removido.`);
    } catch (e: any) {
      notifyError(`Erro ao remover: ${e.message || e}`);
    }
  }

  async function handleSelectVps(server: any) {
    try {
      await ConfigService.SetActiveServer(server.id);
      await loadServers();
      notifySuccess(`Servidor '${server.name}' ativado com sucesso!`);
    } catch (e: any) {
      data.activeVps = server;
      notifySuccess(`Servidor '${server.name}' selecionado como ativo.`);
    }
  }

  function doChangeLocale(locale: string) {
    if (data.activeProfile) {
      data.activeProfile.locale = locale;
    }
    notifySuccess(`Idioma alterado para ${locale === "pt-BR" ? "Português" : "English"}`);
  }
</script>

<div class="space-y-8">
  <!-- Top Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div
      class="inline-flex items-center px-5 py-2.5 rounded-2xl bg-linear-to-br from-violet-100 to-fuchsia-100 dark:from-violet-950/40 dark:to-fuchsia-950/40 border border-violet-200/50 dark:border-violet-800/50 self-start shadow-sm"
    >
      <h1
        class="text-2xl font-bold bg-linear-to-r from-violet-600 to-fuchsia-500 bg-clip-text text-transparent m-0 flex items-center gap-2"
      >
        {t("sidebar.configs")}
      </h1>
    </div>
  </div>

  <!-- Status Alerts -->
  <StatusBanner />

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 items-start">
    <!-- VPS Servers List (col-span-2) -->
    <div
      class="lg:col-span-2 bg-white dark:bg-[#0b0f19] border border-slate-200/70 dark:border-slate-800/80 p-6 rounded-2xl shadow-sm space-y-6"
    >
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h2
            class="text-lg font-bold text-slate-850 dark:text-white flex items-center gap-2"
          >
            {t("config.servers_title")}
          </h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
            {t("config.servers_desc")}
          </p>
        </div>

        <BrandButton onclick={openCreateModal}>
          <span class="text-sm font-normal">+</span>
          {t("config.add_server_btn")}
        </BrandButton>
      </div>

      <!-- Servers List (Rows) -->
      <div class="flex flex-col gap-3">
        {#each data.servers as server (server.id)}
          <VpsCard
            {server}
            isActive={data.activeVps?.id === server.id}
            onSelect={handleSelectVps}
            onTest={runDiagnostic}
            onEdit={openEditModal}
            onDelete={handleDeleteVps}
          />
        {:else}
          <div
            class="text-center py-12 border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-slate-50/30 italic text-slate-500"
          >
            {t("config.empty_servers")}
          </div>
        {/each}
      </div>
    </div>

    <!-- General Settings (col-span-1) -->
    <div class="space-y-6">
      <!-- Language Card -->
      <div
        class="bg-white dark:bg-[#0b0f19] border border-slate-200/70 dark:border-slate-800/80 p-5 rounded-2xl shadow-sm space-y-4"
      >
        <h3
          class="text-sm font-bold text-slate-800 dark:text-slate-200 uppercase tracking-wider flex items-center gap-2"
        >
          {t("config.lang_tab")}
        </h3>

        <div class="flex flex-col gap-2">
          <button
            type="button"
            class="w-full flex items-center justify-between p-3 rounded-xl border text-sm font-semibold transition-all cursor-pointer {activeLocale ===
            'pt-BR'
              ? 'border-violet-500 bg-violet-500/5 text-violet-750 dark:text-violet-300'
              : 'border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-900/30 text-slate-700 dark:text-slate-300'}"
            onclick={() => doChangeLocale("pt-BR")}
          >
            <span>🇧🇷 Português (Brasil)</span>
            {#if activeLocale === "pt-BR"}✓{/if}
          </button>

          <button
            type="button"
            class="w-full flex items-center justify-between p-3 rounded-xl border text-sm font-semibold transition-all cursor-pointer {activeLocale ===
            'en-US'
              ? 'border-violet-500 bg-violet-500/5 text-violet-750 dark:text-violet-300'
              : 'border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-900/30 text-slate-700 dark:text-slate-300'}"
            onclick={() => doChangeLocale("en-US")}
          >
            <span>🇺🇸 English (US)</span>
            {#if activeLocale === "en-US"}✓{/if}
          </button>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- Modal: Criar / Editar VPS -->
<VpsModal
  bind:show={showVpsModal}
  bind:form={formData}
  onSave={handleSaveVps}
  onTest={runDiagnostic}
/>

<!-- Modal: Diagnóstico de Conexão VPS -->
<DiagnosticModal
  bind:show={showDiagnosticModal}
  title={diagnosticTitle}
  loading={diagnosticLoading}
  result={diagnosticResult}
/>
