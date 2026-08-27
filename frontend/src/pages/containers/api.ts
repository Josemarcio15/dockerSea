import type { VpsServer } from "../../../bindings/go-walis/internal/core/db/models.js";
import {
  executeContainerAction,
  getContainerLogs,
  listContainers,
  containerWailsApi,
} from "$lib/domains/containers";
import type { ContainerActionResult, ContainerActionType } from "./types";
export type { VpsServer } from "../../../bindings/go-walis/internal/core/db/models.js";
export const api = {
  list: (server: VpsServer) => listContainers(server, true),
  action: (
    server: VpsServer,
    action: ContainerActionType,
    names: string[],
  ): Promise<ContainerActionResult> =>
    executeContainerAction(server, action, names),
  logs: (server: VpsServer, name: string) =>
    getContainerLogs(server, name, 200),
  events: containerWailsApi,
};
