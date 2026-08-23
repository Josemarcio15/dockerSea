import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type { NetworkActionResult } from "../domain/network.types.js";
import { networkWailsApi } from "../infrastructure/network.wails.js";

export async function pruneNetworks(
  server: VpsServer,
): Promise<NetworkActionResult> {
  return await networkWailsApi.pruneNetworks(server);
}
