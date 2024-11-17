import { create } from 'zustand'
import { NumberRange } from './state'
import { useSearchParams } from 'react-router-dom'
import React, { useMemo } from 'react'
import {
  DefinedTagDto,
  httpTagsGetByName,
  parseSearchBooksRequest,
  searchBookRequestToURLSearchParams,
  SearchBooksRequest,
} from '../../api'
import isEqual from 'lodash.isequal'
import { isSearchQueryEqual, parseQueryStringArray, stringArray } from '@/modules/common/api'
import { useShallow } from 'zustand/react/shallow'

export type SearchParams = {
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

export type SearchParamsState = {
  params: SearchParams
  activeParams: SearchParams
  ready: boolean

  setFavorites(range: NumberRange | null): void
  setWords(range: NumberRange | null): void
  setChapters(range: NumberRange | null): void
  setWordsPerChapter(range: NumberRange | null): void
  setIncludeTags(tags: DefinedTagDto[]): void
  setExcludeTags(tags: DefinedTagDto[]): void
  initFromSearchBooksRequest(sp: SearchBooksRequest): Promise<void>
  applyChanges(params?: SearchParams): void
  hasChanges(): boolean
}

const defaultParams: () => SearchParams = () => ({
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
})

export const useSearchParamsState = create<SearchParamsState>()((set, get) => ({
  params: defaultParams(),
  activeParams: defaultParams(),
  ready: false,

  applyChanges(params) {
    set({ activeParams: params ?? get().params })
  },

  setChapters(range) {
    set((state) => ({
      params: {
        ...state.params,
        chapters: range,
      },
    }))
  },
  setWords(range) {
    set((state) => ({
      params: {
        ...state.params,
        words: range,
      },
    }))
  },
  setWordsPerChapter(range) {
    set((state) => ({
      params: {
        ...state.params,
        wordsPerChapter: range,
      },
    }))
  },
  setFavorites(range) {
    set((state) => ({
      params: {
        ...state.params,
        favorites: range,
      },
    }))
  },
  setIncludeTags(tags) {
    set((state) => ({
      params: {
        ...state.params,
        include: {
          tags,
        },
      },
    }))
  },
  setExcludeTags(tags) {
    set((state) => ({
      params: {
        ...state.params,
        exclude: {
          ...state.params.exclude,
          tags,
        },
      },
    }))
  },
  hasChanges() {
    const { activeParams, params } = get()
    return !isEqual(params, activeParams)
  },
  async initFromSearchBooksRequest(req) {
    const paramsFromUrl = await normalizeRequest(req)

    set({
      params: paramsFromUrl,
      activeParams: paramsFromUrl,
      ready: true,
    })
  },
}))

export function useBookSearchParams() {
  const {
    params,
    setChapters,
    setFavorites,
    setWords,
    setWordsPerChapter,
    setIncludeTags,
    setExcludeTags,
    hasChanges,
    applyChanges,
  } = useSearchParamsState()

  useInitializeSearchBooksRequestFromSearchQuery()

  return {
    params,
    setChapters,
    setFavorites,
    setWords,
    setWordsPerChapter,
    setIncludeTags,
    setExcludeTags,
    applyChanges,
    hasChanges,
  }
}

function useInitializeSearchBooksRequestFromSearchQuery() {
  const [sp, setSp] = useSearchParams()
  const searchRequest = useMemo(() => parseSearchBooksRequest(sp), [sp])
  const initialized = React.useRef(false)
  React.useEffect(() => {
    if (initialized.current) return
    initialized.current = true
    useSearchParamsState.getState().initFromSearchBooksRequest(searchRequest)
  }, [searchRequest])

  const { activeParams, ready } = useSearchParamsState(
    useShallow((s) => ({
      activeParams: s.activeParams,
      ready: s.ready,
    })),
  )
  React.useEffect(() => {
    if (!ready) return
    const newSP = searchBookRequestToURLSearchParams(searchParamsToBookSearchRequest(activeParams))

    if (!isSearchQueryEqual(newSP, sp)) {
      setSp(newSP)
    }
  }, [sp, setSp, activeParams, ready])
}

async function normalizeRequest(req: SearchBooksRequest): Promise<SearchParams> {
  const includeTagNames = parseQueryStringArray(req.it)
  const excludeTagNames = parseQueryStringArray(req.et)

  const allTags = [...new Set([...includeTagNames, ...excludeTagNames])]

  const definedTags = allTags.length === 0 ? [] : await httpTagsGetByName(allTags)

  const includeTags = definedTags.filter((x) => includeTagNames.includes(x.name))
  const excludeTags = definedTags.filter((x) => excludeTagNames.includes(x.name))

  function numberRange(min?: string, max?: string): NumberRange | null {
    if (!min && !max) return null

    let minNumber = min ? +min : null
    if (minNumber !== null && Number.isNaN(minNumber)) minNumber = null

    let maxNumber = max ? +max : null
    if (maxNumber !== null && Number.isNaN(maxNumber)) maxNumber = null

    if (minNumber === null && maxNumber === null) return null
    return {
      min: minNumber,
      max: maxNumber,
    }
  }

  return {
    words: numberRange(req['w.min'], req['w.max']),
    chapters: numberRange(req['c.min'], req['c.max']),
    wordsPerChapter: numberRange(req['wc.min'], req['wc.max']),
    favorites: numberRange(req['f.min'], req['f.max']),
    include: {
      tags: includeTags,
    },
    exclude: {
      tags: excludeTags,
    },
  }
}

export function searchParamsToBookSearchRequest(params: SearchParams): SearchBooksRequest {
  const numStr = (v?: number) => (v === undefined ? undefined : v + '')

  function removeUndefinedValues<T extends object>(obj: T): T {
    // @ts-expect-error it works
    return Object.fromEntries(Object.entries(obj).filter(([_, v]) => v != null))
  }

  return removeUndefinedValues({
    'w.min': numStr(params.words?.min ?? undefined),
    'w.max': numStr(params.words?.max ?? undefined),
    'c.min': numStr(params.chapters?.min ?? undefined),
    'c.max': numStr(params.chapters?.max ?? undefined),
    'wc.min': numStr(params.wordsPerChapter?.min ?? undefined),
    'wc.max': numStr(params.wordsPerChapter?.max ?? undefined),
    'f.min': numStr(params.favorites?.min ?? undefined),
    'f.max': numStr(params.favorites?.max ?? undefined),
    it: stringArray(params.include.tags.map((x) => x.name)),
    et: stringArray(params.exclude.tags.map((x) => x.name)),
  })
}
