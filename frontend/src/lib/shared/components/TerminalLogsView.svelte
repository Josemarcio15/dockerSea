<script lang="ts">
  import { tick } from "svelte";

  let {
    logs = [],
    id = "terminal-logs-view",
    maxHeight = "max-h-[55vh]",
    class: customClass = "",
  }: {
    logs: string[];
    id?: string;
    maxHeight?: string;
    class?: string;
  } = $props();

  let container = $state<HTMLDivElement | null>(null);

  $effect(() => {
    if (logs.length > 0) {
      tick().then(() => {
        if (container) {
          container.scrollTop = container.scrollHeight;
        }
      });
    }
  });

  function getLogColor(line: string): string {
    const upper = line.toUpperCase();
    if (
      upper.includes("FATAL") ||
      upper.includes("ERROR") ||
      upper.includes("ERRO") ||
      upper.includes("FAIL") ||
      upper.includes("PANIC") ||
      upper.includes("✗") ||
      upper.includes("❌")
    ) {
      return "text-red-400 font-semibold";
    }
    if (
      upper.includes("WARN") ||
      upper.includes("WARNING") ||
      upper.includes("DETAIL:")
    ) {
      return "text-amber-300";
    }
    if (
      upper.includes("SUCCESS") ||
      upper.includes("SUCESSO") ||
      upper.includes("✓") ||
      upper.includes("COMPLETE") ||
      upper.includes("CONSTRUÍDA COM SUCESSO") ||
      upper.includes("READY TO ACCEPT CONNECTIONS") ||
      upper.includes("OK")
    ) {
      return "text-emerald-400 font-medium";
    }
    if (
      upper.includes("STEP") ||
      upper.includes("BUILDING") ||
      upper.includes("COMPILING") ||
      upper.includes("DOWNLOADING") ||
      upper.includes("UPLOADING") ||
      upper.includes("PREPARING") ||
      upper.includes("LOG:") ||
      upper.includes("INFO")
    ) {
      return "text-sky-300";
    }
    return "text-slate-300";
  }
</script>

<div
  bind:this={container}
  {id}
  class="bg-[#090d16] text-slate-200 p-4 rounded-xl text-xs font-mono overflow-auto flex-1 shadow-inner border border-slate-800/80 space-y-1 select-text scrollbar-thin {maxHeight} {customClass}"
>
  {#if logs.length === 0}
    <div class="text-slate-500 italic py-6 text-center select-none">
      (Nenhum log disponível)
    </div>
  {:else}
    {#each logs as log, index}
      <div
        class="flex items-start gap-3 leading-relaxed break-all font-mono hover:bg-slate-800/30 px-1.5 py-0.5 rounded transition-colors"
      >
        <span
          class="text-slate-600 dark:text-slate-500 text-[11px] shrink-0 select-none text-right w-7 font-mono pr-1 border-r border-slate-800/60"
        >
          {index + 1}
        </span>
        <span class="{getLogColor(log)} flex-1">{log}</span>
      </div>
    {/each}
  {/if}
</div>
