<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import Modal from "$lib/components/Modal.svelte";

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
  title={t("volumes.prune_confirm_title")}
  buttons={[
    {
      label: t("common.confirm"),
      variant: "success",
      onclick: onConfirm,
    },
  ]}
>
  <p class="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
    {t("volumes.prune_confirm_msg")}
  </p>
</Modal>
