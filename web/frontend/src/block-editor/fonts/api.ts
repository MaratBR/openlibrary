import { UserDataApi } from '@/features/api/user-data'
import { Context, Effect, Layer, Schema } from 'effect'

const USER_DATA_FONTS_KEY = 'bm:fonts'

const EditorFontsData = Schema.Struct({
  favorite: Schema.Array(Schema.String),
})

export class EditorFontsApi extends Context.Service<
  EditorFontsApi,
  {
    readonly getFavoriteFonts: () => Effect.Effect<ReadonlyArray<string>, unknown>
    readonly saveFavoriteFonts: (fonts: ReadonlyArray<string>) => Effect.Effect<void, unknown>
  }
>()('openlibrary/EditorFontsApi') {
  static readonly layer = Layer.effect(
    this,
    Effect.gen(function* () {
      const userDataApi = yield* UserDataApi

      const getFavoriteFonts = Effect.fn('EditorFontsApi.getFavoriteFonts')(function* () {
        const data = yield* userDataApi.getUserData(USER_DATA_FONTS_KEY, {
          schema: EditorFontsData,
        })
        return data?.favorite ?? []
      })

      const saveFavoriteFonts = Effect.fn('EditorFontsApi.saveFavoriteFonts')(function* (
        fonts: ReadonlyArray<string>,
      ) {
        yield* userDataApi.setUserData(USER_DATA_FONTS_KEY, { favorite: fonts })
      })

      return EditorFontsApi.of({ getFavoriteFonts, saveFavoriteFonts })
    }),
  )
}
