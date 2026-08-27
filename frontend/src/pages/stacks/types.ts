import type { ServerCapabilities, StackItem } from "$lib/domains/stacks";

export type { ServerCapabilities, StackItem } from "$lib/domains/stacks";

export type StackSourceType = "editor" | "folder";

export interface StackEditorState {
  id: string;
  name: string;
  projectName: string;
  sourceType: StackSourceType;
  folderPath: string;
  yaml: string;
}

export interface StackPageState {
  stacks: StackItem[];
  capabilities: ServerCapabilities | null;
  loading: boolean;
  searchQuery: string;
}
