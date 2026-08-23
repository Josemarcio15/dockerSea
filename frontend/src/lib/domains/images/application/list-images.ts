import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type { DockerImage } from "../domain/image.types.js";
import { imageWailsApi } from "../infrastructure/image.wails.js";

export async function listImages(server: VpsServer): Promise<DockerImage[]> {
  return await imageWailsApi.listImages(server);
}
