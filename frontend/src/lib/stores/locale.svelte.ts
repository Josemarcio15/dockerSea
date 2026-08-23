import ptBR from "$lib/locales/pt-BR.json";
import enUS from "$lib/locales/en-US.json";

const translations: Record<string, any> = {
  "pt-BR": ptBR,
  "en-US": enUS,
};

// Reactive locale state — Svelte 5 $state rune via .svelte.ts extension
let currentLocale = $state("pt-BR");

export function setLocale(locale: string) {
  if (translations[locale]) {
    currentLocale = locale;
  } else {
    currentLocale = "en-US";
  }
}

export function getLocale(): string {
  return currentLocale;
}

/**
 * Translate a dot-notation key using the current reactive locale.
 * Components that call t() in templates will re-render automatically
 * when setLocale() changes currentLocale (because it's a $state rune).
 */
export function t(
  key: string,
  params?: Record<string, string | number>,
): string {
  const parts = key.split(".");
  let value: any = translations[currentLocale] || translations["en-US"];

  for (const part of parts) {
    if (value && typeof value === "object" && part in value) {
      value = value[part];
    } else {
      // Fallback to English
      let fallbackValue: any = translations["en-US"];
      for (const fbPart of parts) {
        if (
          fallbackValue &&
          typeof fallbackValue === "object" &&
          fbPart in fallbackValue
        ) {
          fallbackValue = fallbackValue[fbPart];
        } else {
          fallbackValue = key;
          break;
        }
      }
      value = fallbackValue;
      break;
    }
  }

  if (typeof value !== "string") {
    return key;
  }

  if (params) {
    let result = value;
    for (const [k, v] of Object.entries(params)) {
      result = result.replace(new RegExp(`{${k}}`, "g"), String(v));
    }
    return result;
  }

  return value;
}
