<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import CodeEditor from "$shared/components/CodeEditor.svelte";

  let {
    site = $bindable(),
    content = $bindable(),
    tab,
    editorKey,
  }: {
    site: string;
    content: string;
    tab: string;
    editorKey: number;
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
    {#key editorKey}<CodeEditor
        value={content}
        mode="nginx"
        onchange={(value) => (content = value)}
      />{/key}
  </div>
</div>
