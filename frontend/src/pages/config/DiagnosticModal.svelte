<script lang="ts">
  import FormModal from "$shared/components/FormModal.svelte";

  let {
    show = $bindable(false),
    title = "Diagnóstico de Conexão",
    loading = false,
    result = null,
  }: {
    show: boolean;
    title: string;
    loading: boolean;
    result: {
      success: boolean;
      message: string;
      steps: Array<{
        name: string;
        status: "success" | "error" | "warning";
        message: string;
      }>;
    } | null;
  } = $props();
</script>

<FormModal bind:show {title}>
  {#if loading}
    <div class="py-8 flex flex-col items-center justify-center space-y-3">
      <div
        class="w-8 h-8 border-4 border-violet-500/30 border-t-violet-500 rounded-full animate-spin"
      ></div>
      <p class="text-xs font-semibold text-slate-600 dark:text-slate-400">
        Testando conexão SSH, ambiente do OS e API do Docker...
      </p>
    </div>
  {:else if result}
    <div class="space-y-4 max-h-[70vh] overflow-y-auto pr-1">
      <!-- Status Banner Resumo -->
      <div
        class="p-4 rounded-2xl border flex items-start gap-3 text-sm font-semibold {result.success
          ? 'bg-emerald-50 dark:bg-emerald-950/30 text-emerald-700 dark:text-emerald-300 border-emerald-200 dark:border-emerald-900/50'
          : 'bg-red-50 dark:bg-red-950/30 text-red-700 dark:text-red-300 border-red-200 dark:border-red-900/50'}"
      >
        <span class="text-xl shrink-0">{result.success ? "✅" : "❌"}</span>
        <div class="space-y-1">
          <div class="font-bold">{result.message}</div>
        </div>
      </div>

      <!-- Etapas do Diagnóstico -->
      <div class="space-y-3">
        <h4
          class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
        >
          Etapas do Diagnóstico
        </h4>

        {#each result.steps as step}
          <div
            class="p-4 rounded-xl border bg-slate-50 dark:bg-slate-900/50 border-slate-200 dark:border-slate-800 space-y-2"
          >
            <div class="flex items-center justify-between">
              <span
                class="font-bold text-sm text-slate-800 dark:text-slate-200 flex items-center gap-2"
              >
                {#if step.status === "success"}
                  <span class="text-emerald-500 font-bold">✓</span>
                {:else if step.status === "warning"}
                  <span class="text-amber-500 font-bold">⚠️</span>
                {:else}
                  <span class="text-red-500 font-bold">✕</span>
                {/if}
                {step.name}
              </span>
              <span
                class="px-2 py-0.5 text-[10px] font-bold rounded-md uppercase {step.status ===
                'success'
                  ? 'bg-emerald-500/10 text-emerald-500 border border-emerald-500/20'
                  : step.status === 'warning'
                    ? 'bg-amber-500/10 text-amber-500 border border-amber-500/20'
                    : 'bg-red-500/10 text-red-500 border border-red-500/20'}"
              >
                {step.status === "success"
                  ? "Sucesso"
                  : step.status === "warning"
                    ? "Atenção"
                    : "Falha"}
              </span>
            </div>
            <p
              class="text-xs text-slate-600 dark:text-slate-300 leading-relaxed whitespace-pre-wrap font-mono bg-white dark:bg-[#0c1220] p-3 rounded-lg border border-slate-200/60 dark:border-slate-800/80"
            >
              {step.message}
            </p>
          </div>
        {/each}
      </div>
    </div>
  {/if}
</FormModal>
