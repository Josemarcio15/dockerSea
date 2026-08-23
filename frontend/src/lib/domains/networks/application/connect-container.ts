import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type { NetworkActionResult } from "../domain/network.types.js";
import { networkWailsApi } from "../infrastructure/network.wails.js";

export async function connectContainer(
  server: VpsServer,
  networkName: string,
  containerName: string,
): Promise<NetworkActionResult> {
  if (!networkName?.trim() || !containerName?.trim()) {
    return {
      success: false,
      message: "Rede e Container são obrigatórios",
    };
  }

  return await networkWailsApi.connectContainer(
    server,
    networkName.trim(),
    containerName.trim(),
  );
}

export async function disconnectContainer(
  server: VpsServer,
  networkName: string,
  containerName: string,
): Promise<NetworkActionResult> {
  if (!networkName?.trim() || !containerName?.trim()) {
    return {
      success: false,
      message: "Rede e Container são obrigatórios",
    };
  }

  return await networkWailsApi.disconnectContainer(
    server,
    networkName.trim(),
    containerName.trim(),
  );
}
