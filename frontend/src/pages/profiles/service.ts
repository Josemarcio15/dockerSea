export function deserialize(text: string): any {
  try {
    return JSON.parse(text);
  } catch {
    return { success: false, message: text };
  }
}

export function parseActionResult(result: any): any {
  if (!result || typeof result !== "object") return null;
  const rawData = result.data;
  if (!rawData) return result;
  try {
    return typeof rawData === "string" ? JSON.parse(rawData) : rawData;
  } catch {
    return rawData;
  }
}

export function actionError(result: any, parsed: any): string {
  return (
    parsed?.message ||
    parsed?.error ||
    result?.error?.message ||
    result?.error ||
    result?.message ||
    (result ? JSON.stringify(result).slice(0, 300) : "Resposta vazia")
  );
}
