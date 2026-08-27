import type { VpsServer } from "./types";
import type { Container, DockerNetwork } from "./types";
import {
  listNetworks,
  createNetwork,
  removeNetworks,
  connectContainer,
  disconnectContainer,
  pruneNetworks,
} from "$lib/domains/networks";

export {
  listNetworks,
  createNetwork,
  removeNetworks,
  connectContainer,
  disconnectContainer,
  pruneNetworks,
};

export async function list(
  server: VpsServer,
): Promise<{ networks: DockerNetwork[]; containers: Container[] }> {
  return listNetworks(server);
}
