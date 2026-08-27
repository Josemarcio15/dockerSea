import * as ExtraBinding from "../../../bindings/go-walis/internal/extras/extraservice.js";

export function listNginxSites(activeVps: any): Promise<any> {
  return ExtraBinding.ListNginxSites(activeVps);
}

export function readNginxSite(
  activeVps: any,
  site: string,
  tab: string,
): Promise<string> {
  return ExtraBinding.ReadNginxSite(activeVps, site, tab);
}

export function deleteNginxSite(
  activeVps: any,
  site: string,
  tab: string,
): Promise<any> {
  return ExtraBinding.DeleteNginxSite(activeVps, site, tab);
}

export function saveNginxSite(
  activeVps: any,
  site: string,
  content: string,
): Promise<any> {
  return ExtraBinding.SaveNginxSite(activeVps, site, content);
}

export function enableNginxSite(activeVps: any, site: string): Promise<any> {
  return ExtraBinding.EnableNginxSite(activeVps, site);
}

export function testNginxConfig(activeVps: any): Promise<any> {
  return ExtraBinding.TestNginxConfig(activeVps);
}

export function restartNginx(activeVps: any): Promise<any> {
  return ExtraBinding.RestartNginx(activeVps);
}

export function listDeployTempFiles(activeVps: any): Promise<any> {
  return ExtraBinding.ListDeployTempFiles(activeVps);
}

export function listDeployTempFilesAt(activeVps: any, path: string): Promise<any> {
  return ExtraBinding.ListDeployTempFilesAt(activeVps, path);
}

export function cleanDeployTempFiles(activeVps: any): Promise<any> {
  return ExtraBinding.CleanDeployTempFiles(activeVps);
}

export function deleteDeployTempPath(activeVps: any, path: string): Promise<any> {
  return ExtraBinding.DeleteDeployTempPath(activeVps, path);
}
