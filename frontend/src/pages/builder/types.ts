export type BuilderStatus = "idle" | "ready" | "building" | "success" | "error";

export interface BuilderFolder {
  name: string;
  path: string;
}

export interface BrowseResult {
  currentPath: string;
  parentPath?: string | null;
  folders?: BuilderFolder[];
  hasDockerfile: boolean;
  hasDockerignore: boolean;
}

export interface SavedPath {
  path: string;
}

export interface BuildResult {
  success?: boolean;
  image?: string;
  message?: string;
}

export interface BuilderStore {
  readonly currentPath: string;
  readonly parentPath: string | null;
  readonly folders: BuilderFolder[];
  readonly hasDockerfile: boolean;
  readonly hasDockerignore: boolean;
  readonly loading: boolean;
  readonly status: BuilderStatus;
  readonly logs: string[];
  readonly builtImage: string;
  readonly errorMsg: string;
  readonly savedPaths: string[];
  customTag: string;
  readonly folderName: string;
  readonly defaultTag: string;
  readonly effectiveTag: string;
  readonly canBuild: boolean;
  loadSavedPaths(): Promise<void>;
  browse(path?: string): Promise<void>;
  saveCurrentPath(): Promise<void>;
  removeSavedPath(path: string): Promise<void>;
  build(): Promise<void>;
  appendLog(line: string): void;
  completeBuild(result: BuildResult): void;
}