import { notifySuccess } from "$shared/stores/notification.svelte";
import * as api from "./api";
import { activeServer } from "./service";
import type { ServersStore } from "./types";

export function createServersStore(getData: () => any): ServersStore {
  let usageCache = $state<Record<string, any>>({});
  let loadingUsage = $state<Record<string, boolean>>({});

  async function fetchUsage(server: any, force = false) {
    if (!server?.id || (!force && usageCache[server.id])) return;
    loadingUsage[server.id] = true;
    try {
      const usage = await api.getSystemUsage(server);
      if (usage) usageCache[server.id] = usage;
    } catch (error) {
      console.warn("Não foi possível obter recursos do servidor:", error);
    } finally {
      loadingUsage[server.id] = false;
    }
  }

  async function activate(server: any) {
    const data = getData();
    try {
      await api.setActiveServer(server.id);
      data.activeVps = server;
      data.servers = (await api.listServers()) || data.servers;
      notifySuccess(`Servidor '${server.name}' ativado com sucesso!`);
    } catch {
      data.activeVps = server;
      notifySuccess(`Servidor '${server.name}' selecionado.`);
    }
    fetchUsage(server);
  }

  async function load() {
    const data = getData();
    try {
      const servers = await api.listServers();
      if (!servers) return;
      data.servers = servers;
      const active = activeServer(servers);
      if (active) {
        data.activeVps = active;
        fetchUsage(active);
      }
    } catch (error) {
      console.warn("Erro ao carregar servidores do SQLite:", error);
    }
  }

  return {
    get usageCache() {
      return usageCache;
    },
    get loadingUsage() {
      return loadingUsage;
    },
    load,
    fetchUsage,
    activate,
  };
}
