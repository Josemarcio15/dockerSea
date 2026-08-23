<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import type { Container } from "../../../../bindings/go-walis/internal/containers/models.js";

  let {
    show = $bindable(false),
    networkName = "",
    selectedContainer = $bindable(""),
    containers = [],
    onSubmit = () => {},
  }: {
    show?: boolean;
    networkName?: string;
    selectedContainer?: string;
    containers?: Container[];
    onSubmit?: () => void;
  } = $props();
</script>

<Modal
  bind:show
  title={t("networks.connect_modal_title")}
  buttons={[
    {
      label: t("networks.connect_btn"),
      variant: "primary",
      onclick: onSubmit,
      disabled: !selectedContainer,
    },
  ]}
>
  <!-- Network info label -->
  <div
    class="flex items-center gap-2 px-4 py-3 rounded-xl bg-slate-50 dark:bg-slate-900/30 border border-slate-200/60 dark:border-slate-800/60"
  >
    <span
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider shrink-0"
    >
      {t("networks.select_network_label")}
    </span>
    <span class="text-sm font-semibold text-violet-600 dark:text-violet-400">
      {networkName}
    </span>
  </div>

  <div class="flex flex-col gap-1.5">
    <label
      for="select-cont"
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >
      {t("networks.select_container_label")}
    </label>
    <select
      id="select-cont"
      class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
      bind:value={selectedContainer}
    >
      {#each containers as container}
        <option value={container.name}>
          {container.name} ({container.image})
        </option>
      {:else}
        <option value="" disabled>{t("networks.no_active_containers")}</option>
      {/each}
    </select>
  </div>
</Modal>
