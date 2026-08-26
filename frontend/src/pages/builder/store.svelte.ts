import { getLocale } from "$shared/stores/locale.svelte";
import * as api from "./api";
import { folderNameFromPath, sanitizeTag, tagFromPath } from "./service";
import type { BuildResult, BuilderFolder, BuilderStatus, BuilderStore } from "./types";

let currentPath = $state("");
let parentPath = $state<string | null>(null);
let folders = $state<BuilderFolder[]>([]);
let hasDockerfile = $state(false);
let hasDockerignore = $state(false);
let loading = $state(false);
let status = $state<BuilderStatus>("idle");
let logs = $state<string[]>([]);
let builtImage = $state("");
let errorMsg = $state("");
let savedPaths = $state<string[]>([]);
let customTag = $state("");

export const builderStore: BuilderStore = {
  get currentPath() { return currentPath; },
  get parentPath() { return parentPath; },
  get folders() { return folders; },
  get hasDockerfile() { return hasDockerfile; },
  get hasDockerignore() { return hasDockerignore; },
  get loading() { return loading; },
  get status() { return status; },
  get logs() { return logs; },
  get builtImage() { return builtImage; },
  get errorMsg() { return errorMsg; },
  get savedPaths() { return savedPaths; },
  get customTag() { return customTag; },
  set customTag(value: string) { customTag = value; },
  get folderName() { return folderNameFromPath(currentPath); },
  get defaultTag() { return currentPath ? tagFromPath(currentPath) : ""; },
  get effectiveTag() { return customTag.trim() ? sanitizeTag(customTag) : this.defaultTag; },
  get canBuild() { return this.effectiveTag.length > 0 && hasDockerfile && status === "ready"; },

  async loadSavedPaths() {
    try { savedPaths = (await api.listSavedPaths() || []).map((item) => item.path); } catch { /* optional data */ }
  },

  async browse(path = "") {
    loading = true;
    errorMsg = "";
    status = "idle";
    customTag = "";
    try {
      const result = await api.browse(path);
      currentPath = result.currentPath;
      parentPath = result.parentPath || null;
      folders = result.folders || [];
      hasDockerfile = result.hasDockerfile;
      hasDockerignore = result.hasDockerignore;
      status = result.hasDockerfile ? "ready" : "idle";
    } catch {
      currentPath = "";
      parentPath = null;
      folders = [];
      hasDockerfile = false;
      hasDockerignore = false;
      errorMsg = "Erro ao navegar para a pasta.";
    }
    finally { loading = false; }
  },

  async saveCurrentPath() {
    if (!currentPath || savedPaths.includes(currentPath)) return;
    await api.savePath(currentPath, folderNameFromPath(currentPath));
    savedPaths = [...savedPaths, currentPath];
  },

  async removeSavedPath(path: string) {
    await api.removeSavedPath(path);
    savedPaths = savedPaths.filter((item) => item !== path);
  },

  async build() {
    if (!this.canBuild) return;
    status = "building";
    logs = [];
    errorMsg = "";
    builtImage = "";
    try { await api.build(currentPath, this.effectiveTag, getLocale()); }
    catch (error: any) { status = "error"; errorMsg = error?.message || "Build falhou"; }
  },

  appendLog(line: string) { logs.push(line); },

  completeBuild(result: BuildResult) {
    if (result.success) { builtImage = result.image || this.effectiveTag; status = "success"; }
    else { errorMsg = result.message || "Build falhou"; status = "error"; }
  },
};