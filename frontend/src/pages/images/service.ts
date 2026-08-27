import type { DockerImage } from "./types";

export function imageKey(image: DockerImage): string {
  return image.repo && image.repo !== "<none>"
    ? `${image.repo}:${image.tag || "latest"}`
    : image.id;
}
