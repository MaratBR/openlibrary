import { ManagerBookDetailsDto } from '@/api/bm/book'
import SanitizeHTML from '@/common/SanitizeHTML'
import { DashboardContent } from '@/components/dashboard-layout-components'
import { formatNumber, formatNumberK } from '@/util/fmt'
import { NavLink } from 'react-router'

export function BookGeneral({ book }: { book: ManagerBookDetailsDto }) {
  return (
    <>
      <DashboardContent.Card className="mt-4">
        <div className="flex gap-2">
          <div className="bg-gradient-to-b from-lime-500 to-lime-300 p-6 rounded-2xl text-3xl font-semibold text-white max-w-64">
            <div>{window._('book.words', { count: formatNumberK(book.words) })}</div>
            <div className="text-lg opacity-80 leading-5">
              {window._('book.wordsPerChapter', { count: formatNumberK(book.wordsPerChapter) })}
            </div>
          </div>

          <div className="bg-gradient-to-b p-6 rounded-2xl text-3xl font-semibold max-w-64 bg-highlight">
            <div>{window._('book.chapters', { count: formatNumber(book.chapters.length) })}</div>
          </div>
        </div>

        <div className="flex gap-1 my-4">
          <NavLink className="btn btn--outline" to={`/books/${book.id}/edit`}>
            <i className="fa-solid fa-pen mr-2" />
            {window._('common.edit')}
          </NavLink>
        </div>

        <dl className="dl">
          <dt className="dt">{window._('bookManager.edit.name')}</dt>
          <dd className="dd">{book.name}</dd>

          <dt className="dt">{window._('bookManager.edit.summary')}</dt>
          <dd className="dd">
            <SanitizeHTML value={book.summary} />
          </dd>
        </dl>
      </DashboardContent.Card>
    </>
  )
}
