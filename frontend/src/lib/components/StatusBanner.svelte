<script lang="ts">
  import {
    getNotification,
    dismissNotification,
  } from "$lib/stores/notification.svelte";

  let notif = $derived(getNotification());

  let styleClasses = $derived(
    notif.type === "error"
      ? "bg-rose-950/90 text-rose-200 border-rose-500/30 shadow-rose-950/50"
      : notif.type === "warning"
        ? "bg-amber-950/90 text-amber-200 border-amber-500/30 shadow-amber-950/50"
        : "bg-emerald-950/90 text-emerald-200 border-emerald-500/30 shadow-emerald-950/50",
  );

  let iconText = $derived(
    notif.type === "error"
      ? "✗"
      : notif.type === "warning"
        ? "!"
        : "✓",
  );

  let iconBg = $derived(
    notif.type === "error"
      ? "bg-rose-500/20 text-rose-400 border-rose-500/30"
      : notif.type === "warning"
        ? "bg-amber-500/20 text-amber-400 border-amber-500/30"
        : "bg-emerald-500/20 text-emerald-400 border-emerald-500/30",
  );
</script>

{#if notif.message}
  <!-- Toast flutuante com z-index alto e backdrop-blur para sobrepor a tela sem empurrar layout -->
  <div
    class="fixed top-6 right-6 z-9999 max-w-md w-full pointer-events-auto animate-slideIn"
  >
    <div
      class="flex items-center gap-3 p-4 rounded-2xl border backdrop-blur-md {styleClasses} text-sm shadow-2xl transition-all"
    >
      <div
        class="w-7 h-7 rounded-xl flex items-center justify-center shrink-0 border font-bold text-xs {iconBg}"
      >
        {iconText}
      </div>

      <span class="flex-1 whitespace-pre-line text-xs font-semibold leading-relaxed">
        {notif.message}
      </span>

      <button
        type="button"
        class="w-6 h-6 rounded-lg bg-white/5 hover:bg-white/10 flex items-center justify-center cursor-pointer text-xs opacity-70 hover:opacity-100 transition-opacity border-none text-inherit"
        onclick={dismissNotification}
      >
        ✕
      </button>
    </div>
  </div>
{/if}
