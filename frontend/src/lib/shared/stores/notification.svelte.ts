/**
 * Global notification banner store.
 *
 * Use `notifySuccess(msg)`, `notifyWarning(msg)`, `notifyError(msg)` em
 * any page or component. The banner automatically dismisses after 5s.
 *
 * @example
 * ```svelte
 * import { notifySuccess, notifyWarning, notifyError } from "$shared/stores/notification.svelte";
 *
 * notifySuccess("VPS selected successfully");
 * notifyWarning("No image selected");
 * notifyError("Error deleting images");
 *
 * // For caught errors (e.message):
 * notifyError(e?.message ?? "Unknown error");
 * ```
 */

type NotificationType = "success" | "warning" | "error";

let _message = $state("");
let _type = $state<NotificationType | null>(null);
let _timeout: ReturnType<typeof setTimeout> | null = null;

function _clear() {
  _message = "";
  _type = null;
}

function _show(msg: string, type: NotificationType) {
  if (_timeout) clearTimeout(_timeout);
  _message = msg;
  _type = type;
  _timeout = setTimeout(_clear, 5000);
}

export function notifySuccess(msg: string) {
  _show(msg, "success");
}

export function notifyWarning(msg: string) {
  _show(msg, "warning");
}

export function notifyError(msg: string) {
  _show(msg, "error");
}

export function getNotification() {
  return { message: _message, type: _type };
}

export function dismissNotification() {
  _clear();
}
