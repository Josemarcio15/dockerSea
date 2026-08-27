import type { VpsServer } from "../../../bindings/go-walis/internal/core/db/models.js";
import type { ImageHistoryItem } from "../../../bindings/go-walis/internal/core/db/models.js";
import { imageWailsApi } from "$lib/domains/images";
import * as configWailsApi from "../../../bindings/go-walis/internal/config/configservice.js";

export const api = imageWailsApi;
export type { VpsServer, ImageHistoryItem };
export function list(server: VpsServer) {
  return api.listImages(server);
}

export function pull(server: VpsServer, imageName: string, profileId: string) {
  return api.pullImage(server, imageName, profileId);
}

export function remove(server: VpsServer, imageIds: string[]) {
  return api.deleteImages(server, imageIds);
}

export function listHistory(profileId: string) {
  return api.listHistory(profileId);
}

export function deleteHistory(ids: string[]) {
  return api.deleteHistory(ids);
}

export function clearHistory(profileId: string) {
  return api.clearHistory(profileId);
}

export function transferImages(
  source: VpsServer,
  destination: VpsServer,
  imageIds: string[],
) {
  return api.transferImages(source, destination, imageIds);
}

export function listContainerConfigs(profileId: string) {
  return (configWailsApi as any).ListContainerConfigs(profileId);
}

export function saveContainerConfig(config: Record<string, unknown>) {
  return (configWailsApi as any).SaveContainerConfig(config);
}

export function deleteContainerConfig(id: string) {
  return (configWailsApi as any).DeleteContainerConfig(id);
}
