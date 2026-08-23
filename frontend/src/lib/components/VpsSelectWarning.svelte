<script lang="ts">
  import { t } from "$lib/stores/locale.svelte";
</script>

<div
  class="flex flex-col items-center justify-center min-h-[65vh] p-6 text-center animate-fadeIn"
>
  <div
    class="max-w-md w-full bg-white dark:bg-[#0c1220] border border-slate-200 dark:border-slate-800 rounded-3xl p-8 shadow-xl flex flex-col items-center space-y-5"
  >
    <!-- Icon Container -->
    <div
      class="w-16 h-16 rounded-2xl bg-amber-500/10 border border-amber-500/20 text-amber-500 flex items-center justify-center shadow-inner"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="1.75"
        class="w-8 h-8"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M5.25 14.25h13.5m-13.5 0a3 3 0 0 1-3-3m3 3a3 3 0 0 0 3 3h7.5a3 3 0 0 0 3-3m-13.5 0V6a3 3 0 0 1 3-3h7.5a3 3 0 0 1 3 3v5.25m-13.5 0h13.5m-13.5 0a3 3 0 0 0-3 3v3a3 3 0 0 0 3 3h13.5a3 3 0 0 0 3-3v-3a3 3 0 0 0-3-3m-9.75-6.75h.008v.008H9.75V7.5Zm0 9h.008v.008H9.75v-.008Zm0 3h.008v.008H9.75v-.008Z"
        />
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M12 9v3.75m0 3.75h.008v.008H12v-.008Z"
          class="stroke-amber-500 stroke-[2.5]"
        />
      </svg>
    </div>

    <!-- Text -->
    <div class="space-y-2">
      <h2 class="text-xl font-bold text-slate-800 dark:text-slate-100">
        {t("common.no_vps_title")}
      </h2>
      <p class="text-sm text-slate-500 dark:text-slate-400 leading-relaxed">
        {t("common.no_vps_desc")}
      </p>
    </div>

    <!-- Action Button -->
    <a
      href="/config"
      class="inline-flex items-center gap-2 px-6 py-3 rounded-xl bg-violet-600 hover:bg-violet-700 text-white font-bold text-sm transition-all shadow-md shadow-violet-500/20 active:scale-95 no-underline"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        class="w-4 h-4"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M5.25 14.25h13.5m-13.5 0a3 3 0 0 1-3-3m3 3a3 3 0 0 0 3 3h7.5a3 3 0 0 0 3-3m-13.5 0V6a3 3 0 0 1 3-3h7.5a3 3 0 0 1 3 3v5.25m-13.5 0h13.5m-13.5 0a3 3 0 0 0-3 3v3a3 3 0 0 0 3 3h13.5a3 3 0 0 0 3-3v-3a3 3 0 0 0-3-3"
        />
      </svg>
      <span>{t("common.go_to_servers")}</span>
    </a>
  </div>
</div>
