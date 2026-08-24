<script lang="ts">
  import { getLocale, t } from "$lib/stores/locale.svelte";
  import { ButtonOrange } from "$lib/components/buttons";
  import { tick } from "svelte";

  interface LayerProgress {
    id: string;
    status: string;
    current_mb: number | null;
    total_mb: number | null;
    percent: number | null;
  }

  let {
    show = $bindable(false),
    title = "",
    loading = false,
    logs = [],
    console_id = "global-terminal-console",
  }: {
    show: boolean;
    title: string;
    loading: boolean;
    logs: string[];
    console_id?: string;
  } = $props();

  // Auto-scroll effect when logs update
  $effect(() => {
    if (logs.length > 0) {
      tick().then(() => {
        const el = document.getElementById(console_id);
        if (el) {
          el.scrollTop = el.scrollHeight;
          let parent = el.parentElement;
          while (parent && !parent.classList.contains("fixed")) {
            parent.scrollTop = parent.scrollHeight;
            parent = parent.parentElement;
          }
        }
      });
    }
  });

  // Parse layers
  let parsedData = $derived.by(() => {
    let layers_order: string[] = [];
    let layers_map: Record<string, LayerProgress> = {};
    let filtered_lines: Array<{ line: string; colorClass: string }> = [];
    let has_layers = false;

    for (const line of logs) {
      const trimmed = line.trim();
      if (!trimmed) continue;

      let is_layer_progress = false;
      let is_repetitive = false;

      if (trimmed.startsWith("[")) {
        const close_idx = trimmed.indexOf("]");
        if (close_idx !== -1) {
          const id = trimmed.substring(1, close_idx);
          const rest = trimmed.substring(close_idx + 1).trim();

          let status = rest;
          let current_mb: number | null = null;
          let total_mb: number | null = null;
          let percent: number | null = null;

          if (rest.startsWith("Downloading")) {
            is_layer_progress = true;
            is_repetitive = true;
            status = "Downloading";
            const mb_part = rest.replace("Downloading ", "");
            const parts = mb_part.split("/");
            if (parts.length === 2) {
              const cur_str = parts[0].replace("MB", "").trim();
              const tot_str = parts[1].replace("MB", "").trim();
              const cur = parseFloat(cur_str);
              const tot = parseFloat(tot_str);
              if (!isNaN(cur) && !isNaN(tot)) {
                current_mb = cur;
                total_mb = tot;
                if (tot > 0) {
                  percent = Math.min((cur / tot) * 100, 100);
                }
              }
            }
          } else if (rest.startsWith("Extracting")) {
            is_layer_progress = true;
            is_repetitive = true;
            status = "Extracting";
          } else if (rest === "Pull complete" || rest === "Download complete") {
            is_layer_progress = true;
            status = rest;
            percent = 100;
          } else if (rest === "Waiting") {
            is_layer_progress = true;
            status = "Waiting";
          } else if (rest === "Pulling fs layer") {
            is_layer_progress = true;
            status = "Pulling fs layer";
          }

          if (is_layer_progress) {
            has_layers = true;
            if (!layers_map[id]) {
              layers_order.push(id);
            }
            layers_map[id] = {
              id,
              status,
              current_mb,
              total_mb,
              percent,
            };
          }
        }
      }

      if (is_repetitive) {
        continue;
      }

      let colorClass = "text-green-400";
      if (
        trimmed.includes("complete") ||
        trimmed.includes("Sucesso") ||
        trimmed.includes("Ok") ||
        trimmed.includes("Downloaded")
      ) {
        colorClass = "text-emerald-400";
      } else if (
        trimmed.includes("Erro") ||
        trimmed.includes("Error") ||
        trimmed.includes("❌")
      ) {
        colorClass = "text-red-400";
      }

      filtered_lines.push({ line, colorClass });
    }

    return { layers_order, layers_map, filtered_lines, has_layers };
  });
</script>

{#if show}
  <div
    class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
  >
    <div
      class="bg-white dark:bg-[#111827] border border-slate-200 dark:border-slate-800 rounded-2xl shadow-2xl w-200 max-w-full max-h-[90vh] flex flex-col text-slate-800 dark:text-slate-100 overflow-hidden"
    >
      <!-- Modal Header -->
      <div
        class="flex justify-between items-center px-6 py-4 border-b border-slate-200 dark:border-slate-850 bg-slate-50 dark:bg-slate-900/50"
      >
        <h2
          class="text-base font-bold text-slate-800 dark:text-slate-100 flex items-center gap-2"
        >
          {title}
        </h2>
        <button
          class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 text-lg bg-transparent border-none cursor-pointer transition-colors"
          onclick={() => (show = false)}
        >
          ✕
        </button>
      </div>

      <!-- Modal Body -->
      <div class="p-6 overflow-y-auto flex-1 flex flex-col min-h-0 space-y-4">
        {#if loading && logs.length === 0}
          <div class="flex flex-col items-center justify-center py-12 px-4 space-y-3">
            <div class="w-10 h-10 rounded-full border-4 border-violet-500/20 border-t-violet-600 animate-spin"></div>
            <span class="text-sm font-medium text-slate-500 dark:text-slate-400 animate-pulse">
              {t("terminal_modal.starting_process")}
            </span>
          </div>
        {:else}
          {#if loading}
            <div
              class="flex items-center gap-2 text-xs text-blue-600 dark:text-blue-400 font-semibold animate-pulse"
            >
              <span
                class="animate-spin inline-block w-4 h-4 border-2 border-current border-t-transparent rounded-full"
              ></span>
              {t("terminal_modal.running_realtime")}
            </div>
          {/if}

          <!-- Layer progress panel -->
          {#if parsedData.has_layers}
            <div
              class="bg-slate-950 p-4 rounded-xl border border-slate-800 flex flex-col gap-3 font-mono text-xs text-slate-350 shadow-inner"
            >
              <div
                class="font-semibold border-b border-slate-800 pb-2 text-slate-400 flex justify-between items-center"
              >
                <span>{t("terminal_modal.layers_progress")}</span>
                <span class="text-blue-400 font-normal">
                  {t("terminal_modal.layers_count", {
                    count: parsedData.layers_order.length,
                  })}
                </span>
              </div>
              <div
                class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-2 max-h-40 overflow-y-auto pr-1"
              >
                {#each parsedData.layers_order as id (id)}
                  {@const layer = parsedData.layers_map[id]}
                  {@const pct = layer.percent ?? 0}
                  {@const is_done =
                    layer.status === "Pull complete" ||
                    layer.status === "Download complete"}
                  {@const bar_setup =
                    layer.status === "Extracting"
                      ? {
                          color:
                            "bg-linear-to-r from-amber-500 to-orange-500 animate-pulse",
                          width: "100%",
                          text: "Extracting...",
                        }
                      : is_done
                        ? {
                            color: "bg-emerald-500",
                            width: "100%",
                            text: "Complete",
                          }
                        : layer.status === "Downloading"
                          ? {
                              color:
                                "bg-linear-to-r from-blue-500 to-indigo-500",
                              width: `${pct}%`,
                              text: "",
                            }
                          : layer.status === "Pulling fs layer"
                            ? {
                                color: "bg-blue-600/30 animate-pulse",
                                width: "15%",
                                text: "Pulling...",
                              }
                            : layer.status === "Waiting"
                              ? {
                                  color: "bg-slate-700",
                                  width: "5%",
                                  text: "Waiting...",
                                }
                              : {
                                  color: "bg-slate-600",
                                  width: "0%",
                                  text: layer.status,
                                }}

                  <div
                    class="flex flex-col gap-1 p-2 rounded bg-slate-900/50 border border-slate-800"
                  >
                    <div class="flex justify-between items-center text-[10px]">
                      <span class="font-semibold text-blue-400">{id}</span>
                      <span class="text-slate-400"
                        >{bar_setup.text || layer.status}</span
                      >
                    </div>

                    <!-- Progress bar track -->
                    <div
                      class="w-full bg-slate-950 h-1.5 rounded overflow-hidden border border-slate-800"
                    >
                      <div
                        class="h-full rounded transition-all duration-150 {bar_setup.color}"
                        style="width: {bar_setup.width}"
                      ></div>
                    </div>

                    <!-- Detail metrics -->
                    <div
                      class="flex justify-between items-center text-[9px] text-slate-500"
                    >
                      {#if layer.current_mb !== null && layer.total_mb !== null}
                        <span
                          >{layer.current_mb.toFixed(
                            1,
                          )}/{layer.total_mb.toFixed(1)} MB</span
                        >
                      {:else}
                        <span></span>
                      {/if}
                      {#if layer.status === "Downloading"}
                        <span>{pct.toFixed(0)}%</span>
                      {:else}
                        <span></span>
                      {/if}
                    </div>
                  </div>
                {/each}
              </div>
            </div>
          {/if}

          <!-- Terminal Logs block -->
          <pre
            id={console_id}
            class="bg-[#020817] text-green-400 p-4 rounded-xl text-xs font-mono overflow-auto whitespace-pre-wrap max-h-[50vh] flex-1 shadow-inner border border-slate-950">
                        {#each parsedData.filtered_lines as item}
              <div class={item.colorClass}>→ {item.line}</div>
            {/each}
                    </pre>
        {/if}
      </div>

      <!-- Modal Footer -->
      <div
        class="flex justify-end px-6 py-4 border-t border-slate-200 dark:border-slate-850 bg-slate-50 dark:bg-slate-900/50"
      >
        <ButtonOrange onclick={() => (show = false)}>
          {t("common.close")}
        </ButtonOrange>
      </div>
    </div>
  </div>
{/if}
