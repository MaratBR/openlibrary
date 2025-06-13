const i18nShowKeys = new URLSearchParams(window.location.search).has('i18n.show-keys')

export function translate(key: string, args?: Record<string, string>) {
  if (i18nShowKeys) {
    return key
  }

  if (!window.i18n) {
    return key
  }

  let translation = window.i18n[key]

  if (!translation) {
    return key
  }

  if (args) {
    // replace {{key}} with args[key]
    for (const [key, value] of Object.entries(args)) {
      translation = translation.replace(`{{${key}}}`, value)
      translation = translation.replace(`{${key}}`, value)
      translation = translation.replace(`{{.${key}}}`, value)
    }
  }

  return translation
}

declare global {
  interface Window {
    i18n?: Record<string, string>
    _: typeof translate
  }
}

window._ = translate
