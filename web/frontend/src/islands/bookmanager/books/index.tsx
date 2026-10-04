import { BMBookAPI, ManagerBookDto } from '@/api/bm/book'
import { DashboardContent } from '@/components/dashboard-layout-components'
import { Pagination } from '@/components/Pagination'
import { getPage } from '@/lib/url'
import { useMutation } from '@tanstack/react-query'
import { useState } from 'react'
import { LoaderFunctionArgs, NavLink, useLoaderData, useRevalidator } from 'react-router'
import { Counts, CoverImage, ManagerDialog, RequestError, usePageTitle } from '../ui'

export const booksRouteLoader = async ({ request }: LoaderFunctionArgs) => {
  const resp = await BMBookAPI.getInstance().getBooks({
    size: 20,
    page: getPage(request.url),
    search: '',
  })
  resp.throwIfError()
  return { booksResponse: resp }
}

export function Books() {
  const { booksResponse } = useLoaderData<Awaited<ReturnType<typeof booksRouteLoader>>>()
  const [notice, setNotice] = useState('')
  usePageTitle(window._('bookManager.books.title'))
  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader title={window._('bookManager.books.title')}>
        <NavLink to="/books/new" className="Btn Btn--lg Btn--primary">
          + {window._('bookManager.books.addBook')}
        </NavLink>
      </DashboardContent.StickyHeader>
      <p className="text-secondary-foreground mb-6">
        {window._('bookManager.ui.libraryDescription')}
      </p>
      <p role="status">{notice}</p>
      {booksResponse.data.books.length ? (
        <div className="BM-grid">
          {booksResponse.data.books.map((book) => (
            <BookCard key={`${book.id}-${book.isTrashed}`} book={book} onChanged={setNotice} />
          ))}
        </div>
      ) : (
        <div className="Card space-y-4">
          <p>{window._('bookManager.ui.emptyLibrary')}</p>
          <NavLink className="Btn Btn--primary" to="/books/new">
            {window._('bookManager.books.addBook')}
          </NavLink>
        </div>
      )}
      {booksResponse.data.totalPages > 1 && (
        <div className="mt-6 flex flex-wrap">
          <Pagination.Facade
            page={booksResponse.data.page}
            size={5}
            totalPages={booksResponse.data.totalPages}
          />
        </div>
      )}
    </DashboardContent.Root>
  )
}

function BookCard({
  book,
  onChanged,
}: {
  book: ManagerBookDto
  onChanged: (message: string) => void
}) {
  const [confirm, setConfirm] = useState(false)
  const revalidator = useRevalidator()
  const mutation = useMutation({
    mutationFn: async () => {
      const response = await BMBookAPI.getInstance().trashBook({
        id: book.id,
        trash: !book.isTrashed,
      })
      response.throwIfError()
    },
    onSuccess: async () => {
      setConfirm(false)
      onChanged(
        window._(
          book.isTrashed
            ? 'bookManager.books.restoreBook.trashedBookNotif'
            : 'bookManager.books.trashBook.trashedBookNotif',
        ),
      )
      window.toast({
        title: window._('common.operationSuccessful'),
        text: window._(
          book.isTrashed
            ? 'bookManager.books.restoreBook.trashedBookNotif'
            : 'bookManager.books.trashBook.trashedBookNotif',
        ),
      })
      await revalidator.revalidate()
    },
  })
  const action = window._(book.isTrashed ? 'bookManager.ui.restore' : 'common.trash')
  const title = window._(
    book.isTrashed ? 'bookManager.ui.restoreTitle' : 'bookManager.ui.trashTitle',
    { name: book.name },
  )
  return (
    <article className="BM-bookCard">
      <NavLink to={`/books/${book.id}`} className="BM-bookArt" aria-label={book.name}>
        <CoverImage cover={book.cover} />
      </NavLink>
      <details
        className="BM-actions"
        onKeyDown={(e) => {
          if (e.key === 'Escape') {
            e.currentTarget.open = false
            e.currentTarget.querySelector('summary')?.focus()
          }
        }}
      >
        <summary
          className="Btn Btn--icon Btn--outline"
          aria-label={window._('bookManager.ui.actions', { name: book.name })}
        >
          ···
        </summary>
        <div>
          <button
            className="Btn Btn--ghost"
            onClick={(e) => {
              const menu = e.currentTarget.closest('details')
              menu?.removeAttribute('open')
              menu?.querySelector('summary')?.focus()
              mutation.reset()
              setConfirm(true)
            }}
          >
            {action}
          </button>
        </div>
      </details>
      <div className="BM-bookBody">
        <p className="text-sm text-secondary-foreground">
          {window._(
            book.isTrashed
              ? 'bookManager.books.trashed'
              : book.isBanned
                ? 'bookManager.books.banned'
                : book.isPubliclyVisible
                  ? 'bookManager.ui.public'
                  : 'bookManager.ui.hidden',
          )}{' '}
          · {book.ageRating}
        </p>
        <h2 className="BM-bookTitle">
          <NavLink to={`/books/${book.id}`} title={book.name}>
            {book.name}
          </NavLink>
        </h2>
        <p className="text-sm text-secondary-foreground">
          <Counts chapters={book.chapters} words={book.words} />
        </p>
      </div>
      <footer className="BM-cardFooter">
        <NavLink className="Link" to={`/books/${book.id}/edit`}>
          {window._('bookManager.ui.editDetails')}
        </NavLink>
        <NavLink className="Link" to={`/books/${book.id}?t=chapters`}>
          {window._('bookManager.edit.chapters')}
        </NavLink>
      </footer>
      {confirm && (
        <ManagerDialog
          title={title}
          onClose={() => {
            if (!mutation.isPending) setConfirm(false)
          }}
        >
          <h2 className="text-xl">{title}</h2>
          <p className="my-4">
            {window._(
              book.isTrashed
                ? 'bookManager.books.restoreBook.description'
                : 'bookManager.books.trashBook.description',
            )}
          </p>
          {mutation.isError && <RequestError />}
          <div className="flex gap-2">
            <button
              className="Btn Btn--outline"
              disabled={mutation.isPending}
              onClick={() => setConfirm(false)}
            >
              {window._('common.cancel')}
            </button>
            <button
              className="Btn Btn--destructive"
              disabled={mutation.isPending}
              onClick={() => mutation.mutate()}
            >
              {mutation.isPending ? window._('bookManager.ui.saving') : action}
            </button>
          </div>
        </ManagerDialog>
      )}
    </article>
  )
}
