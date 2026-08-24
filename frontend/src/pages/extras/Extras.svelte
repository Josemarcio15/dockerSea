<script lang="ts">
  import { onMount } from "svelte";
  import { t } from "$lib/stores/locale.svelte";
  import { useRefreshKey, triggerRefresh } from "$lib/stores/refresh.svelte";
  import {
    notifySuccess,
    notifyWarning,
    notifyError,
  } from "$lib/stores/notification.svelte";
  import DockerseaLoading from "$lib/components/DockerseaLoading.svelte";
  import StatusBanner from "$lib/components/StatusBanner.svelte";
  import VpsSelectWarning from "$lib/components/VpsSelectWarning.svelte";
  import CodeEditor from "$lib/components/CodeEditor.svelte";
  import PortsPanel from "$lib/components/PortsPanel.svelte";
  import NginxLogsModal from "$lib/components/NginxLogsModal.svelte";
  import {
    ButtonPurple,
    ButtonPink,
    ButtonGreen,
    ButtonYellow,
    ButtonRed,
    ButtonBlue,
    ButtonCyan,
    Button,
  } from "$lib/components/buttons";
  import * as ExtraService from "../../../bindings/go-walis/internal/extras/extraservice.js";

  let { data } = $props();

  let site = $state("");
  let busy = $state<string | null>(null);
  let activeTab = $state<"available" | "enabled">("available");
  let mainTab = $state<"nginx" | "ports">("nginx");
  let available = $state<string[]>([]);
  let enabled = $state<string[]>([]);
  let loadingSites = $state(true);
  let fetchError = $state("");
  let editorContent = $state("");
  let editorKey = $state(0);
  let showDeleteDialog = $state(false);
  let showLogs = $state(false);

  // Carrega lista de arquivos Nginx do servidor ativo
  async function loadSites(silent = false) {
    if (!data?.activeVps) {
      loadingSites = false;
      return;
    }
    if (!silent) {
      loadingSites = true;
    }
    fetchError = "";
    try {
      const res = await ExtraService.ListNginxSites(data.activeVps);
      available = res?.available || [];
      enabled = res?.enabled || [];
      if (!site && available.length > 0) {
        site = available[0];
        openSite(available[0]);
      }
    } catch (error: any) {
      const message =
        error?.message || String(error) || "Não foi possível listar os arquivos do Nginx.";
      fetchError = message;
      if (!silent) {
        available = [];
        enabled = [];
      }
    } finally {
      loadingSites = false;
    }
  }

  // Abre arquivo de configuração selecionado
  async function openSite(filename: string) {
    if (!data?.activeVps || !filename) return;
    site = filename;
    try {
      const content = await ExtraService.ReadNginxSite(
        data.activeVps,
        filename,
        activeTab,
      );
      editorContent = content || "";
      editorKey += 1;
    } catch (error: any) {
      const message =
        error?.message || String(error) || "Não foi possível abrir o arquivo.";
      notifyError(message);
    }
  }

  // Cria template para novo site
  function newSite() {
    site = "";
    editorContent = `server {
    listen 80;
    listen 443 ssl;
    server_name example.com;

    # Configure os certificados antes de ativar o SSL:
    # ssl_certificate /etc/letsencrypt/live/example.com/fullchain.pem;
    # ssl_certificate_key /etc/letsencrypt/live/example.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}`;
    editorKey += 1;
  }

  // Executa exclusão de site
  async function deleteSite() {
    if (!site.trim()) {
      notifyWarning(t("extras.select_file_warn"));
      return;
    }
    if (!data?.activeVps) {
      notifyError("Nenhuma VPS ativa selecionada.");
      return;
    }

    showDeleteDialog = false;
    busy = "delete";
    try {
      const result = await ExtraService.DeleteNginxSite(
        data.activeVps,
        site,
        activeTab,
      );
      if (result && result.success) {
        site = "";
        editorContent = "";
        editorKey += 1;
        notifySuccess(result.message || t("extras.delete_success"));
        await loadSites();
      } else {
        notifyError(result?.message || "Falha ao apagar arquivo.");
      }
    } catch (error: any) {
      notifyError(error?.message || String(error) || "Erro ao apagar arquivo.");
    } finally {
      busy = null;
    }
  }

  // Executa comandos Nginx (enable, test, restart, save)
  async function run(action: "enable" | "test" | "restart" | "save") {
    if (!data?.activeVps) {
      notifyError("Nenhuma VPS ativa selecionada.");
      return;
    }

    if ((action === "enable" || action === "save") && !site.trim()) {
      notifyWarning(t("extras.input_name_warn"));
      return;
    }

    busy = action;
    try {
      let result: any;
      switch (action) {
        case "save":
          result = await ExtraService.SaveNginxSite(
            data.activeVps,
            site,
            editorContent,
          );
          break;
        case "enable":
          result = await ExtraService.EnableNginxSite(data.activeVps, site);
          break;
        case "test":
          result = await ExtraService.TestNginxConfig(data.activeVps);
          break;
        case "restart":
          result = await ExtraService.RestartNginx(data.activeVps);
          break;
      }

      if (result && result.success) {
        notifySuccess(result.message || result.output || "Comando executado com sucesso.");
        if (action === "save" || action === "enable") {
          await loadSites(true);
        }
      } else {
        notifyError(
          result?.message || result?.output || "Falha ao executar ação no Nginx.",
        );
      }
    } catch (error: any) {
      notifyError(error?.message || String(error) || "Erro ao executar comando.");
    } finally {
      busy = null;
    }
  }

  $effect(() => {
    useRefreshKey();
    if (data?.activeVps) {
      loadSites();
    } else {
      loadingSites = false;
    }
  });
</script>

{#if !data.activeVps}
  <VpsSelectWarning />
{:else}
  <div class="space-y-6">
    <!-- Top Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div
        class="inline-flex items-center px-5 py-2.5 rounded-2xl bg-linear-to-br from-violet-100 to-fuchsia-100 dark:from-violet-950/40 dark:to-fuchsia-950/40 border border-violet-200/50 dark:border-violet-800/50 self-start shadow-sm"
      >
        <h1
          class="text-2xl font-bold bg-linear-to-r from-violet-600 to-fuchsia-500 bg-clip-text text-transparent m-0 flex items-center gap-2"
        >
          {t("extras.title")}
        </h1>
      </div>

      <div class="flex items-center gap-3">
        <!-- Main Navigation Pills (Nginx / Ports) -->
        <div
          class="inline-flex rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-[#0b0f19] p-1 shadow-xs"
        >
          <button
            type="button"
            class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {mainTab === 'nginx'
              ? 'bg-violet-600 text-white shadow-xs'
              : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}"
            onclick={() => (mainTab = "nginx")}
          >
            {t("extras.tab_nginx")}
          </button>
          <button
            type="button"
            class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {mainTab === 'ports'
              ? 'bg-violet-600 text-white shadow-xs'
              : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}"
            onclick={() => (mainTab = "ports")}
          >
            {t("extras.tab_ports")}
          </button>
        </div>

        <ButtonPink
          size="sm"
          title={t("common.refresh")}
          onclick={() => {
            if (mainTab === "nginx") {
              loadSites();
            } else {
              triggerRefresh();
            }
          }}
        >
          {t("common.refresh")}
        </ButtonPink>
      </div>
    </div>

    <!-- Status Alerts -->
    <StatusBanner />

    {#if mainTab === "ports"}
      <!-- Tab 1: Ports Manager -->
      <section
        class="w-full bg-white dark:bg-[#0b0f19] border border-slate-200/80 dark:border-slate-800/80 rounded-2xl p-5 shadow-sm"
      >
        <PortsPanel activeVps={data.activeVps} />
      </section>
    {:else}
      {#if loadingSites}
        <DockerseaLoading message={t("common.loading")} />
      {:else if fetchError}
        <div
          class="p-6 rounded-2xl border border-red-200 dark:border-red-900/60 bg-red-50 dark:bg-red-950/20 text-red-700 dark:text-red-400"
        >
          <h3 class="font-bold text-sm mb-1">Erro de Conexão Nginx</h3>
          <p class="text-xs whitespace-pre-wrap">{fetchError}</p>
        </div>
      {:else}
        <!-- Tab 2: Nginx Manager -->
        <section
          class="space-y-6 rounded-2xl border border-slate-200/80 dark:border-slate-800/80 bg-white dark:bg-[#0b0f19] p-6 shadow-sm"
        >
          <!-- Nginx Action Toolbar -->
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div class="inline-flex rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-100 dark:bg-slate-900 p-1">
              <button
                type="button"
                class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {activeTab ===
                'available'
                  ? 'bg-violet-600 text-white shadow-xs'
                  : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}"
                onclick={() => (activeTab = "available")}
              >
                {t("extras.sites_available")} ({available.length})
              </button>
              <button
                type="button"
                class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {activeTab ===
                'enabled'
                  ? 'bg-violet-600 text-white shadow-xs'
                  : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}"
                onclick={() => (activeTab = "enabled")}
              >
                {t("extras.sites_enabled")} ({enabled.length})
              </button>
            </div>

            <div class="flex flex-wrap items-center gap-2">
              <ButtonGreen
                size="sm"
                loading={busy === "test"}
                onclick={() => run("test")}
              >
                {busy === "test" ? t("extras.testing") : t("extras.test_nginx")}
              </ButtonGreen>
              <ButtonYellow
                size="sm"
                loading={busy === "restart"}
                onclick={() => run("restart")}
              >
                {busy === "restart" ? t("extras.restarting") : t("extras.restart_nginx")}
              </ButtonYellow>
              <ButtonCyan
                size="sm"
                onclick={() => {
                  showLogs = true;
                }}
              >
                {t("extras.view_logs")}
              </ButtonCyan>
              <ButtonRed
                size="sm"
                disabled={!site.trim() || !!busy}
                onclick={() => {
                  if (site.trim()) {
                    showDeleteDialog = true;
                  } else {
                    notifyWarning(t("extras.select_file_warn"));
                  }
                }}
              >
                {t("extras.delete_file")}
              </ButtonRed>
              <ButtonBlue
                size="sm"
                onclick={newSite}
              >
                {t("extras.new_site")}
              </ButtonBlue>
            </div>
          </div>

          <!-- Site Pills List -->
          <div class="min-h-12 flex items-center">
            {#if (activeTab === "available" ? available : enabled).length === 0}
              <p class="text-xs text-slate-400 dark:text-slate-500 italic">
                {t("extras.empty_sites")}
              </p>
            {:else}
              <div class="flex flex-wrap gap-2">
                {#each activeTab === "available" ? available : enabled as filename}
                  <button
                    type="button"
                    class="px-3 py-1.5 rounded-xl text-xs font-mono border transition-all cursor-pointer {site ===
                    filename
                      ? 'border-violet-500 bg-violet-50 dark:bg-violet-950/40 text-violet-700 dark:text-violet-300 font-bold shadow-xs'
                      : 'border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-400 hover:border-slate-300 dark:hover:border-slate-700'}"
                    onclick={() => openSite(filename)}
                  >
                    {filename}
                  </button>
                {/each}
              </div>
            {/if}
          </div>

          <!-- Input File Name -->
          <div class="space-y-1.5">
            <label
              for="site"
              class="text-[11px] font-bold uppercase tracking-wider text-slate-500 dark:text-slate-400"
            >
              {t("extras.site_file_label")}
            </label>
            <input
              id="site"
              bind:value={site}
              placeholder={t("extras.site_file_placeholder")}
              class="w-full px-3.5 py-2 text-xs rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-[#0c101b] text-slate-900 dark:text-white font-mono focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
            />
            <p class="text-[11px] text-slate-400 font-mono">
              {activeTab === "available"
                ? "/etc/nginx/sites-available"
                : "/etc/nginx/sites-enabled"}/{site || "..."}
            </p>
          </div>

          <!-- Code Editor Block -->
          <div class="space-y-1.5">
            <label
              for="nginx-editor"
              class="text-[11px] font-bold uppercase tracking-wider text-slate-500 dark:text-slate-400"
            >
              {t("extras.site_content_label")}
            </label>
            {#key editorKey}
              <CodeEditor
                value={editorContent}
                mode="nginx"
                onchange={(value) => (editorContent = value)}
              />
            {/key}
          </div>

          <!-- Bottom Action Buttons -->
          <div class="flex flex-wrap gap-2.5 pt-2">
            <ButtonPurple
              size="md"
              loading={busy === "enable"}
              onclick={() => run("enable")}
            >
              {busy === "enable" ? "Executando..." : t("extras.btn_enable")}
            </ButtonPurple>
            <ButtonGreen
              size="md"
              loading={busy === "save"}
              onclick={() => run("save")}
            >
              {busy === "save" ? "Salvando..." : t("extras.btn_save")}
            </ButtonGreen>
          </div>
        </section>
      {/if}
    {/if}
  </div>
{/if}

<!-- Nginx Logs Modal -->
<NginxLogsModal bind:show={showLogs} activeVps={data.activeVps} />

<!-- Confirm Delete Dialog Modal -->
{#if showDeleteDialog}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4 animate-fadeIn">
    <div
      class="w-full max-w-md rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-[#0c1220] p-6 shadow-2xl space-y-4"
    >
      <h2 class="text-base font-bold text-slate-800 dark:text-slate-100">
        {t("extras.delete_confirm_title")}
      </h2>
      <p class="text-xs text-slate-600 dark:text-slate-300 leading-relaxed">
        {t("extras.delete_confirm_msg", { site })}
      </p>
      <div class="flex justify-end gap-2.5 pt-2">
        <Button
          size="sm"
          onclick={() => {
            showDeleteDialog = false;
          }}
        >
          {t("common.cancel")}
        </Button>
        <ButtonRed
          size="sm"
          onclick={deleteSite}
        >
          {t("common.delete")}
        </ButtonRed>
      </div>
    </div>
  </div>
{/if}
