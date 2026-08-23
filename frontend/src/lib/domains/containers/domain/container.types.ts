export interface PortMapping {
  ip?: string;
  privatePort: number;
  publicPort?: number;
  type: string; // "tcp" | "udp"
}

export interface NetworkEndpoint {
  networkId?: string;
  ipAddress: string;
  gateway: string;
  macAddress?: string;
}

export interface MountInfo {
  type: string; // "volume" | "bind" | "tmpfs"
  name: string;
  source: string;
  destination: string;
  readOnly: boolean;
}

export interface Container {
  id: string;
  names: string[] | null;
  name: string;
  image: string;
  imageId: string;
  command: string;
  created: number;
  state: string; // "running" | "exited" | "paused" | "created"
  status: string; // "Up 2 hours", "Exited (0) 5 minutes ago"
  ports: string;
  portList: PortMapping[] | null;
  networks: { [key: string]: NetworkEndpoint | undefined } | null;
  mounts: MountInfo[] | null;
  restartPolicy: string;
  labels: { [key: string]: string | undefined } | null;
  env?: string[] | null;
}

export interface ContainerActionResult {
  success: boolean;
  message: string;
  errors?: string[] | null;
}

export type ContainerActionType = "start" | "stop" | "restart" | "rm";
