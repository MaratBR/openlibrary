import { create } from 'zustand'

export type NumberRange = {
  max: number | null
  min: number | null
}

export type SearchState = {
  extremes: {
    chapters: NumberRange
    words: NumberRange
    wordsPerChapter: NumberRange
    favorites: NumberRange
  }
  loading: boolean
  setLoading(value: boolean): void
  setExtremes(extremes: SearchState['extremes']): void
}

export const useSearchState = create<SearchState>()((set) => ({
  loading: false,
  extremes: {
    chapters: { max: null, min: null },
    words: { max: null, min: null },
    wordsPerChapter: { max: null, min: null },
    favorites: { max: null, min: null },
  },
  setLoading(value) {
    set({
      loading: value,
    })
  },
  setExtremes(extremes) {
    set((state) => {
      function normalizeRange(r: NumberRange): NumberRange {
        if (r.max === r.min && r.min !== null) {
          return {
            min: r.min,
            max: r.min + 1
          }
        }
        return r
      }

      return {
        ...state,
        extremes: {
          words: normalizeRange(extremes.words),
          wordsPerChapter: normalizeRange(extremes.wordsPerChapter),
          favorites: normalizeRange(extremes.favorites),
          chapters: normalizeRange(extremes.chapters),
        }
      }
    })
  },
}))
