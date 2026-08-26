export function folderNameFromPath(path: string): string {
  return path.split("/").pop()?.split("\\").pop() || "";
}

export function sanitizeTag(input: string): string {
  return input
    .toLowerCase()
    .replace(/[^a-z0-9._-]/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-|-$/g, "");
}

export function tagFromPath(path: string): string {
  return sanitizeTag(folderNameFromPath(path));
}