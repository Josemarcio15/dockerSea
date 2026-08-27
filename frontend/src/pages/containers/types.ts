import type {
  Container,
  ContainerActionResult,
  ContainerActionType,
} from "$lib/domains/containers";
export type { Container, ContainerActionResult, ContainerActionType };
export interface ContainersStore {
  readonly selectedNames: string[];
  searchQuery: string;
  readonly containers: Container[];
  readonly loading: boolean;
  readonly fetchError: string;
  readonly filteredContainers: Container[];
  showLogs: boolean;
  readonly logsTitle: string;
  readonly logsLoading: boolean;
  readonly logsContent: string[];
  fetchAll(silent?: boolean): Promise<void>;
  toggleAll(): void;
  handleToggleSelect(name: string): void;
  doActionSelected(action: ContainerActionType): Promise<void>;
  openLogs(containerName: string): Promise<void>;
  setupEventsStream(): () => void;
}
