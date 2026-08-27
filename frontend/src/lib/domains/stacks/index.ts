export * from "./domain/stack.types.js";
export * from "./infrastructure/stack.wails.js";
export * from "./application/stack-actions.js";
import { stackWailsApi } from "./infrastructure/stack.wails.js";

export const browseStackFolder = (path: string) => stackWailsApi.browse(path);
