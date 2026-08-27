import { notifySuccess, notifyError } from "$shared/stores/notification.svelte";
import { triggerRefresh } from "$shared/stores/refresh.svelte";
import * as api from "./api";

export function createProfilesStore() {
  let form = $state({ id: "", name: "" });
  let showModal = $state(false);
  let showDeleteConfirm = $state(false);
  let profileToDelete = $state<{ id: string; name: string } | null>(null);

  async function action(operation: () => Promise<unknown>) {
    try {
      await operation();
      notifySuccess("Operação realizada com sucesso");
      triggerRefresh();
    } catch (error: any) {
      notifyError(error?.message || String(error));
    }
  }
  return {
    get form() {
      return form;
    },
    get showModal() {
      return showModal;
    },
    set showModal(v: boolean) {
      showModal = v;
    },
    get showDeleteConfirm() {
      return showDeleteConfirm;
    },
    set showDeleteConfirm(v: boolean) {
      showDeleteConfirm = v;
    },
    get profileToDelete() {
      return profileToDelete;
    },
    openCreate() {
      form = { id: "", name: "" };
      showModal = true;
    },
    openEdit(profile: any) {
      form = { id: profile.id, name: profile.name };
      showModal = true;
    },
    async save() {
      if (!form.name.trim()) return;
      const profile = { id: form.id, name: form.name, locale: "pt-BR" };
      showModal = false;
      await action(() => api.saveProfile(profile));
    },
    async select(id: string) {
      await action(() => api.setActiveProfile(id));
    },
    requestDelete(id: string, name: string) {
      profileToDelete = { id, name };
      showDeleteConfirm = true;
    },
    async confirmDelete() {
      if (!profileToDelete) return;
      const profileId = profileToDelete.id;
      await action(() => api.deleteProfile(profileId));
      showDeleteConfirm = false;
    },
  };
}
