<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import {
    ButtonGreen,
    ButtonYellow,
    ButtonBlue,
    ButtonPurple,
    ButtonRed,
    ButtonOrange,
  } from "$shared/components/buttons";
  import type { StackItem } from "$lib/domains/stacks";

  let {
    stack,
    on_deploy = () => {},
    on_stop = () => {},
    on_logs = () => {},
    on_edit = () => {},
    on_delete_local = () => {},
    on_remove_remote = () => {},
  }: {
    stack: StackItem;
    on_deploy?: () => void | Promise<void>;
    on_stop?: () => void | Promise<void>;
    on_logs?: () => void | Promise<void>;
    on_edit?: () => void | Promise<void>;
    on_delete_local?: () => void | Promise<void>;
    on_remove_remote?: () => void | Promise<void>;
  } = $props();

  let formattedDate = $derived(
    stack.updatedAt ? new Date(stack.updatedAt).toLocaleString() : "",
  );

  let formattedDeployDate = $derived(
    stack.lastDeployedAt ? new Date(stack.lastDeployedAt).toLocaleString() : "",
  );
</script>

<div
  class="flex flex-col xl:flex-row xl:items-center justify-between p-4 bg-white dark:bg-[#0c1220] border-2 border-slate-300/80 dark:border-slate-700 rounded-2xl gap-4 hover:border-violet-500/50 dark:hover:border-violet-500/40 transition-all shadow-md hover:shadow-lg"
>
  <!-- Informações Principais da Stack -->
  <div class="flex items-start gap-4 min-w-0 flex-1">
    <div
      class="w-11 h-11 rounded-xl bg-violet-500/10 text-violet-600 dark:text-violet-400 flex items-center justify-center font-bold text-xs shrink-0 border border-violet-500/20 mt-0.5 shadow-inner"
    >
      {stack.sourceType === "folder" ? "DIR" : "YAML"}
    </div>

    <div class="min-w-0 flex-1 space-y-2">
      <!-- Linha 1: Título e Badges de Status -->
      <div class="flex items-center gap-2.5 flex-wrap">
        <h3
          class="font-bold text-slate-900 dark:text-white text-base leading-tight m-0"
          title={stack.name}
        >
          {stack.name}
        </h3>

        {#if stack.sourceType === "folder"}
          <span class="px-2 py-0.5 text-[11px] font-semibold rounded-md bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/20 shrink-0">
            {t("stacks.source_type_folder")}
          </span>
        {:else}
          <span class="px-2 py-0.5 text-[11px] font-semibold rounded-md bg-blue-500/10 text-blue-600 dark:text-blue-400 border border-blue-500/20 shrink-0">
            {t("stacks.source_type_editor")}
          </span>
        {/if}

        {#if stack.lastDeployedAt}
          <span class="px-2 py-0.5 text-[11px] font-semibold rounded-md bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20 flex items-center gap-1.5 shrink-0">
            <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
            {t("stacks.last_deployed_label")} {formattedDeployDate}
          </span>
        {:else}
          <span class="px-2 py-0.5 text-[11px] font-medium rounded-md bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 border border-slate-200 dark:border-slate-750 shrink-0">
            {t("stacks.never_deployed")}
          </span>
        {/if}
      </div>

      <!-- Linha 2: Metadados da Stack -->
      <div
        class="flex flex-wrap items-center gap-x-3 gap-y-1.5 text-xs text-slate-500 dark:text-slate-400"
      >
        <span
          class="font-mono bg-slate-100 dark:bg-slate-800/80 px-2 py-0.5 rounded-md text-xs border border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-300 font-medium"
        >
          {stack.projectName || stack.id}
        </span>

        {#if stack.folderPath}
          <span class="text-slate-300 dark:text-slate-700">•</span>
          <span class="font-mono text-xs text-slate-500 dark:text-slate-400" title={stack.folderPath}>
            {stack.folderPath}
          </span>
        {/if}

        {#if stack.createdAt}
          <span class="text-slate-300 dark:text-slate-700">•</span>
          <span>
            {t("stacks.created_at")}: {new Date(stack.createdAt).toLocaleDateString()}
          </span>
        {/if}

        {#if formattedDate}
          <span class="text-slate-300 dark:text-slate-700">•</span>
          <span class="text-slate-400 dark:text-slate-500">
            {t("stacks.updated_at")}: {formattedDate}
          </span>
        {/if}
      </div>
    </div>
  </div>

  <!-- Ações da Stack -->
  <div class="flex flex-wrap items-center gap-2 shrink-0 pt-2 xl:pt-0 border-t xl:border-t-0 border-slate-100 dark:border-slate-800/60">
    <!-- Deploy (Verde) -->
    <ButtonGreen size="sm" onclick={on_deploy}>
      {#snippet icon()}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="currentColor"
          class="w-3.5 h-3.5"
        >
          <path
            fill-rule="evenodd"
            d="M4.5 5.653c0-1.427 1.529-2.33 2.779-1.643l11.54 6.347c1.295.712 1.295 2.573 0 3.286L7.28 19.99c-1.25.687-2.779-.217-2.779-1.643V5.653Z"
            clip-rule="evenodd"
          />
        </svg>
      {/snippet}
      {t("stacks.deploy_btn")}
    </ButtonGreen>

    <!-- Parar / Remover Container (Amarelo / Âmbar) -->
    <ButtonYellow size="sm" onclick={on_stop}>
      {#snippet icon()}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="currentColor"
          class="w-3.5 h-3.5"
        >
          <rect width="12" height="12" x="6" y="6" rx="2" />
        </svg>
      {/snippet}
      {t("stacks.stop_btn")}
    </ButtonYellow>

    <!-- Logs (Azul) -->
    <ButtonBlue size="sm" onclick={on_logs}>
      {#snippet icon()}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="w-3.5 h-3.5"
        >
          <polyline points="4 17 10 11 4 5" />
          <line x1="12" x2="20" y1="19" y2="19" />
        </svg>
      {/snippet}
      {t("stacks.logs_btn")}
    </ButtonBlue>

    <!-- Editar (Roxo / Púrpura) -->
    <ButtonPurple size="sm" onclick={on_edit}>
      {#snippet icon()}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="w-3.5 h-3.5"
        >
          <path
            d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"
          />
          <path d="m15 5 4 4" />
        </svg>
      {/snippet}
      {t("stacks.edit_btn")}
    </ButtonPurple>

    <!-- Remover da VPS / Down (Laranja) -->
    <ButtonOrange size="sm" onclick={on_remove_remote}>
      {#snippet icon()}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="w-3.5 h-3.5"
        >
          <path d="M12 2v10" />
          <path d="m16 8-4 4-4-4" />
          <path d="M2 17h20" />
          <path d="M6 21h12" />
        </svg>
      {/snippet}
      {t("stacks.remove_remote_btn")}
    </ButtonOrange>

    <!-- Excluir Definição Local (Vermelho / Rose) -->
    <ButtonRed size="sm" onclick={on_delete_local}>
      {#snippet icon()}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="w-3.5 h-3.5"
        >
          <path d="M3 6h18" />
          <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
          <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
        </svg>
      {/snippet}
      {t("stacks.delete_local_btn")}
    </ButtonRed>
  </div>
</div>
