<script lang="ts">
  import Modal from "$lib/components/Modal.svelte";
  import Input from "$lib/components/Input.svelte";
  import { SecondaryButton } from "$lib/components/buttons";
  import { t } from "$lib/stores/locale.svelte";
  import { notifySuccess, notifyError } from "$lib/stores/notification.svelte";
  import * as ConfigService from "../../../bindings/go-walis/internal/config/configservice.js";

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
  }: {
    show: boolean;
    form: VpsFormData;
    onSave: (data: VpsFormData) => void;
    onTest: (data: VpsFormData) => void;
  } = $props();

  import { Dialogs } from "@wailsio/runtime";

  let modalTab = $state<"connection" | "docker">("connection");
  let isDetecting = $state(false);

  async function pickSshKey() {
    try {
      const selected = await Dialogs.OpenFile({
        title: "Selecione a chave privada SSH",
        canChooseFiles: true,
        canChooseDirectories: false,
        allowsMultipleSelection: false,
      });
      if (selected && typeof selected === "string") {
        form.sshKeyPath = selected;
        return;
      }
    } catch (e) {
      console.warn("Wails runtime dialog fallback:", e);
    }
    document.getElementById("modal-file-ssh-key")?.click();
  }

  async function handleAutoDetect() {
    isDetecting = true;
    try {
      const payload: any = {
        id: form.id,
        name: form.name,
        connectionType: form.connectionType,
        host: form.host,
        port: parseInt(form.port) || 22,
        username: form.username,
        authType: form.authType,
        sshKeyPath: form.sshKeyPath,
        sshKeyPassphrase: form.sshKeyPassphrase,
        sshPassword: form.sshPassword,
        sudoPassword: form.sudoPassword,
        dockerSocketPath: form.dockerSocketPath,
        dockerPath: form.dockerPath,
        dockerComposePath: form.dockerComposePath,
        isActive: false,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      const res = await ConfigService.AutoDetectDocker(payload);
      if (res && res.success) {
        form.dockerSocketPath = res.dockerSocketPath || "/var/run/docker.sock";
        form.dockerPath = res.dockerPath || "/usr/bin/docker";
        form.dockerComposePath = res.dockerComposePath || "docker compose";
        notifySuccess(res.message || "Ambiente Docker detectado!");
      } else {
        notifyError(res?.message || "Não foi possível autodetectar. Verifique a conexão.");
      }
    } catch (e: any) {
      notifyError(`Erro ao autodetectar: ${e.message || e}`);
    } finally {
      isDetecting = false;
    }
  }
</script>

<Modal
  bind:show
  title={form.id ? t("config.edit_title") : t("config.new_title")}
  buttons={[
    {
      label: "🔍 " + t("config.test_conn_btn"),
      variant: "secondary",
      onclick: () => onTest(form),
      disabled: !form.name.trim() || (form.connectionType === "ssh" && !form.host.trim()),
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
      <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 21a9.004 9.004 0 0 0 8.716-6.747M12 21a9.004 9.004 0 0 1-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 0 1 7.843 4.582M12 3a8.997 8.997 0 0 0-7.843 4.582" />
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
      <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M5.25 14.25h13.5m-13.5 0a3 3 0 0 1-3-3m3 3a3 3 0 1 0 0 6h13.5a3 3 0 1 0 0-6m-16.5-3a3 3 0 0 1 3-3h13.5a3 3 0 0 1 3 3m-19.5 0a4.5 4.5 0 0 1 .9-2.7L5.75 5.1a3 3 0 0 1 2.4-1.1h7.7a3 3 0 0 1 2.4 1.1l2.1 3.45a4.5 4.5 0 0 1 .9 2.7" />
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

      <div class="space-y-4 pt-2 border-t border-slate-100 dark:border-slate-800">
        <h4 class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider">
          Autenticação SSH
        </h4>

        <!-- Chave Privada -->
        <Input
          label={t("config.ssh_key_path")}
          placeholder="Ex: ~/.ssh/id_ed25519 (opcional se usar senha)"
          bind:value={form.sshKeyPath}
        >
          {#snippet trailing()}
            <button
              type="button"
              class="px-3.5 py-2.5 text-xs font-bold rounded-xl border border-violet-200 dark:border-violet-800/60 text-violet-700 dark:text-violet-300 bg-violet-50 dark:bg-violet-950/40 hover:bg-violet-100 dark:hover:bg-violet-900/60 cursor-pointer transition-colors shrink-0"
              onclick={pickSshKey}
            >
              {t("config.select_btn")}
            </button>
            <input
              id="modal-file-ssh-key"
              type="file"
              class="hidden"
              onchange={(e) => {
                const file = (e.target as HTMLInputElement).files?.[0];
                if (file) {
                  const fullPath = (file as any).path;
                  form.sshKeyPath = fullPath || `~/.ssh/${file.name}`;
                }
                (e.target as HTMLInputElement).value = "";
              }}
            />
          {/snippet}
        </Input>

        {#if form.sshKeyPath.trim()}
          <Input
            label={t("config.ssh_key_pass")}
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
        <h5 class="text-xs font-bold text-slate-800 dark:text-slate-200 flex items-center gap-1.5">
          <span>⚡</span> Descoberta Automática de Ambiente
        </h5>
        <p class="text-[11px] text-slate-500 dark:text-slate-400">
          Descobrir automaticamente o socket UNIX, docker e docker compose nesta VPS.
        </p>
      </div>

      <SecondaryButton
        size="sm"
        loading={isDetecting}
        disabled={form.connectionType === "ssh" && !form.host.trim()}
        onclick={handleAutoDetect}
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="m3.75 13.5 10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75Z" />
        </svg>
        Autodetectar
      </SecondaryButton>
    </div>

    <Input
      label={t("config.docker_socket")}
      placeholder={t("config.placeholder_socket")}
      help="Caminho do socket UNIX (Padrão: /var/run/docker.sock)"
      bind:value={form.dockerSocketPath}
    />

    <Input
      label={t("config.docker_path")}
      placeholder={t("config.placeholder_docker_path")}
      help={t("config.help_docker_path")}
      bind:value={form.dockerPath}
    />

    <Input
      label={t("config.docker_compose_path")}
      placeholder={t("config.placeholder_docker_compose_path")}
      help={t("config.help_docker_compose_path")}
      bind:value={form.dockerComposePath}
    />
  {/if}
</Modal>
