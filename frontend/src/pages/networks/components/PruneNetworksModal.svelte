<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import Modal from "$shared/components/Modal.svelte";

  let {
    show = $bindable(false),
    onConfirm = () => {},
  }: {
    show?: boolean;
    onConfirm?: () => void;
  } = $props();
</script>

<Modal
  bind:show
  title={t("networks.prune_confirm_title")}
  buttons={[
    {
      label: t("common.confirm"),
      variant: "success",
      onclick: onConfirm,
    },
  ]}
>
  <p class="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
    {t("networks.prune_confirm_msg")}
  </p>
</Modal>
