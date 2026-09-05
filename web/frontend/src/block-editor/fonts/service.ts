import { Context, Effect, Layer } from 'effect'
import { Font } from '@/features/fonts-loader/api'
import { FontsLoader } from '@/features/fonts-loader/loader'
import { EditorFontsApi } from './api'

export type EditorFontsState = {
  readonly fonts: ReadonlyArray<Readonly<Font>>
  readonly favoriteFonts: ReadonlyArray<string>
  readonly initializing: boolean
  readonly error: unknown
}

const initialState: EditorFontsState = {
  fonts: [],
  favoriteFonts: [],
  initializing: false,
  error: undefined,
}

export class EditorFonts extends Context.Service<
  EditorFonts,
  {
    readonly getState: () => EditorFontsState
    readonly subscribe: (listener: () => void) => () => void
    readonly initialize: () => Effect.Effect<void, unknown>
    readonly setFavorite: (font: string, favorite: boolean) => Effect.Effect<void, unknown>
  }
>()('openlibrary/EditorFonts') {
  static readonly layer = Layer.effect(
    this,
    Effect.gen(function* () {
      const api = yield* EditorFontsApi
      const fontsLoader = yield* FontsLoader
      const listeners = new Set<() => void>()
      let state = initialState
      let initialized = false

      const getState = () => state
      const subscribe = (listener: () => void) => {
        listeners.add(listener)
        return () => listeners.delete(listener)
      }
      const updateState = (update: Partial<EditorFontsState>) =>
        Effect.sync(() => {
          state = { ...state, ...update }
          listeners.forEach((listener) => listener())
        })

      const initialize = Effect.fn('EditorFonts.initialize')(function* () {
        if (initialized || state.initializing) return

        yield* updateState({ initializing: true, error: undefined })
        yield* Effect.gen(function* () {
          const [fonts, favoriteFonts] = yield* Effect.all([
            fontsLoader.fetchFonts(),
            api.getFavoriteFonts(),
          ])
          yield* updateState({ fonts, favoriteFonts: [...favoriteFonts] })
          initialized = true
        }).pipe(
          Effect.tapError((error) => updateState({ error })),
          Effect.ensuring(updateState({ initializing: false })),
        )
      })

      const setFavorite = Effect.fn('EditorFonts.setFavorite')(function* (
        font: string,
        favorite: boolean,
      ) {
        if (state.favoriteFonts.includes(font) === favorite) return

        const favoriteFonts = favorite
          ? [...state.favoriteFonts, font].sort()
          : state.favoriteFonts.filter((currentFont) => currentFont !== font)

        // TODO handle concurrency
        yield* updateState({ error: undefined })
        yield* api
          .saveFavoriteFonts(favoriteFonts)
          .pipe(Effect.tapError((error) => updateState({ error })))
        yield* updateState({ favoriteFonts })
      })

      return EditorFonts.of({ getState, subscribe, initialize, setFavorite })
    }),
  )
}
