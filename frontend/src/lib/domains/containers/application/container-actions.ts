import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type {
  ContainerActionResult,
  ContainerActionType,
} from "../domain/container.types.js";
import { containerWailsApi } from "../infrastructure/container.wails.js";

export async function executeContainerAction(
  server: VpsServer,
  actionType: ContainerActionType,
  containerNames: string[],
): Promise<ContainerActionResult> {
  const names = containerNames.map((n) => n.trim()).filter((n) => n.length > 0);

  if (names.length === 0) {
    return {
      success: false,
      message: "Nenhum container selecionado",
    };
  }

  return await containerWailsApi.executeAction(server, actionType, names);
}

export async function startContainers(
  server: VpsServer,
  containerNames: string[],
): Promise<ContainerActionResult> {
  return executeContainerAction(server, "start", containerNames);
}

export async function stopContainers(
  server: VpsServer,
  containerNames: string[],
): Promise<ContainerActionResult> {
  return executeContainerAction(server, "stop", containerNames);
}

export async function restartContainers(
  server: VpsServer,
  containerNames: string[],
): Promise<ContainerActionResult> {
  return executeContainerAction(server, "restart", containerNames);
}

export async function removeContainers(
  server: VpsServer,
  containerNames: string[],
): Promise<ContainerActionResult> {
  return executeContainerAction(server, "rm", containerNames);
}
