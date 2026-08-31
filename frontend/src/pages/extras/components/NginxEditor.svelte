<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import CodeEditor from "$shared/components/CodeEditor.svelte";

  let {
    site = $bindable(),
    content = $bindable(),
    tab,
    editorKey,
    loading = false,
  }: {
    site: string;
    content: string;
    tab: string;
    editorKey: number;
    loading?: boolean;
  } = $props();
</script>

<div class="space-y-4">
  <div class="space-y-1.5">
    <label
      for="site"
      class="text-[11px] font-bold uppercase tracking-wider text-slate-500 dark:text-slate-400"
      >{t("extras.site_file_label")}</label
    >
    <input
      id="site"
      bind:value={site}
      placeholder={t("extras.site_file_placeholder")}
      class="w-full px-3.5 py-2 text-xs rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-[#0c101b] text-slate-900 dark:text-white font-mono focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all"
    />
    <p class="text-[11px] text-slate-400 font-mono">
      {tab === "available"
        ? "/etc/nginx/sites-available"
        : "/etc/nginx/sites-enabled"}/{site || "..."}
    </p>
  </div>
  <div class="space-y-1.5">
    <label
      for="nginx-editor"
      class="text-[11px] font-bold uppercase tracking-wider text-slate-500 dark:text-slate-400"
      >{t("extras.site_content_label")}</label
    >
    {#if loading}
      <div
        class="h-[320px] w-full rounded-2xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-[#0c101b]/50 flex flex-col items-center justify-center gap-3 text-slate-400"
      >
        <div class="w-6 h-6 border-2 border-violet-500 border-t-transparent rounded-full animate-spin"></div>
        <span class="text-xs font-mono">Carregando {site || "arquivo"}...</span>
      </div>
    {:else}
      {#key editorKey}
        <CodeEditor
          value={content}
          mode="nginx"
          onchange={(value) => (content = value)}
        />
      {/key}
    {/if}
  </div>
</div>
