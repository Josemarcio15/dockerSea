import * as BuilderBinding from "../../../bindings/go-walis/internal/builder/service.js";
import type { BrowseResult, SavedPath } from "./types";

export function browse(path: string): Promise<BrowseResult> {
  return BuilderBinding.Browse(path) as Promise<BrowseResult>;
}

export function listSavedPaths(): Promise<SavedPath[]> {
  return BuilderBinding.ListSavedPaths() as Promise<SavedPath[]>;
}

export function savePath(path: string, label: string): Promise<void> {
  return BuilderBinding.SavePath(path, label);
}

export function removeSavedPath(path: string): Promise<void> {
  return BuilderBinding.RemoveSavedPath(path);
}

export function build(path: string, tag: string, locale: string): Promise<void> {
  return BuilderBinding.Build(path, tag, locale);
}