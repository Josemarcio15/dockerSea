export interface DockerNetworkContainer {
  name: string;
  endpointId: string;
  macAddress: string;
  ip: string;
  ipv6Address?: string;
}

export interface DockerNetwork {
  id: string;
  name: string;
  driver: string;
  scope: string;
  subnet: string;
  gateway: string;
  internal: boolean;
  attachable: boolean;
  containers: DockerNetworkContainer[] | null;
  labels: { [key: string]: string | undefined } | null;
}

export interface NetworkActionResult {
  success: boolean;
  message: string;
  count?: number;
  errors?: string[] | null;
}

export interface CreateNetworkParams {
  name: string;
  driver?: string;
  subnet?: string;
  gateway?: string;
  labels?: { [key: string]: string | undefined } | null;
}
