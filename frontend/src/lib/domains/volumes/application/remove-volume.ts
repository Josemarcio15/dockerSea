import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type { VolumeActionResult } from "../domain/volume.types.js";
import { volumeWailsApi } from "../infrastructure/volume.wails.js";

export async function removeVolumes(
  server: VpsServer,
  names: string[],
): Promise<VolumeActionResult> {
  const filtered = names
    .map((n) => n.trim())
    .filter((n) => n.length > 0);

  if (filtered.length === 0) {
    return {
      success: false,
      message: "Nenhum volume selecionado",
    };
  }

  return await volumeWailsApi.deleteVolumes(server, filtered);
}
