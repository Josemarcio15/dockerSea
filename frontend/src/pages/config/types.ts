export interface VpsFormData {
  id: string;
  name: string;
  connectionType: "local" | "ssh";
  host: string;
  port: string;
  username: string;
  authType: "key" | "password";
  sshKeyPath: string;
  sshKeyPassphrase: string;
  sshPassword: string;
  sudoPassword: string;
  dockerSocketPath: string;
  dockerPath: string;
  dockerComposePath: string;
}

export interface VpsServer extends Omit<VpsFormData, "port"> {
  port: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface DiagnosticResult {
  success: boolean;
  message: string;
  steps: Array<{
    name: string;
    status: "success" | "error" | "warning";
    message: string;
  }>;
}
