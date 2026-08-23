import * as ContainerService from "../../../../../bindings/go-walis/internal/containers/containerservice.js";
import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type {
  Container,
  ContainerActionResult,
  ContainerActionType,
} from "../domain/container.types.js";
import { Events } from "@wailsio/runtime";

export const containerWailsApi = {
  async listContainers(server: VpsServer, all: boolean = true): Promise<Container[]> {
    const list = await ContainerService.ListContainers(server, all);
    return (list as Container[]) || [];
  },

  async executeAction(
    server: VpsServer,
    actionType: ContainerActionType,
    containerNames: string[],
  ): Promise<ContainerActionResult> {
    return await ContainerService.ExecuteAction(server, actionType, containerNames);
  },

  async getLogs(server: VpsServer, containerName: string, tail: number = 200): Promise<string> {
    return await ContainerService.GetLogs(server, containerName, tail);
  },

  async startEventsStream(server: VpsServer): Promise<void> {
    return await ContainerService.StartEventsStream(server);
  },

  async stopEventsStream(serverId: string): Promise<void> {
    return await ContainerService.StopEventsStream(serverId);
  },

  subscribeToEvents(callback: (event: any) => void): () => void {
    return Events.On("docker:container:event", callback);
  },
};
