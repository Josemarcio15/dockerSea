import { getLocale, setLocale, t } from "$shared/stores/locale.svelte";
import { notifySuccess, notifyError } from "$shared/stores/notification.svelte";
import * as api from "./api";
import { dockerDetectionPayload, serverPayload } from "./service";
import type { DiagnosticResult, VpsFormData, VpsServer } from "./types";

export function createConfigStore(getData: () => any) {
  const defaultForm: VpsFormData = {
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
  let form = $state<VpsFormData>({ ...defaultForm });
  let showVpsModal = $state(false);
  let showDiagnosticModal = $state(false);
  let diagnosticTitle = $state("");
  let diagnosticLoading = $state(false);
  let diagnosticResult = $state<DiagnosticResult | null>(null);
  let showDeleteConfirm = $state(false);
  let serverToDelete = $state<VpsServer | null>(null);

  async function load() {
    try {
      const servers = await api.listServers();
      if (servers?.length) {
        const data = getData();
        data.servers = servers;
        const active = servers.find((server) => server.isActive);
        if (active) data.activeVps = active;
      }
    } catch (error) {
      console.warn("Wails runtime ainda não disponível:", error);
    }
  }
  function openCreate() {
    form = { ...defaultForm };
    showVpsModal = true;
  }
  function openEdit(server: VpsServer) {
    form = { ...defaultForm, ...server, port: String(server.port || "22") };
    showVpsModal = true;
  }
  async function test(value = form) {
    diagnosticLoading = true;
    showDiagnosticModal = true;
    diagnosticResult = null;
    diagnosticTitle = `Diagnóstico de Conexão: ${value.name || value.host || "Servidor"}`;
    try {
      diagnosticResult = (await api.testConnection(
        serverPayload(value),
      )) as DiagnosticResult;
    } catch (error: any) {
      diagnosticResult = {
        success: false,
        message: error?.message || String(error),
        steps: [
          {
            name: "Execução do Diagnóstico",
            status: "error",
            message: error?.message || String(error),
          },
        ],
      };
    } finally {
      diagnosticLoading = false;
    }
  }
  async function autoDetect(value: VpsFormData) {
    return api.autoDetectDocker(dockerDetectionPayload(value));
  }
  async function save(value: VpsFormData) {
    if (!value.name.trim()) return;
    try {
      await api.saveServer(serverPayload(value));
      await load();
      notifySuccess("Servidor salvo no SQLite com sucesso!");
    } catch (error: any) {
      notifyError(`Erro ao salvar no banco: ${error?.message || error}`);
    }
    showVpsModal = false;
  }
  function requestDelete(server: VpsServer) {
    serverToDelete = server;
    showDeleteConfirm = true;
  }
  async function confirmDelete() {
    if (!serverToDelete) return;
    try {
      await api.deleteServer(serverToDelete.id);
      await load();
      notifySuccess(`Servidor '${serverToDelete.name}' removido.`);
    } catch (error: any) {
      notifyError(`Erro ao remover: ${error?.message || error}`);
    }
    showDeleteConfirm = false;
  }
  async function select(server: VpsServer) {
    try {
      await api.setActiveServer(server.id);
      await load();
      notifySuccess(`Servidor '${server.name}' ativado com sucesso!`);
    } catch {
      getData().activeVps = server;
      notifySuccess(`Servidor '${server.name}' selecionado como ativo.`);
    }
  }
  function changeLocale(locale: string) {
    const data = getData();
    setLocale(locale);
    if (data.activeProfile) data.activeProfile.locale = locale;
    notifySuccess(
      t("config.language_changed", {
        language: locale === "pt-BR" ? "Português" : "English",
      }),
    );
  }
  return {
    get form() {
      return form;
    },
    set form(value: VpsFormData) {
      form = value;
    },
    get showVpsModal() {
      return showVpsModal;
    },
    set showVpsModal(value: boolean) {
      showVpsModal = value;
    },
    get showDiagnosticModal() {
      return showDiagnosticModal;
    },
    set showDiagnosticModal(value: boolean) {
      showDiagnosticModal = value;
    },
    get diagnosticTitle() {
      return diagnosticTitle;
    },
    get diagnosticLoading() {
      return diagnosticLoading;
    },
    get diagnosticResult() {
      return diagnosticResult;
    },
    get showDeleteConfirm() {
      return showDeleteConfirm;
    },
    set showDeleteConfirm(value: boolean) {
      showDeleteConfirm = value;
    },
    get serverToDelete() {
      return serverToDelete;
    },
    get activeLocale() {
      return getData().activeProfile?.locale || getLocale();
    },
    load,
    openCreate,
    openEdit,
    test,
    autoDetect,
    save,
    requestDelete,
    confirmDelete,
    select,
    changeLocale,
  };
}
