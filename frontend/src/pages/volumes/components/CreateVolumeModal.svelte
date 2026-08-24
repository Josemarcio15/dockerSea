<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import { ButtonRed, Button } from "$lib/components/buttons";

  let {
    show = $bindable(false),
    name = $bindable(""),
    driver = $bindable(""),
    labelEntries = $bindable([]),
    onAddLabel = () => {},
    onRemoveLabel = (index: number) => {},
    onSubmit = () => {},
  }: {
    show?: boolean;
    name?: string;
    driver?: string;
    labelEntries?: Array<{ key: string; value: string }>;
    onAddLabel?: () => void;
    onRemoveLabel?: (index: number) => void;
    onSubmit?: () => void;
  } = $props();
</script>

<Modal
  bind:show
  title={t("volumes.modal_title")}
  buttons={[
    {
      label: t("volumes.create_btn"),
      variant: "primary",
      onclick: onSubmit,
      disabled: !name.trim(),
    },
  ]}
>
  <!-- Nome -->
  <div class="flex flex-col gap-1.5">
    <label
      for="vol-name"
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >
      {t("volumes.field_name")}
    </label>
    <input
      id="vol-name"
      type="text"
      class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
      placeholder={t("volumes.placeholder_name")}
      bind:value={name}
    />
  </div>

  <!-- Driver -->
  <div class="flex flex-col gap-1.5">
    <label
      for="vol-driver"
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >
      {t("volumes.field_driver")}
    </label>
    <input
      id="vol-driver"
      type="text"
      class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
      placeholder="local"
      bind:value={driver}
    />
  </div>

  <!-- Labels -->
  <div class="flex flex-col gap-2">
    <div class="flex items-center justify-between">
      <span
        class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
      >
        {t("volumes.field_labels")}
      </span>
      <Button
        size="xs"
        onclick={onAddLabel}
      >
        {t("volumes.add_label")}
      </Button>
    </div>

    {#if labelEntries.length === 0}
      <p class="text-xs text-slate-400 dark:text-slate-600 italic">
        {t("volumes.no_labels")}
      </p>
    {/if}

    {#each labelEntries as entry, i}
      <div class="flex gap-2 items-center">
        <input
          type="text"
          placeholder="Chave"
          class="flex-1 px-3 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-1 focus:ring-violet-500 transition-all"
          bind:value={entry.key}
        />
        <input
          type="text"
          placeholder="Valor"
          class="flex-1 px-3 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-1 focus:ring-violet-500 transition-all"
          bind:value={entry.value}
        />
        <ButtonRed
          size="xs"
          onclick={() => onRemoveLabel(i)}
        >
          ✕
        </ButtonRed>
      </div>
    {/each}
  </div>
</Modal>
