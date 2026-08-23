import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type { Container } from "../../../../../bindings/go-walis/internal/containers/models.js";
import type { DockerNetwork } from "../domain/network.types.js";
import { networkWailsApi } from "../infrastructure/network.wails.js";

export async function listNetworks(
  server: VpsServer,
): Promise<{ networks: DockerNetwork[]; containers: Container[] }> {
  const [networks, containers] = await Promise.all([
    networkWailsApi.listNetworks(server),
    networkWailsApi.listContainers(server),
  ]);
  return { networks, containers };
}
