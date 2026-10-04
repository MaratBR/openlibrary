import { ManagerBookDetailsDto } from '@/api/bm/book'
import SanitizeHTML from '@/common/SanitizeHTML'
import { Counts, CoverImage } from '../ui'

export function BookGeneral({ book }: { book: ManagerBookDetailsDto }) {
  return (
    <div className="Card space-y-6">
      <div className="flex gap-6 flex-wrap items-center">
        <CoverImage cover={book.cover} />
        <div className="space-y-3 min-w-0">
          <p>
            {window._(
              book.isBanned
                ? 'bookManager.books.banned'
                : book.isPubliclyVisible
                  ? 'bookManager.ui.public'
                  : 'bookManager.ui.hidden',
            )}{' '}
            · {book.ageRating}
          </p>
          <div className="flex flex-wrap gap-2">
            {book.tags.map((tag) => (
              <span key={tag.id} className="Tag">
                {tag.name}
              </span>
            ))}
          </div>
          <a className="Link" href={`/book/${book.id}`} target="_blank" rel="noreferrer">
            {window._('bookManager.ui.publicPage')} ↗
          </a>
        </div>
      </div>
      <div className="bg-secondary rounded-xl p-4 space-y-2">
        <Counts chapters={book.chapters.length} words={book.words} />
        <p className="text-secondary-foreground">
          {window._('bookManager.ui.wordsPerChapter', {
            count: book.wordsPerChapter.toLocaleString(),
          })}
        </p>
      </div>
      <section>
        <h2 className="text-xl mb-3">{window._('bookManager.edit.summary')}</h2>
        {book.summary
          .replace(/<[^>]*>/g, '')
          .replace(/&nbsp;/g, '')
          .trim() ? (
          <SanitizeHTML value={book.summary} />
        ) : (
          <p className="text-secondary-foreground">{window._('bookManager.ui.noSummary')}</p>
        )}
      </section>
    </div>
  )
}
