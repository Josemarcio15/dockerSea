import { filterVolumes } from "$lib/domains/volumes";
import type { DockerVolume } from "$lib/domains/volumes";

export { filterVolumes };

export function filterVolumeList(
  volumes: DockerVolume[],
  query: string,
): DockerVolume[] {
  return filterVolumes(volumes, query);
}
