<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLInputAttributes } from "svelte/elements";

  let {
    id,
    label,
    type = "text",
    value = $bindable(""),
    placeholder = "",
    help = "",
    error = "",
    disabled = false,
    required = false,
    autocomplete = "off",
    class: customClass = "",
    trailing,
  }: {
    id?: string;
    label?: string;
    type?: string;
    value?: any;
    placeholder?: string;
    help?: string;
    error?: string;
    disabled?: boolean;
    required?: boolean;
    autocomplete?: HTMLInputAttributes["autocomplete"];
    class?: string;
    trailing?: Snippet;
  } = $props();

  let generatedId = $derived(id || (label ? "input-" + label.toLowerCase().replace(/\s+/g, "-") : undefined));
</script>

<div class="flex flex-col gap-1.5 {customClass}">
  {#if label}
    <label
      for={generatedId}
      class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >
      {label}
      {#if required}<span class="text-rose-500">*</span>{/if}
    </label>
  {/if}

  <div class="flex gap-2 items-center">
    <input
      id={generatedId}
      {type}
      {placeholder}
      {disabled}
      {required}
      {autocomplete}
      bind:value
      class="w-full px-3.5 py-2.5 text-sm border rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all font-semibold disabled:opacity-50 disabled:cursor-not-allowed {error
        ? 'border-rose-500/80 focus:border-rose-500 focus:ring-rose-500/20'
        : 'border-slate-200 dark:border-slate-800'}"
    />
    {#if trailing}
      {@render trailing()}
    {/if}
  </div>

  {#if error}
    <span class="text-[11px] text-rose-500 font-medium">{error}</span>
  {:else if help}
    <span class="text-[11px] text-slate-400 dark:text-slate-500">{help}</span>
  {/if}
</div>
