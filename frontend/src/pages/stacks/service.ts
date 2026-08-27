import type { StackItem } from "$lib/domains/stacks";

export const defaultYaml = `services:
  web:
    image: nginx:alpine
    ports:
      - "8080:80"
    restart: always
`;

export function filterStacks(stacks: StackItem[], query: string): StackItem[] {
  const normalizedQuery = query.toLowerCase();
  return stacks.filter(
    (stack) =>
      stack.name.toLowerCase().includes(normalizedQuery) ||
      (stack.projectName || "").toLowerCase().includes(normalizedQuery),
  );
}

export function folderName(path: string): string {
  return path.split("/").pop()?.split("\\").pop() || "";
}
