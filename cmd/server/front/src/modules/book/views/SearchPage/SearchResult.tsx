import { useQuery } from '@tanstack/react-query'
import { useSearchState } from './state'
import { getPreloadedBookSearchResult, httpSearchBooks } from '../../api'
import React, { useMemo } from 'react'
import {
  SearchParams,
  searchParamsToBookSearchRequest,
  useSearchParamsState,
} from './search-params'
import SearchBookCard from './SearchBookCard'
import Spinner from '@/components/spinner'
import { useDebounce } from '@/lib/react-utils'

export default function SearchResults() {
  const { data, isLoading } = useDebounceSearchResultsQuery()
  const debouncedIsLoading = useDebounce(isLoading, 10)

  return (
    <div>
      <div id="search-books" className="space-y-2 relative">
        {data?.books.map((book) => <SearchBookCard key={book.id} book={book} />)}
        {debouncedIsLoading && (
          <div className="bg-background/50 absolute inset-0 z-10 flex justify-center pt-8">
            <Spinner className="text-foreground" />
          </div>
        )}
      </div>

      {data && (
        <div id="search-nerd-stats" className="mt-5 text-foreground/30">
          <div className="text-xs font-mono">
            <strong className="font-semibold">Nerd stuff</strong>
            <br />
            {`Took: ${data.took}us (${data.took / 1000}ms)`}
            <br />
            {`Cache status: ${data.cache.hit ? 'hit' : 'miss'}`}
            <br />
            {`Cache key: ${data.cache.key}`}
          </div>
        </div>
      )}
    </div>
  )
}

function useDebounceSearchResultsQuery() {
  const { activeParams: params, ready } = useSearchParamsState()
  const searchParams = useMemo(() => searchParamsToBookSearchRequest(params), [params])
  const searchKey = useMemo(() => getSearchKey(params), [params])

  return useQuery({
    queryKey: ['search', searchKey],
    enabled: ready,
    queryFn: async () => {
      try {
        useSearchState.getState().setLoading(true)
        const result = await httpSearchBooks(searchParams)
        return result
      } finally {
        useSearchState.getState().setLoading(false)
      }
    },
    initialData: getPreloadedBookSearchResult(searchParams),
  })
}

function getSearchKey(params: SearchParams) {
  return JSON.stringify(searchParamsToBookSearchRequest(params))
}
