<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";

  let {
    stack,
    on_deploy = () => {},
    on_stop = () => {},
    on_logs = () => {},
    on_edit = () => {},
    on_delete = () => {},
  }: {
    stack: {
      id: string;
      name: string;
      yamlContent: string;
      createdAt: string;
      updatedAt: string;
    };
    on_deploy?: () => void | Promise<void>;
    on_stop?: () => void | Promise<void>;
    on_logs?: () => void | Promise<void>;
    on_edit?: () => void | Promise<void>;
    on_delete?: () => void | Promise<void>;
  } = $props();

  // Format datetime locally
  const formattedDate = $derived.by(() => {
    try {
      const d = new Date(stack.updatedAt || stack.createdAt);
      return d.toLocaleString("pt-BR");
    } catch (e) {
      return stack.updatedAt || stack.createdAt;
    }
  });

  // Extract project name from YAML `name:` field
  const projectName = $derived.by(() => {
    const match = (stack.yamlContent || "").match(/^name:\s*(.+)/m);
    return match ? match[1].trim() : "";
  });
</script>

<div
  data-stack-id={stack.id}
  title={stack.yamlContent ? "YAML Configurado" : ""}
  class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-4.5 rounded-2xl bg-white dark:bg-[#0c1220] border border-slate-200 dark:border-slate-800/80 hover:border-slate-300 dark:hover:border-slate-700 transition-all duration-300 shadow-md dark:shadow-lg dark:shadow-black/40"
>
  <!-- Stack Title & Icon -->
  <div class="flex items-center gap-3.5 min-w-0">
    <div
      class="w-11 h-11 rounded-xl bg-linear-to-br from-rose-400 to-pink-500 flex items-center justify-center text-lg shadow-lg shadow-pink-500/20 shrink-0 text-white font-bold"
    >
      {stack.name.charAt(0).toUpperCase()}
    </div>
    <div class="flex flex-col min-w-0">
      <span
        class="font-bold text-slate-800 dark:text-slate-100 truncate text-base"
        >{stack.name}</span
      >
      <div class="flex items-center gap-2 mt-0.5">
        {#if projectName}
          <span
            class="text-[10px] px-1.5 py-0.5 rounded-md bg-violet-50 dark:bg-violet-950/30 text-violet-600 dark:text-violet-400 font-mono font-semibold border border-violet-200 dark:border-violet-900/40"
          >
            {projectName}
          </span>
        {/if}
        <span
          class="text-[10px] text-slate-400 dark:text-slate-500 font-semibold uppercase tracking-wider"
        >
          {t("stacks.updated_at")}
          {formattedDate}
        </span>
      </div>
    </div>
  </div>

  <!-- Action buttons -->
  <div class="flex flex-wrap items-center gap-2">
    <button
      type="button"
      class="px-3.5 py-2 text-xs font-bold text-emerald-700 dark:text-white bg-emerald-50 dark:bg-emerald-600 hover:bg-emerald-100 dark:hover:bg-emerald-700 border border-emerald-200 dark:border-none rounded-xl cursor-pointer transition-all duration-200 shadow-md shadow-emerald-500/5 flex items-center gap-1.5"
      onclick={on_deploy}
    >
      {t("stacks.deploy_btn")}
    </button>

    <button
      type="button"
      class="px-3.5 py-2 text-xs font-bold text-amber-700 dark:text-white bg-amber-50 dark:bg-amber-600 hover:bg-amber-100 dark:hover:bg-amber-750 border border-amber-200 dark:border-none rounded-xl cursor-pointer transition-all duration-200 shadow-md shadow-amber-500/5 flex items-center gap-1.5"
      onclick={on_stop}
    >
      {t("stacks.stop_btn")}
    </button>

    <button
      type="button"
      class="px-3.5 py-2 text-xs font-bold text-white bg-blue-600 hover:bg-blue-700 border-none rounded-xl cursor-pointer transition-all duration-200 shadow-md shadow-blue-500/20 flex items-center gap-1.5"
      onclick={on_logs}
    >
      {t("stacks.logs_btn")}
    </button>

    <button
      type="button"
      class="px-3.5 py-2 text-xs font-bold text-white bg-purple-600 hover:bg-purple-700 border-none rounded-xl cursor-pointer transition-all duration-200 shadow-md shadow-purple-500/20 flex items-center gap-1.5"
      onclick={on_edit}
    >
      {t("stacks.edit_btn")}
    </button>

    <button
      type="button"
      class="px-3.5 py-2 text-xs font-bold text-white bg-red-600 hover:bg-red-700 border-none rounded-xl cursor-pointer transition-all duration-200 shadow-md shadow-red-500/20 flex items-center gap-1.5"
      onclick={on_delete}
    >
      {t("stacks.delete_btn")}
    </button>
  </div>
</div>
