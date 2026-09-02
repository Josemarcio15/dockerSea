<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { useRefreshKey, triggerRefresh } from "$shared/stores/refresh.svelte";
  import { notifySuccess, notifyError } from "$shared/stores/notification.svelte";
  import * as ContainerService from "../../../bindings/go-walis/internal/containers/containerservice.js";
  import DockerseaLoading from "$shared/components/DockerseaLoading.svelte";
  import StatusBanner from "$shared/components/StatusBanner.svelte";
  import TerminalModal from "$shared/components/TerminalModal.svelte";
  import ConfigModal from "$shared/components/ConfigModal.svelte";
  import TaskProgressModal from "$shared/components/TaskProgressModal.svelte";
  import VpsSelectWarning from "$shared/components/VpsSelectWarning.svelte";

  import ImageToolbar from "./components/ImageToolbar.svelte";
  import ImageCard from "./components/ImageCard.svelte";
  import ConfirmDialog from "$shared/components/ConfirmDialog.svelte";
  import {
    ButtonBlue,
    ButtonGreen,
    ButtonYellow,
    ButtonPurple,
    ButtonCyan,
    ButtonRed,
    ButtonPink,
    ButtonOrange,
  } from "$shared/components/buttons";
  import { createImagesStore } from "./store.svelte";
  import { viewModeStore } from "$shared/stores/viewMode.svelte";

  let { data } = $props();
  const imgState = createImagesStore(
    () => data?.activeVps,
    () => data?.activeProfile?.id || "default",
    () => data?.servers || [],
  );

  $effect(() => {
    useRefreshKey();
    if (data?.activeVps?.id) {
      imgState.fetchImages();
      imgState.fetchHistory();
    }
  });

  $effect(() => {
    const profileId = data?.activeProfile?.id || "default";
    void imgState.fetchSavedConfigs();
  });

  async function handleCreateContainer(config: any) {
    if (!data.activeVps) return notifyError("Nenhum servidor selecionado.");
    const result = await ContainerService.CreateContainer(data.activeVps, config);
    if (!result?.success) return notifyError(result?.message || "Não foi possível criar o container.");
    notifySuccess(result.message || "Container criado com sucesso!");
    triggerRefresh();
  }

  async function handleSaveProfile(profile: any) {
    await imgState.saveProfile(profile);
  }

  async function handleDeleteProfile(id: string) {
    await imgState.deleteProfile(id);
  }

  // Deletion confirm states
  let showDeleteImagesConfirm = $state(false);
  let showDeleteSingleImageConfirm = $state(false);
  let showPruneImagesConfirm = $state(false);
  let imageToDelete = $state<any>(null);

  function requestDeleteSelectedImages() {
    if (imgState.selectedImageIds.length === 0) return;
    showDeleteImagesConfirm = true;
  }

  async function confirmDeleteSelectedImages() {
    await imgState.handleDeleteSelected();
  }

  function requestPruneUnusedImages() {
    showPruneImagesConfirm = true;
  }

  async function confirmPruneUnusedImages() {
    await imgState.handlePruneUnused();
  }

  function requestDeleteSingleImage(img: any) {
    imageToDelete = img;
    showDeleteSingleImageConfirm = true;
  }

  async function confirmDeleteSingleImage() {
    if (!imageToDelete) return;
    await imgState.handleDeleteSingle(imageToDelete.id);
  }
</script>

{#if !data.activeVps}
  <VpsSelectWarning />
{:else}
  <div class="space-y-6">
    <!-- Top Header -->
    <div
      class="flex flex-col sm:flex-row sm:items-center justify-between gap-4"
    >
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
        class="inline-flex rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-[#0b0f19] p-1 shadow-xs"
      >
        <button
          type="button"
          class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {imgState.activeTab ===
          'my_images'
            ? 'bg-violet-600 text-white shadow-xs'
            : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}"
          onclick={() => (imgState.activeTab = "my_images")}
        >
          {t("images.my_images_tab")}
        </button>
        <button
          type="button"
          class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {imgState.activeTab ===
          'download'
            ? 'bg-violet-600 text-white shadow-xs'
            : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}"
          onclick={() => (imgState.activeTab = "download")}
        >
          {t("images.download_tab")}
        </button>
        <button
          type="button"
          class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {imgState.activeTab ===
          'transfer'
            ? 'bg-violet-600 text-white shadow-xs'
            : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}"
          onclick={() => (imgState.activeTab = "transfer")}
        >
          {t("images.transfer_tab")}
        </button>
      </div>
    </div>

    <!-- Status Alert Banner -->
    <StatusBanner />

    {#if imgState.loading}
      <DockerseaLoading message={t("common.loading")} />
    {:else if imgState.fetchError}
      <div
        class="p-6 rounded-2xl border border-red-200 dark:border-red-900/60 bg-red-50 dark:bg-red-950/20 text-red-700 dark:text-red-400"
      >
        <h3 class="font-bold text-sm mb-1">Erro de Conexão</h3>
        <p class="text-xs whitespace-pre-wrap">{imgState.fetchError}</p>
      </div>
    {:else}
      <!-- Tab 1: My Images -->
      {#if imgState.activeTab === "my_images"}
        <div class="space-y-4">
          <ImageToolbar
            diskUsage={imgState.totalDiskUsage}
            totalCount={imgState.images.length}
            selectedCount={imgState.selectedImageIds.length}
            allSelected={imgState.selectedImageIds.length ===
              imgState.filteredImages.length &&
              imgState.filteredImages.length > 0}
            bind:imageFilter={imgState.imageFilter}
            bind:searchQuery={imgState.searchQuery}
            countInUse={imgState.countInUse}
            countUnused={imgState.countUnused}
            countDangling={imgState.countDangling}
            onToggleAll={imgState.toggleAll}
            onPrune={requestPruneUnusedImages}
            onDeleteSelected={requestDeleteSelectedImages}
          />

          <!-- Images Cards Grid / List -->
          {#if imgState.filteredImages.length === 0}
            <div
              class="text-sm text-slate-400 dark:text-slate-500 py-16 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-slate-50/20"
            >
              {t("images.empty_images")}
            </div>
          {:else}
            <div class={viewModeStore.getGridClass()}>
              {#each imgState.filteredImages as img (img.id + "_" + img.repo + "_" + img.tag)}
                <ImageCard
                  {img}
                  checked={imgState.selectedImageIds.includes(img.id)}
                  on_toggle={() => imgState.toggleChecked(img.id)}
                  on_build={() => imgState.openConfig(img.repo, img.tag)}
                />
              {/each}
            </div>
          {/if}
        </div>
      {/if}

      <!-- Tab 2: Download -->
      {#if imgState.activeTab === "download"}
        <div class="space-y-6">
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
                bind:value={imgState.downloadQuery}
                onkeydown={(e) =>
                  e.key === "Enter" &&
                  imgState.handlePull(imgState.downloadQuery)}
              />
              <ButtonGreen
                size="md"
                onclick={() => imgState.handlePull(imgState.downloadQuery)}
              >
                {t("images.pull_btn")}
              </ButtonGreen>
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
                {#if imgState.imageHistory && imgState.imageHistory.length > 0}
                  <ButtonYellow size="xs" onclick={imgState.toggleAllHistory}>
                    {imgState.selectedHistoryIds.length ===
                    imgState.imageHistory.length
                      ? t("common.deselect_all")
                      : t("common.select_all")}
                  </ButtonYellow>
                  {#if imgState.selectedHistoryIds.length > 0}
                    <span
                      class="text-xs font-semibold text-red-500 px-1 animate-pulse"
                    >
                      {imgState.selectedHistoryIds.length}
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
                {#if imgState.selectedHistoryIds.length > 0}
                  <ButtonRed
                    size="xs"
                    onclick={imgState.handleDeleteHistorySelected}
                  >
                    {t("common.delete")}
                  </ButtonRed>
                {/if}
                {#if imgState.imageHistory && imgState.imageHistory.length > 0}
                  <ButtonRed size="xs" onclick={imgState.handleClearHistory}>
                    {t("images.clear_history")}
                  </ButtonRed>
                {/if}
              </div>
            </div>

            {#if !imgState.imageHistory || imgState.imageHistory.length === 0}
              <div class="text-xs text-slate-400 italic text-center py-6">
                {t("images.empty_history")}
              </div>
            {:else}
              <div class="divide-y divide-slate-100 dark:divide-slate-900">
                {#each imgState.imageHistory as hist (hist.id)}
                  {@const isHistChecked = imgState.selectedHistoryIds.includes(
                    hist.id,
                  )}
                  <div
                    class="flex items-center justify-between py-3 text-xs gap-3"
                  >
                    <div class="flex items-center gap-2.5 min-w-0 flex-1">
                      <button
                        type="button"
                        class="w-4 h-4 rounded border flex items-center justify-center cursor-pointer transition-all duration-150 shrink-0 {isHistChecked
                          ? 'bg-violet-600 border-violet-600 text-white'
                          : 'border-slate-300 dark:border-slate-700 bg-transparent hover:border-violet-500'}"
                        onclick={() => imgState.toggleCheckedHistory(hist.id)}
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
                    <ButtonCyan
                      size="xs"
                      onclick={() => imgState.handlePull(hist.imageName)}
                    >
                      {t("images.repull_btn")}
                    </ButtonCyan>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        </div>
      {/if}

      <!-- Tab 3: Transfer -->
      {#if imgState.activeTab === "transfer"}
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
                  bind:value={imgState.transferSourceId}
                  onchange={imgState.fetchSourceImages}
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
                  bind:value={imgState.transferDestId}
                >
                  <option value="">{t("images.transfer_select_dest")}</option>
                  {#each data.servers as server (server.id)}
                    <option value={server.id}>{server.name}</option>
                  {/each}
                </select>
              </div>
            </div>

            <!-- Source Images -->
            {#if imgState.transferSourceId}
              <div
                class="bg-white dark:bg-[#0b0f19] border border-slate-200/80 dark:border-slate-800/80 p-5 rounded-2xl shadow-sm space-y-4"
              >
                <div class="flex items-center justify-between gap-3">
                  <span
                    class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
                  >
                    {t("images.transfer_select_images")}
                  </span>
                  {#if imgState.sourceImages.length > 0}
                    <ButtonYellow
                      size="xs"
                      onclick={imgState.toggleAllTransfer}
                    >
                      {imgState.selectedTransferIds.length ===
                      imgState.sourceImages.length
                        ? t("common.deselect_all")
                        : t("common.select_all")}
                    </ButtonYellow>
                  {/if}
                </div>

                {#if imgState.sourceLoading}
                  <DockerseaLoading
                    message={t("images.transfer_loading_images")}
                  />
                {:else if imgState.sourceImages.length === 0}
                  <div
                    class="text-sm text-slate-400 dark:text-slate-500 py-8 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-2xl"
                  >
                    {t("images.transfer_no_images")}
                  </div>
                {:else}
                  <div
                    class="divide-y divide-slate-100 dark:divide-slate-900 max-h-80 overflow-y-auto"
                  >
                    {#each imgState.sourceImages as img (imgState.getTransferImageKey(img))}
                      {@const itemKey = imgState.getTransferImageKey(img)}
                      {@const isChecked =
                        imgState.selectedTransferIds.includes(itemKey)}
                      <div class="flex items-center gap-3 py-2.5 text-xs">
                        <button
                          type="button"
                          class="w-5 h-5 rounded-lg border-2 flex items-center justify-center cursor-pointer transition-all duration-150 shrink-0 {isChecked
                            ? 'bg-violet-600/20 border-violet-500 text-violet-500'
                            : 'border-slate-300 dark:border-slate-700 bg-transparent hover:border-violet-500'}"
                          onclick={() => imgState.toggleTransferImage(itemKey)}
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
                          <span class="text-slate-400 shrink-0"
                            >({img.size})</span
                          >
                        </div>
                        <span
                          class="text-[10px] font-mono text-slate-400 dark:text-slate-500 shrink-0"
                          >{img.id}</span
                        >
                      </div>
                    {/each}
                  </div>

                  {#if imgState.selectedTransferIds.length > 0}
                    <div
                      class="flex items-center justify-between pt-2 border-t border-slate-100 dark:border-slate-900"
                    >
                      <span class="text-xs font-semibold text-violet-500">
                        {imgState.selectedTransferIds.length}
                        {t("images.selected_count")}
                      </span>
                      <ButtonPurple
                        size="md"
                        disabled={!imgState.transferDestId ||
                          imgState.transferSourceId ===
                            imgState.transferDestId ||
                          imgState.transferInProgress}
                        loading={imgState.transferInProgress}
                        onclick={imgState.handleTransfer}
                      >
                        {imgState.transferInProgress
                          ? t("images.transferring")
                          : t("images.transfer_action_btn")}
                      </ButtonPurple>
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
  bind:show={imgState.showTerminal}
  title={imgState.terminalTitle}
  loading={imgState.terminalLoading}
  logs={imgState.terminalLogs}
/>

<!-- Config Modal to start container -->
<ConfigModal
  bind:show={imgState.showConfig}
  image={imgState.configTargetImage}
  savedConfigs={imgState.savedConfigs}
  serverId={data.activeVps?.id || ""}
  activeVps={data.activeVps}
  onsubmit={handleCreateContainer}
  onsaveprofile={handleSaveProfile}
  ondeleteprofile={handleDeleteProfile}
/>

<!-- Modal Confirmação Limpar Não Usadas (Prune) -->
<ConfirmDialog
  bind:show={showPruneImagesConfirm}
  title="Limpar Imagens Não Utilizadas"
  message={`Tem certeza de que deseja remover todas as imagens não utilizadas (${imgState.countUnused + imgState.countDangling})?\nImagens associadas a contêineres ativos ou parados serão mantidas.`}
  confirmText="Limpar Imagens"
  type="danger"
  onConfirm={confirmPruneUnusedImages}
/>

<!-- Modal Confirmação Exclusão em Lote -->
<ConfirmDialog
  bind:show={showDeleteImagesConfirm}
  title="Remover Imagens"
  message={`Tem certeza de que deseja remover as ${imgState.selectedImageIds.length} imagem(ns) selecionada(s)?\nEssa ação liberará o espaço em disco correspondente.`}
  confirmText="Remover Imagens"
  type="danger"
  onConfirm={confirmDeleteSelectedImages}
/>

<!-- Modal Confirmação Exclusão Individual -->
<ConfirmDialog
  bind:show={showDeleteSingleImageConfirm}
  title="Remover Imagem"
  message={`Tem certeza de que deseja remover a imagem '${imageToDelete?.repo || ""}:${imageToDelete?.tag || ""}'?`}
  confirmText="Remover Imagem"
  type="danger"
  onConfirm={confirmDeleteSingleImage}
/>

<!-- Pull Progress modal -->
<TaskProgressModal
  bind:show={imgState.showPullProgress}
  title={imgState.pullImageTargetName ? `Download de Imagem: ${imgState.pullImageTargetName}` : "Download de Imagem"}
  eventPrefix="docker:image:pull"
  oncomplete={async () => {
    notifySuccess("Processo de download finalizado!");
    await imgState.fetchImages(true);
    await imgState.fetchHistory();
    triggerRefresh();
  }}
/>

<!-- Carregamento da Transferência com Progresso em Tempo Real Modal -->
{#if imgState.transferInProgress}
  <div
    class="fixed inset-0 bg-black/80 backdrop-blur-md flex items-center justify-center z-50 p-4 animate-fadeIn"
  >
    <div
      class="bg-[#0b101d] border border-slate-800 rounded-3xl p-6 sm:p-8 shadow-2xl max-w-lg w-full flex flex-col items-center text-center gap-5 text-slate-200"
    >
      <div
        class="w-12 h-12 rounded-2xl bg-linear-to-tr from-violet-600 to-fuchsia-600 flex items-center justify-center text-white text-xl shadow-lg shadow-violet-500/30 animate-pulse"
      >
        🚀
      </div>

      <div class="space-y-1">
        <h2 class="text-lg font-bold text-white tracking-wide">
          {t("images.transfer_tab")} entre Servidores
        </h2>
        <p class="text-xs text-slate-400 font-mono">
          {imgState.selectedTransferIds.length} imagem(ns) selecionada(s)
        </p>
      </div>

      <!-- progresso abstrata 0 - 100% Bar -->
      <div
        class="w-full bg-[#070a12] border border-slate-800 rounded-2xl p-5 space-y-3.5 shadow-inner text-left"
      >
        <div class="flex items-center justify-between text-xs font-bold">
          <span class="text-white text-sm">Progresso da Transferência</span>
          <span
            class="font-mono text-base font-bold bg-linear-to-r from-violet-400 to-emerald-400 bg-clip-text text-transparent"
          >
            {imgState.transferPercent}%
          </span>
        </div>

        <div
          class="w-full h-3 bg-slate-900 rounded-full overflow-hidden border border-slate-800 p-0.5"
        >
          <div
            class="h-full rounded-full transition-all duration-300 bg-linear-to-r from-violet-600 via-fuchsia-500 to-emerald-500 shadow-sm shadow-emerald-500/30"
            style="width: {imgState.transferPercent}%"
          ></div>
        </div>

        <div
          class="flex items-center justify-between text-xs font-mono text-slate-400 pt-1"
        >
          <span class="text-emerald-400 font-semibold">
            {imgState.transferSpeed !== "0 B/s"
              ? `⚡ ${imgState.transferSpeed}`
              : "⚡ Conectando stream..."}
          </span>
          <span class="font-bold text-white">
            {imgState.transferFormattedBytes}{imgState.transferFormattedTotalBytes
              ? ` / ${imgState.transferFormattedTotalBytes}`
              : ""}
          </span>
        </div>
      </div>

      <!-- Footer com tempo decorrido e status -->
      <div
        class="flex items-center justify-between w-full text-xs font-mono text-slate-400 pt-2 border-t border-slate-800/80"
      >
        <div class="flex items-center gap-1.5 text-emerald-400">
          <span>⚡</span>
          <span>Pipeline streaming sem buffer em disco</span>
        </div>
        <span class="font-semibold text-slate-300">
          {imgState.transferElapsedSeconds}s
        </span>
      </div>
    </div>
  </div>
{/if}
