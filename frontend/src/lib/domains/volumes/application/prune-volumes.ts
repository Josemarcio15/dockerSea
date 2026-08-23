import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type { VolumeActionResult } from "../domain/volume.types.js";
import { volumeWailsApi } from "../infrastructure/volume.wails.js";

export async function pruneVolumes(
  server: VpsServer,
): Promise<VolumeActionResult> {
  return await volumeWailsApi.pruneVolumes(server);
}
