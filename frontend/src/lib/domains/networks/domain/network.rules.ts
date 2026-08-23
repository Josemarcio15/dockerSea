/**
 * Redes padrão do Docker que possuem restrições especiais
 */
export const DEFAULT_DOCKER_NETWORKS = ["bridge", "host", "none"] as const;

/**
 * Verifica se uma rede é considerada padrão/do sistema
 */
export function isDefaultNetwork(name: string): boolean {
  return DEFAULT_DOCKER_NETWORKS.includes(name as any);
}

/**
 * Validação simples de formato de nome de rede Docker
 */
export function isValidNetworkName(name: string): boolean {
  if (!name || name.trim().length === 0) return false;
  // Docker network names can contain alphanumeric, underscores, dots, and hyphens
  const validRegex = /^[a-zA-Z0-9][a-zA-Z0-9_.-]*$/;
  return validRegex.test(name.trim());
}

/**
 * Filtra redes customizadas (exclui redes padrão e aplica busca)
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
