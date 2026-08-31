import * as ConfigBinding from "../../../bindings/go-walis/internal/config/configservice.js";
import * as ServerBinding from "../../../bindings/go-walis/internal/servers/service.js";

export const listServers = () => ServerBinding.ListServers();
export const saveServer = (server: any) => ServerBinding.SaveServer(server);
export const deleteServer = (id: string) => ServerBinding.DeleteServer(id);
export const setActiveServer = (id: string) =>
  ServerBinding.SetActiveServer(id);
export const testConnection = (server: any) =>
  ServerBinding.TestConnection(server);
export const autoDetectDocker = (server: any) =>
  ConfigBinding.AutoDetectDocker(server);
export const getDatabaseInfo = () => ConfigBinding.GetDatabaseInfo();
export const exportDatabaseBackup = (destinationPath: string) =>
  ConfigBinding.ExportDatabaseBackup(destinationPath);
export const restoreDatabaseBackup = (sourcePath: string) =>
  ConfigBinding.RestoreDatabaseBackup(sourcePath);
export const resetDatabase = () => ConfigBinding.ResetDatabase();


