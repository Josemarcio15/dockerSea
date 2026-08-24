<script lang="ts">
  import { onMount, tick } from "svelte";
  import { t, getLocale } from "$lib/stores/locale.svelte";
  import StatusBanner from "$lib/components/StatusBanner.svelte";
  import VpsSelectWarning from "$lib/components/VpsSelectWarning.svelte";
  import { notifySuccess, notifyError } from "$lib/stores/notification.svelte";
  import { Events } from "@wailsio/runtime";
  import {
    ButtonBlue,
    ButtonGreen,
    ButtonYellow,
    ButtonPurple,
    ButtonCyan,
    ButtonRed,
    ButtonPink,
    ButtonOrange,
  } from "$lib/components/buttons";
  import * as BuilderService from "../../../bindings/go-walis/internal/builder/service.js";

  let { data } = $props();

  let currentPath = $state("");
  let parentPath = $state<string | null>(null);
  let folders = $state<{ name: string; path: string }[]>([]);
  let hasDockerfile = $state(false);
  let hasDockerignore = $state(false);
  let loading = $state(false);
  let status = $state<"idle" | "ready" | "building" | "success" | "error">(
    "idle",
  );
  let logs = $state<string[]>([]);
  let builtImage = $state("");
  let errorMsg = $state("");
  let savedPaths = $state<string[]>([]);
  let logContainer = $state<HTMLDivElement | null>(null);
  let copied = $state(false);
  let unsubscribeProgress: (() => void) | null = null;
  let unsubscribeComplete: (() => void) | null = null;

  // Auto scroll logs when new entries arrive
  $effect(() => {
    if (logs.length > 0) {
      tick().then(() => {
        if (logContainer) {
          logContainer.scrollTop = logContainer.scrollHeight;
        }
      });
    }
  });

  async function copyLogs() {
    if (logs.length === 0) return;
    try {
      await navigator.clipboard.writeText(logs.join("\n"));
      copied = true;
      notifySuccess(t("builder.logs_copied") || "Logs copiados para a área de transferência!");
      setTimeout(() => {
        copied = false;
      }, 2000);
    } catch {
      const textarea = document.createElement("textarea");
      textarea.value = logs.join("\n");
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand("copy");
      document.body.removeChild(textarea);
      copied = true;
      notifySuccess(t("builder.logs_copied") || "Logs copiados para a área de transferência!");
      setTimeout(() => {
        copied = false;
      }, 2000);
    }
  }

  const isLocal = $derived(data?.activeVps?.connectionType === "local");

  async function loadSavedPaths() {
    try {
      const paths = await BuilderService.ListSavedPaths();
      savedPaths = (paths || []).map((p) => p.path);
    } catch {
      /* ignore */
    }
  }

  async function saveCurrentPath() {
    if (!currentPath || savedPaths.includes(currentPath)) return;
    const label = currentPath.split("/").pop()!.split("\\").pop()!;
    try {
      await BuilderService.SavePath(currentPath, label);
      savedPaths = [...savedPaths, currentPath];
      notifySuccess("Pasta salva com sucesso!");
    } catch (e: any) {
      notifyError(`Erro ao salvar pasta: ${e?.message || e}`);
    }
  }

  async function removeSavedPath(path: string, e?: Event) {
    if (e) {
      e.stopPropagation();
    }
    try {
      await BuilderService.RemoveSavedPath(path);
      savedPaths = savedPaths.filter((p) => p !== path);
      notifySuccess("Pasta removida dos favoritos!");
    } catch (e: any) {
      notifyError(`Erro ao remover pasta: ${e?.message || e}`);
    }
  }

  const folderName = $derived(
    currentPath ? currentPath.split("/").pop()!.split("\\").pop()! : "",
  );

  function folderNameFromPath(path: string) {
    return path.split("/").pop()!.split("\\").pop()!;
  }

  let customTag = $state("");

  function sanitizeTag(input: string): string {
    return input
      .toLowerCase()
      .replace(/[^a-z0-9._-]/g, "-")
      .replace(/-+/g, "-")
      .replace(/^-|-$/g, "");
  }

  const defaultTag = $derived(
    currentPath
      ? sanitizeTag(currentPath.split("/").pop()!.split("\\").pop()!)
      : "",
  );

  const effectiveTag = $derived(
    customTag.trim() ? sanitizeTag(customTag) : defaultTag,
  );

  const canBuild = $derived(
    effectiveTag.length > 0 && hasDockerfile && status === "ready",
  );

  async function browse(path: string = "") {
    loading = true;
    errorMsg = "";
    status = "idle";
    customTag = "";
    try {
      const data = await BuilderService.Browse(path);
      currentPath = data.currentPath;
      parentPath = data.parentPath;
      folders = data.folders || [];
      hasDockerfile = data.hasDockerfile;
      hasDockerignore = data.hasDockerignore;
      status = data.hasDockerfile ? "ready" : "idle";
    } catch {
      errorMsg = "Erro ao navegar para a pasta.";
    } finally {
      loading = false;
    }
  }

  function goHome() {
    browse();
  }
  function goUp() {
    if (parentPath) browse(parentPath);
  }
  function goToFolder(path: string) {
    browse(path);
  }

  async function doBuild() {
    if (!canBuild) return;
    status = "building";
    logs = [];
    errorMsg = "";
    builtImage = "";
    try {
      await BuilderService.Build(currentPath, effectiveTag, getLocale());
    } catch (error: any) {
      status = "error";
      errorMsg = error?.message || "Build falhou";
    }
  }

  function reset() {
    currentPath = "";
    parentPath = null;
    folders = [];
    hasDockerfile = false;
    status = "idle";
    loading = false;
    logs = [];
    builtImage = "";
    errorMsg = "";
    browse();
  }

  function goToImages() {
    window.location.href = `/images?highlight=${encodeURIComponent(builtImage)}`;
  }

  onMount(() => {
    unsubscribeProgress = Events.On("builder:progress", (event) => {
      if (event.data?.line) logs.push(event.data.line);
    });
    unsubscribeComplete = Events.On("builder:complete", (event) => {
      const result = event.data as { success?: boolean; image?: string; message?: string };
      if (result.success) {
        builtImage = result.image || effectiveTag;
        status = "success";
      } else {
        errorMsg = result.message || "Build falhou";
        status = "error";
      }
    });
    loadSavedPaths();
    browse();

    return () => {
      unsubscribeProgress?.();
      unsubscribeComplete?.();
    };
  });
</script>

{#if !data.activeVps}
  <VpsSelectWarning />
{:else}
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div
      class="inline-flex items-center px-5 py-2.5 rounded-2xl bg-linear-to-br from-violet-100 to-fuchsia-100 dark:from-violet-950/40 dark:to-fuchsia-950/40 border border-violet-200/50 dark:border-violet-800/50 self-start shadow-sm"
    >
      <h1
        class="text-2xl font-bold bg-linear-to-r from-violet-600 to-fuchsia-500 bg-clip-text text-transparent m-0 flex items-center gap-2"
      >
        {t("sidebar.builder")}
      </h1>
    </div>
  </div>

  <StatusBanner />

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
    <!-- Folder Browser Panel -->
      <div class="space-y-4">
        <div
          class="bg-white dark:bg-[#0b0f19] border border-slate-200/70 dark:border-slate-800/80 p-5 rounded-2xl shadow-sm space-y-4"
        >
          <h3
            class="text-sm font-bold text-slate-800 dark:text-slate-200 uppercase tracking-wider"
          >
            {t("builder.nav_title")}
          </h3>

          <!-- Navigation Buttons -->
          <div class="flex items-center gap-2 flex-wrap">
            <ButtonPurple
              size="xs"
              onclick={goHome}
            >
              {t("builder.nav_home")}
            </ButtonPurple>
            {#if parentPath}
              <ButtonYellow
                size="xs"
                onclick={goUp}
              >
                {t("builder.nav_up")}
              </ButtonYellow>
            {/if}
            {#each savedPaths as path}
              <div class="relative inline-block group">
                <ButtonGreen
                  size="xs"
                  onclick={() => browse(path)}
                  title={path}
                >
                  <span class="mr-1">📌</span>{folderNameFromPath(path)}
                </ButtonGreen>
                <span
                  class="absolute -top-1.5 -right-1.5 w-3.5 h-3.5 rounded-full bg-red-500 text-white text-[8px] flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer z-10"
                  role="button"
                  tabindex="0"
                  onclick={(e) => removeSavedPath(path, e)}
                  onkeydown={(e) => {
                    if (e.key === "Enter" || e.key === " ") {
                      removeSavedPath(path, e);
                    }
                  }}
                  title="Remover atalho"
                >
                  ✕
                </span>
              </div>
            {/each}
          </div>

          <!-- Current Path -->
          {#if currentPath}
            <div
              class="px-3 py-2 rounded-xl bg-slate-100 dark:bg-slate-900/60 border border-slate-200 dark:border-slate-800 text-xs font-mono text-slate-600 dark:text-slate-400 truncate"
              title={currentPath}
            >
              {currentPath}
            </div>
          {/if}

          <!-- Folder List -->
          <div
            class="border border-slate-200 dark:border-slate-800 rounded-xl max-h-64 overflow-y-auto bg-slate-50 dark:bg-slate-900/20"
          >
            {#if loading}
              <div class="flex items-center justify-center py-8">
                <span
                  class="animate-spin inline-block w-5 h-5 border-2 border-violet-500 border-t-transparent rounded-full"
                ></span>
              </div>
            {:else if folders.length === 0 && !hasDockerfile}
              <div class="text-center py-8 text-xs text-slate-400 italic">
                Nenhuma subpasta encontrada.
              </div>
            {:else}
              {#each folders as folder}
                <button
                  type="button"
                  class="w-full flex items-center gap-2.5 px-4 py-2.5 text-xs text-left hover:bg-slate-100 dark:hover:bg-slate-800/60 border-b border-slate-100 dark:border-slate-800/40 last:border-none cursor-pointer transition-colors"
                  onclick={() => goToFolder(folder.path)}
                >
                  <span class="text-base shrink-0">📂</span>
                  <span
                    class="font-mono font-semibold text-slate-700 dark:text-slate-300 truncate"
                    >{folder.name}</span
                  >
                </button>
              {/each}
            {/if}
          </div>

          <!-- Dockerfile & Dockerignore Badges -->
          {#if hasDockerfile}
            <div class="flex flex-wrap items-center gap-2">
              <span
                class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-bold bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20"
              >
                ✓ {t("builder.detect_dockerfile")}
              </span>
              {#if hasDockerignore}
                <span
                  class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-bold bg-blue-500/15 text-blue-600 dark:text-blue-400 border border-blue-500/20"
                >
                  ✓ {t("builder.detect_dockerignore")}
                </span>
              {:else}
                <span
                  class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-medium bg-amber-500/15 text-amber-600 dark:text-amber-400 border border-amber-500/20"
                  title={t("builder.warn_no_dockerignore")}
                >
                  ⚠️ {t("builder.warn_no_dockerignore")}
                </span>
              {/if}
            </div>

            <!-- Image Name / Tag Input -->
            <div
              class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-slate-50 dark:bg-slate-800/80 border border-slate-200 dark:border-slate-700/80 focus-within:border-violet-500 transition-colors"
            >
              <span class="text-xs font-semibold text-slate-500 dark:text-slate-400 select-none whitespace-nowrap">
                🏷️ {t("builder.tag_label")}
              </span>
              <input
                type="text"
                class="bg-transparent text-xs font-mono font-bold text-violet-600 dark:text-violet-400 focus:outline-none flex-1 placeholder:text-slate-400"
                placeholder={defaultTag}
                bind:value={customTag}
                disabled={status === "building"}
              />
              <span class="text-xs font-mono text-slate-400 select-none">:latest</span>
              {#if customTag.trim()}
                <ButtonYellow
                  size="xs"
                  onclick={() => (customTag = "")}
                  title="Restaurar nome padrão da pasta"
                >
                  Restaurar
                </ButtonYellow>
              {/if}
            </div>
          {:else if currentPath && !loading}
            <span
              class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-bold bg-slate-100 dark:bg-slate-800 text-slate-500 border border-slate-200 dark:border-slate-700"
            >
              {t("builder.no_dockerfile")}
            </span>
          {/if}

          <!-- Feedback messages -->
          {#if status === "building"}
            <span
              class="text-xs text-amber-500 font-semibold flex items-center gap-1.5"
            >
              <span
                class="animate-spin inline-block w-3.5 h-3.5 border-2 border-amber-500 border-t-transparent rounded-full"
              ></span>
              {t("builder.building")}
            </span>
          {:else if errorMsg}
            <span class="text-xs text-red-500 font-semibold"
            >
              {errorMsg}</span
            >
          {/if}

          <!-- Build + Save Path Buttons -->
          <div class="flex items-center gap-2">
            <ButtonGreen
              size="md"
              class="flex-1"
              disabled={!canBuild}
              loading={status === "building"}
              onclick={doBuild}
            >
              {t("builder.build_btn")}
            </ButtonGreen>

            {#if currentPath}
              <ButtonPurple
                size="md"
                disabled={savedPaths.includes(currentPath)}
                onclick={saveCurrentPath}
              >
                {t("builder.save_path")}
              </ButtonPurple>
            {/if}
          </div>

          {#if status === "success"}
            <div
              class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400"
            >
              <p class="text-xs font-bold">
                {t("builder.image_ready", { name: builtImage })}
              </p>
            </div>
            <ButtonBlue
              size="md"
              class="w-full"
              onclick={goToImages}
            >
              {t("builder.transfer_btn")}
            </ButtonBlue>
          {/if}

          {#if status === "error" && errorMsg}
            <div
              class="p-3 rounded-xl bg-red-500/10 border border-red-500/20 text-red-500"
            >
              <p class="text-xs font-bold">{errorMsg}</p>
            </div>
          {/if}
        </div>
      </div>

      <!-- Log Panel -->
      <div
        class="bg-[#070b15] border border-[#1e293b] rounded-2xl overflow-hidden flex flex-col h-120 shadow-xl"
      >
        <!-- Log Header -->
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
            {#if status === "building"}
              <span
                class="text-[10px] px-2 py-0.5 rounded-full bg-amber-500/10 border border-amber-500/20 text-amber-400 font-semibold animate-pulse"
              >
                {t("builder.building")}
              </span>
            {/if}
          </div>

          {#if logs.length > 0}
            <button
              type="button"
              class="flex items-center gap-1.5 px-3 py-1 text-xs font-semibold rounded-xl bg-slate-800/80 hover:bg-slate-700 text-slate-200 border border-slate-700/60 transition-all cursor-pointer shadow-sm active:scale-95"
              onclick={copyLogs}
            >
              {#if copied}
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="w-3.5 h-3.5 text-emerald-400"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2.5"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="m4.5 12.75 6 6 9-13.5"
                  /></svg
                >
                <span class="text-emerald-400">{t("common.copied")}</span>
              {:else}
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="w-3.5 h-3.5 text-slate-400"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M15.75 17.25v3.375c0 .621-.504 1.125-1.125 1.125h-9.75a1.125 1.125 0 0 1-1.125-1.125V7.875c0-.621.504-1.125 1.125-1.125H6.75a9.06 9.06 0 0 1 1.5.124m7.5 10.376h3.375c.621 0 1.125-.504 1.125-1.125V11.25c0-4.46-3.243-8.161-7.5-8.876a9.06 9.06 0 0 0-1.5-.124H9.375c-.621 0-1.125.504-1.125 1.125v3.5m7.5 10.375H9.375a1.125 1.125 0 0 1-1.125-1.125v-9.25c0-.621.504-1.125 1.125-1.125h5.25c.621 0 1.125.504 1.125 1.125v9.25c0 .621-.504 1.125-1.125 1.125Z"
                  /></svg
                >
                <span>{t("builder.copy_logs")}</span>
              {/if}
            </button>
          {/if}
        </div>

        <!-- Log Content Container -->
        <div
          bind:this={logContainer}
          class="p-4 overflow-y-auto flex-1 font-mono text-xs text-emerald-400 space-y-1 scrollbar-thin"
        >
          {#if logs.length === 0}
            <div class="text-slate-500 italic">
              {#if status === "building"}{t("builder.waiting_build")}{:else}{t(
                  "builder.select_folder_build",
                )}{/if}
            </div>
          {:else}
            {#each logs as log}
              <div class="leading-relaxed whitespace-pre-wrap">{log}</div>
            {/each}
          {/if}
        </div>
      </div>
  </div>
</div>
{/if}
