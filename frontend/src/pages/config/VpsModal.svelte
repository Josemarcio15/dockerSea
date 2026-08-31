<script lang="ts">
  import FormModal from "$shared/components/FormModal.svelte";
  import Input from "$shared/components/Input.svelte";
  import { ButtonYellow, ButtonPink } from "$shared/components/buttons";
  import { t } from "$shared/stores/locale.svelte";
  import {
    notifySuccess,
    notifyError,
  } from "$shared/stores/notification.svelte";

  import { Dialogs } from "@wailsio/runtime";

  export interface VpsFormData {
    id: string;
    name: string;
    connectionType: "local" | "ssh";
    host: string;
    port: string;
    username: string;
    authType: "key" | "password";
    sshKeyPath: string;
    sshKeyPassphrase: string;
    sshPassword: string;
    sudoPassword: string;
    dockerSocketPath: string;
    dockerPath: string;
    dockerComposePath: string;
  }

  let {
    show = $bindable(false),
    form = $bindable(),
    onSave,
    onTest,
    onAutoDetect,
  }: {
    show: boolean;
    form: VpsFormData;
    onSave: (data: VpsFormData) => void;
    onTest: (data: VpsFormData) => void;
    onAutoDetect: (data: VpsFormData) => Promise<any>;
  } = $props();

  let modalTab = $state<"connection" | "docker">("connection");
  let isDetecting = $state(false);

  // Lista de caminhos descobertos
  let discoveredSockets = $state<string[]>([]);
  let discoveredBins = $state<string[]>([]);
  let discoveredComposes = $state<string[]>([]);

  // Limpa o estado temporário e opções descobertas quando o modal abre ou fecha
  $effect(() => {
    if (!show) {
      discoveredSockets = [];
      discoveredBins = [];
      discoveredComposes = [];
      modalTab = "connection";
      isDetecting = false;
    }
  });

  async function pickSshKey() {
    try {
      const selected = await Dialogs.OpenFile({
        Title: "Selecione a chave privada SSH",
        CanChooseFiles: true,
        CanChooseDirectories: false,
        AllowsMultipleSelection: false,
      });
      if (selected) {
        const path = Array.isArray(selected) ? selected[0] : selected;
        if (path && typeof path === "string") {
          form.sshKeyPath = path;
          return;
        }
      }
    } catch (e) {
      console.warn("Wails Dialogs.OpenFile fallback:", e);
    }
    document.getElementById("modal-file-ssh-key")?.click();
  }

  async function handleAutoDetect() {
    isDetecting = true;
    try {
      const res: any = await onAutoDetect(form);
      if (res && res.success) {
        discoveredSockets = res.availableSockets || [
          res.dockerSocketPath || "/var/run/docker.sock",
        ];
        discoveredBins = res.availableBins || [
          res.dockerPath || "/usr/bin/docker",
        ];
        discoveredComposes = res.availableComposes || [
          res.dockerComposePath || "docker compose",
        ];

        if (
          !form.dockerSocketPath ||
          !discoveredSockets.includes(form.dockerSocketPath)
        ) {
          form.dockerSocketPath =
            res.dockerSocketPath ||
            discoveredSockets[0] ||
            "/var/run/docker.sock";
        }
        if (!form.dockerPath || !discoveredBins.includes(form.dockerPath)) {
          form.dockerPath =
            res.dockerPath || discoveredBins[0] || "/usr/bin/docker";
        }
        if (
          !form.dockerComposePath ||
          !discoveredComposes.includes(form.dockerComposePath)
        ) {
          form.dockerComposePath =
            res.dockerComposePath || discoveredComposes[0] || "docker compose";
        }
        notifySuccess(res.message || "Ambiente Docker detectado!");
      } else {
        notifyError(
          res?.message || "Não foi possível autodetectar. Verifique a conexão.",
        );
      }
    } catch (e: any) {
      notifyError(`Erro ao autodetectar: ${e.message || e}`);
    } finally {
      isDetecting = false;
    }
  }
</script>

<FormModal
  bind:show
  title={form.id ? t("servers.modal_edit_title") : t("servers.modal_new_title")}
  buttons={[
    {
      label: t("config.test_conn_btn"),
      variant: "secondary",
      onclick: () => onTest(form),
      disabled:
        !form.name.trim() ||
        (form.connectionType === "ssh" && !form.host.trim()),
    },
    {
      label: t("common.save"),
      variant: "primary",
      onclick: () => onSave(form),
      disabled: !form.name.trim(),
    },
  ]}
>
  <!-- Tabs Header -->
  <div
    class="flex items-center gap-2 p-1 rounded-xl bg-slate-100 dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800"
  >
    <button
      type="button"
      class="flex-1 py-2 text-xs font-bold rounded-lg transition-all cursor-pointer flex items-center justify-center gap-2 {modalTab ===
      'connection'
        ? 'bg-violet-600 text-white shadow-sm'
        : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-200'}"
      onclick={() => (modalTab = "connection")}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="w-4 h-4"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        stroke-width="2"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M12 21a9.004 9.004 0 0 0 8.716-6.747M12 21a9.004 9.004 0 0 1-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 0 1 7.843 4.582M12 3a8.997 8.997 0 0 0-7.843 4.582"
        />
      </svg>
      Conexão & Acesso
    </button>
    <button
      type="button"
      class="flex-1 py-2 text-xs font-bold rounded-lg transition-all cursor-pointer flex items-center justify-center gap-2 {modalTab ===
      'docker'
        ? 'bg-violet-600 text-white shadow-sm'
        : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-200'}"
      onclick={() => (modalTab = "docker")}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="w-4 h-4"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        stroke-width="2"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M5.25 14.25h13.5m-13.5 0a3 3 0 0 1-3-3m3 3a3 3 0 1 0 0 6h13.5a3 3 0 1 0 0-6m-16.5-3a3 3 0 0 1 3-3h13.5a3 3 0 0 1 3 3m-19.5 0a4.5 4.5 0 0 1 .9-2.7L5.75 5.1a3 3 0 0 1 2.4-1.1h7.7a3 3 0 0 1 2.4 1.1l2.1 3.45a4.5 4.5 0 0 1 .9 2.7"
        />
      </svg>
      Docker Engine
    </button>
  </div>

  {#if modalTab === "connection"}
    <!-- ABA 1: CONEXÃO -->
    <Input
      label={t("config.field_name")}
      placeholder={t("config.placeholder_name")}
      bind:value={form.name}
      required
    />

    <div class="flex flex-col gap-1.5">
      <label
        for="vps-conn"
        class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
        >{t("config.conn_type")}</label
      >
      <select
        id="vps-conn"
        class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all font-semibold"
        bind:value={form.connectionType}
      >
        <option value="local">{t("config.conn_local")}</option>
        <option value="ssh">{t("config.conn_ssh")}</option>
      </select>
    </div>

    {#if form.connectionType === "ssh"}
      <div class="grid grid-cols-3 gap-3">
        <div class="col-span-2">
          <Input
            label={t("config.field_host")}
            placeholder="192.168.1.100"
            bind:value={form.host}
            required
          />
        </div>
        <div class="col-span-1">
          <Input
            label={t("config.field_port")}
            type="number"
            placeholder="22"
            bind:value={form.port}
          />
        </div>
      </div>

      <Input
        label={t("config.field_user")}
        placeholder="root"
        bind:value={form.username}
      />

      <div
        class="space-y-4 pt-2 border-t border-slate-100 dark:border-slate-800"
      >
        <h4
          class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
        >
          Autenticação SSH
        </h4>

        <!-- Chave Privada -->
        <Input
          label={t("config.ssh_key_path")}
          placeholder="Ex: ~/.ssh/id_ed25519 (opcional se usar senha)"
          bind:value={form.sshKeyPath}
        >
          {#snippet trailing()}
            <ButtonYellow size="sm" onclick={pickSshKey}>
              {t("config.select_btn")}
            </ButtonYellow>
            <input
              id="modal-file-ssh-key"
              type="file"
              class="hidden"
              onchange={(e) => {
                const file = (e.target as HTMLInputElement).files?.[0];
                if (file) {
                  form.sshKeyPath = (file as any).path || file.name;
                }
              }}
            />
          {/snippet}
        </Input>

        {#if form.authType === "key"}
          <!-- Passphrase -->
          <Input
            label={t("config.ssh_key_passphrase")}
            type="password"
            placeholder={t("config.placeholder_pass")}
            help="Preencha apenas se a chave privada tiver passphrase."
            bind:value={form.sshKeyPassphrase}
          />
        {/if}

        <!-- Senha SSH -->
        <Input
          label="Senha do usuário SSH"
          type="password"
          placeholder="•••••••• (opcional se usar chave privada)"
          help="Deixe em branco se a VPS utilizar exclusivamente chave pública."
          bind:value={form.sshPassword}
        />

        <!-- Senha Sudo -->
        <Input
          label="Senha do Sudo (opcional)"
          type="password"
          placeholder="Em branco = sem senha / NOPASSWD"
          help="Deixe em branco caso o usuário tenha sudo sem senha (ex: Oracle Cloud, AWS, etc.) ou se for root."
          bind:value={form.sudoPassword}
        />
      </div>
    {/if}
  {:else}
    <!-- ABA 2: DOCKER ENGINE -->
    <div
      class="p-3.5 rounded-xl bg-violet-500/5 border border-violet-500/20 flex flex-col sm:flex-row sm:items-center justify-between gap-3"
    >
      <div class="space-y-0.5">
        <h5
          class="text-xs font-bold text-slate-800 dark:text-slate-200 flex items-center gap-1.5"
        >
          <span>⚡</span> Descoberta Automática de Ambiente
        </h5>
        <p class="text-[11px] text-slate-500 dark:text-slate-400">
          Descobrir automaticamente o socket UNIX, docker e docker compose nesta
          VPS.
        </p>
      </div>

      <ButtonPink
        size="sm"
        loading={isDetecting}
        disabled={form.connectionType === "ssh" && !form.host.trim()}
        onclick={handleAutoDetect}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="w-3.5 h-3.5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="2"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="m3.75 13.5 10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75Z"
          />
        </svg>
        Autodetectar
      </ButtonPink>
    </div>

    <!-- Socket Selector -->
    <div class="space-y-1.5">
      <Input
        label={t("config.docker_socket")}
        placeholder={t("config.placeholder_socket")}
        help="Caminho do socket UNIX (Padrão: /var/run/docker.sock)"
        bind:value={form.dockerSocketPath}
      />
      {#if discoveredSockets.length > 1}
        <div
          class="p-2 rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200 dark:border-slate-800 space-y-1.5"
        >
          <span
            class="text-[10px] font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider block"
          >
            Selecione o socket desejado ({discoveredSockets.length} encontrados):
          </span>
          <div class="flex flex-col gap-1.5">
            {#each discoveredSockets as sock}
              {@const isSelected = form.dockerSocketPath === sock}
              <button
                type="button"
                class="w-full px-3 py-2 text-left rounded-lg border transition-all cursor-pointer flex items-center justify-between gap-2 {isSelected
                  ? 'bg-violet-500/15 border-violet-500 text-violet-700 dark:text-violet-300 font-bold shadow-xs'
                  : 'bg-white dark:bg-[#0c1220] border-slate-200 dark:border-slate-800 hover:border-violet-300 text-slate-700 dark:text-slate-300'}"
                onclick={() => (form.dockerSocketPath = sock)}
              >
                <div class="flex items-center gap-2 min-w-0 flex-1">
                  <span class="text-xs shrink-0"
                    >{sock.includes("/user/") ? "👤" : "⚙️"}</span
                  >
                  <div class="flex flex-col min-w-0">
                    <span
                      class="text-[10px] font-semibold text-slate-400 dark:text-slate-500"
                    >
                      {sock.includes("/user/")
                        ? "Rootless (Ambiente de Usuário)"
                        : "Root / Sistema (Ambiente Padrão)"}
                    </span>
                    <span class="text-xs font-mono truncate">{sock}</span>
                  </div>
                </div>
                <div
                  class="w-4 h-4 rounded-full border flex items-center justify-center shrink-0 {isSelected
                    ? 'border-violet-500 bg-violet-500 text-white'
                    : 'border-slate-300 dark:border-slate-700'}"
                >
                  {#if isSelected}
                    <div class="w-1.5 h-1.5 rounded-full bg-white"></div>
                  {/if}
                </div>
              </button>
            {/each}
          </div>
        </div>
      {/if}
    </div>

    <!-- Docker Bin Selector -->
    <div class="space-y-1.5">
      <Input
        label={t("config.docker_path")}
        placeholder={t("config.placeholder_docker_path")}
        help={t("config.help_docker_path")}
        bind:value={form.dockerPath}
      />
      {#if discoveredBins.length > 1}
        <div
          class="p-2 rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200 dark:border-slate-800 space-y-1.5"
        >
          <span
            class="text-[10px] font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider block"
          >
            Selecione o binário do Docker ({discoveredBins.length} encontrados):
          </span>
          <div class="flex flex-col gap-1.5">
            {#each discoveredBins as bin}
              {@const isSelected = form.dockerPath === bin}
              <button
                type="button"
                class="w-full px-3 py-2 text-left rounded-lg border transition-all cursor-pointer flex items-center justify-between gap-2 {isSelected
                  ? 'bg-violet-500/15 border-violet-500 text-violet-700 dark:text-violet-300 font-bold shadow-xs'
                  : 'bg-white dark:bg-[#0c1220] border-slate-200 dark:border-slate-800 hover:border-violet-300 text-slate-700 dark:text-slate-300'}"
                onclick={() => (form.dockerPath = bin)}
              >
                <div class="flex items-center gap-2 min-w-0 flex-1">
                  <span class="text-xs shrink-0"
                    >{bin.includes("/home/") ? "👤" : "⚙️"}</span
                  >
                  <div class="flex flex-col min-w-0">
                    <span
                      class="text-[10px] font-semibold text-slate-400 dark:text-slate-500"
                    >
                      {bin.includes("/home/")
                        ? "Binário do Usuário (Rootless)"
                        : "Binário do Sistema"}
                    </span>
                    <span class="text-xs font-mono truncate">{bin}</span>
                  </div>
                </div>
                <div
                  class="w-4 h-4 rounded-full border flex items-center justify-center shrink-0 {isSelected
                    ? 'border-violet-500 bg-violet-500 text-white'
                    : 'border-slate-300 dark:border-slate-700'}"
                >
                  {#if isSelected}
                    <div class="w-1.5 h-1.5 rounded-full bg-white"></div>
                  {/if}
                </div>
              </button>
            {/each}
          </div>
        </div>
      {/if}
    </div>

    <!-- Docker Compose Selector -->
    <div class="space-y-1.5">
      <Input
        label={t("config.docker_compose_path")}
        placeholder={t("config.placeholder_docker_compose_path")}
        help={t("config.help_docker_compose_path")}
        bind:value={form.dockerComposePath}
      />
      {#if discoveredComposes.length > 1}
        <div
          class="p-2 rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200 dark:border-slate-800 space-y-1.5"
        >
          <span
            class="text-[10px] font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider block"
          >
            Selecione o Compose ({discoveredComposes.length} encontrados):
          </span>
          <div class="flex flex-col gap-1.5">
            {#each discoveredComposes as comp}
              {@const isSelected = form.dockerComposePath === comp}
              <button
                type="button"
                class="w-full px-3 py-2 text-left rounded-lg border transition-all cursor-pointer flex items-center justify-between gap-2 {isSelected
                  ? 'bg-violet-500/15 border-violet-500 text-violet-700 dark:text-violet-300 font-bold shadow-xs'
                  : 'bg-white dark:bg-[#0c1220] border-slate-200 dark:border-slate-800 hover:border-violet-300 text-slate-700 dark:text-slate-300'}"
                onclick={() => (form.dockerComposePath = comp)}
              >
                <div class="flex items-center gap-2 min-w-0 flex-1">
                  <span class="text-xs shrink-0">📦</span>
                  <div class="flex flex-col min-w-0">
                    <span
                      class="text-[10px] font-semibold text-slate-400 dark:text-slate-500"
                    >
                      {comp.includes("compose") &&
                      !comp.includes("docker-compose")
                        ? "Plugin V2"
                        : "Executável Standalone"}
                    </span>
                    <span class="text-xs font-mono truncate">{comp}</span>
                  </div>
                </div>
                <div
                  class="w-4 h-4 rounded-full border flex items-center justify-center shrink-0 {isSelected
                    ? 'border-violet-500 bg-violet-500 text-white'
                    : 'border-slate-300 dark:border-slate-700'}"
                >
                  {#if isSelected}
                    <div class="w-1.5 h-1.5 rounded-full bg-white"></div>
                  {/if}
                </div>
              </button>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  {/if}
</FormModal>
