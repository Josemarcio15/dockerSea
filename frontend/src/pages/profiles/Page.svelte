<script lang="ts">
  import StatusBanner from "$shared/components/StatusBanner.svelte";
  import { t } from "$shared/stores/locale.svelte";
  import ConfirmDialog from "$shared/components/ConfirmDialog.svelte";
  import { ButtonBlue } from "$shared/components/buttons";
  import ProfileCard from "./components/ProfileCard.svelte";
  import ProfileForm from "./components/ProfileForm.svelte";
  import { createProfilesStore } from "./store.svelte";

  let { data } = $props();
  const store = createProfilesStore();
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

    <ButtonBlue size="sm" onclick={store.openCreate}>
      {t("profiles.new_profile")}
    </ButtonBlue>
  </div>

  <!-- Status Alerts -->
  <StatusBanner />

  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
    {#each data.profiles as profile (profile.id)}
      <ProfileCard
        {profile}
        active={data.activeProfile?.id === profile.id}
        canDelete={data.profiles.length > 1}
        onSelect={() => store.select(profile.id)}
        onEdit={() => store.openEdit(profile)}
        onDelete={() => store.requestDelete(profile.id, profile.name)}
      />
    {/each}
  </div>
</div>

<!-- Modal: Adicionar/Editar Perfil -->
<ProfileForm
  bind:show={store.showModal}
  id={store.form.id}
  bind:name={store.form.name}
  onSave={store.save}
/>

<!-- Modal de Confirmação de Exclusão de Perfil -->
<ConfirmDialog
  bind:show={store.showDeleteConfirm}
  title="Remover Perfil"
  message={`Tem certeza de que deseja remover o perfil '${store.profileToDelete?.name || ""}'?\nEssa ação não poderá ser desfeita.`}
  confirmText="Remover Perfil"
  type="danger"
  onConfirm={store.confirmDelete}
/>
