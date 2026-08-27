export type Route =
  | "servers"
  | "containers"
  | "images"
  | "volumes"
  | "networks"
  | "stacks"
  | "builder"
  | "config"
  | "extras"
  | "profiles";

export const ROUTES: readonly Route[] = [
  "servers",
  "containers",
  "images",
  "volumes",
  "networks",
  "stacks",
  "builder",
  "config",
  "extras",
  "profiles",
];
