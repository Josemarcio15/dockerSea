import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type { NetworkActionResult } from "../domain/network.types.js";
import { isDefaultNetwork } from "../domain/network.rules.js";
import { networkWailsApi } from "../infrastructure/network.wails.js";

export async function removeNetworks(
  server: VpsServer,
  names: string[],
): Promise<NetworkActionResult> {
  const filteredNames = names
    .map((n) => n.trim())
    .filter((n) => n.length > 0 && !isDefaultNetwork(n));

  if (filteredNames.length === 0) {
    return {
      success: false,
      message: "Nenhuma rede válida selecionada para exclusão",
    };
  }

  return await networkWailsApi.deleteNetworks(server, filteredNames);
}
