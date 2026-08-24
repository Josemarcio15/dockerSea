<script lang="ts">
  import { onDestroy } from "svelte";
  import { ButtonPink } from "$lib/components/buttons";
  import * as ExtraService from "../../../bindings/go-walis/internal/extras/extraservice.js";

  let {
    show = $bindable(false),
    activeVps,
  }: { show: boolean; activeVps?: any } = $props();
  let logs = $state("");
  let loading = $state(false);
  let error = $state("");
  let timer: ReturnType<typeof setInterval> | undefined;

  async function refresh() {
    if (!activeVps) {
      error = "Nenhum servidor selecionado.";
      return;
    }
    if (loading) return;
    loading = true;
    error = "";
    try {
      const result = await ExtraService.GetNginxLogs(activeVps, 150);
      logs = result || "Nenhum log encontrado.";
    } catch (e: any) {
      error = e?.message || "Não foi possível carregar os logs.";
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if (show) {
      refresh();
      timer = setInterval(refresh, 3000);
    } else if (timer) {
      clearInterval(timer);
      timer = undefined;
    }
  });

  onDestroy(() => timer && clearInterval(timer));
</script>

{#if show}
  <div class="fixed inset-0 z-60 flex items-center justify-center bg-black/70 p-4">
    <div class="flex max-h-[90vh] w-full max-w-6xl flex-col overflow-hidden rounded-2xl border border-slate-700 bg-[#080d16] shadow-2xl">
      <div class="flex items-center justify-between border-b border-slate-800 px-5 py-4">
        <div>
          <h2 class="text-lg font-bold text-slate-100">Logs do Nginx</h2>
          <p class="text-xs text-slate-400">/var/log/nginx/access.log e error.log</p>
        </div>
        <div class="flex items-center gap-2">
          <ButtonPink size="xs" onclick={refresh} loading={loading}>
            Atualizar
          </ButtonPink>
          <button type="button" class="px-2 text-xl text-slate-400 hover:text-white cursor-pointer" onclick={() => (show = false)}>×</button>
        </div>
      </div>
      <div class="min-h-0 flex-1 overflow-auto p-4">
        {#if loading && !logs}<p class="text-sm text-slate-400">Carregando logs...</p>
        {:else if error}<pre class="whitespace-pre-wrap text-sm text-red-300">{error}</pre>
        {:else}<pre class="whitespace-pre-wrap text-xs leading-relaxed text-emerald-300">{logs}</pre>{/if}
      </div>
    </div>
  </div>
{/if}
