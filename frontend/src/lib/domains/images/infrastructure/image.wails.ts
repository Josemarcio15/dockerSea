import * as ImageService from "../../../../../bindings/go-walis/internal/images/imageservice.js";
import type {
  VpsServer,
  ImageHistoryItem,
} from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type { DockerImage, ImageActionResult } from "../domain/image.types.js";

export const imageWailsApi = {
  async listImages(server: VpsServer): Promise<DockerImage[]> {
    const list = await ImageService.ListImages(server);
    return (list as DockerImage[]) || [];
  },

  async deleteImages(
    server: VpsServer,
    imageIds: string[],
  ): Promise<ImageActionResult> {
    return await ImageService.DeleteImages(server, imageIds);
  },

  async pullImage(
    server: VpsServer,
    imageName: string,
    profileId: string,
  ): Promise<ImageActionResult> {
    return await ImageService.PullImage(server, imageName, profileId);
  },

  async listHistory(profileId: string): Promise<ImageHistoryItem[]> {
    const hist = await ImageService.ListHistory(profileId);
    return hist || [];
  },

  async deleteHistory(ids: string[]): Promise<void> {
    return await ImageService.DeleteHistory(ids);
  },

  async clearHistory(profileId: string): Promise<void> {
    return await ImageService.ClearHistory(profileId);
  },

  async transferImages(
    srcServer: VpsServer,
    dstServer: VpsServer,
    imageIds: string[],
  ): Promise<ImageActionResult> {
    return await ImageService.TransferImages(srcServer, dstServer, imageIds);
  },
};
