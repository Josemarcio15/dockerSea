export interface DockerVolume {
  name: string;
  driver: string;
  mountpoint: string;
  createdAt: string;
  labels: { [key: string]: string | undefined } | null;
  scope: string;
  inUse: boolean;
  containers: (string[] | null)[] | null;
  size?: string;
}

export interface VolumeActionResult {
  success: boolean;
  message: string;
  count?: number;
  errors?: string[] | null;
}

export interface CreateVolumeParams {
  name: string;
  driver?: string;
  labels?: { [key: string]: string | undefined } | null;
}
