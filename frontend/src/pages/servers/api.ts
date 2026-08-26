import * as ServerBinding from "../../../bindings/go-walis/internal/dashboard/service.js";

export function listServers(): Promise<any[]> {
  return ServerBinding.ListServers() as Promise<any[]>;
}

export function setActiveServer(id: string): Promise<void> {
  return ServerBinding.SetActiveServer(id);
}

export function getSystemUsage(server: any): Promise<any> {
  return ServerBinding.GetSystemUsage(server) as Promise<any>;
}