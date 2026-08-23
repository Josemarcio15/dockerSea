import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type {
  CreateNetworkParams,
  NetworkActionResult,
} from "../domain/network.types.js";
import { isValidNetworkName } from "../domain/network.rules.js";
import { networkWailsApi } from "../infrastructure/network.wails.js";

export async function createNetwork(
  server: VpsServer,
  params: CreateNetworkParams,
): Promise<NetworkActionResult> {
  if (!isValidNetworkName(params.name)) {
    return {
      success: false,
      message: "Nome de rede inválido",
    };
  }

  return await networkWailsApi.createNetwork(server, {
    name: params.name.trim(),
    driver: params.driver || "bridge",
    subnet: params.subnet?.trim() || "",
    gateway: params.gateway?.trim() || "",
    labels: params.labels || null,
  });
}
