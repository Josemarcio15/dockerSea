<script lang="ts">
  import { onMount } from "svelte";
  import { EditorView, basicSetup } from "codemirror";
  import { EditorState } from "@codemirror/state";
  import { StreamLanguage } from "@codemirror/language";
  import { oneDark } from "@codemirror/theme-one-dark";

  import { yaml } from "@codemirror/lang-yaml";

  type EditorMode = "nginx" | "json" | "yaml" | "text";
  let {
    value = "",
    mode = "text",
    minHeight = "400px",
    maxHeight = "560px",
    onchange = (_value: string) => {},
  } = $props<{
    value: string;
    mode?: EditorMode;
    minHeight?: string;
    maxHeight?: string;
    onchange?: (value: string) => void;
  }>();

  let container = $state<HTMLDivElement | null>(null);
  let view = $state<EditorView | null>(null);

  const nginxMode = StreamLanguage.define({
    startState: () => ({}),
    token(stream) {
      if (stream.sol() && stream.match(/^\s*#.*$/)) return "comment";
      if (stream.eatSpace()) return null;
      if (stream.match(/^#.*/)) return "comment";
      if (stream.match(/^\$[a-zA-Z_][\w]*/)) return "variableName";
      if (stream.match(/^\b\d+(?:\.\d+)?(?:[kKmMgG]|ms|s)?\b/)) return "number";
      if (stream.match(/^"(?:[^"\\]|\\.)*"|^'(?:[^'\\]|\\.)*'/)) return "string";
      if (stream.match(/^(?:server|location|upstream|events|http|stream|if)\b/)) return "keyword";
      if (stream.match(/^(?:listen|server_name|root|index|proxy_pass|include|return|rewrite|ssl_certificate|ssl_certificate_key|access_log|error_log|client_max_body_size|deny|allow|proxy_set_header)\b/)) return "propertyName";
      if (stream.match(/^(?:on|off|default|auto|localhost|none|all)\b/)) return "bool";
      stream.next();
      return null;
    },
  });

  const jsonMode = StreamLanguage.define({
    startState: () => ({}),
    token(stream) {
      if (stream.eatSpace()) return null;
      if (stream.match(/^"(?:[^"\\]|\\.)*"\s*:/)) return "propertyName";
      if (stream.match(/^"(?:[^"\\]|\\.)*"/)) return "string";
      if (stream.match(/^(?:true|false|null)\b/)) return "bool";
      if (stream.match(/^-?\d+(?:\.\d+)?/)) return "number";
      stream.next();
      return null;
    },
  });

  onMount(() => {
    if (!container) return;
    view = new EditorView({
      state: EditorState.create({
        doc: value,
        extensions: [
          basicSetup,
          ...(mode === "yaml"
            ? [yaml()]
            : mode === "nginx"
              ? [nginxMode]
              : mode === "json"
                ? [jsonMode]
                : []),
          oneDark,
          EditorView.theme({
            "&": { minHeight, maxHeight },
            ".cm-scroller": {
              overflow: "auto",
              fontFamily: "ui-monospace, SFMono-Regular, Menlo, monospace",
            },
            ".cm-content": { minHeight, padding: "14px 0" },
          }),
          EditorView.updateListener.of((update) => {
            if (update.docChanged) onchange(update.state.doc.toString());
          }),
        ],
      }),
      parent: container,
    });
    return () => view?.destroy();
  });

  $effect(() => {
    if (!view) return;
    const current = view.state.doc.toString();
    if (current !== value) view.dispatch({ changes: { from: 0, to: current.length, insert: value } });
  });
</script>

<div bind:this={container} class="w-full overflow-hidden rounded-xl border border-slate-200 dark:border-slate-800 focus-within:ring-2 focus-within:ring-violet-500/20"></div>
