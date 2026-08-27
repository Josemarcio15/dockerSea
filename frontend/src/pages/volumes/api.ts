export {
  listVolumes,
  createVolume,
  removeVolumes,
  pruneVolumes,
} from "$lib/domains/volumes";

export type { DockerVolume } from "$lib/domains/volumes";
export type { VpsServer } from "../../../bindings/go-walis/internal/core/db/models.js";
