<script lang="ts">
  import { tick } from "svelte";

  let {
    logs = [],
    id = "terminal-logs-view",
    maxHeight = "max-h-[55vh]",
    class: customClass = "",
  }: {
    logs: string[];
    id?: string;
    maxHeight?: string;
    class?: string;
  } = $props();

  let container = $state<HTMLDivElement | null>(null);

  $effect(() => {
    if (logs.length > 0) {
      tick().then(() => {
        if (container) {
          container.scrollTop = container.scrollHeight;
        }
      });
    }
  });

  const ANSI_COLOR_MAP: Record<number, string> = {
    // Normal Foreground
    30: "text-slate-900 dark:text-slate-100",
    31: "text-red-600 dark:text-red-400",
    32: "text-emerald-600 dark:text-emerald-400",
    33: "text-amber-600 dark:text-amber-300",
    34: "text-blue-600 dark:text-blue-400",
    35: "text-purple-600 dark:text-purple-400",
    36: "text-cyan-600 dark:text-cyan-400",
    37: "text-slate-700 dark:text-slate-200",
    // Bright Foreground
    90: "text-slate-500 dark:text-slate-400",
    91: "text-rose-600 dark:text-rose-400",
    92: "text-green-600 dark:text-green-300",
    93: "text-yellow-600 dark:text-yellow-200",
    94: "text-sky-600 dark:text-sky-300",
    95: "text-fuchsia-600 dark:text-fuchsia-300",
    96: "text-cyan-600 dark:text-cyan-300",
    97: "text-slate-900 dark:text-white",
  };

  function escapeHtml(text: string): string {
    return text
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#039;");
  }

  function parseAnsi(text: string): string {
    if (!text) return "";

    // Check if line has any ANSI sequences
    const hasAnsi = /(?:\u001b|\x1b|\[)\[?[0-9;]*m/.test(text);

    if (!hasAnsi) {
      const clean = escapeHtml(text.replace(/[\u0000-\u001F\u007F-\u009F]/g, ""));
      return `<span class="${getFallbackColor(text)}">${clean}</span>`;
    }

    let result = "";
    let lastIndex = 0;
    let currentColorClass = "";
    let isBold = false;
    let match: RegExpExecArray | null;

    // Standard ANSI regex
    const standardAnsi = /\x1b\[([0-9;]*)m/g;
    let textToParse = text;

    // If text contains stripped escapes like `[36m` directly without \x1b
    if (!text.includes("\x1b") && !text.includes("\u001b")) {
      textToParse = text.replace(/\[([0-9]{1,2}(?:;[0-9]{1,2})*)m/g, "\x1b[$1m");
    }

    while ((match = standardAnsi.exec(textToParse)) !== null) {
      const chunk = textToParse.slice(lastIndex, match.index);
      if (chunk) {
        const classes = [currentColorClass, isBold ? "font-bold" : ""].filter(Boolean).join(" ");
        result += classes ? `<span class="${classes}">${escapeHtml(chunk)}</span>` : escapeHtml(chunk);
      }

      const codes = (match[1] || "0").split(";").map((c) => parseInt(c, 10) || 0);

      for (const code of codes) {
        if (code === 0) {
          currentColorClass = "";
          isBold = false;
        } else if (code === 1) {
          isBold = true;
        } else if (ANSI_COLOR_MAP[code]) {
          currentColorClass = ANSI_COLOR_MAP[code];
        }
      }

      lastIndex = standardAnsi.lastIndex;
    }

    const remaining = textToParse.slice(lastIndex);
    if (remaining) {
      const classes = [currentColorClass, isBold ? "font-bold" : ""].filter(Boolean).join(" ");
      result += classes ? `<span class="${classes}">${escapeHtml(remaining)}</span>` : escapeHtml(remaining);
    }

    // Clean up any stray unprintable characters
    let html = result.replace(/[\u0000-\u0008\u000B\u000C\u000E-\u001F]/g, "");

    // Realce especial em rosa para SQL e queries SELECT / INSERT / UPDATE / DELETE / JOIN
    html = highlightSql(html);

    return html;
  }

  function highlightSql(html: string): string {
    // Se a linha tiver [SQL...] ou começar com keywords SQL, realça a tag e os comandos SQL em rosa/pink
    html = html.replace(
      /(\[SQL(?::[^\]]+)?\])/gi,
      '<span class="text-pink-600 dark:text-pink-400 font-bold">$1</span>',
    );

    // Realce de queries SQL: SELECT, INSERT INTO, UPDATE, DELETE, FROM, WHERE, LEFT JOIN, INNER JOIN, ORDER BY, GROUP BY, LIMIT
    html = html.replace(
      /\b(SELECT|INSERT\s+INTO|UPDATE|DELETE\s+FROM|FROM|WHERE|LEFT\s+JOIN|RIGHT\s+JOIN|INNER\s+JOIN|JOIN|ORDER\s+BY|GROUP\s+BY|LIMIT|OFFSET|COALESCE|COUNT|SUM|AVG|CASE|WHEN|THEN|ELSE|END)\b/g,
      '<span class="text-pink-600 dark:text-pink-400 font-bold">$1</span>',
    );

    // Realce para tags HTTP (GET, POST, PUT, DELETE, PATCH)
    html = html.replace(
      /(\[HTTP\])/gi,
      '<span class="text-purple-600 dark:text-purple-400 font-bold">$1</span>',
    );

    return html;
  }

  function getFallbackColor(line: string): string {
    const upper = line.toUpperCase();
    if (
      upper.includes("FATAL") ||
      upper.includes("ERROR") ||
      upper.includes("ERRO") ||
      upper.includes("FAIL") ||
      upper.includes("PANIC") ||
      upper.includes("✗") ||
      upper.includes("❌")
    ) {
      return "text-red-600 dark:text-red-400 font-semibold";
    }
    if (
      upper.includes("WARN") ||
      upper.includes("WARNING") ||
      upper.includes("DETAIL:")
    ) {
      return "text-amber-600 dark:text-amber-300";
    }
    if (
      upper.includes("SUCCESS") ||
      upper.includes("SUCESSO") ||
      upper.includes("✓") ||
      upper.includes("COMPLETE") ||
      upper.includes("CONSTRUÍDA COM SUCESSO") ||
      upper.includes("READY TO ACCEPT CONNECTIONS") ||
      upper.includes("OK")
    ) {
      return "text-emerald-600 dark:text-emerald-400 font-medium";
    }
    if (upper.includes("SELECT") || upper.includes("[SQL")) {
      return "text-pink-600 dark:text-pink-400";
    }
    if (
      upper.includes("STEP") ||
      upper.includes("BUILDING") ||
      upper.includes("COMPILING") ||
      upper.includes("DOWNLOADING") ||
      upper.includes("UPLOADING") ||
      upper.includes("PREPARING") ||
      upper.includes("LOG:") ||
      upper.includes("INFO")
    ) {
      return "text-sky-600 dark:text-sky-300";
    }
    return "text-slate-700 dark:text-slate-300";
  }
</script>

<div
  bind:this={container}
  {id}
  class="bg-slate-100/90 dark:bg-[#090d16] text-slate-800 dark:text-slate-200 p-4 rounded-xl text-xs font-mono overflow-auto flex-1 shadow-inner border border-slate-300/80 dark:border-slate-800/80 space-y-1 select-text scrollbar-thin {maxHeight} {customClass}"
>
  {#if logs.length === 0}
    <div class="text-slate-400 dark:text-slate-500 italic py-6 text-center select-none">
      (Nenhum log disponível)
    </div>
  {:else}
    {#each logs as log, index}
      <div
        class="flex items-start gap-3 leading-relaxed break-all font-mono hover:bg-slate-200/60 dark:hover:bg-slate-800/30 px-1.5 py-0.5 rounded transition-colors"
      >
        <span
          class="text-slate-400 dark:text-slate-500 text-[11px] shrink-0 select-none text-right w-7 font-mono pr-1 border-r border-slate-300 dark:border-slate-800/60"
        >
          {index + 1}
        </span>
        <!-- eslint-disable-next-line svelte/no-at-html-tags -->
        <span class="flex-1">{@html parseAnsi(log)}</span>
      </div>
    {/each}
  {/if}
</div>
