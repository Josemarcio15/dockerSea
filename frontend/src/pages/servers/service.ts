export function activeServer(servers: any[]): any | undefined {
  return servers.find((server) => server.isActive);
}
