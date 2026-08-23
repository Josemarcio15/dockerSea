import type { Container } from "./container.types.js";

export function isContainerRunning(container: Container): boolean {
  return container.status.includes("Up") || container.state === "running";
}

export function isContainerStopped(container: Container): boolean {
  return (
    container.status.includes("Exited") ||
    container.state === "exited" ||
    container.state === "created"
  );
}

export function filterContainers(
  containers: Container[],
  searchQuery: string = "",
): Container[] {
  const query = searchQuery.trim().toLowerCase();
  if (!query) return containers || [];

  return (containers || []).filter(
    (c) =>
      (c.name && c.name.toLowerCase().includes(query)) ||
      (c.image && c.image.toLowerCase().includes(query)) ||
      (c.id && c.id.toLowerCase().includes(query)),
  );
}
