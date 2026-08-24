import * as VolumeService from "../../../../../bindings/go-walis/internal/volumes/volumeservice.js";
import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type {
  DockerVolume,
  VolumeActionResult,
  CreateVolumeParams,
} from "../domain/volume.types.js";

export const volumeWailsApi = {
  async listVolumes(server: VpsServer): Promise<DockerVolume[]> {
    const list = await VolumeService.ListVolumes(server);
    return (list as DockerVolume[]) || [];
  },

  async createVolume(
    server: VpsServer,
    params: CreateVolumeParams,
  ): Promise<VolumeActionResult> {
    return await VolumeService.CreateVolume(server, {
      name: params.name,
      driver: params.driver,
      labels: params.labels || undefined,
    });
  },

  async deleteVolumes(
    server: VpsServer,
    names: string[],
  ): Promise<VolumeActionResult> {
    return await VolumeService.DeleteVolumes(server, names);
  },

  async getVolumeSize(server: VpsServer, name: string): Promise<string> {
    return await VolumeService.GetVolumeSize(server, name);
  },

  async pruneVolumes(server: VpsServer): Promise<VolumeActionResult> {
    return await VolumeService.PruneVolumes(server);
  },
};
