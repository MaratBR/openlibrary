import { create } from 'zustand'
import { useSearchParams } from 'react-router-dom'
import React from 'react'
import {
  BookSearchItem,
  DefinedTagDto,
  httpSearchBooks,
  parseSearchBooksRequest,
  SearchBooksResponse,
  tryGetPreloadedSearchResult,
} from '../../api'
import { isNotFalsy, toDictionaryByProperty } from '@/lib/utils'
import { useQuery } from '@tanstack/react-query'
import { SearchBooksRequest, searchBooksRequestToURLSearchParams } from '../../api/search-request'

export type NumberRange = {
  max: number | null
  min: number | null
}

export namespace SearchFilters {
  export type Type = {
    page: number
    words: NumberRange | null
    chapters: NumberRange | null
    wordsPerChapter: NumberRange | null
    favorites: NumberRange | null
    include: {
      tags: DefinedTagDto[]
    }
    exclude: {
      tags: DefinedTagDto[]
    }
  }

  export function toSearchBooksRequest(params: Type): SearchBooksRequest {
    const numStr = (v?: number) => (v === undefined ? undefined : v + '')

    return {
      'w.min': numStr(params.words?.min ?? undefined),
      'w.max': numStr(params.words?.max ?? undefined),
      'c.min': numStr(params.chapters?.min ?? undefined),
      'c.max': numStr(params.chapters?.max ?? undefined),
      'wc.min': numStr(params.wordsPerChapter?.min ?? undefined),
      'wc.max': numStr(params.wordsPerChapter?.max ?? undefined),
      'f.min': numStr(params.favorites?.min ?? undefined),
      'f.max': numStr(params.favorites?.max ?? undefined),
      it: params.include.tags.map((x) => x.id),
      et: params.exclude.tags.map((x) => x.id),
      p: params.page,
    }
  }
}

export type BookSearchItemState = Omit<BookSearchItem, 'tags'> & {
  tags: DefinedTagDto[]
}

export type SearchParamsState = {
  filters: SearchFilters.Type
  activeFilters: SearchFilters.Type
  ready: boolean
  isLoading: boolean
  results: {
    page: number
    pageSize: number
    totalPages: number
    books: BookSearchItemState[]
  }
  nerdStuff?: {
    cacheKey: string
    cacheHit: boolean
    cacheTook: number
  }
  extremes: {
    chapters: NumberRange
    words: NumberRange
    wordsPerChapter: NumberRange
    favorites: NumberRange
  }

  setExtremes(extremes: SearchParamsState['extremes']): void
  setFavorites(range: NumberRange | null): void
  setWords(range: NumberRange | null): void
  setChapters(range: NumberRange | null): void
  setWordsPerChapter(range: NumberRange | null): void
  setIncludeTags(tags: DefinedTagDto[]): void
  setExcludeTags(tags: DefinedTagDto[]): void
  applyChanges(params?: SearchFilters.Type): void
  search(sp: SearchBooksRequest): Promise<void>
  setResponse(response: SearchBooksResponse): void
}

const defaultParams: () => SearchFilters.Type = () => ({
  words: null,
  chapters: null,
  wordsPerChapter: null,
  favorites: null,
  include: {
    tags: [],
  },
  exclude: {
    tags: [],
  },
  page: 1,
})

export const useSearchState = create<SearchParamsState>()((set, get) => ({
  filters: defaultParams(),
  activeFilters: defaultParams(),
  ready: false,
  isLoading: false,

  results: {
    books: [],
    page: 1,
    pageSize: 20,
    totalPages: -1,
  },

  extremes: {
    chapters: { max: null, min: null },
    words: { max: null, min: null },
    wordsPerChapter: { max: null, min: null },
    favorites: { max: null, min: null },
  },

  setExtremes(extremes) {
    set({ extremes })
  },

  applyChanges(params) {
    set({ activeFilters: params ?? get().filters })
  },

  setChapters(range) {
    set((state) => ({
      filters: {
        ...state.filters,
        chapters: range,
      },
    }))
  },
  setWords(range) {
    set((state) => ({
      filters: {
        ...state.filters,
        words: range,
      },
    }))
  },
  setWordsPerChapter(range) {
    set((state) => ({
      filters: {
        ...state.filters,
        wordsPerChapter: range,
      },
    }))
  },
  setFavorites(range) {
    set((state) => ({
      filters: {
        ...state.filters,
        favorites: range,
      },
    }))
  },
  setIncludeTags(tags) {
    set((state) => ({
      filters: {
        ...state.filters,
        include: {
          tags,
        },
      },
    }))
  },
  setExcludeTags(tags) {
    set((state) => ({
      filters: {
        ...state.filters,
        exclude: {
          ...state.filters.exclude,
          tags,
        },
      },
    }))
  },

  setResponse(response) {
    const tagsById = toDictionaryByProperty(response.tags, 'id')

    set((s) => {
      const filters: SearchFilters.Type = {
        ...s.activeFilters,
        include: {
          tags: (response.filter?.includeTags || [])
            .map((tagId) => {
              const tag = tagsById[tagId]
              if (!tag) return null
              return tag
            })
            .filter(isNotFalsy),
        },
        exclude: {
          tags: (response.filter?.excludeTags || [])
            .map((tagId) => {
              const tag = response.tags.find((t) => t.id === tagId)
              if (!tag) return null
              return tag
            })
            .filter(isNotFalsy),
        },
      }

      return {
        results: {
          books: response.books.map((book) => {
            return {
              ...book,
              tags: book.tags.map((tagId) => tagsById[tagId]).filter(isNotFalsy),
            }
          }),
          page: response.page,
          pageSize: response.pageSize,
          totalPages: response.totalPages,
        },
        isLoading: false,
        filters,
        activeFilters: filters,
      }
    })
  },

  async search(req) {
    set({ isLoading: true })

    try {
      const start = performance.now()
      const response = await httpSearchBooks(req)
      const took = performance.now() - start
      if (took < 500) {
        await new Promise((r) => setTimeout(r, 400 - took))
      }
      this.setResponse(response)
    } catch {
      set({ isLoading: false })
    }
  },
}))

export function useBookSearchParams() {
  const {
    filters: params,
    setChapters,
    setFavorites,
    setWords,
    setWordsPerChapter,
    setIncludeTags,
    setExcludeTags,
    applyChanges: applyChangesInState,
  } = useSearchState()

  const [searchRequest, setSearchRequest] = useSearchBooksRequest()

  useQuery({
    queryKey: ['search', searchRequest],
    queryFn: async () => {
      const state = useSearchState.getState()
      const preloadedResult = tryGetPreloadedSearchResult()
      if (preloadedResult) {
        state.setResponse(preloadedResult)
      } else {
        await state.search(searchRequest)
      }
      return 'OK'
    },
    staleTime: 0,
    gcTime: 0,
  })

  const applyChanges = React.useCallback(
    (params?: SearchFilters.Type) => {
      applyChangesInState(params)
      setSearchRequest(SearchFilters.toSearchBooksRequest(useSearchState.getState().activeFilters))
    },
    [applyChangesInState, setSearchRequest],
  )

  return {
    params,
    setChapters,
    setFavorites,
    setWords,
    setWordsPerChapter,
    setIncludeTags,
    setExcludeTags,
    applyChanges,
  }
}

function useSearchBooksRequest(): [
  value: SearchBooksRequest,
  set: (value: SearchBooksRequest) => void,
] {
  const [sp, setSp] = useSearchParams()
  const searchRequest = React.useMemo(() => parseSearchBooksRequest(sp), [sp])
  const setSearchRequest = React.useCallback(
    (value: SearchBooksRequest) => {
      setSp(searchBooksRequestToURLSearchParams(value))
    },
    [setSp],
  )
  return [searchRequest, setSearchRequest]
}
