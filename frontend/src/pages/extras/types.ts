export type ExtrasMainTab = "nginx" | "ports" | "deploy_temp";
export type NginxSiteTab = "available" | "enabled";
export type NginxBusyAction = "delete" | "enable" | "test" | "restart" | "save";

export interface NginxSites {
  available?: string[];
  enabled?: string[];
}
