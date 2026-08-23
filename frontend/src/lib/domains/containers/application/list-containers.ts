import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type { Container } from "../domain/container.types.js";
import { containerWailsApi } from "../infrastructure/container.wails.js";

export async function listContainers(
  server: VpsServer,
  all: boolean = true,
): Promise<Container[]> {
  return await containerWailsApi.listContainers(server, all);
}
