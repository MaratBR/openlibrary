export function translate(key: string) {
  if (window.i18n) {
    return window.i18n[key] ?? key
  }

  return key
}

declare global {
  interface Window {
    i18n?: Record<string, string>
    _: typeof translate
  }
}

window._ = translate
