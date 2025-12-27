import { JSX } from 'preact/jsx-runtime'

const i18nShowKeys = new URLSearchParams(window.location.search).has('i18n.show-keys')

type TranslateArg = string | JSX.Element
type TranslateArgs = Record<string, TranslateArg>

export function translate(key: string, args?: TranslateArgs): string | JSX.Element {
  if (i18nShowKeys) return key
  if (!window.i18n) return key

  let translation = window.i18n[key]
  if (!translation) return key

  if (!args) return translation

  const hasJSX = Object.values(args).some((v) => typeof v !== 'string')

  // Fast path: strings only (original behavior)
  if (!hasJSX) {
    for (const [argKey, value] of Object.entries(args)) {
      const v = value as string
      translation = translation.replace(`{{${argKey}}}`, v)
      translation = translation.replace(`{${argKey}}`, v)
      translation = translation.replace(`{{.${argKey}}}`, v)
    }
    return translation
  }

  // JSX path
  const parts: Array<string | JSX.Element> = [translation]

  for (const [argKey, value] of Object.entries(args)) {
    const patterns = [`{{${argKey}}}`, `{${argKey}}`, `{{.${argKey}}}`]

    for (const pattern of patterns) {
      for (let i = 0; i < parts.length; i++) {
        const part = parts[i]
        if (typeof part !== 'string') continue

        const split = part.split(pattern)
        if (split.length === 1) continue

        const next: Array<string | JSX.Element> = []
        split.forEach((s, idx) => {
          next.push(s)
          if (idx < split.length - 1) {
            next.push(value)
          }
        })

        parts.splice(i, 1, ...next)
        i += next.length - 1
      }
    }
  }

  return <>{parts}</>
}

declare global {
  interface Window {
    i18n?: Record<string, string>
    _: typeof translate
  }
}

window._ = translate
