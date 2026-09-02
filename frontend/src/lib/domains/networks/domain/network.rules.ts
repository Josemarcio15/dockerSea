/**
 * Default Docker networks that have special restrictions
 */
export const DEFAULT_DOCKER_NETWORKS = ["bridge", "host", "none"] as const;

/**
 * Checks if a network is considered default/system
 */
export function isDefaultNetwork(name: string): boolean {
  return DEFAULT_DOCKER_NETWORKS.includes(name as any);
}

/**
 * Simple Docker network name format validation
 */
export function isValidNetworkName(name: string): boolean {
  if (!name || name.trim().length === 0) return false;
  // Docker network names can contain alphanumeric, underscores, dots, and hyphens
  const validRegex = /^[a-zA-Z0-9][a-zA-Z0-9_.-]*$/;
  return validRegex.test(name.trim());
}

/**
 * Filters custom networks (excludes default networks and applies search)
 */
export function filterCustomNetworks<T extends { name: string }>(
  networks: T[],
  searchQuery: string = "",
): T[] {
  const query = searchQuery.trim().toLowerCase();
  return (networks || [])
    .filter((n) => !isDefaultNetwork(n.name))
    .filter((n) => (query ? n.name.toLowerCase().includes(query) : true));
}
