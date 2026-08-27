import type { VpsFormData } from "./types";

export function serverPayload(form: VpsFormData): any {
  return {
    ...form,
    port: parseInt(form.port) || 22,
    isActive: false,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  };
}

export function dockerDetectionPayload(form: VpsFormData): any {
  return serverPayload(form);
}
