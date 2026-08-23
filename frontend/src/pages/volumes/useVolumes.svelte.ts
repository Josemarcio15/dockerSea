import { t } from "$lib/stores/locale.svelte";
import { triggerRefresh } from "$lib/stores/refresh.svelte";
import { notifySuccess, notifyError } from "$lib/stores/notification.svelte";
import type { VpsServer } from "../../../bindings/go-walis/internal/core/db/models.js";
import {
  type DockerVolume,
  filterVolumes,
  listVolumes,
  createVolume,
  removeVolumes,
  pruneVolumes,
} from "$lib/domains/volumes";

export function createVolumesState(getServer: () => VpsServer | undefined) {
  let selectedNames = $state<string[]>([]);
  let searchQuery = $state("");
  let volumes = $state<DockerVolume[]>([]);
  let loading = $state(true);
  let fetchError = $state("");

  // Modal: Novo Volume
  let showCreateModal = $state(false);
  let modalName = $state("");
  let modalDriver = $state("");
  let labelEntries = $state<Array<{ key: string; value: string }>>([]);

  // Modal: Confirmar Limpeza
  let showPruneModal = $state(false);

  let filteredVolumes = $derived(
    filterVolumes(volumes, searchQuery),
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
      const list = await listVolumes(server);
      volumes = list || [];
    } catch (e: any) {
      fetchError = e.message || String(e);
      if (!silent) {
        volumes = [];
      }
    } finally {
      loading = false;
    }
  }

  function toggleChecked(name: string) {
    if (selectedNames.includes(name)) {
      selectedNames = selectedNames.filter((n) => n !== name);
    } else {
      selectedNames = [...selectedNames, name];
    }
  }

  function toggleAll() {
    if (selectedNames.length === filteredVolumes.length) {
      selectedNames = [];
    } else {
      selectedNames = filteredVolumes.map((v) => v.name);
    }
  }

  function openCreateModal() {
    modalName = "";
    modalDriver = "";
    labelEntries = [];
    showCreateModal = true;
  }

  function addLabel() {
    labelEntries = [...labelEntries, { key: "", value: "" }];
  }

  function removeLabel(index: number) {
    labelEntries = labelEntries.filter((_, i) => i !== index);
  }

  async function doCreate() {
    const server = getServer();
    if (!modalName.trim() || !server) return;

    showCreateModal = false;
    notifySuccess(t("volumes.creating_status").replace("{name}", modalName));

    const labels: Record<string, string> = {};
    for (const entry of labelEntries) {
      if (entry.key.trim()) {
        labels[entry.key.trim()] = entry.value.trim();
      }
    }

    try {
      const res = await createVolume(server, {
        name: modalName.trim(),
        driver: modalDriver.trim() || undefined,
        labels: Object.keys(labels).length > 0 ? labels : null,
      });

      if (res.success) {
        notifySuccess(res.message);
        modalName = "";
        modalDriver = "";
        labelEntries = [];
      } else {
        notifyError(res.message || "Erro ao criar volume");
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
      t("volumes.deleting_status").replace("{name}", namesList.join(", ")),
    );

    try {
      const res = await removeVolumes(server, namesList);
      if (res.success) {
        notifySuccess(res.message);
        selectedNames = selectedNames.filter((n) => !namesList.includes(n));
      } else {
        notifyError(res.message || "Erro ao remover volume");
      }
      triggerRefresh();
    } catch (e: any) {
      notifyError(e.message || String(e));
    }
  }

  async function doDeleteSelected() {
    if (selectedNames.length === 0) return;
    await doDelete(selectedNames);
  }

  async function doPrune() {
    const server = getServer();
    if (!server) return;

    showPruneModal = false;
    notifySuccess(t("volumes.prune_confirm_msg"));

    try {
      const res = await pruneVolumes(server);
      if (res.success) {
        notifySuccess(res.message);
      } else {
        notifyError(res.message || "Erro ao limpar volumes");
      }
      triggerRefresh();
    } catch (e: any) {
      notifyError(e.message || String(e));
    }
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
    get volumes() {
      return volumes;
    },
    get loading() {
      return loading;
    },
    get fetchError() {
      return fetchError;
    },
    get filteredVolumes() {
      return filteredVolumes;
    },

    // Modal Create
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
    get labelEntries() {
      return labelEntries;
    },
    set labelEntries(val: Array<{ key: string; value: string }>) {
      labelEntries = val;
    },

    // Modal Prune
    get showPruneModal() {
      return showPruneModal;
    },
    set showPruneModal(val: boolean) {
      showPruneModal = val;
    },

    // Actions
    fetchAll,
    toggleChecked,
    toggleAll,
    openCreateModal,
    addLabel,
    removeLabel,
    doCreate,
    doDelete,
    doDeleteSelected,
    doPrune,
  };
}
