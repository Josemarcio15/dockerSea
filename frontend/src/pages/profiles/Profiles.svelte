<script lang="ts">
  import StatusBanner from "$shared/components/StatusBanner.svelte";
  import { notifySuccess, notifyError } from "$shared/stores/notification.svelte";
  import { t } from "$shared/stores/locale.svelte";
  import { triggerRefresh } from "$shared/stores/refresh.svelte";
  import Modal from "$shared/components/Modal.svelte";
  import ConfirmDialog from "$shared/components/ConfirmDialog.svelte";
  import {
    ButtonBlue,
    ButtonPurple,
    ButtonRed,
    EditButtonIcon,
    TrashButtonIcon,
  } from "$shared/components/buttons";

  function deserialize(text: string): any {
    try {
      return JSON.parse(text);
    } catch {
      return { success: false, message: text };
    }
  }

  async function invalidateAll() {
    triggerRefresh();
  }

  let { data } = $props();

  // Profile form modal
  let showModal = $state(false);
  let modalId = $state("");
  let modalName = $state("");

  // Action runner
  async function runAction(action: string, formData: FormData) {
    try {
      const res = await fetch(`?/${action}`, {
        method: "POST",
        body: formData,
      });
      const result = deserialize(await res.text()) as any;

      let parsedData: any = null;
      if (result && typeof result === "object") {
        const rawData = result.data;
        if (rawData) {
          try {
            parsedData =
              typeof rawData === "string" ? JSON.parse(rawData) : rawData;
          } catch (e) {
            parsedData = rawData;
          }
        } else {
          parsedData = result;
        }
      }

      if (parsedData && parsedData.success) {
        notifySuccess(parsedData.message);
      } else {
        notifyError(
          parsedData?.message ||
            parsedData?.error ||
            result?.error?.message ||
            result?.error ||
            result?.message ||
            (result ? JSON.stringify(result).slice(0, 300) : "Resposta vazia"),
        );
      }

      await invalidateAll();
    } catch (e: any) {
      notifyError(e.message);
    }
  }

  async function doSave() {
    if (!modalName.trim()) return;

    const formData = new FormData();
    if (modalId) formData.append("id", modalId);
    formData.append("name", modalName);

    showModal = false;
    await runAction("save", formData);
  }

  // Delete Profile confirmation
  let showDeleteConfirm = $state(false);
  let profileToDelete = $state<{ id: string; name: string } | null>(null);

  function doDelete(id: string, name: string) {
    profileToDelete = { id, name };
    showDeleteConfirm = true;
  }

  async function confirmDeleteProfile() {
    if (!profileToDelete) return;
    const formData = new FormData();
    formData.append("id", profileToDelete.id);
    await runAction("delete", formData);
  }

  async function doSelect(id: string) {
    const formData = new FormData();
    formData.append("id", id);
    await runAction("select", formData);
  }

  function openCreateModal() {
    modalId = "";
    modalName = "";
    showModal = true;
  }

  function openEditModal(profile: any) {
    modalId = profile.id;
    modalName = profile.name;
    showModal = true;
  }
</script>

<div class="space-y-6">
  <!-- Top Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div
      class="inline-flex items-center px-5 py-2.5 rounded-2xl bg-linear-to-br from-violet-100 to-fuchsia-100 dark:from-violet-950/40 dark:to-fuchsia-950/40 border border-violet-200/50 dark:border-violet-800/50 self-start shadow-sm"
    >
      <h1
        class="text-2xl font-bold bg-linear-to-r from-violet-600 to-fuchsia-500 bg-clip-text text-transparent m-0 flex items-center gap-2"
      >
        {t("profiles.title")}
      </h1>
    </div>

    <ButtonBlue size="sm" onclick={openCreateModal}>
      {t("profiles.new_profile")}
    </ButtonBlue>
  </div>

  <!-- Status Alerts -->
  <StatusBanner />

  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
    {#each data.profiles as profile (profile.id)}
      <div
        class="flex flex-col justify-between p-5 rounded-2xl bg-white dark:bg-[#0c1220] border transition-all duration-300 shadow-md dark:shadow-lg dark:shadow-black/40 hover:-translate-y-1 hover:border-slate-300 dark:hover:border-slate-700 {data
          .activeProfile?.id === profile.id
          ? 'border-violet-500 ring-2 ring-violet-500/10'
          : 'border-slate-200 dark:border-slate-800/80'}"
      >
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            {#if data.activeProfile?.id === profile.id}
              <span
                class="px-2 py-0.5 rounded-full text-[9px] font-bold bg-emerald-500/15 text-emerald-500 border border-emerald-500/20"
                >{t("profiles.active_badge")}</span
              >
            {/if}
          </div>

          <div>
            <h3
              class="font-bold text-slate-800 dark:text-white text-base truncate"
            >
              {profile.name}
            </h3>
          </div>
        </div>

        <div
          class="flex items-center gap-2 mt-6 pt-4 border-t border-slate-100 dark:border-slate-800/60 shrink-0"
        >
          {#if data.activeProfile?.id !== profile.id}
            <ButtonPurple
              size="xs"
              class="flex-1"
              onclick={() => doSelect(profile.id)}
            >
              {t("profiles.select_btn")}
            </ButtonPurple>
          {/if}

          <ButtonPurple
            size="xs"
            onclick={() => openEditModal(profile)}
          >
            {t("common.edit")}
          </ButtonPurple>

          {#if data.profiles.length > 1}
            <ButtonRed
              size="xs"
              onclick={() => doDelete(profile.id, profile.name)}
            >
              {t("common.delete")}
            </ButtonRed>
          {/if}
        </div>
      </div>
    {/each}
  </div>
</div>

<!-- Modal: Adicionar/Editar Perfil -->
<Modal
  bind:show={showModal}
  title={modalId ? t("profiles.edit_title") : t("profiles.create_title")}
  buttons={[
    {
      label: t("common.save"),
      variant: "primary",
      onclick: doSave,
      disabled: !modalName.trim(),
    },
  ]}
>
  <div class="flex flex-col gap-1.5">
    <label
      for="profile-name"
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >
      {t("profiles.field_name")}
    </label>
    <input
      id="profile-name"
      type="text"
      class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all font-semibold"
      placeholder={t("profiles.placeholder_name")}
      bind:value={modalName}
    />
  </div>
</Modal>

<!-- Modal de Confirmação de Exclusão de Perfil -->
<ConfirmDialog
  bind:show={showDeleteConfirm}
  title="Remover Perfil"
  message={`Tem certeza de que deseja remover o perfil '${profileToDelete?.name || ""}'?\nEssa ação não poderá ser desfeita.`}
  confirmText="Remover Perfil"
  type="danger"
  onConfirm={confirmDeleteProfile}
/>
