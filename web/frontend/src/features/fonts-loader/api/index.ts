import { HttpClient } from '@/features/http-client'
import { Context, Effect, Layer, Schema } from 'effect'

export const Font = Schema.Struct({
  name: Schema.String,
  source: Schema.String,
  license: Schema.String,
  includeCss: Schema.String,
  externalLink: Schema.String,
})

export type Font = Schema.Schema.Type<typeof Font>

export class FontsApi extends Context.Service<
  FontsApi,
  {
    readonly getGoogleFonts: () => Effect.Effect<ReadonlyArray<Font>, unknown>
  }
>()('openlibrary/FontsApi') {
  static readonly layer = Layer.effect(
    this,
    Effect.gen(function* () {
      const httpClient = yield* HttpClient
      const getGoogleFonts = Effect.fn('FontsApi.getGoogleFonts')(function* () {
        yield* Effect.logDebug('Requesting Google fonts').pipe(
          Effect.annotateLogs('service', 'FontsApi'),
        )
        const response = yield* Effect.tryPromise(() =>
          httpClient.get('/_api/fonts/google').then((r) => r.json()),
        )
        const fonts = yield* Schema.decodeUnknownEffect(Schema.Array(Font))(response)
        yield* Effect.logDebug('Received Google fonts').pipe(
          Effect.annotateLogs({ service: 'FontsApi', count: fonts.length }),
        )
        return fonts
      })
      return FontsApi.of({ getGoogleFonts })
    }),
  )
}
