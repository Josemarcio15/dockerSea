<script lang="ts">
  import { tick } from "svelte";
  import { t } from "$lib/stores/locale.svelte";
  import { useRefreshKey, triggerRefresh } from "$lib/stores/refresh.svelte";
  import {
    notifySuccess,
    notifyWarning,
    notifyError,
  } from "$lib/stores/notification.svelte";
  import DockerseaLoading from "$lib/components/DockerseaLoading.svelte";
  import StatusBanner from "$lib/components/StatusBanner.svelte";
  import TerminalModal from "$lib/components/TerminalModal.svelte";
  import ConfigModal from "$lib/components/ConfigModal.svelte";
  import PullProgressModal from "$lib/components/PullProgressModal.svelte";
  import ImageCard from "$lib/components/ImageCard.svelte";
  import VpsSelectWarning from "$lib/components/VpsSelectWarning.svelte";

  let { data } = $props();

  let activeTab = $state<"my_images" | "download" | "transfer">("my_images");
  let searchQuery = $state("");
  let downloadQuery = $state("");
  let selectedImageIds = $state<string[]>([]);

  // Transfer state
  let transferSourceId = $state("");
  let transferDestId = $state("");
  let sourceImages = $state<any[]>([]);
  let sourceLoading = $state(false);
  let selectedTransferIds = $state<string[]>([]);
  let transferInProgress = $state(false);
  let transferElapsedSeconds = $state(0);

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
  let images = $state<any[]>([]);
  let loading = $state(true);
  let fetchError = $state("");

  // Reactive fetch: runs on mount and whenever triggerRefresh() is called
  async function fetchImages() {
    if (!data?.activeVps) {
      loading = false;
      return;
    }
    try {
      const res = await fetch("/api/images");
      const json = await res.json();
      images = json.images || [];
      fetchError = json.error || "";
    } catch (e: any) {
      fetchError = e.message;
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    useRefreshKey();
    fetchImages();
  });

  // Aquecer conexão do destino quando selecionado
  $effect(() => {
    const destId = transferDestId;
    if (destId) warmupConnection(destId);
  });

  // Filter states and categorization
  let imageFilter = $state<"all" | "in_use" | "unused" | "dangling">("all");

  // Calculate total disk usage
  const totalDiskUsage = $derived(() => {
    const totalBytes = (images || []).reduce((acc, img) => acc + (img.rawSizeBytes || 0), 0);
    if (totalBytes === 0) return "0 B";
    const units = ["B", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(totalBytes) / Math.log(1024));
    return (totalBytes / Math.pow(1024, i)).toFixed(1) + " " + units[i];
  });

  // Count metrics
  const countInUse = $derived((images || []).filter(img => (img.containersUsing && img.containersUsing.length > 0)).length);
  const countUnused = $derived((images || []).filter(img => (!img.containersUsing || img.containersUsing.length === 0) && img.repo !== "<none>").length);
  const countDangling = $derived((images || []).filter(img => img.repo === "<none>").length);

  // Filter downloaded images with quick filters
  let filteredImages = $derived(
    (images || []).filter((img) => {
      const matchesSearch =
        img.repo.toLowerCase().includes(searchQuery.toLowerCase()) ||
        img.id.toLowerCase().includes(searchQuery.toLowerCase()) ||
        img.tag.toLowerCase().includes(searchQuery.toLowerCase());

      if (!matchesSearch) return false;

      if (imageFilter === "in_use") return img.containersUsing && img.containersUsing.length > 0;
      if (imageFilter === "unused") return (!img.containersUsing || img.containersUsing.length === 0) && img.repo !== "<none>";
      if (imageFilter === "dangling") return img.repo === "<none>";
      return true;
    }),
  );

  // Checked helper
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

  let selectedHistoryIds = $state<string[]>([]);

  function toggleCheckedHistory(id: string) {
    if (selectedHistoryIds.includes(id)) {
      selectedHistoryIds = selectedHistoryIds.filter((x) => x !== id);
    } else {
      selectedHistoryIds = [...selectedHistoryIds, id];
    }
  }

  function toggleAllHistory() {
    if (selectedHistoryIds.length === (data.imageHistory || []).length) {
      selectedHistoryIds = [];
    } else {
      selectedHistoryIds = (data.imageHistory || []).map((h: any) => h.id);
    }
  }

  async function handleDeleteHistorySelected() {
    if (selectedHistoryIds.length === 0) return;
    const formData = new FormData();
    formData.append("ids", selectedHistoryIds.join(","));

    try {
      await fetch("?/deleteHistory", {
        method: "POST",
        body: formData,
      });
      selectedHistoryIds = [];
      await invalidateAll();
    } catch (e) {}
  }

  // Pull image
  async function handlePull(imageName: string) {
    if (!imageName.trim()) {
      notifyWarning(t("images.pull_loading").replace("{image}", ""));
      return;
    }
    pullImageTargetName = imageName;
    showPullProgress = true;
  }

  // Delete selected images
  async function handleDeleteSelected() {
    if (selectedImageIds.length === 0) {
      notifyWarning("Selecione pelo menos uma imagem para deletar.");
      return;
    }

    notifySuccess(t("images.delete_progress", { count: selectedImageIds.length }));

    const formData = new FormData();
    formData.append("ids", selectedImageIds.join(","));

    try {
      const response = await fetch("/api/images/delete", {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ids: selectedImageIds }),
      });
      const result = await response.json();

      if (result.success) {
        notifySuccess(t("images.delete_success", { count: result.count || 1 }));
        selectedImageIds = [];
      } else {
        notifyError(t("images.delete_error", { error: result.error || "Unknown" }));
      }
      await invalidateAll();
      triggerRefresh();
    } catch (e: any) {
      notifyError(e.message);
    }
  }

  // Clear history
  async function handleClearHistory() {
    try {
      const response = await fetch("?/clearHistory", {
        method: "POST",
      });
      const result = deserialize(await response.text()) as any;
      await invalidateAll();
    } catch (e) {}
  }

  // Create container
  async function handleCreateContainer(config: any) {
    notifySuccess(`Criando container '${config.containerName}'...`);

    const formData = new FormData();
    formData.append("config", JSON.stringify(config));

    try {
      const response = await fetch("?/createContainer", {
        method: "POST",
        body: formData,
      });
      const result = deserialize(await response.text()) as any;

      if (result.type === "success" && result.data?.success) {
        notifySuccess(result.data.message);
      } else {
        const errorMsg =
          result.type === "success" && result.data?.message
            ? result.data.message
            : result.type === "failure" && result.data?.message
              ? result.data.message
              : "Erro ao iniciar container";
        notifyError(errorMsg);
      }
    } catch (e: any) {
      notifyError(e.message);
    }
  }

  // Save profile config
  async function handleSaveProfile(profile: any) {
    const formData = new FormData();
    formData.append("config", JSON.stringify(profile));

    try {
      const response = await fetch("?/saveProfile", {
        method: "POST",
        body: formData,
      });
      const result = deserialize(await response.text()) as any;

      if (result.type === "success" && result.data?.success) {
        notifySuccess(result.data.message);
      } else {
        const errorMsg =
          result.type === "success" && result.data?.message
            ? result.data.message
            : result.type === "failure" && result.data?.message
              ? result.data.message
              : "Erro ao salvar perfil";
        notifyError(errorMsg);
      }
      const keepConfigOpen = showConfig;
      await invalidateAll();
      await tick();
      showConfig = keepConfigOpen;
    } catch (e) {}
  }

  // Delete profile config
  async function handleDeleteProfile(profileId: string) {
    const formData = new FormData();
    formData.append("profileId", profileId);

    try {
      const response = await fetch("?/deleteProfile", {
        method: "POST",
        body: formData,
      });
      const result = deserialize(await response.text()) as any;
      const keepConfigOpen = showConfig;
      await invalidateAll();
      await tick();
      showConfig = keepConfigOpen;
    } catch (e) {}
  }

  function openConfig(imageRepo: string, imageTag: string) {
    configTargetImage = `${imageRepo}:${imageTag}`;
    showConfig = true;
  }

  // Transfer helpers
  async function fetchSourceImages() {
    if (!transferSourceId) {
      sourceImages = [];
      return;
    }
    sourceLoading = true;
    selectedTransferIds = [];
    try {
      const res = await fetch(`/api/images?serverId=${transferSourceId}`);
      const json = await res.json();
      sourceImages = json.images || [];
    } catch (e: any) {
      notifyError(`Error loading source images: ${e.message}`);
      sourceImages = [];
    } finally {
      sourceLoading = false;
    }
    // Aquecer conexão em background
    warmupConnection(transferSourceId);
  }

  /**
   * Dispara uma conexão leve com o servidor via `docker.ping()`
   * para que o túnel SSH já fique estabelecido e em cache.
   */
  async function warmupConnection(serverId: string) {
    if (!serverId) return;
    try {
      const formData = new FormData();
      formData.append("serverId", serverId);
      await fetch("/api/images/transfer/warmup", {
        method: "POST",
        body: formData,
      });
    } catch {
      // silencioso — é apenas um warmup
    }
  }

  function toggleTransferImage(id: string) {
    if (selectedTransferIds.includes(id)) {
      selectedTransferIds = selectedTransferIds.filter((x) => x !== id);
    } else {
      selectedTransferIds = [...selectedTransferIds, id];
    }
  }

  function toggleAllTransfer() {
    if (selectedTransferIds.length === sourceImages.length) {
      selectedTransferIds = [];
    } else {
      selectedTransferIds = sourceImages.map((img) => img.id);
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

    transferInProgress = true;
    transferElapsedSeconds = 0;
    const timer = setInterval(() => {
      transferElapsedSeconds += 1;
    }, 1000);

    notifySuccess(
      t("images.transfer_in_progress").replace(
        "{count}",
        String(selectedTransferIds.length),
      ),
    );

    try {
      const formData = new FormData();
      formData.append("sourceId", transferSourceId);
      formData.append("destId", transferDestId);
      formData.append("imageIds", selectedTransferIds.join(","));

      const response = await fetch("/api/images/transfer", {
        method: "POST",
        body: formData,
      });
      const result = await response.json();

      if (result.success) {
        notifySuccess(
          t("images.transfer_success").replace(
            "{count}",
            String(result.count || selectedTransferIds.length),
          ),
        );
        selectedTransferIds = [];
      } else {
        notifyError(
          t("images.transfer_error").replace(
            "{error}",
            result.error || result.message || "Unknown",
          ),
        );
      }
    } catch (e: any) {
      notifyError(t("images.transfer_error").replace("{error}", e.message));
    } finally {
      clearInterval(timer);
      transferInProgress = false;
    }
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
          {t("images.title")}
        </h1>
      </div>

      <!-- Tabs Switcher -->
      <div
        class="flex items-center gap-1.5 p-1 bg-slate-100 dark:bg-slate-900 border border-slate-200/50 dark:border-slate-800 rounded-xl"
      >
        <button
          type="button"
          class="px-4 py-2 text-xs font-bold rounded-lg cursor-pointer transition-all {activeTab ===
          'my_images'
            ? 'bg-white dark:bg-[#0b0f19] text-violet-600 dark:text-violet-400 shadow-sm'
            : 'text-slate-500 hover:text-slate-850'}"
          onclick={() => (activeTab = "my_images")}
        >
          {t("images.my_images_tab")}
        </button>
        <button
          type="button"
          class="px-4 py-2 text-xs font-bold rounded-lg cursor-pointer transition-all {activeTab ===
          'download'
            ? 'bg-white dark:bg-[#0b0f19] text-violet-600 dark:text-violet-400 shadow-sm'
            : 'text-slate-500 hover:text-slate-850'}"
          onclick={() => (activeTab = "download")}
        >
          {t("images.download_tab")}
        </button>
        <button
          type="button"
          class="px-4 py-2 text-xs font-bold rounded-lg cursor-pointer transition-all {activeTab ===
          'transfer'
            ? 'bg-white dark:bg-[#0b0f19] text-violet-600 dark:text-violet-400 shadow-sm'
            : 'text-slate-500 hover:text-slate-850'}"
          onclick={() => (activeTab = "transfer")}
        >
          {t("images.transfer_tab")}
        </button>
      </div>
    </div>

    <!-- Status Alert Banner -->
    <StatusBanner />

    {#if loading}
      <DockerseaLoading message={t("common.loading")} />
    {:else if fetchError}
    <div
      class="p-6 rounded-2xl border border-red-200 dark:border-red-900/60 bg-red-50 dark:bg-red-950/20 text-red-700 dark:text-red-400"
    >
      <h3 class="font-bold text-sm mb-1">Erro de Conexão</h3>
      <p class="text-xs whitespace-pre-wrap">{fetchError}</p>
    </div>
  {:else}
    <!-- Tab 1: My Images -->
    {#if activeTab === "my_images"}
      <div class="space-y-4">
        <!-- Toolbar & Quick Filters -->
        <div
          class="flex flex-col gap-3 bg-white dark:bg-[#0b0f19] border border-slate-200/80 dark:border-slate-800/80 p-4 rounded-2xl shadow-sm"
        >
          <!-- Top Row: Disk usage stats & Main actions -->
          <div class="flex flex-wrap items-center justify-between gap-3">
            <div class="flex items-center gap-3">
              <div class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-violet-500/10 border border-violet-500/20 text-violet-600 dark:text-violet-400 text-xs font-bold">
                <span>📊 {t("images.disk_space")}</span>
                <span class="font-mono text-sm">{totalDiskUsage()}</span>
              </div>
              <span class="text-xs text-slate-400">{t("images.registered_images", { count: String(images.length) })}</span>
            </div>

            <div class="flex items-center gap-2">
              <button
                type="button"
                class="px-3.5 py-2 rounded-xl border-none cursor-pointer text-xs font-bold text-white bg-amber-500 hover:bg-amber-600 transition-colors shadow-md shadow-amber-500/20"
                onclick={toggleAll}
              >
                {selectedImageIds.length === filteredImages.length
                  ? t("common.deselect_all")
                  : t("common.select_all")}
              </button>
              {#if selectedImageIds.length > 0}
                <button
                  type="button"
                  class="px-3.5 py-2 rounded-xl border-none cursor-pointer text-xs font-bold text-white bg-red-500 hover:bg-red-600 transition-colors shadow-md shadow-red-500/20"
                  onclick={handleDeleteSelected}
                >
                  {t("images.delete_selected")} ({selectedImageIds.length})
                </button>
              {/if}
            </div>
          </div>

          <!-- Bottom Row: Filter chips & Search input -->
          <div class="flex flex-wrap items-center justify-between gap-3 pt-2 border-t border-slate-100 dark:border-slate-800/50">
            <div class="flex items-center gap-1.5 flex-wrap">
              <button
                type="button"
                class="px-3 py-1.5 text-xs font-bold rounded-xl border transition-all cursor-pointer {imageFilter === 'all'
                  ? 'bg-violet-600 text-white border-violet-600 shadow-md shadow-violet-500/20'
                  : 'bg-slate-100 dark:bg-slate-900 text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-800 hover:border-violet-500'}"
                onclick={() => (imageFilter = "all")}
              >
                {t("images.filter_all", { count: String(images.length) })}
              </button>
              <button
                type="button"
                class="px-3 py-1.5 text-xs font-bold rounded-xl border transition-all cursor-pointer {imageFilter === 'in_use'
                  ? 'bg-emerald-600 text-white border-emerald-600 shadow-md shadow-emerald-500/20'
                  : 'bg-slate-100 dark:bg-slate-900 text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-800 hover:border-emerald-500'}"
                onclick={() => (imageFilter = "in_use")}
              >
                {t("images.filter_in_use", { count: String(countInUse) })}
              </button>
              <button
                type="button"
                class="px-3 py-1.5 text-xs font-bold rounded-xl border transition-all cursor-pointer {imageFilter === 'unused'
                  ? 'bg-blue-600 text-white border-blue-600 shadow-md shadow-blue-500/20'
                  : 'bg-slate-100 dark:bg-slate-900 text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-800 hover:border-blue-500'}"
                onclick={() => (imageFilter = "unused")}
              >
                {t("images.filter_unused", { count: String(countUnused) })}
              </button>
              {#if countDangling > 0}
                <button
                  type="button"
                  class="px-3 py-1.5 text-xs font-bold rounded-xl border transition-all cursor-pointer {imageFilter === 'dangling'
                    ? 'bg-amber-600 text-white border-amber-600 shadow-md shadow-amber-500/20'
                    : 'bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/30 hover:border-amber-500'}"
                  onclick={() => (imageFilter = "dangling")}
                >
                  {t("images.filter_dangling", { count: String(countDangling) })}
                </button>
              {/if}
            </div>

            <input
              type="text"
              placeholder={t("common.search")}
              class="px-4 py-2 text-xs rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900 text-slate-850 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:border-violet-500 transition-all w-60"
              bind:value={searchQuery}
            />
          </div>
        </div>

        <!-- Images Cards Grid -->
        {#if filteredImages.length === 0}
          <div
            class="text-sm text-slate-400 dark:text-slate-500 py-16 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-slate-50/20"
          >
            {t("images.empty_images")}
          </div>
        {:else}
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {#each filteredImages as img (img.id + '_' + img.repo + '_' + img.tag)}
              <ImageCard
                {img}
                checked={selectedImageIds.includes(img.id)}
                on_toggle={() => toggleChecked(img.id)}
                on_build={() => openConfig(img.repo, img.tag)}
                on_delete={() => {
                  selectedImageIds = [img.id];
                  handleDeleteSelected();
                }}
              />
            {/each}
          </div>
        {/if}
      </div>
    {/if}

    <!-- Tab 2: Download -->
    {#if activeTab === "download"}
      <div class="space-y-6">
        <!-- Search/Pull Form -->
        <div
          class="bg-white dark:bg-[#0b0f19] border border-slate-200/80 dark:border-slate-800/80 p-5 rounded-2xl shadow-sm flex flex-col gap-4"
        >
          <label
            for="pull-image-input"
            class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
            >{t("images.download_title")}</label
          >
          <div class="flex gap-2">
            <input
              id="pull-image-input"
              type="text"
              placeholder="nginx:latest"
              class="flex-1 px-4 py-2.5 text-sm rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 placeholder-slate-400 focus:outline-none focus:border-violet-500 transition-all font-mono"
              bind:value={downloadQuery}
              onkeydown={(e) => e.key === "Enter" && handlePull(downloadQuery)}
            />
            <button
              type="button"
              class="px-5 py-2.5 rounded-xl border-none cursor-pointer text-sm font-bold text-white bg-violet-600 hover:bg-violet-700 transition-colors shadow-md shadow-violet-500/20 shrink-0"
              onclick={() => handlePull(downloadQuery)}
            >
              {t("images.pull_btn")}
            </button>
          </div>
        </div>

        <!-- History -->
        <div
          class="bg-white dark:bg-[#0b0f19] border border-slate-200/80 dark:border-slate-800/80 p-5 rounded-2xl shadow-sm space-y-4"
        >
          <div
            class="flex flex-wrap items-center justify-between gap-3 pb-3 border-b border-slate-100 dark:border-slate-900"
          >
            <div class="flex items-center gap-3">
              {#if data.imageHistory && data.imageHistory.length > 0}
                <button
                  type="button"
                  class="px-3 py-1.5 rounded-xl border-none cursor-pointer text-xs font-bold text-white bg-amber-500 hover:bg-amber-600 transition-colors shadow-md shadow-amber-500/20"
                  onclick={toggleAllHistory}
                >
                  {selectedHistoryIds.length === data.imageHistory.length
                    ? t("common.deselect_all")
                    : t("common.select_all")}
                </button>
                {#if selectedHistoryIds.length > 0}
                  <span
                    class="text-xs font-semibold text-red-500 px-1 animate-pulse"
                  >
                    {selectedHistoryIds.length}
                    {t("images.selected_count")}
                  </span>
                {/if}
              {/if}
              <span
                class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
                >{t("images.history_title")}</span
              >
            </div>

            <div class="flex items-center gap-2">
              {#if selectedHistoryIds.length > 0}
                <button
                  type="button"
                  class="text-xs text-red-500 hover:text-red-650 font-bold bg-transparent border-none cursor-pointer transition-colors"
                  onclick={handleDeleteHistorySelected}
                >
                  {t("common.delete")}
                </button>
              {/if}
              {#if data.imageHistory && data.imageHistory.length > 0}
                <button
                  type="button"
                  class="text-xs text-red-500 hover:text-red-650 font-bold bg-transparent border-none cursor-pointer transition-colors"
                  onclick={handleClearHistory}
                >
                  {t("images.clear_history")}
                </button>
              {/if}
            </div>
          </div>

          {#if !data.imageHistory || data.imageHistory.length === 0}
            <div class="text-xs text-slate-400 italic text-center py-6">
              {t("images.empty_history")}
            </div>
          {:else}
            <div class="divide-y divide-slate-100 dark:divide-slate-900">
              {#each data.imageHistory as hist (hist.id)}
                {@const isHistChecked = selectedHistoryIds.includes(hist.id)}
                <div
                  class="flex items-center justify-between py-3 text-xs gap-3"
                >
                  <div class="flex items-center gap-2.5 min-w-0 flex-1">
                    <button
                      type="button"
                      class="w-4 h-4 rounded border flex items-center justify-center cursor-pointer transition-all duration-150 shrink-0 {isHistChecked
                        ? 'bg-violet-600 border-violet-600 text-white'
                        : 'border-slate-300 dark:border-slate-700 bg-transparent hover:border-violet-500'}"
                      onclick={() => toggleCheckedHistory(hist.id)}
                    >
                      {#if isHistChecked}
                        <span
                          class="text-white text-[9px] font-bold leading-none"
                          >✓</span
                        >
                      {/if}
                    </button>
                    <div class="flex flex-col min-w-0 gap-0.5">
                      <span
                        class="font-bold text-slate-800 dark:text-slate-200 font-mono truncate"
                        >{hist.imageName}</span
                      >
                      <span
                        class="text-[10px] text-slate-400 dark:text-slate-500"
                      >
                        {t("images.downloaded_at")}{new Date(
                          hist.pulledAt,
                        ).toLocaleString()}
                      </span>
                    </div>
                  </div>
                  <button
                    type="button"
                    class="px-3 py-1.5 rounded-lg border border-slate-200 dark:border-slate-800 text-xs font-semibold text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 cursor-pointer transition-colors shrink-0 bg-white dark:bg-slate-900"
                    onclick={() => handlePull(hist.imageName)}
                  >
                    {t("images.repull_btn")}
                  </button>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      </div>
    {/if}

    <!-- Tab 3: Transfer -->
    {#if activeTab === "transfer"}
      <div class="space-y-6">
        {#if !data.servers || data.servers.length < 2}
          <div
            class="flex flex-col items-center justify-center py-16 px-4 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-slate-50/50 dark:bg-slate-900/10"
          >
            <h3
              class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-1"
            >
              {t("images.transfer_no_servers")}
            </h3>
            <p class="text-sm text-slate-500 dark:text-slate-400">
              <a
                href="/config"
                class="text-violet-600 dark:text-violet-400 font-semibold hover:underline"
                >{t("sidebar.configs")}</a
              >
            </p>
          </div>
        {:else}
          <!-- Server Selection Row -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Source -->
            <div
              class="bg-white dark:bg-[#0b0f19] border border-slate-200/80 dark:border-slate-800/80 p-5 rounded-2xl shadow-sm space-y-3"
            >
              <label
                for="transfer-source-select"
                class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider flex items-center gap-2"
              >
                <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
                {t("images.transfer_source")}
              </label>
              <select
                id="transfer-source-select"
                class="w-full px-4 py-2.5 text-sm rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:outline-none focus:border-violet-500 transition-all"
                bind:value={transferSourceId}
                onchange={fetchSourceImages}
              >
                <option value="">{t("images.transfer_select_source")}</option>
                {#each data.servers as server (server.id)}
                  <option value={server.id}>{server.name}</option>
                {/each}
              </select>
            </div>

            <!-- Destination -->
            <div
              class="bg-white dark:bg-[#0b0f19] border border-slate-200/80 dark:border-slate-800/80 p-5 rounded-2xl shadow-sm space-y-3"
            >
              <label
                for="transfer-dest-select"
                class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider flex items-center gap-2"
              >
                <span class="w-2 h-2 rounded-full bg-blue-500"></span>
                {t("images.transfer_dest")}
              </label>
              <select
                id="transfer-dest-select"
                class="w-full px-4 py-2.5 text-sm rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:outline-none focus:border-violet-500 transition-all"
                bind:value={transferDestId}
              >
                <option value="">{t("images.transfer_select_dest")}</option>
                {#each data.servers as server (server.id)}
                  <option value={server.id}>{server.name}</option>
                {/each}
              </select>
            </div>
          </div>

          <!-- Source Images -->
          {#if transferSourceId}
            <div
              class="bg-white dark:bg-[#0b0f19] border border-slate-200/80 dark:border-slate-800/80 p-5 rounded-2xl shadow-sm space-y-4"
            >
              <div class="flex items-center justify-between gap-3">
                <span
                  class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
                >
                  {t("images.transfer_select_images")}
                </span>
                {#if sourceImages.length > 0}
                  <button
                    type="button"
                    class="px-3 py-1.5 rounded-xl border-none cursor-pointer text-xs font-bold text-white bg-amber-500 hover:bg-amber-600 transition-colors shadow-md shadow-amber-500/20"
                    onclick={toggleAllTransfer}
                  >
                    {selectedTransferIds.length === sourceImages.length
                      ? t("common.deselect_all")
                      : t("common.select_all")}
                  </button>
                {/if}
              </div>

              {#if sourceLoading}
                <DockerseaLoading
                  message={t("images.transfer_loading_images")}
                />
              {:else if sourceImages.length === 0}
                <div
                  class="text-sm text-slate-400 dark:text-slate-500 py-8 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-2xl"
                >
                  {t("images.transfer_no_images")}
                </div>
              {:else}
                <div
                  class="divide-y divide-slate-100 dark:divide-slate-900 max-h-80 overflow-y-auto"
                >
                  {#each sourceImages as img (img.id)}
                    {@const isChecked = selectedTransferIds.includes(img.id)}
                    <div class="flex items-center gap-3 py-2.5 text-xs">
                      <button
                        type="button"
                        class="w-5 h-5 rounded-lg border-2 flex items-center justify-center cursor-pointer transition-all duration-150 shrink-0 {isChecked
                          ? 'bg-violet-600/20 border-violet-500 text-violet-500'
                          : 'border-slate-300 dark:border-slate-700 bg-transparent hover:border-violet-500'}"
                        onclick={() => toggleTransferImage(img.id)}
                      >
                        {#if isChecked}
                          <span
                            class="text-violet-500 text-xs font-bold leading-none"
                            >✓</span
                          >
                        {/if}
                      </button>
                      <div class="flex items-center gap-2 min-w-0 flex-1">
                        <span
                          class="font-mono font-bold text-slate-800 dark:text-slate-200 truncate"
                          >{img.repo}:{img.tag}</span
                        >
                        <span class="text-slate-400 shrink-0">({img.size})</span
                        >
                      </div>
                      <span
                        class="text-[10px] font-mono text-slate-400 dark:text-slate-500 shrink-0"
                        >{img.id}</span
                      >
                    </div>
                  {/each}
                </div>

                {#if selectedTransferIds.length > 0}
                  <div
                    class="flex items-center justify-between pt-2 border-t border-slate-100 dark:border-slate-900"
                  >
                    <span class="text-xs font-semibold text-violet-500">
                      {selectedTransferIds.length}
                      {t("images.selected_count")}
                    </span>
                    <button
                      type="button"
                      class="px-5 py-2.5 rounded-xl border-none cursor-pointer text-sm font-bold text-white bg-violet-600 hover:bg-violet-700 transition-colors shadow-md shadow-violet-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
                      disabled={!transferDestId ||
                        transferSourceId === transferDestId ||
                        transferInProgress}
                      onclick={handleTransfer}
                    >
                      {#if transferInProgress}
                        {t("common.loading")}
                      {:else}
                        {t("images.transfer_btn")}
                      {/if}
                    </button>
                  </div>
                {/if}
              {/if}
            </div>
          {/if}
        {/if}
      </div>
    {/if}
  {/if}
</div>
{/if}

<!-- Logs terminal modal -->
<TerminalModal
  bind:show={showTerminal}
  title={terminalTitle}
  loading={terminalLoading}
  logs={terminalLogs}
  console_id="pull-image-terminal-console"
/>

<!-- Config Modal to start container -->
<ConfigModal
  bind:show={showConfig}
  image={configTargetImage}
  savedConfigs={data.savedConfigs}
  serverId={data.activeVps?.id || ""}
  onsubmit={handleCreateContainer}
  onsaveprofile={handleSaveProfile}
  ondeleteprofile={handleDeleteProfile}
/>

<!-- Pull Progress modal -->
<PullProgressModal
  bind:show={showPullProgress}
  imageName={pullImageTargetName}
  oncomplete={async () => {
    notifySuccess("Processo de download finalizado!");
    await invalidateAll();
    triggerRefresh();
  }}
/>

<!-- Modal de Carregamento da Transferência -->
{#if transferInProgress}
  <div
    class="fixed inset-0 bg-black/60 backdrop-blur-md flex items-center justify-center z-50 p-4 animate-fadeIn"
  >
    <div
      class="bg-white dark:bg-[#0c1220] border border-slate-200 dark:border-slate-800 rounded-3xl p-8 shadow-2xl max-w-md w-full flex flex-col items-center text-center gap-5"
    >
      <DockerseaLoading message={t("images.transfer_modal_sending")} />
      
      <div class="flex flex-col gap-2.5 w-full pt-3 border-t border-slate-100 dark:border-slate-800/80">
        <span class="text-xs font-bold text-violet-600 dark:text-violet-400">
          {t("images.transfer_in_progress", { count: selectedTransferIds.length })}
        </span>

        <!-- Indicador visual de progresso de transmissão de dados -->
        <div class="flex items-center justify-between text-[11px] font-mono px-3 py-2 bg-slate-50 dark:bg-[#070a12] rounded-xl border border-slate-200/80 dark:border-slate-800">
          <div class="flex items-center gap-1.5 text-emerald-600 dark:text-emerald-400">
            <span class="animate-pulse">⚡</span>
            <span class="font-bold">Streaming Docker</span>
          </div>
          <span class="text-slate-500 dark:text-slate-400 font-semibold">
            {transferElapsedSeconds}s decorridos
          </span>
        </div>
      </div>
    </div>
  </div>
{/if}
