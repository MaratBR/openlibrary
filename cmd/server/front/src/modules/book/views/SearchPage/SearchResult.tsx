import React from 'react'
import { useSearchState } from './search-params'
import SearchBookCard from './SearchBookCard'
import Spinner from '@/components/spinner'
import { useDebounce } from '@/lib/react-utils'
import { useShallow } from 'zustand/shallow'
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
} from '@/components/ui/pagination'
import { NavLink, useSearchParams } from 'react-router-dom'

export default function SearchResults() {
  const books = useSearchState((s) => s.results.books)

  if (books.length === 0) {
    return (
      <div className="text-lg text-muted-foreground">
        No results, try to change your search query
      </div>
    )
  }

  return (
    <div>
      <div id="search-books" className="space-y-2 relative">
        <SearchPagination />
        <ul id="search-books-list" className="space-y-2">
          <LoadingOverlay />

          {books.map((book) => (
            <SearchBookCard key={book.id} book={book} />
          ))}
        </ul>
        <SearchPagination />
      </div>

      <NerdStuff />
    </div>
  )
}

function LoadingOverlay() {
  const isLoading = useSearchState((s) => s.isLoading)
  const debouncedIsLoading = useDebounce(isLoading, 10)

  return debouncedIsLoading ? (
    <div className="bg-background/60 absolute inset-0 z-10 flex justify-center pt-8">
      <Spinner className="text-foreground" />
    </div>
  ) : null
}

function SearchPagination() {
  const [sp] = useSearchParams()

  const pagination = useSearchState(
    useShallow((x) => ({
      page: x.results.page,
      totalPages: x.results.totalPages,
      pageSize: x.results.pageSize,
    })),
  )

  const MAX_PAGES = 7

  const links = React.useMemo(() => {
    function getHref(page: number) {
      const params = new URLSearchParams(sp)
      params.set('p', page.toString())
      return '?' + params.toString()
    }

    const links = [
      <PaginationItem key={pagination.page}>
        <PaginationLink isActive>{pagination.page}</PaginationLink>
      </PaginationItem>,
    ]

    let j = 1
    const MAX_ITER = 100

    while (links.length < MAX_PAGES && j < MAX_ITER) {
      const li = pagination.page - j,
        ri = pagination.page + j

      if (li > 0) {
        links.unshift(
          <PaginationItem key={li}>
            <NavLink to={getHref(li)}>
              <PaginationLink>{li}</PaginationLink>
            </NavLink>
          </PaginationItem>,
        )
      }

      if (ri <= pagination.totalPages && links.length < MAX_PAGES) {
        links.push(
          <PaginationItem key={ri}>
            <NavLink to={getHref(ri)}>
              <PaginationLink>{ri}</PaginationLink>
            </NavLink>
          </PaginationItem>,
        )
      }

      j++
    }

    if (j > MAX_ITER) {
      throw new Error('Infinite loop')
    }

    return links
  }, [pagination.page, pagination.totalPages, sp])

  return (
    <Pagination className="justify-start mt-4 mb-6">
      <PaginationContent>{links}</PaginationContent>
    </Pagination>
  )
}

function NerdStuff() {
  const meta = useSearchState((x) => x.nerdStuff)

  if (!meta) return null

  return <pre>{JSON.stringify(meta, null, 2)}</pre>
}
