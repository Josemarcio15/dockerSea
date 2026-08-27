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
            class="text-lg font-bold text-slate-850 dark:text-white flex items-center gap-2"
          >
            {t("config.servers_title")}
          </h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
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
              ? 'border-violet-500 bg-violet-500/5 text-violet-750 dark:text-violet-300'
              : 'border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-900/30 text-slate-700 dark:text-slate-300'}"
            onclick={() => doChangeLocale("pt-BR")}
          >
            <span>🇧🇷 Português (Brasil)</span>
            {#if activeLocale === "pt-BR"}✓{/if}
          </button>

          <button
            type="button"
            class="w-full flex items-center justify-between p-3 rounded-xl border text-sm font-semibold transition-all cursor-pointer {activeLocale ===
            'en-US'
              ? 'border-violet-500 bg-violet-500/5 text-violet-750 dark:text-violet-300'
              : 'border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-900/30 text-slate-700 dark:text-slate-300'}"
            onclick={() => doChangeLocale("en-US")}
          >
            <span>🇺🇸 English (US)</span>
            {#if activeLocale === "en-US"}✓{/if}
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

<!-- Modal de Confirmação de Exclusão -->
<ConfirmDialog
  bind:show={store.showDeleteConfirm}
  title="Remover VPS"
  message={`Tem certeza de que deseja remover o servidor '${serverToDelete?.name || ""}'?\nEssa ação não poderá ser desfeita.`}
  confirmText="Remover Servidor"
  type="danger"
  onConfirm={confirmDeleteVps}
/>
