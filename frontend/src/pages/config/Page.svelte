<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import StatusBanner from "$shared/components/StatusBanner.svelte";
  import { createConfigStore } from "./store.svelte";

  import VpsCard from "./VpsCard.svelte";
  import VpsModal from "./VpsModal.svelte";
  import DiagnosticModal from "./DiagnosticModal.svelte";
  import ConfirmDialog from "$shared/components/ConfirmDialog.svelte";
  import { ButtonPurple } from "$shared/components/buttons";
  import type { VpsFormData, VpsServer } from "./types";

  let { data } = $props();

  const store = createConfigStore(() => data);
  const formData = $derived(store.form);
  const showVpsModal = $derived(store.showVpsModal);
  const showDiagnosticModal = $derived(store.showDiagnosticModal);
  const showDeleteConfirm = $derived(store.showDeleteConfirm);
  const activeLocale = $derived(store.activeLocale);

  // Carregar servidores do SQLite ao abrir a tela
  $effect(() => {
    void store.load();
  });

  const diagnosticTitle = $derived(store.diagnosticTitle);
  const diagnosticLoading = $derived(store.diagnosticLoading);
  const diagnosticResult = $derived(store.diagnosticResult);
  const serverToDelete = $derived(store.serverToDelete);
  const openCreateModal = () => store.openCreate();

  const openEditModal = (server: VpsServer) => store.openEdit(server);

  async function runDiagnostic(vpsData?: VpsFormData) {
    await store.test(vpsData);
  }

  async function handleSaveVps(savedForm: VpsFormData) {
    if (!savedForm.name.trim()) return;

    await store.save(savedForm);
  }

  // Delete VPS confirmation
  const handleDeleteVps = (server: VpsServer) => store.requestDelete(server);

  async function confirmDeleteVps() {
    await store.confirmDelete();
  }

  async function handleSelectVps(server: VpsServer) {
    await store.select(server);
  }

  function doChangeLocale(locale: string) {
    store.changeLocale(locale);
  }
</script>

<div class="space-y-8">
  <!-- Top Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div
      class="inline-flex items-center px-5 py-2.5 rounded-2xl bg-linear-to-br from-violet-100 to-fuchsia-100 dark:from-violet-950/40 dark:to-fuchsia-950/40 border border-violet-200/50 dark:border-violet-800/50 self-start shadow-sm"
    >
      <h1
        class="text-2xl font-bold bg-linear-to-r from-violet-600 to-fuchsia-500 bg-clip-text text-transparent m-0 flex items-center gap-2"
      >
        {t("sidebar.configs")}
      </h1>
    </div>
  </div>

  <!-- Status Alerts -->
  <StatusBanner />

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 items-start">
    <!-- VPS Servers List (col-span-2) -->
    <div
      class="lg:col-span-2 bg-white dark:bg-[#0b0f19] border border-slate-200/70 dark:border-slate-800/80 p-6 rounded-2xl shadow-sm space-y-6"
    >
      <div
        class="flex flex-col sm:flex-row sm:items-center justify-between gap-4"
      >
        <div>
          <h2
            class="text-lg font-bold text-slate-900 dark:text-white flex items-center gap-2"
          >
            {t("config.servers_title")}
          </h2>
          <p class="text-xs text-slate-600 dark:text-slate-400 mt-1">
            {t("config.servers_desc")}
          </p>
        </div>

        <ButtonPurple onclick={openCreateModal}>
          <span class="text-sm font-normal">+</span>
          {t("config.add_server_btn")}
        </ButtonPurple>
      </div>

      <!-- Servers List (Rows) -->
      <div class="flex flex-col gap-3">
        {#each data.servers as server (server.id)}
          <VpsCard
            {server}
            isActive={data.activeVps?.id === server.id}
            onSelect={handleSelectVps}
            onTest={runDiagnostic}
            onEdit={openEditModal}
            onDelete={handleDeleteVps}
          />
        {:else}
          <div
            class="text-center py-12 border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-slate-50/30 italic text-slate-500"
          >
            {t("config.empty_servers")}
          </div>
        {/each}
      </div>
    </div>

    <!-- General Settings (col-span-1) -->
    <div class="space-y-6">
      <!-- Language Card -->
      <div
        class="bg-white dark:bg-[#0b0f19] border border-slate-200/70 dark:border-slate-800/80 p-5 rounded-2xl shadow-sm space-y-4"
      >
        <h3
          class="text-sm font-bold text-slate-800 dark:text-slate-200 uppercase tracking-wider flex items-center gap-2"
        >
          {t("config.lang_tab")}
        </h3>

        <div class="flex flex-col gap-2">
          <button
            type="button"
            class="w-full flex items-center justify-between p-3 rounded-xl border text-sm font-semibold transition-all cursor-pointer {activeLocale ===
            'pt-BR'
              ? 'border-violet-500 bg-violet-50 dark:bg-violet-500/10 text-violet-900 dark:text-violet-300 font-bold shadow-xs'
              : 'border-slate-300 dark:border-slate-800 bg-white dark:bg-transparent hover:bg-slate-100 dark:hover:bg-slate-900/30 text-slate-800 dark:text-slate-300'}"
            onclick={() => doChangeLocale("pt-BR")}
          >
            <span>{t("config.lang_pt_br")}</span>
            {#if activeLocale === "pt-BR"}
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-violet-600 dark:text-violet-400" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 0 1 0 1.414l-8 8a1 1 0 0 1-1.414 0l-4-4a1 1 0 0 1 1.414-1.414L8 12.586l7.293-7.293a1 1 0 0 1 1.414 0Z" clip-rule="evenodd" />
              </svg>
            {/if}
          </button>

          <button
            type="button"
            class="w-full flex items-center justify-between p-3 rounded-xl border text-sm font-semibold transition-all cursor-pointer {activeLocale ===
            'en-US'
              ? 'border-violet-500 bg-violet-50 dark:bg-violet-500/10 text-violet-900 dark:text-violet-300 font-bold shadow-xs'
              : 'border-slate-300 dark:border-slate-800 bg-white dark:bg-transparent hover:bg-slate-100 dark:hover:bg-slate-900/30 text-slate-800 dark:text-slate-300'}"
            onclick={() => doChangeLocale("en-US")}
          >
            <span>{t("config.lang_en_us")}</span>
            {#if activeLocale === "en-US"}
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-violet-600 dark:text-violet-400" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 0 1 0 1.414l-8 8a1 1 0 0 1-1.414 0l-4-4a1 1 0 0 1 1.414-1.414L8 12.586l7.293-7.293a1 1 0 0 1 1.414 0Z" clip-rule="evenodd" />
              </svg>
            {/if}
          </button>
        </div>
      </div>

      <!-- Local Database & Maintenance Card -->
      <div
        class="bg-white dark:bg-[#0b0f19] border border-slate-200/70 dark:border-slate-800/80 p-5 rounded-2xl shadow-sm space-y-4"
      >
        <h3
          class="text-sm font-bold text-slate-900 dark:text-slate-200 uppercase tracking-wider flex items-center gap-2"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-violet-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M20.25 6.375c0 2.278-3.694 4.125-8.25 4.125S3.75 8.653 3.75 6.375m16.5 0c0-2.278-3.694-4.125-8.25-4.125S3.75 4.097 3.75 6.375m16.5 0v11.25c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125V6.375m16.5 5.625c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125m16.5 5.625c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125" />
          </svg>
          {t("config.db_section_title")}
        </h3>
        <p class="text-xs text-slate-600 dark:text-slate-400">
          {t("config.db_section_desc")}
        </p>

        {#if store.dbPath}
          <div class="space-y-1">
            <span class="text-[11px] font-bold text-slate-600 dark:text-slate-400 uppercase">{t("config.db_path_label")}</span>
            <div class="p-2.5 rounded-xl bg-slate-100 dark:bg-slate-900 border border-slate-300 dark:border-slate-800 text-[11px] font-mono text-slate-800 dark:text-slate-400 break-all select-all font-semibold">
              {store.dbPath}
            </div>
          </div>
        {/if}

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
          <button
            type="button"
            class="flex items-center justify-center gap-2 py-2.5 px-3 rounded-xl text-xs font-semibold bg-violet-600 hover:bg-violet-700 text-white dark:bg-violet-500/10 dark:hover:bg-violet-500/20 dark:text-violet-300 border border-transparent dark:border-violet-500/30 transition-all cursor-pointer disabled:opacity-50 shadow-sm"
            disabled={store.isBackingUp || store.isRestoring}
            onclick={store.exportBackup}
          >
            {#if store.isBackingUp}
              <span class="animate-spin text-sm">↻</span> {t("config.db_exporting")}
            {:else}
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3" />
              </svg>
              {t("config.db_backup_btn")}
            {/if}
          </button>

          <button
            type="button"
            class="flex items-center justify-center gap-2 py-2.5 px-3 rounded-xl text-xs font-semibold bg-blue-500/10 hover:bg-blue-500/20 text-blue-700 dark:text-blue-300 border border-blue-500/30 transition-all cursor-pointer disabled:opacity-50"
            disabled={store.isBackingUp || store.isRestoring}
            onclick={store.restoreBackup}
          >
            {#if store.isRestoring}
              <span class="animate-spin text-sm">↻</span> {t("config.db_importing")}
            {:else}
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5m-13.5-9L12 3m0 0 4.5 4.5M12 3v13.5" />
              </svg>
              {t("config.db_restore_btn")}
            {/if}
          </button>
        </div>

        <!-- Danger Zone -->
        <div class="pt-3 border-t border-red-500/20 space-y-2">
          <span class="text-[11px] font-bold text-red-600 dark:text-red-400 uppercase tracking-wider flex items-center gap-1.5">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z" />
            </svg>
            {t("config.db_danger_title")}
          </span>
          <button
            type="button"
            class="w-full flex items-center justify-center gap-2 py-2.5 px-4 rounded-xl text-xs font-semibold bg-red-500/10 hover:bg-red-500/20 text-red-700 dark:text-red-400 border border-red-500/30 transition-all cursor-pointer disabled:opacity-50"
            disabled={store.isResetting}
            onclick={store.requestResetDb}
          >
            {#if store.isResetting}
              <span class="animate-spin text-sm">↻</span> {t("config.db_resetting")}
            {:else}
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
              </svg>
              {t("config.db_reset_btn")}
            {/if}
          </button>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- Modal: Criar / Editar VPS -->
<VpsModal
  bind:show={store.showVpsModal}
  bind:form={store.form}
  onSave={handleSaveVps}
  onTest={runDiagnostic}
  onAutoDetect={store.autoDetect}
/>

<!-- Modal: Diagnóstico de Conexão VPS -->
<DiagnosticModal
  bind:show={store.showDiagnosticModal}
  title={diagnosticTitle}
  loading={diagnosticLoading}
  result={diagnosticResult}
/>

<!-- Modal de Confirmação de Exclusão de VPS -->
<ConfirmDialog
  bind:show={store.showDeleteConfirm}
  title="Remover VPS"
  message={`Tem certeza de que deseja remover o servidor '${serverToDelete?.name || ""}'?\nEssa ação não poderá ser desfeita.`}
  confirmText="Remover Servidor"
  type="danger"
  onConfirm={confirmDeleteVps}
/>

<!-- Modal de Confirmação de Reset do Banco -->
<ConfirmDialog
  bind:show={store.showResetConfirm}
  title={t("config.db_reset_confirm_title")}
  message={t("config.db_reset_confirm_message")}
  confirmText={t("config.db_reset_confirm_action")}
  type="danger"
  onConfirm={store.confirmResetDb}
/>
