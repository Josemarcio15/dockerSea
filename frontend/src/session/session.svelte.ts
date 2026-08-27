import * as ProfileService from "../../bindings/go-walis/internal/profiles/service.js";
import * as ServerService from "../../bindings/go-walis/internal/servers/service.js";
import type { SessionState } from "./session.types";

export const session = $state<SessionState>({
  servers: [],
  activeVps: null,
  profiles: [],
  activeProfile: {
    name: "Perfil Padrão",
    locale: "pt-BR",
  },
});

export async function loadSession(): Promise<void> {
  const servers = await ServerService.ListServers();
  if (servers) {
    session.servers = servers;
    session.activeVps = servers.find((server: any) => server.isActive) || null;
  }

  const profiles = await ProfileService.ListProfiles();
  if (profiles && profiles.length > 0) {
    session.profiles = profiles;
    session.activeProfile =
      profiles.find((profile: any) => profile.isActive) || profiles[0];
  }
}
