import { DashboardContent } from '@/components/dashboard-layout-components'
import { LoaderFunctionArgs, NavLink, useLoaderData, useSearchParams } from 'react-router'
import { Schema } from 'effect'
import { BookGeneral } from './BookGeneral'
import { BookChapters } from './BookChapters'
import { BMBookAPI } from '@/api/bm/book'
import { BackToBooks, BookNavigation, usePageTitle } from '../ui'

export const bookRouteLoader = async ({ params }: LoaderFunctionArgs) => {
  const { bookId } = Schema.decodeUnknownSync(Schema.Struct({ bookId: Schema.NonEmptyString }))(
    params,
  )
  const resp = await BMBookAPI.getInstance().getBook(bookId)
  resp.throwIfError()
  return { bookResponse: resp, bookId }
}

export function Book() {
  const { bookResponse } = useLoaderData<Awaited<ReturnType<typeof bookRouteLoader>>>()
  const book = bookResponse.data
  const [params] = useSearchParams()
  const tab = params.get('t') === 'chapters' ? 'chapters' : 'general'
  usePageTitle(`${book.name} · ${window._(`bookManager.ui.${tab}`)}`)
  return (
    <DashboardContent.Root>
      <BackToBooks />
      <DashboardContent.StickyHeader title={book.name}>
        <NavLink className="Btn Btn--primary" to={`/books/${book.id}/edit`}>
          {window._('bookManager.ui.editDetails')}
        </NavLink>
      </DashboardContent.StickyHeader>
      <BookNavigation id={book.id} active={tab} />
      {tab === 'general' ? <BookGeneral book={book} /> : <BookChapters book={book} />}
    </DashboardContent.Root>
  )
}
