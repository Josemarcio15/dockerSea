import type { Container } from "./types";
export function filterContainers(
  containers: Container[],
  query: string,
): Container[] {
  const normalized = query.trim().toLowerCase();
  if (!normalized) return containers;
  return containers.filter((container) =>
    [container.name, container.image, container.status, container.state].some(
      (value) => value?.toLowerCase().includes(normalized),
    ),
  );
}
