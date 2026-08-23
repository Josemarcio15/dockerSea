<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import CodeEditor from "$lib/components/CodeEditor.svelte";

  interface PortMapping {
    external: string;
    internal: string;
  }

  interface EnvVar {
    name: string;
    value: string;
  }

  interface VolumeMapping {
    host: string;
    container: string;
  }

  let {
    show = $bindable(false),
    image = "",
    savedConfigs = [],
    serverId = "",
    onsubmit = (config: any) => {},
    onsaveprofile = (profile: any) => {},
    ondeleteprofile = (profileId: string) => {},
  }: {
    show: boolean;
    image: string;
    savedConfigs: any[];
    serverId: string;
    onsubmit: (config: any) => void;
    onsaveprofile: (profile: any) => void;
    ondeleteprofile: (profileId: string) => void;
  } = $props();

  // Lazy fetch volumes/networks apenas quando o modal abre
  // (usa as rotas existentes que já descobrem a VPS ativa via cookie)
  let existingVolumes = $state<string[]>([]);
  let existingNetworks = $state<string[]>([]);

  $effect(() => {
    if (!show) return;
    existingVolumes = [];
    existingNetworks = [];
    fetch("/api/volumes")
      .then((r) => r.json())
      .then((d) => {
        existingVolumes = (d.volumes || []).map((v: any) => v.name);
      })
      .catch(() => {});
    fetch("/api/networks")
      .then((r) => r.json())
      .then((d) => {
        existingNetworks = (d.networks || []).map((n: any) => n.name);
      })
      .catch(() => {});
  });

  // Form fields state
  let containerName = $state("");
  let projectName = $state("");
  let ports = $state<PortMapping[]>([]);
  let envs = $state<EnvVar[]>([]);
  let volumes = $state<VolumeMapping[]>([]);
  let network = $state("");
  let restartPolicy = $state("");
  let commands = $state<string[]>([]);
  let description = $state("");
  let profileName = $state("");

  let loadedProfileId = $state<string | null>(null);
  let loadedProfileName = $state("");
  let isModified = $state(false);
  let showConfirmNoVolume = $state(false);
  let importFileInput = $state<HTMLInputElement>();
  let jsonEditor = $state("");
  let jsonMessage = $state("");
  let editorTab = $state<"fields" | "json">("fields");

  const emptyConfigTemplate = {
    name: "",
    image: "",
    containerName: "",
    projectName: "",
    ports: [{ port: "", "port-intern": "" }],
    envs: [{ name: "", value: "" }],
    volumes: [{ host: "", container: "" }],
    network: "",
    restartPolicy: "",
    commands: [""],
    description: "",
  };

  function downloadJson() {
    const config = getEditorConfig();
    const blob = new Blob([JSON.stringify(config, null, 2)], {
      type: "application/json",
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = `${profileName.trim() || "docker-config"}.json`;
    link.click();
    URL.revokeObjectURL(url);
  }

  function getEditorConfig() {
    return {
      name: profileName,
      image,
      containerName,
      projectName,
      ports: ports.length
        ? ports.map((p) => ({ port: p.external, "port-intern": p.internal }))
        : [{ port: "", "port-intern": "" }],
      envs: envs.length ? envs : [{ name: "", value: "" }],
      volumes: volumes.length ? volumes : [{ host: "", container: "" }],
      network,
      restartPolicy,
      commands: commands.length ? commands : [""],
      description,
    };
  }

  function syncJsonEditor() {
    jsonEditor = JSON.stringify(getEditorConfig(), null, 2);
    jsonMessage = "";
  }

  async function copyJson() {
    await navigator.clipboard.writeText(jsonEditor);
    jsonMessage = "JSON copiado.";
  }

  function applyJson() {
    try {
      const parsed = JSON.parse(jsonEditor);
      const cfg = { ...emptyConfigTemplate, ...(parsed.config || parsed) };
      containerName = cfg.containerName || "";
      projectName = cfg.projectName || "";
      network = cfg.network || "";
      restartPolicy = cfg.restartPolicy || "";
      description = cfg.description || "";
      profileName = cfg.name || profileName;
      ports = (Array.isArray(cfg.ports) ? cfg.ports : []).map((p: any) => ({ external: p.port ?? p.external ?? "", internal: p["port-intern"] ?? p.internal ?? "" }));
      envs = (Array.isArray(cfg.envs) ? cfg.envs : []).map((e: any) => ({ name: e.name ?? "", value: e.value ?? "" }));
      volumes = (Array.isArray(cfg.volumes) ? cfg.volumes : []).map((v: any) => ({ host: v.host ?? "", container: v.container ?? "" }));
      commands = Array.isArray(cfg.commands) ? cfg.commands.map((c: any) => String(c ?? "")) : [""];
      loadedProfileId = null;
      loadedProfileName = "";
      isModified = true;
      jsonMessage = "JSON aplicado ao formulário.";
    } catch {
      jsonMessage = "JSON inválido. Verifique a estrutura antes de aplicar.";
    }
  }

  function loadExample() {
    jsonEditor = JSON.stringify({
      name: "meu-perfil",
      image: "nginx:latest",
      containerName: "meu-container",
      projectName: "meu-projeto",
      ports: [{ port: "8080", "port-intern": "80" }],
      envs: [{ name: "APP_ENV", value: "production" }],
      volumes: [{ host: "./data", container: "/app/data" }],
      network: "app-network",
      restartPolicy: "unless-stopped",
      commands: ["nginx -g 'daemon off;'"],
      description: "Exemplo de configuração Docker",
    }, null, 2);
    applyJson();
    editorTab = "json";
  }

  function openImport() {
    importFileInput?.click();
  }

  async function importJson(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;
    try {
      const parsed = JSON.parse(await file.text());
      const cfg = { ...emptyConfigTemplate, ...(parsed.config || parsed) };
      containerName = cfg.containerName || "";
      projectName = cfg.projectName || "";
      network = cfg.network || "";
      restartPolicy = cfg.restartPolicy || "";
      description = cfg.description || "";
      profileName = cfg.name || "";
      loadedProfileId = null;
      loadedProfileName = "";
      ports = (Array.isArray(cfg.ports) ? cfg.ports : emptyConfigTemplate.ports).map((p: any) => ({
        external: p.port ?? p.external ?? "",
        internal: p["port-intern"] ?? p.internal ?? "",
      }));
      envs = (Array.isArray(cfg.envs) ? cfg.envs : emptyConfigTemplate.envs).map((e: any) => ({
        name: e.name ?? "",
        value: e.value ?? "",
      }));
      volumes = (Array.isArray(cfg.volumes) ? cfg.volumes : emptyConfigTemplate.volumes).map((v: any) => ({
        host: v.host ?? "",
        container: v.container ?? "",
      }));
      commands = Array.isArray(cfg.commands)
        ? cfg.commands.map((command: any) => String(command ?? ""))
        : cfg.command
          ? String(cfg.command).split(" && ")
          : [""];
      isModified = true;
    } catch {
      // O formulário permanece intacto quando o arquivo não é JSON válido.
    } finally {
      input.value = "";
    }
  }

  // Helpers
  function getDefaultContainerPath(imgName: string): string {
    const img = imgName.split(":")[0].split("/").pop()?.toLowerCase() || "";
    if (img.includes("postgres")) return "/var/lib/postgresql/data";
    if (img.includes("mysql") || img.includes("mariadb"))
      return "/var/lib/mysql";
    if (img.includes("mongo")) return "/data/db";
    if (img.includes("redis")) return "/data";
    if (img.includes("influx")) return "/var/lib/influxdb";
    if (img.includes("rabbitmq")) return "/var/lib/rabbitmq";
    if (img.includes("elasticsearch")) return "/usr/share/elasticsearch/data";
    if (img.includes("nginx")) return "/usr/share/nginx/html";
    if (img.includes("httpd")) return "/usr/local/apache2/htdocs";
    return "/data";
  }

  // Computed properties
  let hasEmptyVolume = $derived(
    volumes.some((v) => !v.host.trim() || !v.container.trim()),
  );

  let nameAlreadyExists = $derived(
    profileName.trim() !== "" &&
      savedConfigs.some(
        (c) => c.name === profileName && c.id !== loadedProfileId,
      ),
  );

  let isNameChanged = $derived(
    loadedProfileId !== null && profileName !== loadedProfileName,
  );

  let saveDisabled = $derived(
    profileName.trim() === "" || nameAlreadyExists || hasEmptyVolume,
  );

  // Actions
  function loadProfile(cfg: any) {
    containerName = cfg.containerName || "";
    projectName = cfg.projectName || "";
    network = cfg.network || "";
    restartPolicy = cfg.restartPolicy || "no";
    description = cfg.description || "";
    profileName = cfg.name || "";
    loadedProfileId = cfg.id;
    loadedProfileName = cfg.name || "";

    if (cfg.command) {
      commands = cfg.command.split(" && ");
    } else {
      commands = [""];
    }

    // Parse ports "host:container host2:container2"
    ports = [];
    if (cfg.ports) {
      for (const entry of cfg.ports.split(/\s+/)) {
        const parts = entry.split(":");
        if (parts.length === 2) {
          ports.push({ external: parts[0], internal: parts[1] });
        }
      }
    }

    // Parse ENVs
    envs = [];
    if (cfg.env) {
      const rawEnv = cfg.env.trim();
      const hasNewline = rawEnv.includes("\n");
      const tokens = rawEnv.split(/\s+/);
      const isLegacySpaceSeparated =
        !hasNewline &&
        tokens.length > 1 &&
        tokens.every((token: string) => token.includes("="));

      if (isLegacySpaceSeparated) {
        for (const entry of tokens) {
          const idx = entry.indexOf("=");
          if (idx > 0) {
            envs.push({
              name: entry.substring(0, idx),
              value: entry.substring(idx + 1),
            });
          }
        }
      } else {
        const lines = rawEnv.split(/\r?\n/);
        for (const line of lines) {
          const trimmed = line.trim();
          if (!trimmed) continue;
          const idx = trimmed.indexOf("=");
          if (idx > 0) {
            envs.push({
              name: trimmed.substring(0, idx).trim(),
              value: trimmed.substring(idx + 1),
            });
          }
        }
      }
    }

    // Parse volumes "host:container host2:container2"
    volumes = [];
    if (cfg.volumes) {
      for (const entry of cfg.volumes.split(/\s+/)) {
        const parts = entry.split(":");
        if (parts.length === 2) {
          volumes.push({ host: parts[0], container: parts[1] });
        } else if (entry) {
          volumes.push({ host: entry, container: "" });
        }
      }
    }

    isModified = false;
    syncJsonEditor();
  }

  function createNewProfile() {
    loadedProfileId = null;
    loadedProfileName = "";
    profileName = "";
    containerName = "";
    projectName = "";
    ports = [];
    envs = [];
    volumes = [];
    network = "";
    restartPolicy = "";
    commands = [""];
    description = "";
    jsonEditor = "";
    jsonMessage = "Novo perfil pronto para edição.";
    isModified = true;
    editorTab = "fields";
  }

  function addPort() {
    ports = [...ports, { external: "", internal: "" }];
    isModified = true;
  }

  function removePort(idx: number) {
    ports = ports.filter((_, i) => i !== idx);
    isModified = true;
  }

  function addEnv() {
    envs = [...envs, { name: "", value: "" }];
    isModified = true;
  }

  function removeEnv(idx: number) {
    envs = envs.filter((_, i) => i !== idx);
    isModified = true;
  }

  function addVolume() {
    volumes = [...volumes, { host: "", container: "" }];
    isModified = true;
  }

  function removeVolume(idx: number) {
    volumes = volumes.filter((_, i) => i !== idx);
    isModified = true;
  }

  function addCommand() {
    commands = [...commands, ""];
    isModified = true;
  }

  function removeCommand(idx: number) {
    commands = commands.filter((_, i) => i !== idx);
    isModified = true;
  }

  function triggerSave() {
    const portsStr = ports.map((p) => `${p.external}:${p.internal}`).join(" ");
    const envStr = envs
      .filter((e) => e.name.trim() !== "")
      .map((e) => `${e.name}=${e.value}`)
      .join("\n");
    const volStr = volumes
      .map((v) => `${v.host}:${v.container || getDefaultContainerPath(image)}`)
      .join(" ");
    const commandStr = commands.filter((c) => c.trim() !== "").join(" && ");

    onsaveprofile({
      id: isNameChanged ? null : loadedProfileId,
      name: profileName,
      containerName,
      projectName,
      image,
      ports: portsStr || null,
      env: envStr || null,
      volumes: volStr || null,
      network: network || null,
      restartPolicy,
      command: commandStr || null,
      description: description || null,
    });

    isModified = false;
  }

  function validateAndSubmit() {
    if (volumes.length === 0) {
      showConfirmNoVolume = true;
    } else {
      submitForm();
    }
  }

  function submitForm() {
    showConfirmNoVolume = false;
    show = false;

    const volList = volumes.map(
      (v) => `${v.host}:${v.container || getDefaultContainerPath(image)}`,
    );
    const commandStr = commands.filter((c) => c.trim() !== "").join(" && ");

    onsubmit({
      containerName,
      image,
      ports,
      envs,
      volumes: volList,
      network,
      restartPolicy,
      projectName,
      command: commandStr || undefined,
    });
  }

  // Set initial values and reset on close
  $effect(() => {
    if (show) {
      if (containerName === "") {
        const shortName = image.split(":")[0].split("/").pop() || "app";
        containerName = `${shortName}-container`;
      }
      if (commands.length === 0) {
        commands = [""];
      }
    }
    if (!show) {
      loadedProfileId = null;
      loadedProfileName = "";
      isModified = false;
      containerName = "";
      projectName = "";
      ports = [];
      envs = [];
      volumes = [];
      network = "";
      restartPolicy = "";
      commands = [];
      description = "";
      profileName = "";
      jsonEditor = "";
      jsonMessage = "";
      editorTab = "fields";
    }
  });
</script>

{#if show}
  <div
    class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
  >
    <div
      class="bg-white dark:bg-[#111827] border border-slate-200 dark:border-slate-800 rounded-2xl shadow-2xl w-[90vw] h-[90vh] max-w-none max-h-none flex flex-col text-slate-800 dark:text-slate-100 overflow-hidden"
    >
      <!-- Header -->
      <div
        class="flex justify-between items-center px-6 py-4 border-b border-slate-200 dark:border-slate-850 bg-slate-50 dark:bg-slate-900/50"
      >
        <h2 class="text-base font-bold text-slate-850 dark:text-slate-100">
          {t("images.config_title")}
        </h2>
        <button
          type="button"
          class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 text-lg bg-transparent border-none cursor-pointer transition-colors"
          onclick={() => (show = false)}
        >
          ✕
        </button>
      </div>

      <!-- Body -->
      <div
        class="p-6 overflow-y-auto flex-1 grid grid-cols-1 lg:grid-cols-[280px_minmax(0,1fr)] gap-6 min-h-0 pr-3"
      >
        <!-- Saved profiles -->
        {#if savedConfigs.length > 0}
          <div class="flex flex-col gap-2">
            <div class="flex items-center justify-between gap-2">
              <span class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider">{t("images.config_saved_profiles")}</span>
              <button type="button" class="rounded-lg bg-violet-600 px-2.5 py-1.5 text-[10px] font-bold text-white hover:bg-violet-700" onclick={createNewProfile}>+ Novo perfil</button>
            </div>
            <div class="flex flex-col gap-2 max-h-[calc(92vh-180px)] overflow-y-auto pr-1">
              {#each savedConfigs as cfg (cfg.id)}
                <div
                  class="flex items-center gap-2 p-3 rounded-xl bg-slate-50 dark:bg-slate-900 border border-slate-150 dark:border-slate-800 shadow-sm text-xs"
                >
                  <div
                    class="flex-1 text-slate-600 dark:text-slate-350 flex flex-col min-w-0 gap-0.5"
                  >
                    <div>
                      <span class="font-bold text-slate-700 dark:text-slate-300"
                        >{t("images.config_profile_name")}</span
                      >
                      <span
                        class="font-semibold text-violet-600 dark:text-violet-400"
                        >{cfg.name}</span
                      >
                    </div>
                    {#if cfg.description}
                      <div class="truncate text-slate-400 dark:text-slate-500">
                        <span
                          class="font-bold text-slate-700 dark:text-slate-300"
                          >{t("images.config_profile_desc")}</span
                        >
                        <span>{cfg.description}</span>
                      </div>
                    {/if}
                  </div>
                  <button
                    type="button"
                    class="px-2.5 py-1 text-[10px] font-bold rounded-lg text-white bg-violet-600 hover:bg-violet-700 cursor-pointer transition-colors shadow-sm shrink-0"
                    onclick={() => loadProfile(cfg)}
                  >
                    {t("images.config_load")}
                  </button>
                  <button
                    type="button"
                    class="px-2 py-1 text-xs rounded border-none text-white bg-red-500 hover:bg-red-600 cursor-pointer shrink-0 transition-colors"
                    onclick={() => ondeleteprofile(cfg.id)}
                  >
                    ✕
                  </button>
                </div>
              {/each}
            </div>
          </div>
        {/if}

        {#if savedConfigs.length === 0}
          <div class="flex flex-col gap-3 rounded-xl border border-slate-200 dark:border-slate-800 p-3">
            <span class="text-[10px] font-bold uppercase tracking-wider text-slate-400">Perfis salvos</span>
            <button type="button" class="rounded-lg bg-violet-600 px-3 py-2 text-xs font-bold text-white hover:bg-violet-700" onclick={createNewProfile}>+ Novo perfil</button>
          </div>
        {/if}

        <div class="lg:col-start-2 flex flex-col min-w-0 gap-4">
        <div class="flex rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/40 p-1">
          <button type="button" class="flex-1 rounded-lg px-3 py-2 text-xs font-bold {editorTab === 'fields' ? 'bg-white dark:bg-slate-800 text-violet-600 shadow-sm' : 'text-slate-500'}" onclick={() => (editorTab = "fields")}>Campos</button>
          <button type="button" class="flex-1 rounded-lg px-3 py-2 text-xs font-bold {editorTab === 'json' ? 'bg-white dark:bg-slate-800 text-violet-600 shadow-sm' : 'text-slate-500'}" onclick={() => { editorTab = "json"; if (!jsonEditor) syncJsonEditor(); }}>JSON</button>
        </div>

        {#if editorTab === "json"}
        <div class="space-y-2 rounded-xl border border-violet-200 dark:border-violet-900/50 bg-violet-50/40 dark:bg-violet-950/10 p-3">
          <div class="flex items-center justify-between gap-2">
            <div>
              <span class="text-[10px] font-bold uppercase tracking-wider text-violet-700 dark:text-violet-300">JSON do perfil</span>
              <p class="text-[10px] text-slate-500">Selecione um perfil acima, depois edite ou cole as variáveis aqui.</p>
            </div>
            <div class="flex gap-2 shrink-0">
              <button type="button" class="px-2.5 py-1.5 text-[10px] font-bold rounded-lg border border-violet-300 text-violet-700 dark:text-violet-300" onclick={loadExample}>Exemplo</button>
              <button type="button" class="px-2.5 py-1.5 text-[10px] font-bold rounded-lg bg-violet-600 text-white hover:bg-violet-700" onclick={copyJson}>Copiar</button>
              <button type="button" class="px-2.5 py-1.5 text-[10px] font-bold rounded-lg border border-violet-300 text-violet-700 dark:text-violet-300" onclick={applyJson}>Aplicar</button>
            </div>
          </div>
          <CodeEditor value={jsonEditor} mode="json" onchange={(value) => { jsonEditor = value; isModified = true; }} />
          {#if jsonMessage}<p class="text-[10px] font-semibold {jsonMessage.includes('inválido') ? 'text-red-500' : 'text-emerald-600'}">{jsonMessage}</p>{/if}
        </div>
        {/if}

        {#if editorTab === "fields"}
        <!-- Image name -->
        <div class="flex flex-col gap-1">
          <label
            for="config-image"
            class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
            >Imagem</label
          >
          <input
            id="config-image"
            type="text"
            class="w-full px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-100 dark:bg-[#0c101b] text-slate-400 dark:text-slate-500 cursor-not-allowed font-mono"
            value={image}
            readonly
          />
        </div>

        <!-- Container Name -->
        <div class="flex flex-col gap-1">
          <label
            for="config-container-name"
            class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
            >{t("images.config_container_name")}</label
          >
          <input
            id="config-container-name"
            type="text"
            class="w-full px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
            bind:value={containerName}
            oninput={() => (isModified = true)}
            placeholder={t("images.config_placeholder_container")}
          />
        </div>

        <!-- Project Name -->
        <div class="flex flex-col gap-1">
          <label
            for="config-project-name"
            class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
            >{t("images.config_project_name")}</label
          >
          <input
            id="config-project-name"
            type="text"
            class="w-full px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
            bind:value={projectName}
            oninput={() => (isModified = true)}
            placeholder={t("images.config_placeholder_project")}
          />
        </div>

        <!-- Ports mapping -->
        <div class="flex flex-col gap-2">
          <div class="flex items-center justify-between">
            <span
              class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
              >{t("images.config_ports")}</span
            >
            <button
              type="button"
              class="text-xs text-violet-600 hover:text-violet-700 dark:text-violet-400 dark:hover:text-violet-300 font-bold cursor-pointer bg-transparent border-none transition-colors"
              onclick={addPort}
            >
              {t("images.config_add_port")}
            </button>
          </div>
          {#each ports as port, i}
            <div class="flex gap-2 items-center">
              <input
                type="text"
                class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
                placeholder={t("images.config_placeholder_port_ext")}
                bind:value={port.external}
                oninput={() => (isModified = true)}
              />
              <span class="text-slate-400">:</span>
              <input
                type="text"
                class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
                placeholder={t("images.config_placeholder_port_int")}
                bind:value={port.internal}
                oninput={() => (isModified = true)}
              />
              <button
                type="button"
                class="px-2.5 py-1.5 text-xs rounded-xl border-none text-white bg-red-500 hover:bg-red-600 cursor-pointer transition-colors"
                onclick={() => removePort(i)}
              >
                ✕
              </button>
            </div>
          {/each}
        </div>

        <!-- Environment variables -->
        <div class="flex flex-col gap-2">
          <div class="flex items-center justify-between">
            <span
              class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
              >{t("images.config_envs")}</span
            >
            <button
              type="button"
              class="text-xs text-violet-600 hover:text-violet-700 dark:text-violet-400 dark:hover:text-violet-300 font-bold cursor-pointer bg-transparent border-none transition-colors"
              onclick={addEnv}
            >
              {t("images.config_add_env")}
            </button>
          </div>
          {#each envs as env, i}
            <div class="flex gap-2 items-center">
              <input
                type="text"
                class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
                placeholder={t("images.config_placeholder_env_name")}
                bind:value={env.name}
                oninput={() => (isModified = true)}
              />
              <span class="text-slate-400">=</span>
              <input
                type="text"
                class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
                placeholder={t("images.config_placeholder_env_value")}
                bind:value={env.value}
                oninput={() => (isModified = true)}
              />
              <button
                type="button"
                class="px-2.5 py-1.5 text-xs rounded-xl border-none text-white bg-red-500 hover:bg-red-600 cursor-pointer transition-colors"
                onclick={() => removeEnv(i)}
              >
                ✕
              </button>
            </div>
          {/each}
        </div>

        <!-- Volumes mapping -->
        <div class="flex flex-col gap-2">
          <div class="flex items-center justify-between">
            <span
              class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
              >{t("images.config_volumes")}</span
            >
            <button
              type="button"
              class="text-xs text-violet-600 hover:text-violet-700 dark:text-violet-400 dark:hover:text-violet-300 font-bold cursor-pointer bg-transparent border-none transition-colors"
              onclick={addVolume}
            >
              {t("images.config_add_volume")}
            </button>
          </div>
          {#each volumes as vol, i}
            <div class="flex gap-2 items-center">
              <input
                type="text"
                class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
                placeholder={t("images.config_placeholder_vol_host")}
                bind:value={vol.host}
                oninput={() => (isModified = true)}
                list="volumes-datalist"
              />
              <span class="text-slate-400">:</span>
              <input
                type="text"
                class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
                placeholder={vol.host
                  ? getDefaultContainerPath(image)
                  : t("images.config_placeholder_vol_container")}
                bind:value={vol.container}
                oninput={() => (isModified = true)}
              />
              <button
                type="button"
                class="px-2.5 py-1.5 text-xs rounded-xl border-none text-white bg-red-500 hover:bg-red-600 cursor-pointer transition-colors"
                onclick={() => removeVolume(i)}
              >
                ✕
              </button>
            </div>
          {/each}
          <datalist id="volumes-datalist">
            {#each existingVolumes as vol}
              <option value={vol}></option>
            {/each}
          </datalist>
          {#if hasEmptyVolume}
            <div class="text-xs text-red-500 font-semibold mt-1">
              ⚠️ {t("images.config_empty_volume_warning")}
            </div>
          {/if}
        </div>

        <!-- Network selection -->
        <div class="flex flex-col gap-1">
          <label
            for="config-network"
            class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
            >{t("images.config_network")}</label
          >
          <input
            id="config-network"
            type="text"
            class="w-full px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
            placeholder={t("images.config_placeholder_network")}
            bind:value={network}
            oninput={() => (isModified = true)}
            list="networks-datalist"
          />
          <datalist id="networks-datalist">
            {#each existingNetworks as net}
              <option value={net}></option>
            {/each}
          </datalist>
        </div>

        <!-- Restart policy -->
        <div class="flex flex-col gap-1">
          <label
            for="config-restart-policy"
            class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
            >{t("containers.card_restart")}</label
          >
          <select
            id="config-restart-policy"
            class="w-full px-3.5 py-2.5 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors appearance-none"
            bind:value={restartPolicy}
            onchange={() => (isModified = true)}
          >
            <option value="no">{t("containers.card_restart_no")}</option>
            <option value="always">{t("containers.card_restart_always")}</option
            >
            <option value="unless-stopped"
              >{t("containers.card_restart_unless_stopped")}</option
            >
            <option value="on-failure"
              >{t("containers.card_restart_on_failure")}</option
            >
          </select>
        </div>

        <!-- Command -->
        <div class="flex flex-col gap-2">
          <div class="flex items-center justify-between">
            <span
              class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
              >{t("images.config_command")}</span
            >
            <button
              type="button"
              class="text-xs text-violet-600 hover:text-violet-700 dark:text-violet-400 dark:hover:text-violet-300 font-bold cursor-pointer bg-transparent border-none transition-colors"
              onclick={addCommand}
            >
              {t("images.config_add_command")}
            </button>
          </div>

          {#each commands as cmd, i}
            <div class="flex gap-2 items-center">
              <input
                type="text"
                class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors font-mono"
                placeholder={t("images.config_placeholder_command")}
                bind:value={commands[i]}
                oninput={() => (isModified = true)}
              />
              <button
                type="button"
                class="px-2.5 py-1.5 text-xs rounded-xl border-none text-white bg-red-500 hover:bg-red-600 cursor-pointer transition-colors"
                onclick={() => removeCommand(i)}
              >
                ✕
              </button>
            </div>
          {/each}
          <span
            class="text-[10px] text-slate-400 dark:text-slate-500 mt-0.5 leading-relaxed"
          >
            {t("images.config_command_hint")}
          </span>
        </div>

        <!-- Description -->
        <div class="flex flex-col gap-1">
          <label
            for="config-description"
            class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
            >{t("images.config_desc")}</label
          >
          <textarea
            id="config-description"
            class="w-full px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors resize-y min-h-15"
            placeholder={t("images.config_placeholder_desc")}
            bind:value={description}
            oninput={() => (isModified = true)}
          ></textarea>
        </div>

        <!-- Save Profile Card -->
        <div
          class="flex flex-col gap-1.5 p-3.5 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/20"
        >
          <label
            for="config-profile-name"
            class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
            >{t("images.config_save_profile")}</label
          >
          <div class="flex gap-2 items-center">
            <input
              id="config-profile-name"
              type="text"
              class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
              placeholder={t("images.config_placeholder_profile")}
              bind:value={profileName}
            />
            {#if profileName.trim()}
              <button
                type="button"
                class="px-4 py-2.5 rounded-xl border-none cursor-pointer text-xs font-bold text-white transition-colors shadow-sm {saveDisabled
                  ? 'bg-slate-400 cursor-not-allowed'
                  : 'bg-green-600 hover:bg-green-700 shadow-green-500/10'}"
                disabled={saveDisabled}
                onclick={triggerSave}
              >
                {loadedProfileId && !isNameChanged
                  ? t("images.config_btn_update")
                  : t("images.config_btn_create_profile")}
              </button>
            {/if}
          </div>
          {#if nameAlreadyExists}
            <div class="text-[11px] text-red-500 font-semibold mt-1">
              Porta ou Perfil já cadastrado com este nome.
            </div>
          {/if}
        </div>
        {/if}
        </div>
      </div>

      <!-- Footer -->
      <div
        class="flex justify-end gap-3 px-6 py-4 border-t border-slate-200 dark:border-slate-850 bg-slate-50 dark:bg-slate-900/50"
      >
        <button
          type="button"
          class="px-4 py-2 text-xs font-bold rounded-xl border border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 cursor-pointer bg-white dark:bg-slate-900 transition-colors"
          onclick={() => (show = false)}
        >
          {t("common.cancel")}
        </button>
        <button
          type="button"
          class="px-4 py-2 text-xs font-bold rounded-xl border-none text-white transition-colors shadow-md shadow-violet-500/20 {hasEmptyVolume
            ? 'bg-slate-400 cursor-not-allowed'
            : 'bg-violet-600 hover:bg-violet-700 cursor-pointer'}"
          disabled={hasEmptyVolume}
          onclick={validateAndSubmit}
        >
          {t("images.config_create_container")}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Confirmation No Volume dialog -->
{#if showConfirmNoVolume}
  <div
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-60"
  >
    <div
      class="bg-white dark:bg-[#1b2330] border border-slate-200 dark:border-[#3f3f46] rounded-2xl p-6 shadow-2xl max-w-sm w-full mx-4 text-slate-800 dark:text-slate-100"
    >
      <h3 class="text-sm font-bold text-slate-800 dark:text-slate-100 mb-2">
        {t("images.config_confirm_persistence")}
      </h3>
      <p class="text-xs text-red-600 dark:text-red-400 font-semibold mb-6">
        {t("images.config_persistence_warning")}
      </p>
      <div class="flex gap-3 justify-end">
        <button
          type="button"
          class="px-4 py-2 text-xs font-bold rounded-xl border border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-350 hover:bg-slate-100 dark:hover:bg-[#3f3f46] bg-slate-50 dark:bg-[#27272a] transition-colors cursor-pointer"
          onclick={() => (showConfirmNoVolume = false)}
        >
          {t("images.config_no")}
        </button>
        <button
          type="button"
          class="px-4 py-2 text-xs font-bold rounded-xl border-none cursor-pointer text-white bg-emerald-600 hover:bg-emerald-700 transition-colors"
          onclick={submitForm}
        >
          {t("images.config_yes")}
        </button>
      </div>
    </div>
  </div>
{/if}
