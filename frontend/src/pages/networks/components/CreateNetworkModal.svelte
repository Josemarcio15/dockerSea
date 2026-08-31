<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import FormModal from "$shared/components/FormModal.svelte";

  let {
    show = $bindable(false),
    name = $bindable(""),
    driver = $bindable("bridge"),
    subnet = $bindable(""),
    gateway = $bindable(""),
    onSubmit = () => {},
  }: {
    show?: boolean;
    name?: string;
    driver?: string;
    subnet?: string;
    gateway?: string;
    onSubmit?: () => void;
  } = $props();
</script>

<FormModal
  bind:show
  title={t("networks.create_title")}
  buttons={[
    {
      label: t("networks.create_btn"),
      variant: "primary",
      onclick: onSubmit,
      disabled: !name.trim(),
    },
  ]}
>
  <div class="flex flex-col gap-1.5">
    <label
      for="net-name"
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >
      {t("networks.field_name")}
    </label>
    <input
      id="net-name"
      type="text"
      class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
      placeholder={t("networks.placeholder_name")}
      bind:value={name}
    />
  </div>

  <div class="flex flex-col gap-1.5">
    <label
      for="net-driver"
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >
      {t("networks.field_driver")}
    </label>
    <select
      id="net-driver"
      class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
      bind:value={driver}
    >
      <option value="bridge">{t("networks.driver_bridge")}</option>
      <option value="host">host</option>
      <option value="overlay">overlay</option>
      <option value="macvlan">macvlan</option>
    </select>
  </div>

  <div class="flex flex-col gap-1.5">
    <label
      for="net-subnet"
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >
      {t("networks.field_subnet")}
    </label>
    <input
      id="net-subnet"
      type="text"
      class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
      placeholder="ex: 172.20.0.0/16"
      bind:value={subnet}
    />
  </div>

  <div class="flex flex-col gap-1.5">
    <label
      for="net-gateway"
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >
      {t("networks.field_gateway")}
    </label>
    <input
      id="net-gateway"
      type="text"
      class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
      placeholder="ex: 172.20.0.1"
      bind:value={gateway}
    />
  </div>
</FormModal>
