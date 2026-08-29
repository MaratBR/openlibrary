import { Font } from '@/features/fonts-loader/api'
import { FontsLoader } from '@/features/fonts-loader/loader'
import { UserDataApi } from '@/public.api/user-data'
import { Effect } from 'effect'
import { UnknownException } from 'effect/Cause'
import { z } from 'zod'
import { create } from 'zustand'

export type FontsState = {
  initializing: boolean
  error: unknown
  fonts: ReadonlyArray<Readonly<Font>>
  init(): Effect.Effect<void, UnknownException, never>

  favoriteFonts: string[]
  favorite(font: string, favorite: boolean): Effect.Effect<void, never, never>
}

const USER_DATA_FONTS_KEY = 'bm:fonts'

const FontsStateSchema = z.object({
  favorite: z.array(z.string()),
})

export const useFonts = create<FontsState>()((set, get) => ({
  initializing: false,
  error: undefined,
  fonts: [],

  init() {
    return Effect.gen(function* () {
      if (get().initializing) {
        return
      }
      set({ initializing: true })

      try {
        const fonts = yield* FontsLoader.fetchFonts()
        const fontState = yield* UserDataApi.getUserData(USER_DATA_FONTS_KEY, {
          schema: FontsStateSchema,
        })
        set({ fonts, initializing: false, favoriteFonts: fontState?.favorite ?? [] })
      } catch (error: unknown) {
        set({ initializing: false, error })
      }
    })
  },

  favoriteFonts: [],

  favorite(font, favorite) {
    return Effect.sync(() => {
      const favoriteCurrent = get().favoriteFonts.includes(font)

      if (favoriteCurrent === favorite) {
        return
      }

      if (favorite) {
        set({ favoriteFonts: [...get().favoriteFonts, font].sort() })
      } else {
        set({ favoriteFonts: get().favoriteFonts.filter((x) => x !== font) })
      }
    })
  },
}))
