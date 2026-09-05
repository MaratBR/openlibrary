import PQueue from 'p-queue'
import { Context, Effect, Layer } from 'effect'
import { Font, FontsApi } from './api'
import { FontsLoaderConfig } from './config'

export class FontsLoader extends Context.Service<
  FontsLoader,
  {
    readonly fontLoaded: (font: string) => boolean
    readonly getFonts: () => ReadonlyArray<Readonly<Font>>
    readonly fetchFonts: () => Effect.Effect<ReadonlyArray<Readonly<Font>>, unknown>
    readonly addFonts: (fonts: ReadonlyArray<string>) => Effect.Effect<void, unknown>
    readonly ensureFontLoaded: (font: string) => Effect.Effect<void, unknown>
    readonly attachIframe: (iframe: HTMLIFrameElement) => Effect.Effect<void, unknown>
    readonly detachIframe: (iframe: HTMLIFrameElement) => Effect.Effect<void>
  }
>()('openlibrary/FontsLoader') {
  static readonly layer = Layer.effect(
    this,
    Effect.gen(function* () {
      const fontsApi = yield* FontsApi
      const config = yield* FontsLoaderConfig
      const cssQueue = new PQueue({ concurrency: 1 })
      const preloadedFonts = new Set(config.preloadedFonts)
      const loadedFonts = new Set<string>()
      const attachedIframes = new Set<HTMLIFrameElement>()
      let fonts: ReadonlyArray<Font> = []
      let loaded = false

      const fontLoaded = (font: string): boolean =>
        preloadedFonts.has(font) || loadedFonts.has(font)

      const getFonts = (): ReadonlyArray<Readonly<Font>> => fonts

      const fetchFonts = Effect.fn('FontsLoader.fetchFonts')(function* () {
        if (loaded) {
          yield* Effect.logDebug('Returning cached fonts').pipe(
            Effect.annotateLogs({ service: 'FontsLoader', count: fonts.length }),
          )
          return getFonts()
        }

        yield* Effect.logDebug('Fetching fonts').pipe(Effect.annotateLogs('service', 'FontsLoader'))
        fonts = yield* fontsApi.getGoogleFonts()
        loaded = true
        yield* Effect.logDebug('Cached fonts').pipe(
          Effect.annotateLogs({ service: 'FontsLoader', count: fonts.length }),
        )
        return getFonts()
      })

      const loadCss = (css: string, targetDocument: Document) =>
        Effect.tryPromise(() =>
          cssQueue.add(
            () =>
              new Promise<void>((resolve, reject) => {
                const link = targetDocument.createElement('link')
                link.href = css
                link.type = 'text/css'
                link.rel = 'stylesheet'
                link.media = 'screen,print'
                link.dataset.fontLoader = 'true'

                link.addEventListener('load', () => resolve())
                link.addEventListener('error', (event) => reject(event.error))
                targetDocument.head.appendChild(link)
              }),
            { timeout: 60000 },
          ),
        ).pipe(
          Effect.tapError((error) =>
            Effect.logError('Font stylesheet failed to load').pipe(
              Effect.annotateLogs({ service: 'FontsLoader', css, error }),
            ),
          ),
          Effect.asVoid,
        )

      const addFonts = Effect.fn('FontsLoader.addFonts')(function* (
        requestedFonts: ReadonlyArray<string>,
      ) {
        const missingFonts = [...new Set(requestedFonts)].filter((font) => !fontLoaded(font))
        if (missingFonts.length === 0) {
          yield* Effect.logDebug('Requested fonts are already loaded').pipe(
            Effect.annotateLogs({ service: 'FontsLoader', fonts: requestedFonts }),
          )
          return
        }

        if (preloadedFonts.size + loadedFonts.size + missingFonts.length > config.maxFontsCount) {
          return yield* Effect.fail(
            new Error(`FontsLoader cannot load more than ${config.maxFontsCount} fonts`),
          )
        }

        const url = `/_api/fonts/google/include?name=${encodeURIComponent(missingFonts.join(','))}`
        yield* Effect.logDebug('Loading font stylesheet').pipe(
          Effect.annotateLogs({ service: 'FontsLoader', fonts: missingFonts }),
        )
        const iframeDocuments = [...attachedIframes].flatMap((iframe) =>
          iframe.contentDocument === null ? [] : [iframe.contentDocument],
        )
        yield* Effect.all(
          [document, ...iframeDocuments].map((targetDocument) => loadCss(url, targetDocument)),
          { concurrency: 1 },
        )
        missingFonts.forEach((font) => loadedFonts.add(font))
        yield* Effect.logDebug('Loaded font stylesheet').pipe(
          Effect.annotateLogs({
            service: 'FontsLoader',
            fonts: missingFonts,
            documents: iframeDocuments.length + 1,
          }),
        )
      })

      const ensureFontLoaded = Effect.fn('FontsLoader.ensureFontLoaded')(function* (font: string) {
        yield* addFonts([font])
      })

      const attachIframe = Effect.fn('FontsLoader.attachIframe')(function* (
        iframe: HTMLIFrameElement,
      ) {
        const iframeDocument = iframe.contentDocument
        if (iframeDocument === null) {
          return yield* Effect.fail(new Error('Cannot attach an iframe without a content document'))
        }

        attachedIframes.add(iframe)
        const fontsToLoad = [...loadedFonts]
        if (fontsToLoad.length > 0) {
          const url = `/_api/fonts/google/include?name=${encodeURIComponent(fontsToLoad.join(','))}`
          yield* Effect.logDebug('Loading existing fonts into attached iframe').pipe(
            Effect.annotateLogs({ service: 'FontsLoader', fonts: fontsToLoad }),
          )
          yield* loadCss(url, iframeDocument)
        }
        yield* Effect.logDebug('Attached iframe').pipe(
          Effect.annotateLogs({ service: 'FontsLoader', loadedFonts: fontsToLoad.length }),
        )
      })

      const detachIframe = Effect.fn('FontsLoader.detachIframe')(function* (
        iframe: HTMLIFrameElement,
      ) {
        attachedIframes.delete(iframe)
        yield* Effect.logDebug('Detached iframe').pipe(
          Effect.annotateLogs('service', 'FontsLoader'),
        )
      })

      return FontsLoader.of({
        fontLoaded,
        getFonts,
        fetchFonts,
        addFonts,
        ensureFontLoaded,
        attachIframe,
        detachIframe,
      })
    }),
  )
}
