<script lang="ts">
  import { onDestroy } from "svelte";
  import { t } from "$shared/stores/locale.svelte";
  import { Events } from "@wailsio/runtime";
  import { ButtonOrange } from "$shared/components/buttons";

  let {
    show = $bindable(false),
    imageName = "",
    sseUrl = "",
    title = "",
    eventPrefix = "docker:image:pull",
    oncomplete = () => {},
  }: {
    show: boolean;
    imageName?: string;
    sseUrl?: string;
    title?: string;
    eventPrefix?: string;
    oncomplete?: () => void;
  } = $props();

  let logs = $state<string[]>([]);
  let statusMsg = $state("");
  let errorMsg = $state("");
  let finished = $state(false);
  let percent = $state(0);
  let terminalElement = $state<HTMLDivElement | null>(null);

  let eventSource: EventSource | null = null;
  let unsubscribeEvents: (() => void) | null = null;
  let lastLoggedStatus: Record<string, string> = {};

  const isGenericMode = $derived(!!sseUrl);
  const isRunning = $derived(!finished && !errorMsg);
  const modalTitle = $derived(title || (imageName ? `Pull: ${imageName}` : "Progresso"));
  const displayPercent = $derived(finished ? 100 : Math.min(99, percent || (isRunning ? 15 : 0)));

  $effect(() => {
    if (terminalElement && logs.length) {
      terminalElement.scrollTop = terminalElement.scrollHeight;
    }
  });

  function closeListeners() {
    if (eventSource) {
      eventSource.close();
      eventSource = null;
    }
    if (unsubscribeEvents) {
      unsubscribeEvents();
      unsubscribeEvents = null;
    }
  }

  function handleProgressData(data: any) {
    if (data.error) {
      errorMsg = data.error;
      logs.push(`✗ Erro: ${data.error}`);
      closeListeners();
      return;
    }

    if (data.status === "pull_complete" || data.status === "complete") {
      statusMsg = data.message || "Download e extração concluídos!";
      logs.push(`✓ ${data.message || "Imagem baixada com sucesso!"}`);
      finished = true;
      percent = 100;
      closeListeners();
      oncomplete();
      return;
    }

    if (data.id) {
      const layerId = data.id;
      const currentStatus = data.status || "";
      if (lastLoggedStatus[layerId] !== currentStatus) {
        lastLoggedStatus[layerId] = currentStatus;
        if (!["Downloading", "Extracting"].includes(currentStatus)) {
          logs.push(`→ [${layerId}] ${currentStatus}`);
        }
      }
      if (data.progressDetail?.total > 0) {
        const rawPct = Math.round((data.progressDetail.current / data.progressDetail.total) * 100);
        if (rawPct > percent) percent = Math.min(95, rawPct);
      }
    } else if (data.message) {
      statusMsg = data.message;
      const phase = (data.phase || "").toLowerCase();
      const phasePrefix = data.phase ? `[${data.phase.toUpperCase()}] ` : "";
      logs.push(`→ ${phasePrefix}${data.message}`);

      // Progressão gradual de acordo com a fase do deploy
      if (phase === "preparing" && percent < 15) percent = 15;
      else if (phase === "uploading" && percent < 30) percent = 30;
      else if (phase === "validating" && percent < 45) percent = 45;
      else if (phase === "building") {
        if (percent < 50) percent = 50;
        else if (percent < 90) percent += 2; // avança a cada log de build
      } else if (phase === "starting" && percent < 92) percent = 92;
      else if (phase === "checking" && percent < 96) percent = 96;
      else if (phase === "cleaning" && percent < 98) percent = 98;
    } else if (data.line) {
      logs.push(data.line);
      statusMsg = data.line;
      if (percent < 90) percent += 1;
    } else if (data.status) {
      statusMsg = data.status;
      if (!logs.some((l) => l.includes(data.status))) {
        logs.push(`→ ${data.status}`);
      }
    }
  }

  function startGenericSse() {
    errorMsg = "";
    finished = false;
    percent = 5;
    statusMsg = t("terminal_modal.starting");
    logs = [];
    closeListeners();

    eventSource = new EventSource(sseUrl);
    eventSource.addEventListener("progress", (ev: MessageEvent) => {
      try { handleProgressData(JSON.parse(ev.data)); } catch {}
    });
    eventSource.addEventListener("complete", (ev: MessageEvent) => {
      try {
        const data = JSON.parse(ev.data);
        finished = true;
        percent = 100;
        statusMsg = data.message || (data.success ? "Concluído!" : "Falhou!");
        logs.push(data.message || (data.success ? "✓ Concluído!" : "✗ Falhou!"));
        closeListeners();
        oncomplete();
      } catch {}
    });
    eventSource.onerror = () => {
      if (!finished) {
        errorMsg = "Erro de conexão com o servidor.";
        logs.push("✗ Erro de conexão com o servidor.");
        closeListeners();
      }
    };
  }

  function startWailsEventsListener() {
    errorMsg = "";
    finished = false;
    percent = 5;
    statusMsg = t("terminal_modal.starting_process");
    logs = [];
    lastLoggedStatus = {};
    closeListeners();

    const unsubs = [
      Events.On(`${eventPrefix}:started`, (ev: any) => {
        const d = ev?.data ?? ev;
        const msg = d?.projectName ? `Iniciando deploy: ${d.projectName}` : "Deploy iniciado...";
        statusMsg = msg;
        logs.push(`🚀 ${msg}`);
      }),
      Events.On(`${eventPrefix}:progress`, (ev: any) => {
        handleProgressData(ev?.data ?? ev);
      }),
      Events.On(`${eventPrefix}:failed`, (ev: any) => {
        const d = ev?.data ?? ev;
        finished = true;
        errorMsg = d?.message || "Erro na operação";
        logs.push(`✗ [ERRO] ${d?.message || "Falha na operação"}`);
        closeListeners();
      }),
      Events.On(`${eventPrefix}:complete`, (ev: any) => {
        const d = ev?.data ?? ev;
        finished = true;
        percent = 100;
        if (d?.success !== false) {
          statusMsg = d?.message || "Concluído!";
          logs.push(d?.message || "✓ Operação concluída com sucesso!");
          oncomplete();
        } else {
          errorMsg = d?.message || "Erro na execução";
          logs.push(`✗ Erro: ${d?.message || "Falha na operação"}`);
        }
        closeListeners();
      })
    ];

    unsubscribeEvents = () => {
      unsubs.forEach((unsub) => { if (typeof unsub === "function") unsub(); });
    };
  }

  $effect(() => {
    if (show) {
      if (isGenericMode) startGenericSse();
      else startWailsEventsListener();
    }
    return closeListeners;
  });

  onDestroy(closeListeners);
</script>

{#if show}
  <div class="fixed inset-0 bg-black/80 backdrop-blur-md flex items-center justify-center z-50 p-4 animate-fadeIn">
    <div class="bg-[#0b101d] border border-slate-800/90 rounded-3xl w-200 max-w-full p-6 shadow-2xl flex flex-col max-h-[92vh] text-slate-200 gap-4">
      
      <!-- Top Title -->
      <div class="flex justify-between items-center pb-3 border-b border-slate-800/80">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-xl bg-violet-600/20 border border-violet-500/30 flex items-center justify-center text-violet-400 font-bold">
            ⬇️
          </div>
          <div>
            <h2 class="text-base font-bold text-white tracking-wide m-0">{modalTitle}</h2>
            <p class="text-xs text-slate-400 m-0">
              {isRunning ? (statusMsg || "Processando operação...") : finished ? "Processo finalizado com sucesso" : errorMsg ? "Operação falhou" : "Aguardando"}
            </p>
          </div>
        </div>
        <button
          type="button"
          class="w-8 h-8 rounded-xl bg-slate-850 hover:bg-slate-800 border border-slate-700 text-slate-400 hover:text-white flex items-center justify-center transition-colors {isRunning ? 'cursor-not-allowed opacity-30' : 'cursor-pointer'}"
          disabled={isRunning}
          onclick={() => { closeListeners(); show = false; }}
        >
          ✕
        </button>
      </div>

      <!-- Barra de Progresso Geral -->
      <div class="border border-slate-800 rounded-2xl p-4 bg-[#070b14]/80 space-y-2.5 shadow-inner">
        <div class="flex items-center justify-between text-xs font-bold">
          <span class="text-white">Progresso Geral</span>
          <span class="font-mono text-sm font-bold bg-linear-to-r from-violet-400 to-emerald-400 bg-clip-text text-transparent">
            {displayPercent}%
          </span>
        </div>
        <div class="w-full h-3 bg-slate-900 rounded-full overflow-hidden border border-slate-800 p-0.5">
          <div
            class="h-full rounded-full transition-all duration-300 bg-linear-to-r from-violet-600 via-blue-500 to-emerald-500 shadow-sm shadow-emerald-500/30"
            style="width: {displayPercent}%"
          ></div>
        </div>
      </div>

      <!-- Terminal de Logs -->
      <div
        bind:this={terminalElement}
        class="border border-slate-800/80 rounded-2xl p-4 bg-[#050811] font-mono text-xs text-emerald-400 h-64 overflow-y-auto space-y-1.5 scrollbar-thin shadow-inner"
      >
        {#if logs.length === 0}
          <div class="text-slate-500 italic flex items-center gap-2">
            <span class="animate-spin text-sm">⚙️</span>
            {t("terminal_modal.waiting_stream")}
          </div>
        {:else}
          {#each logs as log}
            <div class="leading-relaxed whitespace-pre-wrap">{log}</div>
          {/each}
        {/if}
      </div>

      <!-- Rodapé / Ação -->
      <div class="flex justify-between items-center pt-3 border-t border-slate-800/80">
        <div class="flex items-center gap-2 text-xs font-semibold">
          {#if isRunning}
            <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-violet-500/10 border border-violet-500/20 text-violet-400">
              <span class="w-2 h-2 rounded-full bg-violet-400 animate-pulse"></span>
              {statusMsg || "Executando operação..."}
            </span>
          {:else if finished}
            <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-emerald-500/10 border border-emerald-500/20 text-emerald-400">
              ✓ Processo concluído com sucesso
            </span>
          {:else if errorMsg}
            <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-rose-500/10 border border-rose-500/20 text-rose-400">
              ✗ Falha na execução
            </span>
          {/if}
        </div>

        <ButtonOrange disabled={isRunning} onclick={() => { closeListeners(); show = false; }}>
          {t("common.close")}
        </ButtonOrange>
      </div>

    </div>
  </div>
{/if}
