<script lang="ts">
  import { tick } from "svelte";
  import { t } from "$shared/stores/locale.svelte";
  import { notifySuccess } from "$shared/stores/notification.svelte";
  import type { BuilderStatus } from "../types";

  let { logs, status }: { logs: string[]; status: BuilderStatus } = $props();
  let logContainer = $state<HTMLDivElement | null>(null);
  let copied = $state(false);

  $effect(() => {
    if (logs.length > 0) {
      tick().then(() => {
        if (logContainer) logContainer.scrollTop = logContainer.scrollHeight;
      });
    }
  });

  async function copyLogs() {
    if (logs.length === 0) return;
    try {
      await navigator.clipboard.writeText(logs.join("\n"));
    } catch {
      const textarea = document.createElement("textarea");
      textarea.value = logs.join("\n");
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand("copy");
      document.body.removeChild(textarea);
    }
    copied = true;
    notifySuccess(
      t("builder.logs_copied") || "Logs copiados para a área de transferência!",
    );
    setTimeout(() => (copied = false), 2000);
  }
</script>

<div
  class="bg-[#070b15] border border-[#1e293b] rounded-2xl overflow-hidden flex flex-col h-120 shadow-xl"
>
  <div
    class="px-4 py-3 bg-[#0d1424] border-b border-[#1e293b] flex items-center justify-between shrink-0"
  >
    <div class="flex items-center gap-2.5">
      <span
        class="w-2.5 h-2.5 rounded-full {status === 'building'
          ? 'bg-amber-400 animate-pulse'
          : status === 'success'
            ? 'bg-emerald-400'
            : status === 'error'
              ? 'bg-red-400'
              : 'bg-slate-500'}"
      ></span>
      <span class="text-xs font-bold text-slate-300 tracking-wide uppercase"
        >{t("builder.logs_title")}</span
      >
      {#if status === "building"}<span
          class="text-[10px] px-2 py-0.5 rounded-full bg-amber-500/10 border border-amber-500/20 text-amber-400 font-semibold animate-pulse"
          >{t("builder.building")}</span
        >{/if}
    </div>
    {#if logs.length > 0}<button
        type="button"
        class="flex items-center gap-1.5 px-3 py-1 text-xs font-semibold rounded-xl bg-slate-800/80 hover:bg-slate-700 text-slate-200 border border-slate-700/60 transition-all cursor-pointer shadow-sm active:scale-95"
        onclick={copyLogs}
        >{copied ? t("common.copied") : t("builder.copy_logs")}</button
      >{/if}
  </div>
  <div
    bind:this={logContainer}
    class="p-4 overflow-y-auto flex-1 font-mono text-xs text-emerald-400 space-y-1 scrollbar-thin"
  >
    {#if logs.length === 0}<div class="text-slate-500 italic">
        {status === "building"
          ? t("builder.waiting_build")
          : t("builder.select_folder_build")}
      </div>{:else}{#each logs as log}<div
          class="leading-relaxed whitespace-pre-wrap"
        >
          {log}
        </div>{/each}{/if}
  </div>
</div>
