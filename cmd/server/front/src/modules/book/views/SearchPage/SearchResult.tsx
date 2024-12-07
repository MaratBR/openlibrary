import React from 'react'
import { useSearchState } from './search-params'
import SearchBookCard from './SearchBookCard'
import Spinner from '@/components/spinner'
import { useDebounce } from '@/lib/react-utils'

export default function SearchResults() {
  const books = useSearchState((s) => s.books)
  const isLoading = useSearchState((s) => s.isLoading)
  const debouncedIsLoading = useDebounce(isLoading, 10)

  return (
    <div>
      <div id="search-books" className="space-y-2 relative">
        {books.map((book) => (
          <SearchBookCard key={book.id} book={book} />
        ))}
        {debouncedIsLoading && (
          <div className="bg-background/50 absolute inset-0 z-10 flex justify-center pt-8">
            <Spinner className="text-foreground" />
          </div>
        )}
      </div>
      {/* 
      {data && (
        <div id="search-nerd-stats" className="mt-5 text-foreground/30">
          <div className="text-xs font-mono">
            <strong className="font-semibold">Nerd stuff</strong>
            <br />
            {`Took: ${data.booksTook}us (${data.booksTook / 1000}ms)`}
            <br />
            {`Cache status: ${data.booksMeta.cacheHit ? 'hit' : 'miss'}`}
            <br />
            {`Cache key: ${data.booksMeta.cacheKey}`}
          </div>
        </div>
      )} */}
    </div>
  )
}
