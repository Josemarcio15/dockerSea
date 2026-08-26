import { t } from "$shared/stores/locale.svelte";
import { triggerRefresh } from "$shared/stores/refresh.svelte";
import { notifySuccess, notifyError } from "$shared/stores/notification.svelte";
import type { VpsServer } from "../../../bindings/go-walis/internal/core/db/models.js";
import type { Container } from "../../../bindings/go-walis/internal/containers/models.js";
import {
  type DockerNetwork,
  filterCustomNetworks,
  listNetworks,
  createNetwork,
  removeNetworks,
  connectContainer,
  disconnectContainer,
  pruneNetworks,
} from "$lib/domains/networks";

export function createNetworksState(getServer: () => VpsServer | undefined) {
  let searchQuery = $state("");
  let networks = $state<DockerNetwork[]>([]);
  let containers = $state<Container[]>([]);
  let loading = $state(true);
  let fetchError = $state("");
  let selectedNetworkNames = $state<string[]>([]);

  // Modal: Nova Rede
  let showCreateModal = $state(false);
  let modalName = $state("");
  let modalDriver = $state("bridge");
  let modalSubnet = $state("");
  let modalGateway = $state("");

  // Modal: Conectar Container
  let showConnectModal = $state(false);
  let selectedNetworkName = $state("");
  let selectedContainerName = $state("");

  // Modal: Prune
  let showPruneModal = $state(false);

  // Redes filtradas derivadas
  let filteredNetworks = $derived(
    filterCustomNetworks(networks, searchQuery),
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
      const data = await listNetworks(server);
      networks = data.networks || [];
      containers = data.containers || [];
    } catch (e: any) {
      fetchError = e.message || String(e);
      if (!silent) {
        networks = [];
        containers = [];
      }
    } finally {
      loading = false;
    }
  }

  function toggleAll() {
    if (selectedNetworkNames.length === filteredNetworks.length) {
      selectedNetworkNames = [];
    } else {
      selectedNetworkNames = filteredNetworks.map((n) => n.name);
    }
  }

  function handleToggleSelect(name: string) {
    if (selectedNetworkNames.includes(name)) {
      selectedNetworkNames = selectedNetworkNames.filter((n) => n !== name);
    } else {
      selectedNetworkNames = [...selectedNetworkNames, name];
    }
  }

  function openCreateModal() {
    modalName = "";
    modalDriver = "bridge";
    modalSubnet = "";
    modalGateway = "";
    showCreateModal = true;
  }

  function openConnectModal(networkName: string) {
    selectedNetworkName = networkName;
    selectedContainerName = containers[0]?.name || "";
    showConnectModal = true;
  }

  async function doCreate() {
    const server = getServer();
    if (!modalName.trim() || !server) return;

    notifySuccess(t("networks.creating_status").replace("{name}", modalName));
    showCreateModal = false;

    try {
      const res = await createNetwork(server, {
        name: modalName.trim(),
        driver: modalDriver,
        subnet: modalSubnet.trim(),
        gateway: modalGateway.trim(),
      });

      if (res.success) {
        notifySuccess(res.message);
        modalName = "";
        modalDriver = "bridge";
        modalSubnet = "";
        modalGateway = "";
      } else {
        notifyError(res.message || "Erro ao criar rede");
      }
      triggerRefresh();
    } catch (e: any) {
      notifyError(e.message || String(e));
    }
  }

  async function doDelete(names: string | string[]) {
    const server = getServer();
    if (!server) return;

    const namesList = Array.isArray(names)
      ? names
      : names.split(",").map((n) => n.trim()).filter(Boolean);

    if (namesList.length === 0) return;

    notifySuccess(
      t("networks.deleting_status").replace("{name}", namesList.join(", ")),
    );

    try {
      const res = await removeNetworks(server, namesList);
      if (res.success) {
        notifySuccess(res.message);
        selectedNetworkNames = selectedNetworkNames.filter(
          (n) => !namesList.includes(n),
        );
      } else {
        notifyError(res.message || "Erro ao remover rede");
      }
      triggerRefresh();
    } catch (e: any) {
      notifyError(e.message || String(e));
    }
  }

  async function doDeleteSelected() {
    if (selectedNetworkNames.length === 0) return;
    await doDelete(selectedNetworkNames);
  }

  async function doPrune() {
    const server = getServer();
    if (!server) return;

    showPruneModal = false;
    notifySuccess(t("networks.pruning_status"));

    try {
      const res = await pruneNetworks(server);
      if (res.success) {
        notifySuccess(res.message);
      } else {
        notifyError(res.message || "Erro ao limpar redes");
      }
      triggerRefresh();
    } catch (e: any) {
      notifyError(e.message || String(e));
    }
  }

  async function doConnect() {
    const server = getServer();
    if (!selectedNetworkName || !selectedContainerName || !server) return;

    notifySuccess(
      t("networks.connecting_status")
        .replace("{container}", selectedContainerName)
        .replace("{network}", selectedNetworkName),
    );
    showConnectModal = false;

    try {
      const res = await connectContainer(
        server,
        selectedNetworkName,
        selectedContainerName,
      );
      if (res.success) {
        notifySuccess(res.message);
      } else {
        notifyError(res.message || "Erro ao conectar container");
      }
      triggerRefresh();
    } catch (e: any) {
      notifyError(e.message || String(e));
    }
  }

  async function doDisconnect(networkName: string, containerName: string) {
    const server = getServer();
    if (!networkName || !containerName || !server) return;

    notifySuccess(
      t("networks.disconnecting_status")
        .replace("{container}", containerName)
        .replace("{name}", networkName),
    );

    try {
      const res = await disconnectContainer(
        server,
        networkName,
        containerName,
      );
      if (res.success) {
        notifySuccess(res.message);
      } else {
        notifyError(res.message || "Erro ao desconectar container");
      }
      triggerRefresh();
    } catch (e: any) {
      notifyError(e.message || String(e));
    }
  }

  return {
    get searchQuery() {
      return searchQuery;
    },
    set searchQuery(val: string) {
      searchQuery = val;
    },
    get networks() {
      return networks;
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
    get selectedNetworkNames() {
      return selectedNetworkNames;
    },
    get filteredNetworks() {
      return filteredNetworks;
    },

    // Create Modal State
    get showCreateModal() {
      return showCreateModal;
    },
    set showCreateModal(val: boolean) {
      showCreateModal = val;
    },
    get modalName() {
      return modalName;
    },
    set modalName(val: string) {
      modalName = val;
    },
    get modalDriver() {
      return modalDriver;
    },
    set modalDriver(val: string) {
      modalDriver = val;
    },
    get modalSubnet() {
      return modalSubnet;
    },
    set modalSubnet(val: string) {
      modalSubnet = val;
    },
    get modalGateway() {
      return modalGateway;
    },
    set modalGateway(val: string) {
      modalGateway = val;
    },

    // Connect Modal State
    get showConnectModal() {
      return showConnectModal;
    },
    set showConnectModal(val: boolean) {
      showConnectModal = val;
    },
    get selectedNetworkName() {
      return selectedNetworkName;
    },
    set selectedNetworkName(val: string) {
      selectedNetworkName = val;
    },
    get selectedContainerName() {
      return selectedContainerName;
    },
    set selectedContainerName(val: string) {
      selectedContainerName = val;
    },

    // Prune Modal State
    get showPruneModal() {
      return showPruneModal;
    },
    set showPruneModal(val: boolean) {
      showPruneModal = val;
    },

    // Actions
    fetchAll,
    toggleAll,
    handleToggleSelect,
    openCreateModal,
    openConnectModal,
    doCreate,
    doDelete,
    doDeleteSelected,
    doPrune,
    doConnect,
    doDisconnect,
  };
}
