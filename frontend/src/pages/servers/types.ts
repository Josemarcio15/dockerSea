export interface ServersStore {
  readonly usageCache: Record<string, any>;
  readonly loadingUsage: Record<string, boolean>;
  load(): Promise<void>;
  fetchUsage(server: any, force?: boolean): Promise<void>;
  activate(server: any): Promise<void>;
}
