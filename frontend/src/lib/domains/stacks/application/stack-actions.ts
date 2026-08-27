import { stackWailsApi } from "../infrastructure/stack.wails.js";
import type {
  StackItem,
  ServerCapabilities,
  StackActionResult,
} from "../domain/stack.types.js";

export async function listStacksUseCase(
  profileId: string = "default",
): Promise<StackItem[]> {
  return await stackWailsApi.listStacks(profileId);
}

export async function getServerCapabilitiesUseCase(
  profileId: string = "default",
): Promise<ServerCapabilities | null> {
  try {
    return await stackWailsApi.getServerCapabilities(profileId);
  } catch {
    return null;
  }
}

export async function saveStackUseCase(stack: StackItem): Promise<void> {
  const nowIso = new Date().toISOString();
  await stackWailsApi.saveStack({
    ...stack,
    createdAt: stack.createdAt || nowIso,
    updatedAt: nowIso,
  });
}

export async function deleteStackDefinitionUseCase(
  stackId: string,
): Promise<void> {
  await stackWailsApi.deleteStackDefinition(stackId);
}

export async function removeStackRemoteUseCase(
  profileId: string,
  stackId: string,
  deleteVolumes: boolean = false,
): Promise<StackActionResult> {
  return await stackWailsApi.removeStackRemote(
    profileId,
    stackId,
    deleteVolumes,
  );
}

export async function deployStackUseCase(
  profileId: string,
  stackId: string,
): Promise<StackActionResult> {
  return await stackWailsApi.deployStack(profileId, stackId);
}

export async function stopStackUseCase(
  profileId: string,
  stackId: string,
): Promise<StackActionResult> {
  return await stackWailsApi.stopStack(profileId, stackId);
}

export async function getStackLogsUseCase(
  profileId: string,
  stackId: string,
  tail: number = 200,
): Promise<string> {
  return await stackWailsApi.getLogs(profileId, stackId, tail);
}
