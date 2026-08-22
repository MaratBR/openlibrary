import { searchModerationBooks } from '@/features/moderation/api'
import type { ModerationBookListEntry } from '@/features/moderation/api'
import { Checkbox } from '@/components'
import { DashboardContent } from '@/components/dashboard-layout-components'
import { ErrorDisplay } from '@/components/error'
import { Pagination } from '@/components/Pagination'
import { OLAPIResponse } from '@/http-client'
import { useEffect, useState } from 'react'
import { Form, Link, LoaderFunctionArgs, useLoaderData, useRouteError } from 'react-router'

async function responseData<T>(request: Promise<OLAPIResponse<T>>) {
  const response = await request
  response.throwIfError()
  return response.data
}

export async function moderationBooksLoader({ request }: LoaderFunctionArgs) {
  const params = new URL(request.url).searchParams
  const filters = {
    search: params.get('search') ?? '',
    exact: params.get('exact') === 'true',
    includeBanned: params.get('includeBanned') === 'true',
    includeDeleted: params.get('includeDeleted') === 'true',
  }
  const page = Math.max(1, Number(params.get('page')) || 1)
  return {
    books: await responseData(
      searchModerationBooks(
        filters.search,
        filters.exact,
        filters.includeBanned,
        filters.includeDeleted,
        page,
      ),
    ),
    filters,
  }
}

export default function ModerationBooks() {
  const { books, filters } = useLoaderData<Awaited<ReturnType<typeof moderationBooksLoader>>>()
  const [exact, setExact] = useState(filters.exact)
  const [includeBanned, setIncludeBanned] = useState(filters.includeBanned)
  const [includeDeleted, setIncludeDeleted] = useState(filters.includeDeleted)
  useEffect(() => {
    setExact(filters.exact)
    setIncludeBanned(filters.includeBanned)
    setIncludeDeleted(filters.includeDeleted)
  }, [filters])
  const pagination = {
    search: filters.search,
    exact: String(filters.exact),
    includeBanned: String(filters.includeBanned),
    includeDeleted: String(filters.includeDeleted),
  }
  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader title={window._('moderationPortal.books')} />
      <div className="w-full max-w-360 mx-auto px-4 md:px-8 py-7">
        <p className="text-secondary-foreground mb-6">
          {window._('moderationPortal.booksPage.description')}
        </p>
        <DashboardContent.Card className="card--elevated p-4 md:p-5 mb-6">
          <Form method="get" className="grid gap-4">
            <label>
              <span className="block text-sm font-medium mb-1.5">
                {window._('moderationPortal.booksPage.search')}
              </span>
              <div className="relative">
                <i className="fa-solid fa-magnifying-glass absolute left-3 top-1/2 -translate-y-1/2 text-secondary-foreground" />
                <input
                  className="input w-full pl-10"
                  name="search"
                  defaultValue={filters.search}
                  placeholder={window._('moderationPortal.booksPage.searchPlaceholder')}
                />
              </div>
            </label>
            <div className="flex flex-wrap items-center gap-x-6 gap-y-3">
              <input type="hidden" name="exact" value={String(exact)} />
              <label className="flex items-center gap-2 cursor-pointer" htmlFor="books-exact">
                <Checkbox
                  id="books-exact"
                  checked={exact}
                  onCheckedChange={(checked) => setExact(checked === true)}
                />
                {window._('moderationPortal.booksPage.exact')}
              </label>
              <input type="hidden" name="includeBanned" value={String(includeBanned)} />
              <label
                className="flex items-center gap-2 cursor-pointer"
                htmlFor="books-include-banned"
              >
                <Checkbox
                  id="books-include-banned"
                  checked={includeBanned}
                  onCheckedChange={(checked) => setIncludeBanned(checked === true)}
                />
                {window._('moderationPortal.booksPage.includeBanned')}
              </label>
              <input type="hidden" name="includeDeleted" value={String(includeDeleted)} />
              <label
                className="flex items-center gap-2 cursor-pointer"
                htmlFor="books-include-deleted"
              >
                <Checkbox
                  id="books-include-deleted"
                  checked={includeDeleted}
                  onCheckedChange={(checked) => setIncludeDeleted(checked === true)}
                />
                {window._('moderationPortal.booksPage.includeDeleted')}
              </label>
              <button className="btn btn--primary ml-auto">
                {window._('moderationPortal.booksPage.apply')}
              </button>
            </div>
          </Form>
        </DashboardContent.Card>
        <div className="flex justify-between items-baseline mb-3">
          <h2 className="text-lg font-semibold">
            {window._('moderationPortal.booksPage.results')}
          </h2>
          <span className="text-sm text-secondary-foreground">
            {books.total} {window._('moderationPortal.booksPage.found')}
          </span>
        </div>
        {books.entries.length ? (
          <div className="grid gap-3">
            {books.entries.map((book) => (
              <BookRow key={book.id} book={book} />
            ))}
          </div>
        ) : (
          <div className="card py-16 text-center text-secondary-foreground">
            {window._('moderationPortal.booksPage.empty')}
          </div>
        )}
        {books.totalPages > 1 && (
          <div className="flex justify-center mt-7">
            <Pagination.Facade
              page={books.page}
              totalPages={books.totalPages}
              size={7}
              getTo={(page) => ({
                search: new URLSearchParams({ ...pagination, page: String(page) }).toString(),
              })}
            />
          </div>
        )}
      </div>
    </DashboardContent.Root>
  )
}

function BookRow({ book }: { book: ModerationBookListEntry }) {
  const states = [
    book.isBanned && window._('moderationPortal.booksPage.banned'),
    book.isShadowBanned && window._('moderationPortal.booksPage.shadowBanned'),
    book.isTrashed && window._('moderationPortal.booksPage.trashed'),
    book.isPermanentlyRemoved && window._('moderationPortal.booksPage.deleted'),
  ].filter(Boolean)
  return (
    <Link to={`/books/${book.id}`} className="card card--elevated block hover:bg-foreground/5">
      <div className="flex flex-col md:flex-row md:items-center gap-4">
        <div className="min-w-0 flex-1">
          <div className="flex flex-wrap gap-2 items-center">
            <h3 className="font-semibold text-lg truncate">{book.name}</h3>
            {states.length ? (
              states.map((state) => (
                <span key={state as string} className="chip chip--destructive">
                  {state}
                </span>
              ))
            ) : (
              <span className="chip chip--primary">
                {window._('moderationPortal.booksPage.active')}
              </span>
            )}
          </div>
          <p className="text-sm text-secondary-foreground">
            {book.authorUserName} · #{book.id}
          </p>
        </div>
        <div className="grid grid-cols-3 gap-4 text-sm">
          <Metric
            label={window._('moderationPortal.booksPage.words')}
            value={book.words.toLocaleString()}
          />
          <Metric label={window._('moderationPortal.booksPage.chapters')} value={book.chapters} />
          <Metric
            label={window._('moderationPortal.booksPage.reports')}
            value={book.reportsCount}
          />
        </div>
      </div>
    </Link>
  )
}
function Metric({ label, value }: { label: string; value: string | number }) {
  return (
    <div>
      <div className="text-secondary-foreground">{label}</div>
      <div className="font-medium">{value}</div>
    </div>
  )
}

export function ModerationBooksErrorPage() {
  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader title={window._('moderationPortal.books')} />
      <div className="p-8">
        <div className="card">
          <ErrorDisplay error={useRouteError()} />
        </div>
      </div>
    </DashboardContent.Root>
  )
}
