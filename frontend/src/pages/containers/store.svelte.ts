import { t } from "$shared/stores/locale.svelte";
import { triggerRefresh } from "$shared/stores/refresh.svelte";
import {
  notifySuccess,
  notifyWarning,
  notifyError,
} from "$shared/stores/notification.svelte";
import type { VpsServer } from "./api";
import { api } from "./api";
import { filterContainers } from "./service";
import { getCached, setCached } from "$shared/stores/swr-cache";
import type { Container, ContainerActionType, ContainersStore } from "./types";

export function createContainersStore(
  getServer: () => VpsServer | undefined,
): ContainersStore {
  let selectedNames = $state<string[]>([]);
  let searchQuery = $state("");
  let containers = $state<Container[]>([]);
  let loading = $state(true);
  let fetchError = $state("");
  let showLogs = $state(false);
  let logsTitle = $state("");
  let logsLoading = $state(false);
  let logsContent = $state<string[]>([]);
  let filteredContainers = $derived(filterContainers(containers, searchQuery));

  async function fetchAll(silent = false, forceRefresh = false) {
    const server = getServer();
    if (!server) {
      loading = false;
      return;
    }

    const cacheKey = `containers:${server.id}`;
    const cachedData = forceRefresh ? undefined : getCached<Container[]>(cacheKey);

    // Se temos dados em cache e não é um refresh forçado, renderizamos instantaneamente (0ms)
    if (cachedData && cachedData.length >= 0) {
      containers = cachedData;
      loading = false;
      silent = true; // Searches em background sem travar o layout com spinner
    } else if (!silent) {
      loading = true;
    }

    fetchError = "";
    try {
      const freshData = (await api.list(server)) || [];
      containers = freshData;
      setCached(cacheKey, freshData);
    } catch (error: any) {
      fetchError = error.message || String(error);
      if (!cachedData && !silent) containers = [];
    } finally {
      loading = false;
    }
  }

  function toggleAll() {
    selectedNames =
      selectedNames.length === filteredContainers.length
        ? []
        : filteredContainers.map((container) => container.name);
  }

  function handleToggleSelect(name: string) {
    selectedNames = selectedNames.includes(name)
      ? selectedNames.filter((item) => item !== name)
      : [...selectedNames, name];
  }

  async function doActionSelected(action: ContainerActionType) {
    const server = getServer();
    if (!server || selectedNames.length === 0) return;

    notifyWarning(
      t(`containers.status_${action}`).replace(
        "{name}",
        selectedNames.join(", "),
      ),
    );

    try {
      const result = await api.action(server, action, selectedNames);
      if (result.success) {
        notifySuccess(result.message);
        if (action === "rm") selectedNames = [];
      } else {
        notifyError(result.message || "Erro na execução");
      }
      triggerRefresh();
    } catch (error: any) {
      notifyError(error.message || String(error));
    }
  }

  async function openLogs(name: string) {
    const server = getServer();
    if (!server) return;

    logsTitle = `Logs: ${name}`;
    logsContent = [];
    logsLoading = true;
    showLogs = true;
    try {
      const logs = await api.logs(server, name);
      logsContent = logs ? logs.split("\n") : ["(Sem logs disponíveis)"];
    } catch (error: any) {
      logsContent = [
        `Erro ao carregar logs: ${error.message || String(error)}`,
      ];
    } finally {
      logsLoading = false;
    }
  }

  function setupEventsStream() {
    const server = getServer();
    if (!server?.id) return () => {};

    api.events.startEventsStream(server).catch(console.warn);
    const unsubscribe = api.events.subscribeToEvents((event: any) => {
      const data = event?.data?.[0] || event?.data || event;
      if (!data?.serverId || data.serverId === server.id) fetchAll(true);
    });

    return () => {
      unsubscribe();
      api.events.stopEventsStream(server.id);
    };
  }
  return {
    get selectedNames() {
      return selectedNames;
    },
    get searchQuery() {
      return searchQuery;
    },
    set searchQuery(value) {
      searchQuery = value;
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
    get showLogs() {
      return showLogs;
    },
    set showLogs(value) {
      showLogs = value;
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
    fetchAll,
    toggleAll,
    handleToggleSelect,
    doActionSelected,
    openLogs,
    setupEventsStream,
  };
}
