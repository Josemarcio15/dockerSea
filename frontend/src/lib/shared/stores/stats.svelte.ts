export type ContainerStats = {
  BlockIO: string;
  CPUPerc: string;
  Container: string;
  ID: string;
  IPv4Address?: string;
  MemPerc: string;
  MemUsage: string;
  Name: string;
  NetIO: string;
  PIDs: string;
};

let eventSource: EventSource | null = null;
let activeVpsId: string | null = null;
let subscribers = 0;

export const statsState = $state<{
  stats: ContainerStats[];
  loading: boolean;
  online: boolean;
  checking: boolean;
  error: string | null;
}>({
  stats: [],
  loading: true,
  online: true,
  checking: false,
  error: null,
});

export function reconnectStats() {
  if (activeVpsId) {
    statsState.checking = true;
    startStream(activeVpsId);
  }
}

export function subscribeToStats(vpsId: string): () => void {
  subscribers++;

  if (activeVpsId !== vpsId) {
    startStream(vpsId);
  }

  return () => {
    subscribers--;
    if (subscribers <= 0) {
      stopStream();
      subscribers = 0;
    }
  };
}

function startStream(vpsId: string) {
  if (eventSource) {
    eventSource.close();
  }

  activeVpsId = vpsId;
  statsState.loading = true;
  statsState.checking = true;
  statsState.error = null;

  eventSource = new EventSource(`/api/stats?vpsId=${vpsId}`);

  eventSource.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      statsState.stats = data;
      statsState.loading = false;
      statsState.checking = false;
      statsState.online = true;
      statsState.error = null;
    } catch (e) {
      console.error("Erro ao fazer parse dos stats do docker", e);
    }
  };

  eventSource.onerror = () => {
    statsState.loading = false;
    statsState.checking = false;
    // In the desktop Wails application there is an HTTP /api/stats route;
    // the absence of this route does not mean the VPS is offline.
    if (typeof window !== "undefined" && !["http:", "https:"].includes(window.location.protocol)) {
      statsState.online = true;
      statsState.error = null;
      return;
    }
    statsState.online = false;
    statsState.error = "Conexão perdida";
  };
}

function stopStream() {
  if (eventSource) {
    eventSource.close();
    eventSource = null;
  }
  activeVpsId = null;
  statsState.stats = [];
  statsState.loading = true;
  statsState.checking = false;
  statsState.online = true;
  statsState.error = null;
}
