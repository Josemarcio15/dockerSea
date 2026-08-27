export interface ProfileFormState {
  id: string;
  name: string;
}
export interface ProfileStore {
  readonly form: ProfileFormState;
  readonly showModal: boolean;
  readonly showDeleteConfirm: boolean;
  readonly profileToDelete: { id: string; name: string } | null;
  openCreate(): void;
  openEdit(profile: any): void;
  select(id: string): Promise<void>;
  save(): Promise<void>;
  requestDelete(id: string, name: string): void;
  confirmDelete(): Promise<void>;
}
