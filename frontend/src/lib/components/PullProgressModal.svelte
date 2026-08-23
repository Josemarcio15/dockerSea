<script lang="ts">
  import { onDestroy } from "svelte";
  import { t } from "$lib/stores/locale.svelte";
  import { Events } from "@wailsio/runtime";

  let {
    show = $bindable(false),
    imageName = "",
    sseUrl = "",
    title = "",
    oncomplete = () => {},
  } = $props();

  interface LayerProgress {
    status: string;
    progress: string;
    current: number;
    total: number;
  }

  let layers = $state<Record<string, LayerProgress>>({});
  let logs = $state<string[]>([]);
  let lastLoggedStatus: Record<string, string> = {};

  let statusMsg = $state("");
  let errorMsg = $state("");
  let finished = $state(false);
  let completeResult = $state<{ success: boolean; message: string } | null>(
    null,
  );

  let eventSource: EventSource | null = null;
  let unsubscribeProgress: (() => void) | null = null;
  let unsubscribeComplete: (() => void) | null = null;
  let terminalElement = $state<HTMLDivElement | null>(null);

  // Determine mode: generic SSE or image pull
  const isGenericMode = $derived(!!sseUrl);
  const isRunning = $derived(!finished && !errorMsg);
  const modalTitle = $derived(
    title || (imageName ? `Pull: ${imageName}` : "Progresso"),
  );

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
    if (unsubscribeProgress) {
      unsubscribeProgress();
      unsubscribeProgress = null;
    }
    if (unsubscribeComplete) {
      unsubscribeComplete();
      unsubscribeComplete = null;
    }
  }

  function handleProgressData(data: any) {
    if (data.error) {
      errorMsg = data.error;
      logs.push(`→ Erro: ${data.error}`);
      closeListeners();
      return;
    }

    if (data.status === "pull_complete" || data.status === "complete") {
      statusMsg = data.message || "Download e extração concluídos!";
      logs.push(`→ Status: ${data.message || "Imagem baixada com sucesso!"}`);
      finished = true;
      for (const key of Object.keys(layers)) {
        layers[key].status = "Pull complete";
        layers[key].current = 1;
        layers[key].total = 1;
      }
      closeListeners();
      oncomplete();
      return;
    }

    if (data.id) {
      const layerId = data.id;
      const currentStatus = data.status || "";
      const progressDetail = data.progressDetail || {};
      let current = progressDetail.current || 0;
      let total = progressDetail.total || 0;

      const isRepoPullStart =
        currentStatus.includes("Pulling from") || layerId === "latest";
      if (isRepoPullStart) {
        if (lastLoggedStatus[layerId] !== currentStatus) {
          lastLoggedStatus[layerId] = currentStatus;
          logs.push(`→ [${layerId}] ${currentStatus}`);
        }
        return;
      }

      if (currentStatus === "Already exists") {
        layers[layerId] = {
          status: "Already exists",
          progress: "100%",
          current: 1,
          total: 1,
        };
        if (lastLoggedStatus[layerId] !== currentStatus) {
          lastLoggedStatus[layerId] = currentStatus;
          logs.push(`→ [${layerId}] Camada já existe no cache`);
        }
        return;
      }

      const isComplete = currentStatus === "Pull complete";
      if (isComplete) {
        current = 1;
        total = 1;
      }

      if (!layers[layerId]) {
        layers[layerId] = {
          status: currentStatus,
          progress: data.progress || "",
          current,
          total,
        };
      } else {
        layers[layerId].status = currentStatus;
        if (data.progress) layers[layerId].progress = data.progress;
        if (current || isComplete) layers[layerId].current = current;
        if (total || isComplete) layers[layerId].total = total;
      }

      const skipLogging = ["Downloading", "Extracting"].includes(currentStatus);
      if (!skipLogging && lastLoggedStatus[layerId] !== currentStatus) {
        lastLoggedStatus[layerId] = currentStatus;
        logs.push(`→ [${layerId}] ${currentStatus}`);
      }
    } else if (data.line) {
      logs.push(data.line);
      statusMsg = data.line;
    } else if (data.status) {
      statusMsg = data.status;
      const isLogged = logs.some((l) => l.includes(data.status));
      if (!isLogged) {
        logs.push(`→ ${data.status}`);
      }
    }
  }

  function startGenericSse() {
    errorMsg = "";
    finished = false;
    completeResult = null;
    statusMsg = t("terminal_modal.starting");
    logs = [];
    closeListeners();

    eventSource = new EventSource(sseUrl);
    eventSource.addEventListener("progress", (event: MessageEvent) => {
      try {
        const data = JSON.parse(event.data);
        handleProgressData(data);
      } catch (e) {
        console.error("Error parsing progress event:", e);
      }
    });

    eventSource.addEventListener("complete", (event: MessageEvent) => {
      try {
        const data = JSON.parse(event.data);
        finished = true;
        completeResult = { success: data.success, message: data.message };
        statusMsg = data.message || (data.success ? "Concluído!" : "Falhou!");
        logs.push(data.message || (data.success ? "✓ Concluído!" : "✗ Falhou!"));
        closeListeners();
        oncomplete();
      } catch (e) {
        console.error("Error parsing complete event:", e);
      }
    });

    eventSource.onerror = () => {
      if (!finished) {
        errorMsg = "Erro de conexão com o servidor.";
        logs.push("→ Erro de conexão com o servidor.");
        closeListeners();
      }
    };
  }

  function startWailsEventsListener() {
    errorMsg = "";
    finished = false;
    completeResult = null;
    statusMsg = t("terminal_modal.starting_download");
    layers = {};
    logs = [];
    lastLoggedStatus = {};

    closeListeners();

    unsubscribeProgress = Events.On("docker:image:pull:progress", (event: any) => {
      try {
        const data = event?.data ?? event;
        handleProgressData(data);
      } catch (e) {
        console.error("Erro ao processar evento de progresso:", e);
      }
    });

    unsubscribeComplete = Events.On("docker:image:pull:complete", (event: any) => {
      try {
        const data = event?.data ?? event;
        finished = true;
        completeResult = { success: data.success, message: data.message };
        if (data.success) {
          statusMsg = data.message || "Concluído!";
          logs.push(data.message || "✓ Download e extração concluídos com sucesso!");
          oncomplete();
        } else {
          errorMsg = data.message || "Erro no download da imagem";
          logs.push(`✗ Erro: ${data.message || "Falha no download"}`);
        }
        closeListeners();
      } catch (e) {
        console.error("Erro ao processar evento de conclusão:", e);
      }
    });
  }

  $effect(() => {
    if (show) {
      if (isGenericMode) {
        startGenericSse();
      } else if (imageName) {
        startWailsEventsListener();
      }
    }
    return () => {
      closeListeners();
    };
  });

  onDestroy(() => {
    closeListeners();
  });

  // Cálculo coerente do progresso por camada e progresso global unificado
  function getLayerCalculatedProgress(info: LayerProgress): { percent: number; label: string; badgeColor: string } {
    const s = info.status.toLowerCase();
    if (s === "pull complete" || s === "already exists") {
      return { percent: 100, label: "Concluído", badgeColor: "bg-emerald-500/15 text-emerald-400 border-emerald-500/30" };
    }
    if (s.includes("download complete")) {
      return { percent: 60, label: "Download 100% • Extraindo...", badgeColor: "bg-blue-500/15 text-blue-400 border-blue-500/30" };
    }
    if (s.includes("extracting")) {
      const rawPct = info.total > 0 ? info.current / info.total : 0;
      const combined = Math.min(99, Math.round(60 + rawPct * 40));
      return { percent: combined, label: `Extraindo (${Math.round(rawPct * 100)}%)`, badgeColor: "bg-amber-500/15 text-amber-400 border-amber-500/30" };
    }
    if (s.includes("downloading")) {
      const rawPct = info.total > 0 ? info.current / info.total : 0;
      const combined = Math.min(59, Math.round(rawPct * 60));
      return { percent: combined, label: `Baixando (${Math.round(rawPct * 100)}%)`, badgeColor: "bg-violet-500/15 text-violet-400 border-violet-500/30" };
    }
    if (s.includes("waiting")) {
      return { percent: 0, label: "Na fila", badgeColor: "bg-slate-700/30 text-slate-400 border-slate-700/50" };
    }
    return { percent: 10, label: info.status || "Processando", badgeColor: "bg-slate-700/30 text-slate-300 border-slate-700/50" };
  }

  let layerList = $derived(
    Object.entries(layers).map(([id, info]) => {
      const { percent, label, badgeColor } = getLayerCalculatedProgress(info);
      return {
        id,
        ...info,
        percentage: percent,
        label,
        badgeColor,
      };
    }),
  );

  // Progresso total abstrato e fluido de 0% a 100%
  let calculatedPercent = $state(0);

  $effect(() => {
    if (finished) {
      calculatedPercent = 100;
    } else if (layerList.length > 0) {
      const sum = layerList.reduce((acc, l) => acc + l.percentage, 0);
      const avg = Math.round(sum / layerList.length);
      if (avg > calculatedPercent) {
        calculatedPercent = avg;
      }
    } else if (isRunning && calculatedPercent === 0) {
      calculatedPercent = 5;
    }
  });

  const displayPercent = $derived(finished ? 100 : Math.min(99, calculatedPercent));
</script>

{#if show}
  <div
    class="fixed inset-0 bg-black/80 backdrop-blur-md flex items-center justify-center z-50 p-4 animate-fadeIn"
  >
    <div
      class="bg-[#0b101d] border border-slate-800/90 rounded-3xl w-200 max-w-full p-6 shadow-2xl flex flex-col max-h-[92vh] text-slate-200"
    >
      <!-- Title Bar -->
      <div class="flex justify-between items-center pb-4 border-b border-slate-800/80">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-xl bg-violet-600/20 border border-violet-500/30 flex items-center justify-center text-violet-400 font-bold">
            ⬇️
          </div>
          <div>
            <h2 class="text-base font-bold text-white tracking-wide">
              {modalTitle}
            </h2>
            <p class="text-xs text-slate-400">
              {isRunning ? "Baixando imagem Docker" : finished ? "Processo finalizado com sucesso" : "Aguardando"}
            </p>
          </div>
        </div>
        <button
          type="button"
          class="w-8 h-8 rounded-xl bg-slate-850 hover:bg-slate-800 border border-slate-700 text-slate-400 hover:text-white flex items-center justify-center transition-colors {isRunning
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

      <!-- Content Area -->
      <div class="flex-1 overflow-y-auto py-5 space-y-4 scrollbar-thin">
        
        <!-- Barra de Progresso Geral 0 - 100% -->
        <div class="border border-slate-800 rounded-2xl p-5 bg-[#070b14]/80 space-y-3.5 shadow-inner">
          <div class="flex items-center justify-between text-xs font-bold">
            <span class="text-white text-sm">Progresso Geral</span>
            <span class="font-mono text-base font-bold bg-linear-to-r from-violet-400 to-emerald-400 bg-clip-text text-transparent">
              {displayPercent}%
            </span>
          </div>

          <div class="w-full h-3.5 bg-slate-900 rounded-full overflow-hidden border border-slate-800 p-0.5">
            <div
              class="h-full rounded-full transition-all duration-300 bg-linear-to-r from-violet-600 via-blue-500 to-emerald-500 shadow-sm shadow-emerald-500/30"
              style="width: {displayPercent}%"
            ></div>
          </div>
        </div>

        <!-- 2. Terminal de Logs com visual refinado -->
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
      </div>

      <!-- Footer Bar -->
      <div class="flex justify-between items-center pt-4 border-t border-slate-800/80">
        <div class="flex items-center gap-2 text-xs font-semibold">
          {#if isRunning}
            <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-violet-500/10 border border-violet-500/20 text-violet-400">
              <span class="w-2 h-2 rounded-full bg-violet-400 animate-pulse"></span>
              Baixando e extraindo dados
            </span>
          {:else if finished}
            <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-emerald-500/10 border border-emerald-500/20 text-emerald-400">
              ✓ Processo concluído com sucesso
            </span>
          {/if}
        </div>

        <button
          type="button"
          class="px-6 py-2.5 rounded-xl text-xs font-bold text-white transition-all shadow-md {isRunning
            ? 'bg-slate-700 text-slate-400 cursor-not-allowed opacity-50'
            : 'bg-emerald-600 hover:bg-emerald-500 shadow-emerald-600/20 cursor-pointer'}"
          disabled={isRunning}
          onclick={() => {
            closeListeners();
            show = false;
          }}
        >
          {isRunning ? t("terminal_modal.wait") : t("terminal_modal.close")}
        </button>
      </div>
    </div>
  </div>
{/if}
