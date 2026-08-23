import * as NetworkService from "../../../../../bindings/go-walis/internal/networks/networkservice.js";
import type { VpsServer } from "../../../../../bindings/go-walis/internal/core/db/models.js";
import type { Container } from "../../../../../bindings/go-walis/internal/containers/models.js";
import type {
  DockerNetwork,
  NetworkActionResult,
  CreateNetworkParams,
} from "../domain/network.types.js";

export const networkWailsApi = {
  async listNetworks(server: VpsServer): Promise<DockerNetwork[]> {
    const res = await NetworkService.ListNetworks(server);
    return res || [];
  },

  async listContainers(server: VpsServer): Promise<Container[]> {
    const res = await NetworkService.ListNetworkContainers(server);
    return res || [];
  },

  async createNetwork(
    server: VpsServer,
    params: CreateNetworkParams,
  ): Promise<NetworkActionResult> {
    return await NetworkService.CreateNetwork(server, {
      name: params.name,
      driver: params.driver,
      subnet: params.subnet,
      gateway: params.gateway,
      labels: params.labels,
    });
  },

  async deleteNetworks(
    server: VpsServer,
    names: string[],
  ): Promise<NetworkActionResult> {
    return await NetworkService.DeleteNetworks(server, names);
  },

  async connectContainer(
    server: VpsServer,
    networkName: string,
    containerName: string,
  ): Promise<NetworkActionResult> {
    return await NetworkService.ConnectContainer(
      server,
      networkName,
      containerName,
    );
  },

  async disconnectContainer(
    server: VpsServer,
    networkName: string,
    containerName: string,
  ): Promise<NetworkActionResult> {
    return await NetworkService.DisconnectContainer(
      server,
      networkName,
      containerName,
    );
  },

  async pruneNetworks(server: VpsServer): Promise<NetworkActionResult> {
    return await NetworkService.PruneNetworks(server);
  },
};
