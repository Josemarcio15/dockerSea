import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type { DockerVolume } from "../domain/volume.types.js";
import { volumeWailsApi } from "../infrastructure/volume.wails.js";

export async function listVolumes(server: VpsServer): Promise<DockerVolume[]> {
  return await volumeWailsApi.listVolumes(server);
}

export async function getVolumeSize(
  server: VpsServer,
  name: string,
): Promise<string> {
  return await volumeWailsApi.getVolumeSize(server, name);
}
