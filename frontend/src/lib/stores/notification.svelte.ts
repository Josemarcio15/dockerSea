/**
 * Global notification banner store.
 *
 * Use `notifySuccess(msg)`, `notifyWarning(msg)`, `notifyError(msg)` em
 * qualquer página ou componente. O banner desaparece automaticamente após 5s.
 *
 * @example
 * ```svelte
 * import { notifySuccess, notifyWarning, notifyError } from "$lib/stores/notification.svelte";
 *
 * notifySuccess("VPS selecionada com sucesso");
 * notifyWarning("Nenhuma imagem selecionada");
 * notifyError("Erro ao deletar imagens");
 *
 * // Para erros capturados (e.message):
 * notifyError(e?.message ?? "Erro desconhecido");
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
