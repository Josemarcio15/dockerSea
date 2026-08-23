import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type { ImageActionResult } from "../domain/image.types.js";
import { imageWailsApi } from "../infrastructure/image.wails.js";

export async function removeImages(
  server: VpsServer,
  imageIds: string[],
): Promise<ImageActionResult> {
  const filtered = imageIds
    .map((id) => id.trim())
    .filter((id) => id.length > 0);

  if (filtered.length === 0) {
    return {
      success: false,
      message: "Nenhuma imagem selecionada",
    };
  }

  return await imageWailsApi.deleteImages(server, filtered);
}

export async function pullImage(
  server: VpsServer,
  imageName: string,
  profileId: string = "default",
): Promise<ImageActionResult> {
  if (!imageName?.trim()) {
    return {
      success: false,
      message: "Nome da imagem não informado",
    };
  }

  return await imageWailsApi.pullImage(server, imageName.trim(), profileId);
}
