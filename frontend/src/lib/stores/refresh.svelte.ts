/**
 * Global reactive refresh counter.
 *
 * Use `triggerRefresh()` after any action that modifies Docker data.
 * Pages watching `useRefreshKey()` via `$effect` will re-fetch automatically.
 *
 * @example
 * ```svelte
 * import { useRefreshKey, triggerRefresh } from "$lib/stores/refresh.svelte";
 *
 * $effect(() => {
 *   useRefreshKey(); // creates reactive dependency
 *   if (data?.activeVps) fetchData();
 * });
 *
 * async function doAction() {
 *   await fetch("?/action", ...);
 *   triggerRefresh();
 * }
 * ```
 */

let _refreshKey = $state(0);

/**
 * Increments the refresh counter. Call this after any server-side mutation
 * (pull, delete, create, prune, etc.) to signal pages to re-fetch their data.
 */
export function triggerRefresh(): void {
  _refreshKey++;
}

/**
 * Returns the current refresh key. Use inside `$effect(() => useRefreshKey() && ...)`
 * to reactively re-run side effects whenever `triggerRefresh()` is called.
 */
export function useRefreshKey(): number {
  return _refreshKey;
}
