<script lang="ts">
  import {
    getNotification,
    dismissNotification,
  } from "$lib/stores/notification.svelte";

  let notif = $derived(getNotification());

  let bgClass = $derived(
    notif.type === "error"
      ? "bg-red-50 dark:bg-red-950/30 text-red-700 dark:text-red-400 border-red-200 dark:border-red-900/50"
      : notif.type === "warning"
        ? "bg-amber-50 dark:bg-amber-950/30 text-amber-700 dark:text-amber-400 border-amber-200 dark:border-amber-900/50"
        : "bg-green-50 dark:bg-emerald-950/30 text-green-700 dark:text-emerald-400 border-green-200 dark:border-emerald-900/50",
  );
</script>

{#if notif.message}
  <div
    class="flex items-center justify-between p-4 rounded-xl border mb-4 {bgClass} text-sm transition-all duration-300 shadow-md"
  >
    <span class="whitespace-pre-line">{notif.message}</span>
    <button
      class="bg-transparent border-none cursor-pointer text-sm font-bold opacity-70 hover:opacity-100 ml-4 text-inherit"
      onclick={dismissNotification}
    >
      ✕
    </button>
  </div>
{/if}
