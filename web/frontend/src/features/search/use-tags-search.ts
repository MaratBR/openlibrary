import { DEFAULT_SEARCH_DEBOUNCE } from '@/config'
import { appRuntime } from '@/effect/runtime'
import { Effect } from 'effect'
import { useEffect, useState } from 'react'
import { DefinedTagDto, SearchApi } from './api'

export type TagsSearchOptions = {
  query: string
  fetchDefault?: boolean
  debounceTimeout?: number
  enabled?: boolean
}

export type TagsSearchResult = {
  data: ReadonlyArray<DefinedTagDto> | undefined
  error: unknown
  isLoading: boolean
}

export function useTagsSearch({
  query,
  fetchDefault = false,
  debounceTimeout = DEFAULT_SEARCH_DEBOUNCE,
  enabled = true,
}: TagsSearchOptions): TagsSearchResult {
  const normalizedQuery = query.trim().toLowerCase()
  const [debouncedQuery, setDebouncedQuery] = useState(normalizedQuery)
  const [result, setResult] = useState<TagsSearchResult>({
    data: undefined,
    error: undefined,
    isLoading: false,
  })

  useEffect(() => {
    const timeout = window.setTimeout(() => setDebouncedQuery(normalizedQuery), debounceTimeout)
    return () => window.clearTimeout(timeout)
  }, [normalizedQuery, debounceTimeout])

  useEffect(() => {
    if (!enabled || (debouncedQuery.length === 0 && !fetchDefault)) {
      setResult({ data: undefined, error: undefined, isLoading: false })
      return
    }

    const abortController = new AbortController()
    setResult((current) => ({ ...current, error: undefined, isLoading: true }))
    void appRuntime
      .runPromise(
        Effect.flatMap(SearchApi, (searchApi) => searchApi.searchTags(debouncedQuery)),
        { signal: abortController.signal },
      )
      .then(
        (data) => setResult({ data, error: undefined, isLoading: false }),
        (error: unknown) => {
          if (!abortController.signal.aborted) {
            setResult({ data: undefined, error, isLoading: false })
          }
        },
      )

    return () => abortController.abort()
  }, [debouncedQuery, enabled, fetchDefault])

  return result
}
