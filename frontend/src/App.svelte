<script lang="ts">
  import StatusBanner from "$shared/components/StatusBanner.svelte";
  import SidebarItem from "$shared/components/SidebarItem.svelte";
  import ActiveServerWidget from "$shared/components/ActiveServerWidget.svelte";
  import { setLocale, t } from "$shared/stores/locale.svelte";
  import { onMount } from "svelte";

  // Import page components
  import Servers from "./pages/servers/Page.svelte";
  import Containers from "./pages/containers/Page.svelte";
  import Images from "./pages/images/Page.svelte";
  import Volumes from "./pages/volumes/Page.svelte";
  import Networks from "./pages/networks/Page.svelte";
  import Stacks from "./pages/stacks/Page.svelte";
  import Builder from "./pages/builder/Page.svelte";
  import Config from "./pages/config/Page.svelte";
  import Extras from "./pages/extras/Page.svelte";
  import Profiles from "./pages/profiles/Page.svelte";
  import { loadSession, session } from "$session/session.svelte";
  import { navigation, navigate } from "$navigation/navigation.svelte";

  let { data = {} }: { data?: any } = $props();

  let appData = $state<Record<string, any>>({
    locale: "pt-BR",
    theme: "dark",
    servers: session.servers,
    activeVps: session.activeVps,
    profiles: session.profiles,
    activeProfile: session.activeProfile,
  });

  $effect(() => {
    if (data && typeof data === "object") {
      Object.assign(appData, data);
    }
  });

  $effect(() => {
    appData.servers = session.servers;
    appData.activeVps = session.activeVps;
    appData.profiles = session.profiles;
    appData.activeProfile = session.activeProfile;
  });

  // Sync locale
  $effect(() => {
    if (appData?.locale) {
      setLocale(appData.locale);
    }
  });

  let isDark = $state(true);
  let appVersion = "0.0.5-alpha";
  let hasProfile = $state(true);

  let darkClass = $derived(isDark ? "dark" : "");

  onMount(async () => {
    const storedTheme = localStorage.getItem("theme");
    if (storedTheme) {
      isDark = storedTheme === "dark";
    } else {
      isDark = data?.theme
        ? data.theme === "dark"
        : window.matchMedia("(prefers-color-scheme: dark)").matches;
    }

    try {
      await loadSession();
      if (session.activeProfile?.locale) {
        appData.locale = session.activeProfile.locale;
      }
    } catch (e: any) {
      console.warn("Erro ao buscar servidores/perfis do SQLite no App:", e);
    }
  });

  $effect(() => {
    if (typeof document !== "undefined") {
      if (isDark) {
        document.documentElement.classList.add("dark");
        localStorage.setItem("theme", "dark");
      } else {
        document.documentElement.classList.remove("dark");
        localStorage.setItem("theme", "light");
      }
    }
  });

  function toggleTheme() {
    isDark = !isDark;
  }
</script>

{#snippet iconHome()}
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.5"
    class="w-5 h-5"
    ><path
      stroke-linecap="round"
      stroke-linejoin="round"
      d="m2.25 12 8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25"
    /></svg
  >
{/snippet}

{#snippet iconContainers()}
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.5"
    class="w-5 h-5"
    ><path
      stroke-linecap="round"
      stroke-linejoin="round"
      d="m21 7.5-9-5.25L3 7.5m18 0-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-9v9"
    /></svg
  >
{/snippet}

{#snippet iconImages()}
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.5"
    class="w-5 h-5"
    ><path
      stroke-linecap="round"
      stroke-linejoin="round"
      d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909M3.75 21h16.5A2.25 2.25 0 0 0 22.5 18.75V5.25A2.25 2.25 0 0 0 20.25 3H3.75A2.25 2.25 0 0 0 1.5 5.25v13.5A2.25 2.25 0 0 0 3.75 21Zm16.5-13.5h.008v.008h-.008V7.5Z"
    /></svg
  >
{/snippet}

{#snippet iconVolumes()}
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.5"
    class="w-5 h-5"
    ><path
      stroke-linecap="round"
      stroke-linejoin="round"
      d="m20.25 7.5-.625 10.632a2.25 2.25 0 0 1-2.247 2.118H6.622a2.25 2.25 0 0 1-2.247-2.118L3.75 7.5m8.25 3v6.75m0 0-3-3m3 3 3-3M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125Z"
    /></svg
  >
{/snippet}

{#snippet iconNetworks()}
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.5"
    class="w-5 h-5"
    ><path
      stroke-linecap="round"
      stroke-linejoin="round"
      d="M12 21a9.004 9.004 0 0 0 8.716-6.747M12 21a9.004 9.004 0 0 1-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 0 1 7.843 4.582M12 3a8.997 8.997 0 0 0-7.843 4.582m15.686 0A11.953 11.953 0 0 1 12 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0 1 21 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0 1 12 16.5c-3.162 0-6.133-.815-8.716-2.247m0 0A9.015 9.015 0 0 1 3 12c0-1.605.42-3.113 1.157-4.418"
    /></svg
  >
{/snippet}

{#snippet iconPorts()}
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.5"
    class="w-5 h-5"
    ><path
      stroke-linecap="round"
      stroke-linejoin="round"
      d="M8.25 6.75h7.5M8.25 12h7.5m-7.5 5.25h7.5M5.25 3.75h13.5A2.25 2.25 0 0 1 21 6v12a2.25 2.25 0 0 1-2.25 2.25H5.25A2.25 2.25 0 0 1 3 18V6a2.25 2.25 0 0 1 2.25-2.25Z"
    /></svg
  >
{/snippet}

{#snippet iconStacks()}
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.5"
    class="w-5 h-5"
    ><path
      stroke-linecap="round"
      stroke-linejoin="round"
      d="M3.75 12h16.5m-16.5 3.75h16.5M3.75 19.5h16.5M5.625 4.5h12.75a1.875 1.875 0 0 1 0 3.75H5.625a1.875 1.875 0 0 1 0-3.75Z"
    /></svg
  >
{/snippet}

{#snippet iconBuilder()}
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.5"
    class="w-5 h-5"
    ><path
      stroke-linecap="round"
      stroke-linejoin="round"
      d="M21.75 6.75a4.5 4.5 0 0 1-4.884 4.484c-1.076-.091-2.264.071-2.95.904l-7.152 7.154a2.25 2.25 0 1 1-3.181-3.182l7.152-7.152c.833-.687.995-1.874.904-2.95a4.5 4.5 0 0 1 6.336-4.486l-3.398 3.398a.75.75 0 0 0 0 1.06l2.038 2.038a.75.75 0 0 0 1.06 0l3.398-3.398c.213.565.293 1.19.271 1.832Z"
    /></svg
  >
{/snippet}

{#snippet iconConfig()}
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.5"
    class="w-5 h-5"
    ><path
      stroke-linecap="round"
      stroke-linejoin="round"
      d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.325.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 0 1 1.37.49l1.296 2.247a1.125 1.125 0 0 1-.26 1.431l-1.003.827c-.293.241-.438.613-.43.992a7.723 7.723 0 0 1 0 .255c-.008.378.137.75.43.991l1.004.827c.424.35.534.955.26 1.43l-1.298 2.247a1.125 1.125 0 0 1-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.47 6.47 0 0 1-.22.128c-.331.183-.581.495-.644.869l-.213 1.281c-.09.543-.56.94-1.11.94h-2.594c-.55 0-1.02-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 0 1-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 0 1-1.369-.49l-1.297-2.247a1.125 1.125 0 0 1 .26-1.431l1.004-.827c.292-.24.437-.613.43-.991a6.932 6.932 0 0 1 0-.255c.007-.38-.138-.751-.43-.992l-1.004-.827a1.125 1.125 0 0 1-.26-1.43l1.297-2.247a1.125 1.125 0 0 1 1.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.086.22-.128.332-.183.582-.495.644-.869l.214-1.28Z"
    /></svg
  >
{/snippet}

{#snippet iconProfile()}
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.5"
    class="w-5 h-5"
    ><path
      stroke-linecap="round"
      stroke-linejoin="round"
      d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z"
    /></svg
  >
{/snippet}

<div class="flex h-screen w-screen font-sans overflow-hidden {darkClass}">
  <StatusBanner />
  <div
    class="w-60 h-full bg-linear-to-b from-violet-900 via-violet-800 to-indigo-900 text-white flex flex-col shadow-2xl shrink-0 overflow-y-auto overflow-x-hidden"
  >
    <div
      class="px-6 py-5 text-base font-bold text-white border-b border-white/10 flex items-center gap-3 bg-white/5"
    >
      <div
        class="w-9 h-9 rounded-xl bg-linear-to-br from-violet-400 to-fuchsia-400 flex items-center justify-center shadow-lg shadow-violet-500/30"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 512 512"
          class="w-6 h-6 text-white"
          fill="none"
          stroke="currentColor"
          stroke-width="28"
        >
          <circle cx="256" cy="220" r="80" stroke-width="32" />
          <line x1="256" y1="110" x2="256" y2="330" stroke-width="28" />
          <line x1="146" y1="220" x2="366" y2="220" stroke-width="28" />
          <line x1="178" y1="142" x2="334" y2="298" stroke-width="28" />
          <line x1="178" y1="298" x2="334" y2="142" stroke-width="28" />
          <circle cx="256" cy="220" r="24" fill="currentColor" />
          <path
            fill="currentColor"
            stroke="none"
            d="M 0 420 Q 140 370 270 430 T 512 400 L 512 512 L 0 512 Z"
          />
        </svg>
      </div>
      <div class="flex flex-col">
        <span>DockSea</span>
        <span class="text-[10px] text-violet-300 font-normal tracking-wider"
          >CONTAINER MANAGER</span
        >
      </div>
    </div>

    <div class="flex flex-col gap-1 px-3 pt-5">
      <SidebarItem
        icon={iconHome}
        iconBg="from-emerald-400 to-teal-400"
        label={t("sidebar.devices")}
        active={navigation.currentRoute === "servers"}
        disabled={!hasProfile}
        onclick={() => navigate("servers")}
      />
      <SidebarItem
        icon={iconContainers}
        iconBg="from-blue-400 to-cyan-400"
        label={t("sidebar.containers")}
        active={navigation.currentRoute === "containers"}
        disabled={!hasProfile}
        onclick={() => navigate("containers")}
      />
      <SidebarItem
        icon={iconImages}
        iconBg="from-amber-400 to-orange-400"
        label={t("sidebar.images")}
        active={navigation.currentRoute === "images"}
        disabled={!hasProfile}
        onclick={() => navigate("images")}
      />
      <SidebarItem
        icon={iconVolumes}
        iconBg="from-cyan-400 to-teal-400"
        label={t("sidebar.volumes")}
        active={navigation.currentRoute === "volumes"}
        disabled={!hasProfile}
        onclick={() => navigate("volumes")}
      />
      <SidebarItem
        icon={iconNetworks}
        iconBg="from-green-400 to-teal-400"
        label={t("sidebar.networks")}
        active={navigation.currentRoute === "networks"}
        disabled={!hasProfile}
        onclick={() => navigate("networks")}
      />
      <SidebarItem
        icon={iconStacks}
        iconBg="from-rose-400 to-pink-400"
        label={t("sidebar.stacks")}
        active={navigation.currentRoute === "stacks"}
        disabled={!hasProfile}
        onclick={() => navigate("stacks")}
      />
      <SidebarItem
        icon={iconBuilder}
        iconBg="from-amber-400 to-orange-400"
        label={t("sidebar.builder")}
        active={navigation.currentRoute === "builder"}
        disabled={!hasProfile}
        onclick={() => navigate("builder")}
      />
      <SidebarItem
        icon={iconConfig}
        iconBg="from-slate-400 to-slate-500"
        label={t("sidebar.configs")}
        active={navigation.currentRoute === "config"}
        disabled={!hasProfile}
        onclick={() => navigate("config")}
      />
      <SidebarItem
        icon={iconConfig}
        iconBg="from-cyan-400 to-blue-500"
        label="Extras"
        active={navigation.currentRoute === "extras"}
        disabled={!hasProfile}
        onclick={() => navigate("extras")}
      />
      <SidebarItem
        icon={iconProfile}
        iconBg="from-fuchsia-400 to-violet-500"
        label={t("sidebar.profiles")}
        active={navigation.currentRoute === "profiles"}
        disabled={false}
        onclick={() => navigate("profiles")}
      />
    </div>

    <!-- Active Server Widget -->
    <ActiveServerWidget vps={appData.activeVps} />

    <div
      class="mx-3 mt-2 p-3.5 rounded-xl bg-white/4 border border-white/10 backdrop-blur-md shadow-lg flex flex-col gap-2"
    >
      <span
        class="text-[10px] text-violet-300 font-bold uppercase tracking-wider"
        >{t("app.active_profile")}</span
      >
      <div class="flex items-center gap-2.5">
        <div
          class="w-8 h-8 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center text-violet-300 shrink-0 shadow-inner"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="w-4 h-4"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M17.982 18.725A7.488 7.488 0 0 0 12 15.75a7.488 7.488 0 0 0-5.982 2.975m11.963 0a9 9 0 1 0-11.963 0m11.963 0A8.966 8.966 0 0 1 12 21a8.966 8.966 0 0 1-5.982-2.275M15 9.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
            />
          </svg>
        </div>
        <div class="flex flex-col min-w-0 flex-1">
          <span class="text-sm text-violet-200 font-semibold truncate"
            >{appData?.activeProfile?.name || "Perfil Default"}</span
          >
          <span class="text-[10px] text-violet-300/70 truncate"
            >{appData?.activeProfile?.locale || "pt-BR"}</span
          >
        </div>
      </div>
    </div>

    <div class="flex-1"></div>

    <div
      class="mx-3 mb-2 px-4 py-2 rounded-xl text-[11px] text-orange-400 font-semibold text-center border border-orange-500/20 bg-orange-500/10"
    >
      {t("version", { version: appVersion })}
    </div>

    <button
      class="mx-3 mb-3 px-4 py-3 rounded-xl border border-white/10 flex items-center gap-3 text-sm text-violet-300 hover:text-white hover:bg-white/10 cursor-pointer transition-all duration-200"
      onclick={toggleTheme}
    >
      <span class="text-lg">{isDark ? "☀️" : "🌙"}</span>
      {isDark ? t("app.mode_light") : t("app.mode_dark")}
    </button>
  </div>

  <!-- Main Content Area where routes (pages) are rendered -->
  <div
    class="flex-1 p-6 bg-slate-200 dark:from-slate-900 dark:to-indigo-950 dark:bg-gradient-to-br overflow-y-auto"
  >
    {#if navigation.currentRoute === "servers"}
      <Servers data={appData} {navigate} />
    {:else if navigation.currentRoute === "containers"}
      <Containers data={appData} />
    {:else if navigation.currentRoute === "images"}
      <Images data={appData} />
    {:else if navigation.currentRoute === "volumes"}
      <Volumes data={appData} />
    {:else if navigation.currentRoute === "networks"}
      <Networks data={appData} />
    {:else if navigation.currentRoute === "stacks"}
      <Stacks data={appData} />
    {:else if navigation.currentRoute === "builder"}
      <Builder data={appData} />
    {:else if navigation.currentRoute === "config"}
      <Config data={appData} />
    {:else if navigation.currentRoute === "extras"}
      <Extras data={appData} />
    {:else if navigation.currentRoute === "profiles"}
      <Profiles data={appData} />
    {/if}
  </div>
</div>
