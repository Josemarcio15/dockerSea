<script lang="ts">
  import { onDestroy } from "svelte";
  import { t } from "$shared/stores/locale.svelte";
  import { Events } from "@wailsio/runtime";
  import { ButtonOrange } from "$shared/components/buttons";
  import TerminalLogsView from "$shared/components/TerminalLogsView.svelte";

  let {
    show = $bindable(false),
    title = "",
    eventPrefix = "docker:image:pull",
    sseUrl = "",
    oncomplete = () => {},
  }: {
    show: boolean;
    title?: string;
    eventPrefix?: string;
    sseUrl?: string;
    oncomplete?: () => void;
  } = $props();

  let logs = $state<string[]>([]);
  let statusMsg = $state("");
  let errorMsg = $state("");
  let finished = $state(false);
  let percent = $state(0);

  let eventSource: EventSource | null = null;
  let unsubscribeEvents: (() => void) | null = null;
  let lastLoggedStatus: Record<string, string> = {};

  const isGenericMode = $derived(!!sseUrl);
  const isRunning = $derived(!finished && !errorMsg);
  const modalTitle = $derived(title || "Progresso da Operação");
  const displayPercent = $derived(
    finished ? 100 : Math.min(99, percent || (isRunning ? 15 : 0)),
  );

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
      logs = [...logs, `✗ Erro: ${data.error}`];
      closeListeners();
      return;
    }

    if (data.status === "pull_complete" || data.status === "complete") {
      statusMsg = data.message || "Operação concluída com sucesso!";
      logs = [...logs, `✓ ${data.message || "Operação concluída com sucesso!"}`];
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
          logs = [...logs, `→ [${layerId}] ${currentStatus}`];
        }
      }
      if (data.progressDetail?.total > 0) {
        const rawPct = Math.round(
          (data.progressDetail.current / data.progressDetail.total) * 100,
        );
        if (rawPct > percent) percent = Math.min(95, rawPct);
      }
    } else if (data.message) {
      statusMsg = data.message;
      const phase = (data.phase || "").toLowerCase();
      const phasePrefix = data.phase ? `[${data.phase.toUpperCase()}] ` : "";
      logs = [...logs, `→ ${phasePrefix}${data.message}`];

      // Gradual progression according to the deploy or build phase
      if (phase === "preparing" && percent < 15) percent = 15;
      else if (phase === "uploading" && percent < 30) percent = 30;
      else if (phase === "validating" && percent < 45) percent = 45;
      else if (phase === "building") {
        if (percent < 50) percent = 50;
        else if (percent < 90) percent += 2;
      } else if (phase === "starting" && percent < 92) percent = 92;
      else if (phase === "checking" && percent < 96) percent = 96;
      else if (phase === "cleaning" && percent < 98) percent = 98;
    } else if (data.line) {
      logs = [...logs, data.line];
      statusMsg = data.line;
      if (percent < 90) percent += 1;
    } else if (data.status) {
      statusMsg = data.status;
      if (!logs.some((l) => l.includes(data.status))) {
        logs = [...logs, `→ ${data.status}`];
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
      try {
        handleProgressData(JSON.parse(ev.data));
      } catch {}
    });
    eventSource.addEventListener("complete", (ev: MessageEvent) => {
      try {
        const data = JSON.parse(ev.data);
        finished = true;
        percent = 100;
        statusMsg = data.message || (data.success ? "Concluído!" : "Falhou!");
        logs = [
          ...logs,
          data.message || (data.success ? "✓ Concluído!" : "✗ Falhou!"),
        ];
        closeListeners();
        oncomplete();
      } catch {}
    });
    eventSource.onerror = () => {
      if (!finished) {
        errorMsg = "Erro de conexão com o servidor.";
        logs = [...logs, "✗ Erro de conexão com o servidor."];
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
        const msg = d?.projectName
          ? `Iniciando: ${d.projectName}`
          : "Operação iniciada...";
        statusMsg = msg;
        logs = [...logs, `🚀 ${msg}`];
      }),
      Events.On(`${eventPrefix}:progress`, (ev: any) => {
        handleProgressData(ev?.data ?? ev);
      }),
      Events.On(`${eventPrefix}:failed`, (ev: any) => {
        const d = ev?.data ?? ev;
        finished = true;
        errorMsg = d?.message || "Erro na operação";
        logs = [...logs, `✗ [ERRO] ${d?.message || "Falha na operação"}`];
        closeListeners();
      }),
      Events.On(`${eventPrefix}:complete`, (ev: any) => {
        const d = ev?.data ?? ev;
        finished = true;
        percent = 100;
        if (d?.success !== false) {
          statusMsg = d?.message || "Concluído com sucesso!";
          logs = [
            ...logs,
            d?.message || "✓ Operação concluída com sucesso!",
          ];
          oncomplete();
        } else {
          errorMsg = d?.message || "Erro na execução";
          logs = [...logs, `✗ Erro: ${d?.message || "Falha na operação"}`];
        }
        closeListeners();
      }),
    ];

    unsubscribeEvents = () => {
      unsubs.forEach((unsub) => {
        if (typeof unsub === "function") unsub();
      });
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
  <div
    class="fixed inset-0 bg-black/80 backdrop-blur-md flex items-center justify-center z-50 p-4 animate-fadeIn"
  >
    <div
      class="bg-[#0b101d] border border-slate-800/90 rounded-3xl w-[720px] max-w-full h-[580px] max-h-[92vh] p-6 shadow-2xl flex flex-col text-slate-200 gap-4 overflow-hidden"
    >
      <!-- Top Title -->
      <div
        class="flex justify-between items-center pb-3 border-b border-slate-800/80 shrink-0"
      >
        <div class="flex items-center gap-3 min-w-0">
          <div
            class="w-3 h-3 rounded-full shrink-0 {isRunning
              ? 'bg-amber-400 animate-pulse'
              : finished
                ? 'bg-emerald-400'
                : 'bg-rose-500'}"
          ></div>
          <div class="min-w-0">
            <h3 class="text-base font-bold text-white tracking-wide truncate">
              {modalTitle}
            </h3>
            <p class="text-xs text-slate-400 mt-0.5 truncate">
              {isRunning
                ? statusMsg || "Processando tarefa..."
                : finished
                  ? "Processo finalizado com sucesso"
                  : "Operação interrompida"}
            </p>
          </div>
        </div>
        <button
          type="button"
          class="w-8 h-8 rounded-xl bg-slate-850 hover:bg-slate-800 border border-slate-700 text-slate-400 hover:text-white flex items-center justify-center transition-colors shrink-0 {isRunning
            ? 'cursor-not-allowed opacity-30'
            : 'cursor-pointer'}"
          disabled={isRunning}
          onclick={() => {
            closeListeners();
            show = false;
          }}
        >
          ✕
        </button>
      </div>

      <!-- General Progress Bar 0-100% -->
      <div
        class="border border-slate-800 rounded-2xl p-4 bg-[#070b14]/80 space-y-2.5 shadow-inner shrink-0"
      >
        <div class="flex items-center justify-between text-xs font-bold">
          <span class="text-white">Progresso Geral</span>
          <span
            class="font-mono text-sm font-bold bg-linear-to-r from-violet-400 to-emerald-400 bg-clip-text text-transparent"
          >
            {displayPercent}%
          </span>
        </div>
        <div
          class="w-full h-3 bg-slate-900 rounded-full overflow-hidden border border-slate-800 p-0.5"
        >
          <div
            class="h-full rounded-full transition-all duration-300 bg-linear-to-r from-violet-600 via-blue-500 to-emerald-500 shadow-sm shadow-emerald-500/30"
            style="width: {displayPercent}%"
          ></div>
        </div>
      </div>

      <!-- Reusable Logs Terminal (Occupies full remaining vertical space without jumping) -->
      <div class="flex-1 flex flex-col min-h-0 overflow-hidden">
        <TerminalLogsView {logs} maxHeight="h-full max-h-full" />
      </div>

      <!-- Footer / Action -->
      <div
        class="flex justify-between items-center pt-3 border-t border-slate-800/80"
      >
        <div class="flex items-center gap-2 text-xs font-semibold">
          {#if isRunning}
            <span
              class="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-violet-500/10 border border-violet-500/20 text-violet-400"
            >
              <span
                class="w-2 h-2 rounded-full bg-violet-400 animate-pulse"
              ></span>
              {statusMsg || "Executando operação..."}
            </span>
          {:else if finished}
            <span
              class="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-emerald-500/10 border border-emerald-500/20 text-emerald-400"
            >
              ✓ Processo concluído com sucesso
            </span>
          {:else if errorMsg}
            <span
              class="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-rose-500/10 border border-rose-500/20 text-rose-400"
            >
              ✗ Falha na execução
            </span>
          {/if}
        </div>

        <ButtonOrange
          disabled={isRunning}
          onclick={() => {
            closeListeners();
            show = false;
          }}
        >
          {t("common.close")}
        </ButtonOrange>
      </div>
    </div>
  </div>
{/if}
