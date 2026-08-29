import PQueue from 'p-queue'
import { Font, FontsApi } from './api'
import { Effect } from 'effect'

export namespace FontsLoader {
  const cssQueue = new PQueue({ concurrency: 1 })

  let fonts: ReadonlyArray<Font> = []
  let loaded = false

  export function fontLoaded(font: string): boolean {
    try {
      return document.fonts.check(font)
    } catch {
      return false
    }
  }

  export function getFonts(): ReadonlyArray<Readonly<Font>> {
    return fonts
  }

  export function fetchFonts() {
    return Effect.gen(function* () {
      if (loaded) {
        return getFonts()
      }

      const fontsApi = yield* FontsApi
      const googleFonts = yield* fontsApi.getGoogleFonts()
      fonts = googleFonts

      return getFonts()
    })
  }

  export async function addFonts(fonts: string[]) {
    if (fonts.length === 0) return

    const missingFonts = fonts.filter((font) => !fontLoaded(font))
    if (missingFonts.length === 0) return

    const url = `/_api/fonts/google/include?name=${encodeURIComponent(missingFonts.join(','))}`

    await loadCss(url)
  }

  export async function ensureFontLoaded(font: string) {
    await addFonts([font])
  }

  function loadCss(css: string) {
    return cssQueue.add(
      () => {
        return new Promise<void>((resolve, reject) => {
          const $link = document.createElement('link')
          $link.href = css
          $link.type = 'text/css'
          $link.rel = 'stylesheet'
          $link.media = 'screen,print'
          $link.dataset.fontLoader = 'true'

          $link.addEventListener('load', () => {
            console.debug('[font-loader] css file is loaded', css)
            resolve()
          })
          $link.addEventListener('error', (ev) => {
            // TODO retry policy?
            console.error('[font-loader] css file failed to load', ev.error)
            reject(ev.error)
          })

          document.getElementsByTagName('head')[0].appendChild($link)
        })
      },
      { timeout: 60000 },
    )
  }
}
