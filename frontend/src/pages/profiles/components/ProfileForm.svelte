<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import Modal from "$shared/components/Modal.svelte";
  let { show = $bindable(false), id, name = $bindable(""), onSave } = $props();
</script>

<Modal
  bind:show
  title={id ? t("profiles.edit_title") : t("profiles.create_title")}
  buttons={[
    {
      label: t("common.save"),
      variant: "primary",
      onclick: onSave,
      disabled: !name.trim(),
    },
  ]}
>
  <div class="flex flex-col gap-1.5">
    <label
      for="profile-name"
      class="text-xs font-bold text-slate-400 uppercase tracking-wider"
      >{t("profiles.field_name")}</label
    ><input
      id="profile-name"
      bind:value={name}
      placeholder={t("profiles.placeholder_name")}
      class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500"
    />
  </div>
</Modal>
