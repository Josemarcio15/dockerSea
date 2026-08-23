<script lang="ts">
  import { onDestroy } from "svelte";
  import { t } from "$lib/stores/locale.svelte";

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
  let terminalElement = $state<HTMLDivElement | null>(null);

  // Determine mode: generic SSE or image pull
  const isGenericMode = $derived(!!sseUrl);
  const isRunning = $derived(!finished && !errorMsg);
  const modalTitle = $derived(
    title || (imageName ? `Pull: ${imageName}` : "Progresso"),
  );

  function closeEventSource() {
    if (eventSource) {
      eventSource.close();
      eventSource = null;
    }
  }

  $effect(() => {
    if (terminalElement && logs.length) {
      terminalElement.scrollTop = terminalElement.scrollHeight;
    }
  });

  function startGenericSse() {
    errorMsg = "";
    finished = false;
    completeResult = null;
    statusMsg = t("terminal_modal.starting");
    logs = [];

    closeEventSource();

    eventSource = new EventSource(sseUrl);

    eventSource.addEventListener("progress", (event: MessageEvent) => {
      try {
        const data = JSON.parse(event.data);

        // Dados de progresso de layers (mesmo formato do /api/images/pull)
        if (data.id) {
          const layerId = data.id;
          const currentStatus = data.status || "";
          const progressDetail = data.progressDetail || {};
          let current = progressDetail.current || 0;
          let total = progressDetail.total || 0;

          const isRepoPullStart =
            currentStatus.includes("Pulling from") || layerId === "latest";
          if (currentStatus === "Already exists" || isRepoPullStart) {
            if (layers[layerId]) {
              delete layers[layerId];
            }
            if (lastLoggedStatus[layerId] !== currentStatus) {
              lastLoggedStatus[layerId] = currentStatus;
              logs.push(`→ [${layerId}] ${currentStatus}`);
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

          const skipLogging = ["Downloading", "Extracting"].includes(
            currentStatus,
          );
          if (!skipLogging && lastLoggedStatus[layerId] !== currentStatus) {
            lastLoggedStatus[layerId] = currentStatus;
            logs.push(`→ [${layerId}] ${currentStatus}`);
          }
        } else if (data.line) {
          // Linha de texto simples
          logs.push(data.line);
          statusMsg = data.line;
        } else if (data.status) {
          statusMsg = data.status;
          const isLogged = logs.some((l) => l.includes(data.status));
          if (!isLogged) {
            logs.push(`→ ${data.status}`);
          }
        }
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
        if (data.success) {
          logs.push(`${data.message}`);
        } else {
          logs.push(`${data.message}`);
        }
        closeEventSource();
        oncomplete();
      } catch (e) {
        console.error("Error parsing complete event:", e);
      }
    });

    eventSource.onerror = () => {
      if (!finished) {
        errorMsg = "Erro de conexão com o servidor.";
        logs.push("→ Erro de conexão com o servidor.");
        closeEventSource();
      }
    };
  }

  function startPull() {
    errorMsg = "";
    finished = false;
    completeResult = null;
    statusMsg = t("terminal_modal.starting_download");
    layers = {};
    logs = [];
    lastLoggedStatus = {};

    closeEventSource();

    const url = `/api/images/pull?imageName=${encodeURIComponent(imageName)}`;
    eventSource = new EventSource(url);

    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        if (data.error) {
          errorMsg = data.error;
          logs.push(`→ Erro: ${data.error}`);
          closeEventSource();
          return;
        }

        if (data.status === "pull_complete") {
          statusMsg = data.message || "Completo!";
          logs.push(
            `→ Status: ${data.message || "Imagem baixada com sucesso!"}`,
          );
          finished = true;
          for (const key of Object.keys(layers)) {
            layers[key].status = "Complete";
            layers[key].current = 1;
            layers[key].total = 1;
          }
          closeEventSource();
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
          if (currentStatus === "Already exists" || isRepoPullStart) {
            if (layers[layerId]) {
              delete layers[layerId];
            }
            if (lastLoggedStatus[layerId] !== currentStatus) {
              lastLoggedStatus[layerId] = currentStatus;
              logs.push(`→ [${layerId}] ${currentStatus}`);
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

          const skipLogging = ["Downloading", "Extracting"].includes(
            currentStatus,
          );
          if (!skipLogging && lastLoggedStatus[layerId] !== currentStatus) {
            lastLoggedStatus[layerId] = currentStatus;
            logs.push(`→ [${layerId}] ${currentStatus}`);
          }
        } else if (data.status) {
          statusMsg = data.status;
          const isLogged = logs.some((l) => l.includes(data.status));
          if (!isLogged) {
            logs.push(`→ ${data.status}`);
          }
        }
      } catch (e) {
        console.error("Error parsing progress event:", e);
      }
    };

    eventSource.onerror = () => {
      if (!finished) {
        errorMsg = "Erro de conexão ao buscar progresso de download.";
        logs.push("→ Erro de conexão ao buscar progresso de download.");
        closeEventSource();
      }
    };
  }

  $effect(() => {
    if (show) {
      if (isGenericMode) {
        startGenericSse();
      } else if (imageName) {
        startPull();
      }
    }
    return () => {
      closeEventSource();
    };
  });

  onDestroy(() => {
    closeEventSource();
  });

  let layerList = $derived(
    Object.entries(layers).map(([id, info]) => ({
      id,
      ...info,
      percentage: info.total
        ? Math.round((info.current / info.total) * 100)
        : 0,
    })),
  );
</script>

{#if show}
  <div
    class="fixed inset-0 bg-black/75 backdrop-blur-xs flex items-center justify-center z-50 p-4"
  >
    <div
      class="bg-[#0f172a] border border-[#1e293b] rounded-2xl w-196 max-w-full p-6 shadow-2xl flex flex-col max-h-[90vh] text-slate-200 animate-scaleIn"
    >
      <!-- Title Bar -->
      <div
        class="flex justify-between items-center pb-4 border-b border-[#1e293b]"
      >
        <div class="flex items-center gap-3">
          <h2 class="text-lg font-bold text-white tracking-wide">
            {modalTitle}
          </h2>
        </div>
        <button
          type="button"
          class="text-slate-400 hover:text-white text-xl bg-transparent border-none transition-colors {isRunning
            ? 'cursor-not-allowed opacity-30'
            : 'cursor-pointer'}"
          disabled={isRunning}
          onclick={() => {
            closeEventSource();
            show = false;
          }}
        >
          ✕
        </button>
      </div>

      <!-- Content Area -->
      <div class="flex-1 overflow-y-auto py-5 space-y-5 scrollbar-thin">
        <!-- 1. Layer Progress Card (quando houver layers sendo baixadas) -->
        {#if !isGenericMode || layerList.length > 0 || logs.length === 0}
          <div
            class="border border-[#1e293b] rounded-xl p-4 bg-[#0c101a]/40 space-y-4"
          >
            <div
              class="flex justify-between items-center text-xs font-bold tracking-wider text-slate-400 uppercase"
            >
              <span class="flex items-center gap-1.5 text-white">
                {t("terminal_modal.layers_progress")}
              </span>
              <span class="text-blue-400 font-mono">
                {t("terminal_modal.layers_count", { count: layerList.length })}
              </span>
            </div>

            {#if layerList.length === 0}
              <div
                class="flex flex-col items-center justify-center py-6 space-y-2"
              >
                <div
                  class="animate-spin rounded-full h-6 w-6 border-2 border-violet-500 border-t-transparent"
                ></div>
                <p class="text-xs text-slate-400 font-medium">
                  {statusMsg}
                </p>
              </div>
            {:else}
              <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3">
                {#each layerList as layer (layer.id)}
                  <div
                    class="bg-[#0b0f19] border border-[#1e293b] rounded-lg p-3 flex flex-col justify-between min-h-17"
                  >
                    <div class="flex justify-between items-center text-xs mb-2">
                      <span
                        class="font-mono text-blue-400 font-semibold truncate max-w-25"
                        title={layer.id}
                      >
                        {layer.id}
                      </span>
                      <span
                        class="text-slate-350 font-medium truncate max-w-20"
                        title={layer.status}
                      >
                        {layer.status === "Pull complete"
                          ? "Complete"
                          : layer.status}
                      </span>
                    </div>
                    <div class="space-y-1">
                      <div
                        class="w-full h-1.5 bg-[#1e293b] rounded-full overflow-hidden"
                      >
                        <div
                          class="h-full bg-emerald-500 rounded-full transition-all duration-300"
                          style="width: {layer.status === 'Pull complete'
                            ? 100
                            : layer.percentage}%"
                        ></div>
                      </div>
                    </div>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/if}

        <!-- 2. Terminal Logs Card -->
        <div
          bind:this={terminalElement}
          class="border border-[#1e293b] rounded-xl p-4 bg-[#070b15] font-mono text-xs text-emerald-400 {isGenericMode
            ? 'h-96'
            : 'h-64'} overflow-y-auto space-y-1 scrollbar-thin"
        >
          {#if logs.length === 0}
            <div class="text-slate-500 italic">
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
      <div class="flex justify-end pt-4 border-t border-[#1e293b] items-center">
        {#if !finished && !errorMsg}
          <div
            class="flex items-center gap-2 text-xs text-slate-400 mr-auto font-medium"
          >
            <span class="relative flex h-2 w-2">
              <span
                class="animate-ping absolute inline-flex h-full w-full rounded-full bg-violet-400 opacity-75"
              ></span>
              <span
                class="relative inline-flex rounded-full h-2 w-2 bg-violet-500"
              ></span>
            </span>
            {statusMsg}
          </div>
        {/if}

        <button
          type="button"
          class="px-6 py-2 rounded-lg text-sm font-bold text-white bg-emerald-600 hover:bg-emerald-700 transition-colors shadow-md shadow-emerald-500/20 {isRunning
            ? 'cursor-not-allowed opacity-30'
            : 'cursor-pointer'}"
          disabled={isRunning}
          onclick={() => {
            closeEventSource();
            show = false;
          }}
        >
          {isRunning ? t("terminal_modal.wait") : t("terminal_modal.close")}
        </button>
      </div>
    </div>
  </div>
{/if}
