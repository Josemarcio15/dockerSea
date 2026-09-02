/**
 * SWR (Stale-While-Revalidate) Cache Store
 * Allows instant navigation (0ms) with background synchronization.
 */

type CacheEntry<T> = {
  data: T;
  timestamp: number;
};

const cache = new Map<string, CacheEntry<any>>();

export function getCached<T>(key: string): T | undefined {
  const entry = cache.get(key);
  return entry ? (entry.data as T) : undefined;
}

export function setCached<T>(key: string, data: T): void {
  cache.set(key, {
    data,
    timestamp: Date.now(),
  });
}

export function invalidateCache(keyPrefix?: string): void {
  if (!keyPrefix) {
    cache.clear();
    return;
  }
  for (const key of cache.keys()) {
    if (key.startsWith(keyPrefix)) {
      cache.delete(key);
    }
  }
}
