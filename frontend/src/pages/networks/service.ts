export function normalizeNetworkNames(names: string | string[]): string[] {
  return (Array.isArray(names) ? names : names.split(","))
    .map((name) => name.trim())
    .filter(Boolean);
}

export function createNetworkPayload(
  name: string,
  driver: string,
  subnet: string,
  gateway: string,
) {
  return {
    name: name.trim(),
    driver,
    subnet: subnet.trim(),
    gateway: gateway.trim(),
  };
}
