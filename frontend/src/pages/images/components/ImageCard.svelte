<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import type { DockerImage } from "$lib/domains/images";
  import { ButtonBlue } from "$lib/components/buttons";

  let {
    img,
    checked = false,
    on_toggle = () => {},
    on_build = () => {},
  }: {
    img: DockerImage;
    checked?: boolean;
    on_toggle?: () => void;
    on_build?: () => void;
  } = $props();

  let expanded = $state(false);

  const isInUse = $derived(
    img.containersUsing && img.containersUsing.length > 0,
  );
  const containerCount = $derived(
    img.containersUsing ? img.containersUsing.length : 0,
  );
  const repoLower = $derived(img.repo.toLowerCase());
  const containersHint = $derived(
    (img.containersUsing || []).join(" ").toLowerCase(),
  );

  const iconUrl = $derived.by(() => {
    const searchTarget =
      repoLower === "<none>" ? `${repoLower} ${containersHint}` : repoLower;

    if (searchTarget.includes("wordpress"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/wordpress/wordpress-plain.svg";
    if (searchTarget.includes("postgres") || searchTarget.includes("postgre"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/postgresql/postgresql-original.svg";
    if (searchTarget.includes("mariadb"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mariadb/mariadb-original.svg";
    if (searchTarget.includes("mysql"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mysql/mysql-original.svg";
    if (searchTarget.includes("mongo"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mongodb/mongodb-original.svg";
    if (searchTarget.includes("redis"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/redis/redis-original.svg";
    if (searchTarget.includes("debian"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/debian/debian-original.svg";
    if (searchTarget.includes("ubuntu"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/ubuntu/ubuntu-plain.svg";
    if (searchTarget.includes("alpine"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/alpinejs/alpinejs-original.svg";
    if (searchTarget.includes("nginx"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nginx/nginx-original.svg";
    if (searchTarget.includes("apache") || searchTarget.includes("httpd"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/apache/apache-original.svg";
    if (searchTarget.includes("node"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nodejs/nodejs-original.svg";
    if (searchTarget.includes("python"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/python/python-original.svg";
    if (searchTarget.includes("golang") || searchTarget.includes("go:"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg";
    if (searchTarget.includes("rust"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/rust/rust-plain.svg";
    if (searchTarget.includes("php"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/php/php-original.svg";
    if (searchTarget.includes("traefik"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/traefik/traefik-original.svg";
    if (searchTarget.includes("grafana"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/grafana/grafana-original.svg";
    if (searchTarget.includes("prometheus"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/prometheus/prometheus-original.svg";
    if (searchTarget.includes("rabbitmq"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/rabbitmq/rabbitmq-original.svg";
    if (searchTarget.includes("elasticsearch") || searchTarget.includes("elastic"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/elasticsearch/elasticsearch-original.svg";
    if (searchTarget.includes("docker"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/docker/docker-original.svg";
    return null;
  });

  const createdStr = $derived.by(() => {
    if (!img.created) return "";
    const date = new Date(img.created * 1000);
    const day = String(date.getDate()).padStart(2, "0");
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
  });
</script>

<div
  class="relative rounded-2xl bg-white dark:bg-[#0b101d] border border-slate-200/80 dark:border-slate-800/80 hover:border-violet-500/40 dark:hover:border-violet-500/40 hover:shadow-lg dark:hover:shadow-violet-950/20 transition-all duration-200 flex flex-col justify-between overflow-hidden self-start w-full text-slate-700 dark:text-slate-200 shadow-sm p-3.5 gap-3 group"
>
  <!-- Card Header Principal -->
  <div class="flex items-center gap-2.5 min-w-0">
    <!-- Checkbox Customizado -->
    <button
      type="button"
      class="w-5 h-5 rounded-lg border flex items-center justify-center cursor-pointer transition-all duration-150 shrink-0 {checked
        ? 'bg-violet-600 border-violet-500 text-white shadow-xs'
        : 'border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 hover:border-violet-400'}"
      onclick={on_toggle}
    >
      {#if checked}
        <span class="text-white text-[11px] font-bold leading-none">✓</span>
      {/if}
    </button>

    <!-- Logo / Ícone -->
    <div
      class="w-10 h-10 rounded-xl bg-slate-50 dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 p-2 flex items-center justify-center shrink-0 shadow-inner"
    >
      {#if iconUrl}
        <img
          src={iconUrl}
          alt={img.repo}
          class="w-full h-full object-contain"
        />
      {:else}
        <span class="text-base font-bold text-violet-500">IMG</span>
      {/if}
    </div>

    <!-- Título + Tag + Status Compacto -->
    <div class="flex flex-col min-w-0 flex-1">
      <div class="flex items-center gap-1.5 min-w-0">
        <span
          class="font-bold text-xs text-slate-900 dark:text-white truncate"
          title={img.repo}
        >
          {#if !img.repo || img.repo === "<none>"}
            <span class="text-slate-400 italic font-normal">&lt;sem tag&gt;</span>
          {:else}
            {img.repo}
          {/if}
        </span>
      </div>

      <div class="flex items-center gap-1.5 mt-0.5">
        <span
          class="text-[10px] font-mono px-1.5 py-0.2 rounded bg-violet-50 dark:bg-violet-950/50 text-violet-700 dark:text-violet-300 font-semibold border border-violet-200/50 dark:border-violet-800/40 truncate max-w-[85px]"
          title={img.tag}
        >
          {img.tag || "latest"}
        </span>

        <span
          class="text-[10px] font-medium flex items-center gap-1 {isInUse
            ? 'text-emerald-600 dark:text-emerald-400'
            : 'text-slate-400 dark:text-slate-500'}"
        >
          <span
            class="w-1.5 h-1.5 rounded-full {isInUse
              ? 'bg-emerald-500'
              : 'bg-slate-300 dark:bg-slate-600'} shrink-0"
          ></span>
          {isInUse ? `${containerCount} em uso` : 'Livre'}
        </span>
      </div>
    </div>

    <!-- Botão Expandir / Recolher Detalhes -->
    <button
      type="button"
      class="w-7 h-7 rounded-lg bg-slate-100 dark:bg-slate-800/60 hover:bg-slate-200 dark:hover:bg-slate-800 text-slate-500 dark:text-slate-400 flex items-center justify-center cursor-pointer transition-colors text-[10px] shrink-0 border border-slate-200/50 dark:border-slate-700/50"
      onclick={() => (expanded = !expanded)}
      title={expanded ? "Recolher detalhes" : "Mais detalhes"}
    >
      {expanded ? "▲" : "▼"}
    </button>
  </div>

  <!-- Detalhes Extras (visíveis apenas quando expandido) -->
  {#if expanded}
    <div class="flex flex-col gap-2.5 text-xs pt-1 border-t border-slate-100 dark:border-slate-800/60 animate-fadeIn">
      <!-- ID & Tamanho Grid -->
      <div class="grid grid-cols-2 gap-2 text-[11px]">
        <div class="flex flex-col p-2 rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200/50 dark:border-slate-800/50">
          <span class="text-[9px] text-slate-400 uppercase font-bold tracking-wider">ID</span>
          <span class="font-mono text-blue-600 dark:text-blue-400 font-semibold truncate">{img.id.substring(0, 8)}</span>
        </div>

        <div class="flex flex-col p-2 rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200/50 dark:border-slate-800/50">
          <span class="text-[9px] text-slate-400 uppercase font-bold tracking-wider">
            {t("images.card_size")}
          </span>
          <span class="font-semibold text-slate-700 dark:text-slate-200 truncate">
            {img.size}
          </span>
        </div>
      </div>

      <!-- Data de Criação -->
      <div class="flex items-center justify-between text-[11px] px-1">
        <span class="text-slate-400 text-[10px] uppercase font-bold tracking-wider">
          {t("images.card_created")}
        </span>
        <span class="font-medium text-slate-700 dark:text-slate-300">
          {createdStr || "—"}
        </span>
      </div>

      <!-- Containers Usando -->
      {#if img.containersUsing && img.containersUsing.length > 0}
        <div class="flex flex-col gap-1.5 p-2 rounded-xl bg-purple-50/40 dark:bg-purple-950/20 border border-purple-200/50 dark:border-purple-900/30">
          <span class="text-[9px] text-purple-600 dark:text-purple-400 font-bold uppercase tracking-wider">
            {t("images.card_containers_using")}
          </span>
          <div class="flex flex-wrap gap-1">
            {#each img.containersUsing as containerName}
              <span class="px-2 py-0.5 rounded-md bg-white dark:bg-slate-900 border border-purple-200/60 dark:border-purple-900/50 text-[10px] text-purple-700 dark:text-purple-300 font-medium">
                {containerName}
              </span>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  {/if}

  <!-- Ação Sempre Visível (mesmo recolhido) -->
  <div class="pt-1 border-t border-slate-100 dark:border-slate-800/60">
    <ButtonBlue
      size="xs"
      class="w-full justify-center"
      onclick={on_build}
    >
      {t("images.btn_build_container")}
    </ButtonBlue>
  </div>
</div>
