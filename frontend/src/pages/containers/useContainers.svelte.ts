import { t } from "$shared/stores/locale.svelte";
import { triggerRefresh } from "$shared/stores/refresh.svelte";
import {
  notifySuccess,
  notifyWarning,
  notifyError,
} from "$shared/stores/notification.svelte";
import type { VpsServer } from "../../../bindings/go-walis/internal/core/db/models.js";
import {
  type Container,
  type ContainerActionType,
  filterContainers,
  listContainers,
  executeContainerAction,
  getContainerLogs,
  containerWailsApi,
} from "$lib/domains/containers";

export function createContainersState(getServer: () => VpsServer | undefined) {
  let selectedNames = $state<string[]>([]);
  let searchQuery = $state("");
  let containers = $state<Container[]>([]);
  let loading = $state(true);
  let fetchError = $state("");

  // Logs Modal state
  let showLogs = $state(false);
  let logsTitle = $state("");
  let logsLoading = $state(false);
  let logsContent = $state<string[]>([]);

  let filteredContainers = $derived(
    filterContainers(containers, searchQuery),
  );

  async function fetchAll(silent = false) {
    const server = getServer();
    if (!server) {
      loading = false;
      return;
    }
    if (!silent) {
      loading = true;
    }
    fetchError = "";
    try {
      const list = await listContainers(server, true);
      containers = list || [];
    } catch (e: any) {
      fetchError = e.message || String(e);
      if (!silent) {
        containers = [];
      }
    } finally {
      loading = false;
    }
  }

  function toggleAll() {
    if (selectedNames.length === filteredContainers.length) {
      selectedNames = [];
    } else {
      selectedNames = filteredContainers.map((c) => c.name);
    }
  }

  function handleToggleSelect(name: string) {
    if (selectedNames.includes(name)) {
      selectedNames = selectedNames.filter((n) => n !== name);
    } else {
      selectedNames = [...selectedNames, name];
    }
  }

  async function doAction(action: ContainerActionType, names: string[]) {
    const server = getServer();
    if (!server || names.length === 0) return;

    notifyWarning(
      t(`containers.status_${action}`).replace("{name}", names.join(", ")),
    );

    try {
      const res = await executeContainerAction(server, action, names);
      if (res.success) {
        notifySuccess(res.message);
        if (action === "rm") {
          selectedNames = selectedNames.filter((n) => !names.includes(n));
        }
      } else {
        notifyError(res.message || "Erro na execução");
      }
      triggerRefresh();
    } catch (e: any) {
      notifyError(e.message || String(e));
    }
  }

  async function doActionSelected(action: ContainerActionType) {
    if (selectedNames.length === 0) return;
    await doAction(action, selectedNames);
  }

  async function openLogs(containerName: string) {
    const server = getServer();
    if (!server) return;

    logsTitle = `Logs: ${containerName}`;
    logsContent = [];
    logsLoading = true;
    showLogs = true;

    try {
      const logs = await getContainerLogs(server, containerName, 200);
      logsContent = logs ? logs.split("\n") : ["(Sem logs disponíveis)"];
    } catch (e: any) {
      logsContent = [`Erro ao carregar logs: ${e.message || String(e)}`];
    } finally {
      logsLoading = false;
    }
  }

  function setupEventsStream() {
    const server = getServer();
    if (!server?.id) return () => {};

    containerWailsApi.startEventsStream(server).catch((err) => {
      console.warn("Não foi possível iniciar stream de eventos Docker:", err);
    });

    const unsubscribe = containerWailsApi.subscribeToEvents((event: any) => {
      const evData = event?.data?.[0] || event?.data || event;
      if (!evData?.serverId || evData.serverId === server.id) {
        fetchAll(true);
      }
    });

    return () => {
      unsubscribe();
      if (server?.id) {
        containerWailsApi.stopEventsStream(server.id);
      }
    };
  }

  return {
    get selectedNames() {
      return selectedNames;
    },
    set selectedNames(val: string[]) {
      selectedNames = val;
    },
    get searchQuery() {
      return searchQuery;
    },
    set searchQuery(val: string) {
      searchQuery = val;
    },
    get containers() {
      return containers;
    },
    get loading() {
      return loading;
    },
    get fetchError() {
      return fetchError;
    },
    get filteredContainers() {
      return filteredContainers;
    },

    // Logs state
    get showLogs() {
      return showLogs;
    },
    set showLogs(val: boolean) {
      showLogs = val;
    },
    get logsTitle() {
      return logsTitle;
    },
    get logsLoading() {
      return logsLoading;
    },
    get logsContent() {
      return logsContent;
    },

    // Actions
    fetchAll,
    toggleAll,
    handleToggleSelect,
    doAction,
    doActionSelected,
    openLogs,
    setupEventsStream,
  };
}
