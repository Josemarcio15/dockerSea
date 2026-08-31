import { getLocale, setLocale, t } from "$shared/stores/locale.svelte";
import { notifySuccess, notifyError } from "$shared/stores/notification.svelte";
import { Dialogs } from "@wailsio/runtime";
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
  let showResetConfirm = $state(false);
  let isBackingUp = $state(false);
  let isResetting = $state(false);
  let dbPath = $state("");

  async function load() {
    try {
      const dbInfo = await api.getDatabaseInfo();
      if (dbInfo?.path) {
        dbPath = dbInfo.path;
      }
      const servers = await api.listServers();
      if (servers?.length) {
        const data = getData();
        data.servers = servers;
        const active = servers.find((server) => server.isActive);
        if (active) data.activeVps = active;
      } else {
        const data = getData();
        data.servers = [];
        data.activeVps = null;
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

  let isRestoring = $state(false);

  async function exportBackup() {
    isBackingUp = true;
    try {
      const now = new Date();
      const dateStr = now.toISOString().slice(0, 10);
      const defaultFilename = `docksea_backup_${dateStr}.db`;

      const selected = await Dialogs.SaveFile({
        Title: t("config.db_save_dialog_title"),
        DefaultFilename: defaultFilename,
        Filters: [
          {
            DisplayName: "SQLite Database (*.db)",
            Pattern: "*.db",
          },
        ],
      });

      if (!selected || typeof selected !== "string") {
        return; // Usuário cancelou
      }

      await api.exportDatabaseBackup(selected);
      notifySuccess(t("config.db_backup_success", { path: selected }));
    } catch (error: any) {
      notifyError(`Erro ao exportar backup: ${error?.message || error}`);
    } finally {
      isBackingUp = false;
    }
  }

  async function restoreBackup() {
    isRestoring = true;
    try {
      const selected = await Dialogs.OpenFile({
        Title: t("config.db_open_dialog_title"),
        CanChooseFiles: true,
        CanChooseDirectories: false,
        AllowsMultipleSelection: false,
        Filters: [
          {
            DisplayName: "SQLite Database (*.db)",
            Pattern: "*.db",
          },
        ],
      });

      if (!selected) {
        return; // Usuário cancelou
      }

      const filePath = Array.isArray(selected) ? selected[0] : selected;
      if (!filePath || typeof filePath !== "string") {
        return;
      }

      await api.restoreDatabaseBackup(filePath);
      await load();
      notifySuccess(t("config.db_restore_success"));
    } catch (error: any) {
      notifyError(`Erro ao restaurar backup: ${error?.message || error}`);
    } finally {
      isRestoring = false;
    }
  }

  function requestResetDb() {
    showResetConfirm = true;
  }

  async function confirmResetDb() {
    isResetting = true;
    try {
      await api.resetDatabase();
      await load();
      notifySuccess(t("config.db_reset_success"));
    } catch (error: any) {
      notifyError(`Erro ao resetar banco: ${error?.message || error}`);
    } finally {
      isResetting = false;
      showResetConfirm = false;
    }
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
    get showResetConfirm() {
      return showResetConfirm;
    },
    set showResetConfirm(value: boolean) {
      showResetConfirm = value;
    },
    get isBackingUp() {
      return isBackingUp;
    },
    get isRestoring() {
      return isRestoring;
    },
    get isResetting() {
      return isResetting;
    },
    get dbPath() {
      return dbPath;
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
    exportBackup,
    restoreBackup,
    requestResetDb,
    confirmResetDb,
  };
}
