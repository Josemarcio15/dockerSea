export const EVENTS = {
  container: "docker:container:event",
  imagePullProgress: "docker:image:pull:progress",
  imagePullComplete: "docker:image:pull:complete",
  imageTransferProgress: "docker:image:transfer:progress",
  imageTransferComplete: "docker:image:transfer:complete",
} as const;
