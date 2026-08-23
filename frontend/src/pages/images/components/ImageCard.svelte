<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
  import type { DockerImage } from "$lib/domains/images";
  import { PrimaryButton, DangerButton } from "$lib/components/buttons";

  let {
    img,
    checked = false,
    on_toggle = () => {},
    on_build = () => {},
    on_delete = () => {},
  }: {
    img: DockerImage;
    checked?: boolean;
    on_toggle?: () => void;
    on_build?: () => void;
    on_delete?: () => void;
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
    if (searchTarget.includes("debian"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/debian/debian-original.svg";
    if (searchTarget.includes("alpine"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/alpinejs/alpinejs-original.svg";
    if (searchTarget.includes("redis"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/redis/redis-original.svg";
    if (searchTarget.includes("nginx"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nginx/nginx-original.svg";
    if (searchTarget.includes("node"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nodejs/nodejs-original.svg";
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
  class="relative rounded-2xl bg-[#f0f3f8] dark:bg-[#0c1220] border border-slate-300/80 dark:border-slate-800/80 hover:border-slate-400 dark:hover:border-slate-700 transition-all duration-300 flex flex-col justify-between overflow-hidden self-start w-full text-slate-700 dark:text-slate-200 shadow-md dark:shadow-lg dark:shadow-black/40 p-4 gap-3.5 group"
>
  <!-- Card Header -->
  <div class="flex items-center justify-between gap-3">
    <div class="flex items-center gap-3 min-w-0 flex-1">
      <!-- Checkbox -->
      <button
        type="button"
        class="w-5.5 h-5.5 rounded-lg border-2 flex items-center justify-center cursor-pointer transition-all duration-150 shrink-0 {checked
          ? 'bg-violet-600 border-violet-500 text-white'
          : 'border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-900/60 hover:border-violet-500'}"
        onclick={on_toggle}
      >
        {#if checked}
          <span class="text-white text-xs font-bold leading-none">✓</span>
        {/if}
      </button>

      <!-- Logo / Icon -->
      {#if iconUrl}
        <img
          src={iconUrl}
          alt={img.repo}
          class="w-5 h-5 object-contain shrink-0"
        />
      {:else}
        <div
          class="w-5 h-5 rounded-md bg-slate-200 dark:bg-slate-800 flex items-center justify-center shrink-0 text-[10px] font-bold text-slate-500"
        >
          📦
        </div>
      {/if}

      <!-- Title / Tag Button -->
      <button
        type="button"
        class="flex-1 font-mono font-bold text-sm tracking-tight truncate text-slate-855 dark:text-slate-100 px-4 py-2 rounded-2xl bg-white dark:bg-slate-800/30 border border-slate-200/80 dark:border-slate-700/30 flex items-center justify-between gap-2.5 cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors text-left shadow-xs min-w-0"
        onclick={() => (expanded = !expanded)}
      >
        <div class="flex items-center gap-2 truncate">
          <span
            class="truncate font-semibold text-slate-855 dark:text-white grow"
          >
            {#if img.repo === "<none>"}
              <span class="text-slate-400 dark:text-slate-500 italic"
                >&lt;sem tag&gt;</span
              >
            {:else}
              {img.repo}
            {/if}
          </span>
          <span
            class="text-xs px-2 py-0.5 rounded-lg bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 shrink-0 font-normal"
          >
            {img.tag}
          </span>
        </div>
        <span class="text-[10px] text-slate-400 dark:text-slate-500 shrink-0">
          {expanded ? "▲" : "▼"}
        </span>
      </button>
    </div>

    <!-- Status Badge -->
    <div class="flex items-center gap-2 shrink-0">
      <span
        class="px-2.5 py-1 rounded-xl text-xs font-semibold border flex items-center gap-1.5 {isInUse
          ? 'text-emerald-700 dark:text-emerald-400 bg-emerald-100 dark:bg-emerald-950/30 border-emerald-250 dark:border-emerald-900/50'
          : 'text-slate-600 dark:text-slate-400 bg-slate-100 dark:bg-slate-800/40 border-slate-200/80 dark:border-slate-700/50'}"
      >
        <span
          class="w-1.5 h-1.5 rounded-full {isInUse
            ? 'bg-emerald-500'
            : 'bg-slate-400'} shrink-0"
        ></span>
        {isInUse
          ? t("images.in_use_count", { count: containerCount })
          : t("images.unused")}
      </span>
    </div>
  </div>

  <!-- Expanded Details -->
  {#if expanded}
    <div class="flex flex-col gap-3 text-xs pt-1">
      <!-- ID & Size Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <!-- ID -->
        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-violet-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <div class="grow flex flex-col gap-0.5 min-w-0">
            <span
              class="text-[9px] text-violet-500 font-bold uppercase tracking-wider"
            >
              {t("images.card_image_id")}
            </span>
            <span
              class="font-mono text-xs font-bold text-blue-500 dark:text-blue-400 truncate"
            >
              {img.id.substring(0, 12)}
            </span>
          </div>
        </div>

        <!-- Size -->
        <div
          class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-cyan-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <div class="flex flex-col gap-0.5">
            <span
              class="text-[9px] text-cyan-500 font-bold uppercase tracking-wider"
            >
              {t("images.card_size")}
            </span>
            <span
              class="font-mono font-semibold text-slate-700 dark:text-slate-200 text-xs"
            >
              {img.size}
            </span>
          </div>
        </div>
      </div>

      <!-- Created At -->
      <div
        class="flex items-center justify-between p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-emerald-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
      >
        <div class="flex flex-col gap-0.5">
          <span
            class="text-[9px] text-emerald-500 font-bold uppercase tracking-wider"
          >
            {t("images.card_created")}
          </span>
          <span
            class="font-semibold text-slate-700 dark:text-slate-200 text-xs"
          >
            {createdStr || "—"}
          </span>
        </div>
      </div>

      <!-- Containers Using -->
      {#if img.containersUsing && img.containersUsing.length > 0}
        <div
          class="flex flex-col gap-2 p-3.5 rounded-xl bg-white dark:bg-[#070a12] border border-slate-200 dark:border-slate-800/80 border-l-4 border-l-indigo-500 shadow-md dark:shadow-black/40 hover:-translate-y-0.5 transition-all"
        >
          <span
            class="text-[9px] text-indigo-500 font-bold uppercase tracking-wider"
          >
            {t("images.card_containers_using")}
          </span>
          <div class="flex flex-wrap gap-2">
            {#each img.containersUsing as containerName}
              <span
                class="px-2.5 py-1 rounded-lg bg-slate-50 dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 text-slate-800 dark:text-slate-200 font-medium"
              >
                {containerName}
              </span>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Actions -->
      <div class="flex justify-end gap-2 pt-1">
        <PrimaryButton
          size="sm"
          onclick={on_build}
        >
          {t("images.btn_build_container")}
        </PrimaryButton>
        <DangerButton
          size="sm"
          onclick={on_delete}
        >
          {t("images.btn_delete")}
        </DangerButton>
      </div>
    </div>
  {/if}
</div>
