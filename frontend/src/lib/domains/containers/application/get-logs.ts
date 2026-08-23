import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import { containerWailsApi } from "../infrastructure/container.wails.js";

export async function getContainerLogs(
  server: VpsServer,
  containerName: string,
  tail: number = 200,
): Promise<string> {
  return await containerWailsApi.getLogs(server, containerName, tail);
}
