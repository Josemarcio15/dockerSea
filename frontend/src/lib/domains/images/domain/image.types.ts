export interface DockerImage {
  id: string;
  repo: string;
  tag: string;
  size: string;
  rawSizeBytes: number;
  created: number;
  containersUsing: string[] | null;
  virtualSize?: number;
  sharedSize?: number;
}

export interface ImageActionResult {
  success: boolean;
  message: string;
  count?: number;
  errors?: string[] | null;
}
