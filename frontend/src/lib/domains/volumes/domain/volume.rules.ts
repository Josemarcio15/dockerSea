import type { DockerVolume } from "./volume.types.js";

export function isValidVolumeName(name: string): boolean {
  if (!name || name.trim().length === 0) return false;
  const validRegex = /^[a-zA-Z0-9][a-zA-Z0-9_.-]*$/;
  return validRegex.test(name.trim());
}

export function filterVolumes(
  volumes: DockerVolume[],
  searchQuery: string = "",
): DockerVolume[] {
  const query = searchQuery.trim().toLowerCase();
  if (!query) return volumes || [];

  return (volumes || []).filter(
    (v) =>
      (v.name && v.name.toLowerCase().includes(query)) ||
      (v.driver && v.driver.toLowerCase().includes(query)),
  );
}
