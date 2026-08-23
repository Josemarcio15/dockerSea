import { t } from "$lib/stores/locale.svelte";
import { triggerRefresh } from "$lib/stores/refresh.svelte";
import {
  notifySuccess,
  notifyWarning,
  notifyError,
} from "$lib/stores/notification.svelte";
import type {
  VpsServer,
  ImageHistoryItem,
} from "../../../bindings/go-walis/internal/core/db/models.js";
import {
  type DockerImage,
  filterImages,
  listImages,
  removeImages,
  pullImage,
  imageWailsApi,
} from "$lib/domains/images";
import { Events } from "@wailsio/runtime";

export function createImagesState(
  getServer: () => VpsServer | undefined,
  getProfileId: () => string = () => "default",
  getServersList: () => VpsServer[] = () => [],
) {
  let activeTab = $state<"my_images" | "download" | "transfer">("my_images");
  let searchQuery = $state("");
  let downloadQuery = $state("");
  let selectedImageIds = $state<string[]>([]);
  let imageFilter = $state<"all" | "in_use" | "unused" | "dangling">("all");

  // Transfer state
  let transferSourceId = $state("");
  let transferDestId = $state("");
  let sourceImages = $state<DockerImage[]>([]);
  let sourceLoading = $state(false);
  let selectedTransferIds = $state<string[]>([]);
  let transferInProgress = $state(false);
  let transferElapsedSeconds = $state(0);
  let transferStage = $state("preparing");
  let transferFormattedBytes = $state("0 B");
  let transferFormattedTotalBytes = $state("");
  let transferSpeed = $state("0 B/s");
  let transferMessage = $state("");
  let transferPercent = $state(0);

  // Terminal Modal state
  let showTerminal = $state(false);
  let terminalTitle = $state("");
  let terminalLoading = $state(false);
  let terminalLogs = $state<string[]>([]);

  // Config Modal state
  let showConfig = $state(false);
  let configTargetImage = $state("");

  // Pull Progress Modal state
  let showPullProgress = $state(false);
  let pullImageTargetName = $state("");

  // Client-side fetched data
  let images = $state<DockerImage[]>([]);
  let imageHistory = $state<ImageHistoryItem[]>([]);
  let selectedHistoryIds = $state<string[]>([]);
  let loading = $state(true);
  let fetchError = $state("");

  // Metrics
  let countInUse = $derived(
    (images || []).filter(
      (img) => img.containersUsing && img.containersUsing.length > 0,
    ).length,
  );
  let countUnused = $derived(
    (images || []).filter(
      (img) =>
        (!img.containersUsing || img.containersUsing.length === 0) &&
        img.repo !== "<none>",
    ).length,
  );
  let countDangling = $derived(
    (images || []).filter((img) => img.repo === "<none>").length,
  );

  let totalDiskUsage = $derived(() => {
    const totalBytes = (images || []).reduce(
      (acc, img) => acc + (img.rawSizeBytes || 0),
      0,
    );
    if (totalBytes === 0) return "0 B";
    const units = ["B", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(totalBytes) / Math.log(1024));
    return (totalBytes / Math.pow(1024, i)).toFixed(1) + " " + units[i];
  });

  // Filtered images with search & quick filters
  let filteredImages = $derived(
    (images || []).filter((img) => {
      const matchesSearch =
        img.repo.toLowerCase().includes(searchQuery.toLowerCase()) ||
        img.id.toLowerCase().includes(searchQuery.toLowerCase()) ||
        img.tag.toLowerCase().includes(searchQuery.toLowerCase());

      if (!matchesSearch) return false;

      if (imageFilter === "in_use")
        return img.containersUsing && img.containersUsing.length > 0;
      if (imageFilter === "unused")
        return (
          (!img.containersUsing || img.containersUsing.length === 0) &&
          img.repo !== "<none>"
        );
      if (imageFilter === "dangling") return img.repo === "<none>";
      return true;
    }),
  );

  async function fetchImages(silent = false) {
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
      const list = await listImages(server);
      images = list || [];
    } catch (e: any) {
      fetchError = e.message || String(e);
      if (!silent) {
        images = [];
      }
    } finally {
      loading = false;
    }
  }

  async function fetchHistory() {
    try {
      const profileId = getProfileId();
      const hist = await imageWailsApi.listHistory(profileId);
      imageHistory = hist || [];
    } catch (e) {
      imageHistory = [];
    }
  }

  function toggleChecked(id: string) {
    if (selectedImageIds.includes(id)) {
      selectedImageIds = selectedImageIds.filter((x) => x !== id);
    } else {
      selectedImageIds = [...selectedImageIds, id];
    }
  }

  function toggleAll() {
    if (selectedImageIds.length === filteredImages.length) {
      selectedImageIds = [];
    } else {
      selectedImageIds = filteredImages.map((img) => img.id);
    }
  }

  function toggleCheckedHistory(id: string) {
    if (selectedHistoryIds.includes(id)) {
      selectedHistoryIds = selectedHistoryIds.filter((x) => x !== id);
    } else {
      selectedHistoryIds = [...selectedHistoryIds, id];
    }
  }

  function toggleAllHistory() {
    if (selectedHistoryIds.length === (imageHistory || []).length) {
      selectedHistoryIds = [];
    } else {
      selectedHistoryIds = (imageHistory || []).map((h) => h.id);
    }
  }

  async function handleDeleteHistorySelected() {
    if (selectedHistoryIds.length === 0) return;
    try {
      await imageWailsApi.deleteHistory(selectedHistoryIds);
      selectedHistoryIds = [];
      await fetchHistory();
      notifySuccess("Itens do histórico removidos com sucesso!");
    } catch (e: any) {
      notifyError(e.message || "Erro ao deletar histórico");
    }
  }

  async function handleClearHistory() {
    try {
      const profileId = getProfileId();
      await imageWailsApi.clearHistory(profileId);
      await fetchHistory();
      notifySuccess("Histórico limpo com sucesso!");
    } catch (e: any) {
      notifyError(e.message || "Erro ao limpar histórico");
    }
  }

  async function handlePull(imageName: string) {
    if (!imageName.trim()) {
      notifyWarning(t("images.pull_loading").replace("{image}", ""));
      return;
    }
    const server = getServer();
    if (!server) {
      notifyError("Selecione um servidor VPS ativo.");
      return;
    }

    pullImageTargetName = imageName;
    showPullProgress = true;

    const profileId = getProfileId();
    pullImage(server, imageName, profileId)
      .then((res) => {
        if (!res.success) {
          notifyError(res.message || "Falha ao baixar imagem");
        }
      })
      .catch((err) => {
        notifyError(err.message || String(err));
      });
  }

  async function handleDeleteSelected() {
    if (selectedImageIds.length === 0) {
      notifyWarning("Selecione pelo menos uma imagem para deletar.");
      return;
    }
    const server = getServer();
    if (!server) {
      notifyError("Nenhum servidor VPS ativo.");
      return;
    }

    notifySuccess(
      t("images.delete_progress", { count: selectedImageIds.length }),
    );

    try {
      const result = await removeImages(server, selectedImageIds);
      if (result.success) {
        notifySuccess(
          t("images.delete_success", {
            count: result.count || selectedImageIds.length,
          }),
        );
        selectedImageIds = [];
        await fetchImages(true);
      } else {
        notifyError(
          t("images.delete_error", { error: result.message || "Unknown" }),
        );
        await fetchImages(true);
      }
      triggerRefresh();
    } catch (e: any) {
      notifyError(e.message || String(e));
    }
  }

  function openConfig(imageRepo: string, imageTag: string) {
    configTargetImage = `${imageRepo}:${imageTag}`;
    showConfig = true;
  }

  async function fetchSourceImages() {
    if (!transferSourceId) {
      sourceImages = [];
      return;
    }
    const servers = getServersList();
    const sourceServer = servers.find((s: any) => s.id === transferSourceId);
    if (!sourceServer) {
      sourceImages = [];
      return;
    }

    sourceLoading = true;
    selectedTransferIds = [];
    try {
      const list = await listImages(sourceServer);
      sourceImages = list || [];
    } catch (e: any) {
      notifyError(`Erro ao carregar imagens da origem: ${e.message || e}`);
      sourceImages = [];
    } finally {
      sourceLoading = false;
    }
  }

  function getTransferImageKey(img: DockerImage): string {
    if (img.repo && img.repo !== "<none>") {
      return `${img.repo}:${img.tag || "latest"}`;
    }
    return img.id;
  }

  function toggleTransferImage(key: string) {
    if (selectedTransferIds.includes(key)) {
      selectedTransferIds = selectedTransferIds.filter((x) => x !== key);
    } else {
      selectedTransferIds = [...selectedTransferIds, key];
    }
  }

  function toggleAllTransfer() {
    if (selectedTransferIds.length === sourceImages.length) {
      selectedTransferIds = [];
    } else {
      selectedTransferIds = sourceImages.map((img) =>
        getTransferImageKey(img),
      );
    }
  }

  async function handleTransfer() {
    if (!transferSourceId || !transferDestId) {
      notifyWarning(t("images.transfer_no_selection"));
      return;
    }
    if (transferSourceId === transferDestId) {
      notifyWarning(t("images.transfer_same_server"));
      return;
    }
    if (selectedTransferIds.length === 0) {
      notifyWarning(t("images.transfer_no_selection"));
      return;
    }

    const servers = getServersList();
    const srcServer = servers.find((s: any) => s.id === transferSourceId);
    const dstServer = servers.find((s: any) => s.id === transferDestId);

    if (!srcServer || !dstServer) {
      notifyError("Servidores de origem ou destino inválidos.");
      return;
    }

    transferInProgress = true;
    transferElapsedSeconds = 0;
    transferStage = "preparing";
    transferFormattedBytes = "0 B";
    transferFormattedTotalBytes = "";
    transferSpeed = "0 B/s";
    transferPercent = 0;
    transferMessage = "Iniciando pipeline de transferência...";

    const timer = setInterval(() => {
      transferElapsedSeconds += 1;
    }, 1000);

    const unsubProgress = Events.On(
      "docker:image:transfer:progress",
      (event: any) => {
        const p = event?.data ?? event;
        if (p) {
          if (p.stage) transferStage = p.stage;
          if (p.formattedBytes) transferFormattedBytes = p.formattedBytes;
          if (p.formattedTotalBytes)
            transferFormattedTotalBytes = p.formattedTotalBytes;
          if (p.speed) transferSpeed = p.speed;
          if (typeof p.percent === "number") transferPercent = p.percent;
          if (p.message) transferMessage = p.message;
        }
      },
    );

    const unsubComplete = Events.On(
      "docker:image:transfer:complete",
      (event: any) => {
        const c = event?.data ?? event;
        if (c && c.success) {
          transferPercent = 100;
          transferStage = "complete";
        }
      },
    );

    notifySuccess(
      t("images.transfer_in_progress").replace(
        "{count}",
        String(selectedTransferIds.length),
      ),
    );

    try {
      const result = await imageWailsApi.transferImages(
        srcServer,
        dstServer,
        selectedTransferIds,
      );

      if (result.success) {
        notifySuccess(
          t("images.transfer_success").replace(
            "{count}",
            String(result.count || selectedTransferIds.length),
          ),
        );
        selectedTransferIds = [];
        const active = getServer();
        if (active?.id === dstServer.id) {
          await fetchImages(true);
        }
      } else {
        notifyError(
          t("images.transfer_error").replace(
            "{error}",
            result.message || "Unknown",
          ),
        );
      }
    } catch (e: any) {
      notifyError(
        t("images.transfer_error").replace(
          "{error}",
          e.message || String(e),
        ),
      );
    } finally {
      clearInterval(timer);
      unsubProgress();
      unsubComplete();
      transferInProgress = false;
    }
  }

  return {
    get activeTab() {
      return activeTab;
    },
    set activeTab(val: "my_images" | "download" | "transfer") {
      activeTab = val;
    },
    get searchQuery() {
      return searchQuery;
    },
    set searchQuery(val: string) {
      searchQuery = val;
    },
    get downloadQuery() {
      return downloadQuery;
    },
    set downloadQuery(val: string) {
      downloadQuery = val;
    },
    get selectedImageIds() {
      return selectedImageIds;
    },
    set selectedImageIds(val: string[]) {
      selectedImageIds = val;
    },
    get imageFilter() {
      return imageFilter;
    },
    set imageFilter(val: "all" | "in_use" | "unused" | "dangling") {
      imageFilter = val;
    },
    get images() {
      return images;
    },
    get imageHistory() {
      return imageHistory;
    },
    get selectedHistoryIds() {
      return selectedHistoryIds;
    },
    set selectedHistoryIds(val: string[]) {
      selectedHistoryIds = val;
    },
    get loading() {
      return loading;
    },
    get fetchError() {
      return fetchError;
    },
    get filteredImages() {
      return filteredImages;
    },
    get countInUse() {
      return countInUse;
    },
    get countUnused() {
      return countUnused;
    },
    get countDangling() {
      return countDangling;
    },
    get totalDiskUsage() {
      return totalDiskUsage();
    },

    // Pull Progress Modal State
    get showPullProgress() {
      return showPullProgress;
    },
    set showPullProgress(val: boolean) {
      showPullProgress = val;
    },
    get pullImageTargetName() {
      return pullImageTargetName;
    },

    // Config Modal State
    get showConfig() {
      return showConfig;
    },
    set showConfig(val: boolean) {
      showConfig = val;
    },
    get configTargetImage() {
      return configTargetImage;
    },

    // Terminal Modal State
    get showTerminal() {
      return showTerminal;
    },
    set showTerminal(val: boolean) {
      showTerminal = val;
    },
    get terminalTitle() {
      return terminalTitle;
    },
    get terminalLoading() {
      return terminalLoading;
    },
    get terminalLogs() {
      return terminalLogs;
    },

    // Transfer State
    get transferSourceId() {
      return transferSourceId;
    },
    set transferSourceId(val: string) {
      transferSourceId = val;
    },
    get transferDestId() {
      return transferDestId;
    },
    set transferDestId(val: string) {
      transferDestId = val;
    },
    get sourceImages() {
      return sourceImages;
    },
    get sourceLoading() {
      return sourceLoading;
    },
    get selectedTransferIds() {
      return selectedTransferIds;
    },
    set selectedTransferIds(val: string[]) {
      selectedTransferIds = val;
    },
    get transferInProgress() {
      return transferInProgress;
    },
    get transferElapsedSeconds() {
      return transferElapsedSeconds;
    },
    get transferStage() {
      return transferStage;
    },
    get transferFormattedBytes() {
      return transferFormattedBytes;
    },
    get transferFormattedTotalBytes() {
      return transferFormattedTotalBytes;
    },
    get transferSpeed() {
      return transferSpeed;
    },
    get transferMessage() {
      return transferMessage;
    },
    get transferPercent() {
      return transferPercent;
    },

    // Actions
    fetchImages,
    fetchHistory,
    toggleChecked,
    toggleAll,
    toggleCheckedHistory,
    toggleAllHistory,
    handleDeleteHistorySelected,
    handleClearHistory,
    handlePull,
    handleDeleteSelected,
    openConfig,
    fetchSourceImages,
    getTransferImageKey,
    toggleTransferImage,
    toggleAllTransfer,
    handleTransfer,
  };
}
