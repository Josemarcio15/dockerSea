import type { DockerImage } from "$lib/domains/images";

export type { DockerImage };
export type ImageTab = "my_images" | "download" | "transfer";
export type ImageFilter = "all" | "in_use" | "unused" | "dangling";
