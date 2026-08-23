import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type {
  CreateVolumeParams,
  VolumeActionResult,
} from "../domain/volume.types.js";
import { isValidVolumeName } from "../domain/volume.rules.js";
import { volumeWailsApi } from "../infrastructure/volume.wails.js";

export async function createVolume(
  server: VpsServer,
  params: CreateVolumeParams,
): Promise<VolumeActionResult> {
  if (!isValidVolumeName(params.name)) {
    return {
      success: false,
      message: "Nome de volume inválido",
    };
  }

  return await volumeWailsApi.createVolume(server, {
    name: params.name.trim(),
    driver: params.driver || "local",
    labels: params.labels || null,
  });
}
