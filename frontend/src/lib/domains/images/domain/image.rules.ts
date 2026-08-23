import type { DockerImage } from "./image.types.js";

export function filterImages(
  images: DockerImage[],
  searchQuery: string = "",
): DockerImage[] {
  const query = searchQuery.trim().toLowerCase();
  if (!query) return images || [];

  return (images || []).filter(
    (img) =>
      (img.repo && img.repo.toLowerCase().includes(query)) ||
      (img.tag && img.tag.toLowerCase().includes(query)) ||
      (img.id && img.id.toLowerCase().includes(query)),
  );
}

export function formatImageDisplayName(image: DockerImage): string {
  if (image.repo === "<none>") {
    return `<untagged> (${image.id.substring(0, 12)})`;
  }
  return `${image.repo}:${image.tag}`;
}
