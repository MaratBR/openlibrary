import { Font } from '@/features/fonts-loader/api'
import { FontsLoader } from '@/features/fonts-loader/loader'
import { UserDataApi } from '@/features/api/user-data'
import { Effect } from 'effect'
import { atom, useAtomValue, useSetAtom } from 'jotai'
import { z } from 'zod'

const USER_DATA_FONTS_KEY = 'bm:fonts'

const FontsStateSchema = z.object({
  favorite: z.array(z.string()),
})

const fontsAtom = atom<ReadonlyArray<Readonly<Font>>>([])
const favoriteFontsAtom = atom(Array<string>())
const fontsInitializingAtom = atom(false)
const fontsErrorAtom = atom<unknown>(undefined)

const initializeFontsAtom = atom(null, (get, set) => {
  if (get(fontsInitializingAtom)) return Effect.void

  return Effect.gen(function* () {
    yield* Effect.sync(() => {
      set(fontsInitializingAtom, true)
      set(fontsErrorAtom, undefined)
    })

    yield* Effect.gen(function* () {
      const fonts = yield* FontsLoader.fetchFonts()
      const userDataApi = yield* UserDataApi
      const fontState = yield* userDataApi.getUserData(USER_DATA_FONTS_KEY, {
        schema: FontsStateSchema,
      })
      yield* Effect.sync(() => {
        set(fontsAtom, fonts)
        set(favoriteFontsAtom, fontState?.favorite ?? [])
      })
    }).pipe(
      Effect.tapError((error) => Effect.sync(() => set(fontsErrorAtom, error))),
      Effect.ensuring(Effect.sync(() => set(fontsInitializingAtom, false))),
    )
  })
})

const setFontFavoriteAtom = atom(null, (get, set, font: string, favorite: boolean) =>
  Effect.sync(() => {
    const favoriteFonts = get(favoriteFontsAtom)
    const favoriteCurrent = favoriteFonts.includes(font)

    if (favoriteCurrent === favorite) return

    set(
      favoriteFontsAtom,
      favorite
        ? [...favoriteFonts, font].sort()
        : favoriteFonts.filter((currentFont) => currentFont !== font),
    )
  }),
)

export function useFonts() {
  return {
    initializing: useAtomValue(fontsInitializingAtom),
    error: useAtomValue(fontsErrorAtom),
    fonts: useAtomValue(fontsAtom),
    favoriteFonts: useAtomValue(favoriteFontsAtom),
    init: useSetAtom(initializeFontsAtom),
    favorite: useSetAtom(setFontFavoriteAtom),
  }
}
