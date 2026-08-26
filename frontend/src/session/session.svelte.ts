import * as ConfigService from "../../bindings/go-walis/internal/config/configservice.js";
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
  const servers = await ConfigService.ListServers();
  if (servers) {
    session.servers = servers;
    session.activeVps = servers.find((server: any) => server.isActive) || null;
  }

  const profiles = await ConfigService.ListProfiles();
  if (profiles && profiles.length > 0) {
    session.profiles = profiles;
    session.activeProfile =
      profiles.find((profile: any) => profile.isActive) || profiles[0];
  }
}
