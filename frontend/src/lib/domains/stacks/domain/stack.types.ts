export interface StackItem {
  id: string;
  name: string;
  projectName: string;
  sourceType: "editor" | "folder";
  folderPath: string;
  yamlContent: string;
  profileId: string;
  lastDeployedRemoteDir?: string;
  lastDeployedAt?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface StackActionResult {
  success: boolean;
  message: string;
  logs?: string;
}

export interface StackProgressEvent {
  stackId: string;
  deployId: string;
  phase: string;
  message: string;
  success?: boolean;
}

export interface ServerCapabilities {
  dockerAvailable: boolean;
  composeAvailable: boolean;
  dockerVersion: string;
  composeVersion: string;
  buildxAvailable: boolean;
  architecture: string;
  operatingSystem: string;
}
