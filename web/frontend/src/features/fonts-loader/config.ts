// oxlint-disable require-yield
import { Context, Effect, Layer } from 'effect'

export class FontsLoaderConfig extends Context.Service<
  FontsLoaderConfig,
  {
    readonly preloadedFonts: string[]
    readonly maxFontsCount: number
  }
>()('openlibrary/FontsLoaderConfig') {
  static readonly layer = Layer.effect(
    this,
    Effect.gen(function* () {
      // TODO Load the fonts configuration asynchronously when a remote source is available.
      return FontsLoaderConfig.of({ preloadedFonts: [], maxFontsCount: 200 })
    }),
  )
}
