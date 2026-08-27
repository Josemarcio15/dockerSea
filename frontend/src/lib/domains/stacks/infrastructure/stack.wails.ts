import * as StackService from "../../../../../bindings/go-walis/internal/stacks/stackservice.js";
import type {
  StackItem,
  StackActionResult,
  StackProgressEvent,
  ServerCapabilities,
} from "../domain/stack.types.js";
import { Events } from "@wailsio/runtime";
import { Browse as browseBuilder } from "../../../../../bindings/go-walis/internal/builder/service.js";

export const stackWailsApi = {
  browse(path: string) {
    return browseBuilder(path);
  },
  async listStacks(profileId: string = "default"): Promise<StackItem[]> {
    const list = await (StackService as any).ListStacks(profileId);
    return (list as StackItem[]) || [];
  },

  async getStack(stackId: string): Promise<StackItem | null> {
    return await (StackService as any).GetStack(stackId);
  },

  async saveStack(item: StackItem): Promise<void> {
    return await (StackService as any).SaveStack(item as any);
  },

  async deleteStackDefinition(id: string): Promise<void> {
    if ((StackService as any).DeleteStackDefinition) {
      return await (StackService as any).DeleteStackDefinition(id);
    }
    return await (StackService as any).DeleteStack(id);
  },

  async deleteStack(id: string): Promise<void> {
    return await this.deleteStackDefinition(id);
  },

  async deployStack(
    profileId: string,
    stackId: string,
  ): Promise<StackActionResult> {
    return await (StackService as any).DeployStack(profileId, stackId);
  },

  async stopStack(
    profileId: string,
    stackId: string,
  ): Promise<StackActionResult> {
    return await (StackService as any).StopStack(profileId, stackId);
  },

  async removeStackRemote(
    profileId: string,
    stackId: string,
    deleteVolumes: boolean = false,
  ): Promise<StackActionResult> {
    return await (StackService as any).RemoveStackRemote(
      profileId,
      stackId,
      deleteVolumes,
    );
  },

  async getLogs(
    profileId: string,
    stackId: string,
    tail: number = 200,
  ): Promise<string> {
    return await (StackService as any).GetStackLogs(profileId, stackId, tail);
  },

  async getServerCapabilities(
    profileId: string = "default",
  ): Promise<ServerCapabilities> {
    return await (StackService as any).GetServerCapabilities(profileId);
  },

  subscribeToDeployStarted(callback: (event: any) => void): () => void {
    return Events.On("stacks:deploy:started", (ev) => callback(ev.data));
  },

  subscribeToDeployProgress(
    callback: (event: StackProgressEvent) => void,
  ): () => void {
    return Events.On("stacks:deploy:progress", (ev) =>
      callback(ev.data as StackProgressEvent),
    );
  },

  subscribeToDeployComplete(
    callback: (event: StackProgressEvent) => void,
  ): () => void {
    return Events.On("stacks:deploy:complete", (ev) =>
      callback(ev.data as StackProgressEvent),
    );
  },

  subscribeToDeployFailed(
    callback: (event: StackProgressEvent) => void,
  ): () => void {
    return Events.On("stacks:deploy:failed", (ev) =>
      callback(ev.data as StackProgressEvent),
    );
  },
};
