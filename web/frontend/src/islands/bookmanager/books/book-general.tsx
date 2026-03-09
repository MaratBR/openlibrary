import { ManagerBookDetailsDto } from '@/api/bm/book'
import { formatNumber, formatNumberK } from '@/util/fmt'

export function BookGeneral({ book }: { book: ManagerBookDetailsDto }) {
  return (
    <>
      <div class="card mt-4">
        <div class="flex gap-2">
          <div class="bg-gradient-to-b from-lime-500 to-lime-300 p-6 rounded-2xl shadow-sm shadow-lime-400 text-3xl font-semibold text-white max-w-64">
            <div>{window._('book.words', { count: formatNumberK(book.words) })}</div>
            <div class="text-lg opacity-80 leading-5">
              {window._('book.wordsPerChapter', { count: formatNumberK(book.wordsPerChapter) })}
            </div>
          </div>

          <div class="bg-gradient-to-b p-6 rounded-2xl shadow-sm text-3xl font-semibold max-w-64 bg-highlight">
            <div>{window._('book.chapters', { count: formatNumber(book.chapters.length) })}</div>
          </div>
        </div>
      </div>
    </>
  )
}
