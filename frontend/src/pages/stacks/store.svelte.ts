import { notifyError, notifySuccess } from "$shared/stores/notification.svelte";
import { triggerRefresh } from "$shared/stores/refresh.svelte";
import {
  listStacksUseCase,
  getServerCapabilitiesUseCase,
  saveStackUseCase,
  deleteStackDefinitionUseCase,
  removeStackRemoteUseCase,
  deployStackUseCase,
  stopStackUseCase,
  getStackLogsUseCase,
} from "./api";
import { getCached, setCached } from "$shared/stores/swr-cache";
import { filterStacks } from "./service";
import type { ServerCapabilities, StackItem } from "./api";
import { browse } from "./api";
import { t } from "$shared/stores/locale.svelte";

export function createStacksStore(getData: () => any) {
  let searchQuery = $state("");
  let stacks = $state<StackItem[]>([]);
  let loading = $state(false);
  let capabilities = $state<ServerCapabilities | null>(null);

  async function browseFolder(path = "") {
    try {
      return await browse(path);
    } catch (error: any) {
      notifyError(error?.message || "Erro ao navegar nas pastas");
      return null;
    }
  }

  async function load() {
    const data = getData();
    const profileId = data?.activeProfile?.id || "default";
    const cacheKey = `stacks:${profileId}`;
    const cachedData = getCached<StackItem[]>(cacheKey);

    if (cachedData && cachedData.length >= 0) {
      stacks = cachedData;
      loading = false;
    } else {
      loading = true;
    }

    try {
      const freshStacks = (await listStacksUseCase(profileId)) || [];
      stacks = freshStacks;
      setCached(cacheKey, freshStacks);
      capabilities = data?.activeVps
        ? await getServerCapabilitiesUseCase(profileId)
        : null;
    } catch (error: any) {
      notifyError(error?.message || "Erro ao carregar stacks");
    } finally {
      loading = false;
    }
  }

  const profileId = () => getData()?.activeProfile?.id || "default";

  async function save(
    stack: Omit<Parameters<typeof saveStackUseCase>[0], "profileId">,
  ) {
    try {
      await saveStackUseCase({ ...stack, profileId: profileId() });
      notifySuccess(t("stacks.saving_status").replace("{name}", stack.name));
      await load();
      triggerRefresh();
      return true;
    } catch (error: any) {
      notifyError(error?.message || "Erro ao salvar stack");
      return false;
    }
  }

  async function deleteLocal(id: string) {
    try {
      await deleteStackDefinitionUseCase(id);
      await load();
      triggerRefresh();
      return true;
    } catch (error: any) {
      notifyError(error?.message || "Erro ao excluir stack");
      return false;
    }
  }

  async function removeRemote(id: string, deleteVolumes: boolean) {
    try {
      const result = await removeStackRemoteUseCase(
        profileId(),
        id,
        deleteVolumes,
      );
      if (result.success) notifySuccess(result.message);
      else notifyError(result.message);
      await load();
      triggerRefresh();
      return result;
    } catch (error: any) {
      notifyError(error?.message || "Erro ao remover stack na VPS");
      return { success: false, message: error?.message || String(error) };
    }
  }

  function deploy(id: string) {
    return deployStackUseCase(profileId(), id).catch((error: any) => {
      notifyError(error?.message || "Erro ao iniciar deploy");
      throw error;
    });
  }

  async function stop(id: string) {
    try {
      const result = await stopStackUseCase(profileId(), id);
      if (result.success) notifySuccess(result.message);
      else notifyError(result.message);
      await load();
      triggerRefresh();
      return result;
    } catch (error: any) {
      notifyError(error?.message || "Erro ao parar stack");
      throw error;
    }
  }

  function logs(id: string) {
    return getStackLogsUseCase(profileId(), id, 200);
  }

  return {
    get searchQuery() {
      return searchQuery;
    },
    set searchQuery(value: string) {
      searchQuery = value;
    },
    get stacks() {
      return stacks;
    },
    get loading() {
      return loading;
    },
    get capabilities() {
      return capabilities;
    },
    get filteredStacks() {
      return filterStacks(stacks, searchQuery);
    },
    load,
    browseFolder,
    save,
    deleteLocal,
    removeRemote,
    deploy,
    stop,
    logs,
  };
}
