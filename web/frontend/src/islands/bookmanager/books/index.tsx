import { BMBookAPI, ManagerBookDto } from '@/api/bm/book'
import { BookCover } from '@/components/BookCover'
import { DashboardContent } from '@/components/dashboard-layout-components'
import Modal from '@/components/Modal'
import { Pagination } from '@/components/Pagination'
import { getPage } from '@/lib/url'
import { formatNumberK } from '@/util/fmt'
import { useMutation } from '@tanstack/react-query'
import { useState } from 'preact/hooks'
import { LoaderFunctionArgs, NavLink, useLoaderData } from 'react-router'

export const booksRouteLoader = async ({ params, request }: LoaderFunctionArgs) => {
  const page = getPage(request.url)

  const resp = await BMBookAPI.getInstance().getBooks({
    size: 20,
    page,
    search: '',
  })

  return {
    booksResponse: resp,
  }
}

export function Books() {
  const { booksResponse } = useLoaderData<Awaited<ReturnType<typeof booksRouteLoader>>>()

  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader title={window._('bookManager.books.title')} />

      <div class="card">
        <div class="my-2 ml-4">
          <Pagination.Facade
            page={booksResponse.data.page}
            size={10}
            totalPages={booksResponse.data.totalPages}
          />
        </div>

        <table class="table">
          <tbody>
            {booksResponse.data.books.map((book) => (
              <BookRow key={book.id} book={book} />
            ))}
          </tbody>
        </table>
      </div>
    </DashboardContent.Root>
  )
}

function BookRow({ book }: { book: ManagerBookDto }) {
  const [trashed, setTrashed] = useState(book.isTrashed)

  return (
    <tr>
      <td style={{ width: 166 }}>
        <BookCover cover={book.cover} />
      </td>
      <td>
        <div>
          <span class="text-lg font-medium">{book.name}</span>
        </div>

        <div class="flex gap-1">
          <div class="chip chip--secondary chip--lg">
            {window._('book.chapters', { count: formatNumberK(book.chapters) })}
          </div>

          <div class="chip chip--secondary chip--lg">
            {window._('book.words', { count: formatNumberK(book.words) })}
          </div>
        </div>
      </td>
      <td>
        <div class="flex gap-2">
          <NavLink to={`/books/${book.id}`} className="btn btn--lg primary">
            <i class="fa-solid fa-pen mr-2" />
            {window._('common.edit')}
          </NavLink>
          <TrashBookButton book={book} trashed={trashed} onTrashedChanged={setTrashed} />
        </div>
      </td>
    </tr>
  )
}

function TrashBookButton({
  book,
  onTrashedChanged,
  trashed,
}: {
  book: ManagerBookDto
  trashed: boolean
  onTrashedChanged: (trashed: boolean) => void
}) {
  const [openTrashModal, setOpenTrashModal] = useState(false)
  const [openUntrashModal, setOpenUntrashModal] = useState(false)

  const trashBookMutation = useMutation({
    mutationFn: async (trash: boolean) => {
      const response = await BMBookAPI.getInstance().trashBook({
        trash,
        id: book.id,
      })
      response.throwIfError()
      window.toast({
        title: window._('common.operationSuccessful'),
        text: trash
          ? window._('bookManager.books.trashBook.trashedBookNotif')
          : window._('bookManager.books.restoreBook.trashedBookNotif'),
      })
      setOpenTrashModal(false)
      setOpenUntrashModal(false)
      onTrashedChanged(trash)
    },
  })

  return (
    <>
      <button
        onClick={() => {
          if (trashed) {
            setOpenUntrashModal(true)
          } else {
            setOpenTrashModal(true)
          }
        }}
        class="btn btn--lg btn--outline destructive"
      >
        {trashed ? window._('common.untrash') : window._('common.trash')}
      </button>
      <Modal onClose={() => setOpenTrashModal(false)} open={openTrashModal}>
        <div class="max-w-128">
          <h2 class="text-lg font-semibold">{window._('bookManager.books.trashBook.title')}</h2>
          <p class="my-2">{window._('bookManager.books.trashBook.description')}</p>
          <div class="flex gap-2 mt-4">
            <button onClick={() => setOpenTrashModal(false)} class="btn">
              {window._('common.cancel')}
            </button>
            <button class="btn destructive" onClick={() => trashBookMutation.mutate(true)}>
              {trashBookMutation.isPending && <span class="circle-loader mr-1" />}
              {window._('common.trash')}
            </button>
          </div>
        </div>
      </Modal>
      <Modal onClose={() => setOpenUntrashModal(false)} open={openUntrashModal}>
        <div class="max-w-128">
          <h2 class="text-lg font-semibold">{window._('bookManager.books.restoreBook.title')}</h2>
          <p class="my-2">{window._('bookManager.books.restoreBook.description')}</p>
          <div class="flex gap-2 mt-4">
            <button onClick={() => setOpenUntrashModal(false)} class="btn">
              {window._('common.cancel')}
            </button>
            <button class="btn destructive" onClick={() => trashBookMutation.mutate(false)}>
              {trashBookMutation.isPending && <span class="circle-loader mr-1" />}
              {window._('common.untrash')}
            </button>
          </div>
        </div>
      </Modal>
    </>
  )
}
