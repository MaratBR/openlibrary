import { useParams } from 'react-router'
import { BookDetailsDto, useBookQuery } from '../../api'
import AdultIndicator from '@/components/adult-indicator'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { NavLink } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { LayoutDashboard } from 'lucide-react'
import ChapterCard from './ChapterCard'
import BookInfoCard from './BookInfoCard'
import BookFavoritesCounter from '../../components/book-favorites-counter'

export default function BookPage() {
  const { id } = useParams<{ id: string }>()

  const { data } = useBookQuery(id)

  return (
    <>
      {data && (
        <div className="bg-muted pb-6">
          <div className="container-default">
            <header className="page-header relative">
              <h1 className="page-header-text">
                {data.isAdult && <BookAdultIndicator />}
                {data.name}
              </h1>
              <p>
                by&nbsp;
                <NavLink className="link" to={`/user/${data.author.id}`}>
                  {data.author.name}
                </NavLink>
              </p>

              {data.permissions.canEdit && <QuickEditSection bookId={data.id} />}
            </header>
            <BookInfoCard book={data} />
          </div>
        </div>
      )}
      <main className="container-default relative">
        {data && (
          <>
            <ChaptersList book={data} />
          </>
        )}
      </main>
    </>
  )
}

function ChaptersList({ book }: { book: BookDetailsDto }) {
  return (
    <section id="chapters" className="mt-8">
      <section className="page-section">
        <BookFavoritesCounter bookId={book.id} count={book.favorites} isLiked={book.isFavorite} />
      </section>

      <section className="page-section">
        <h2 className="text-xl font-semibold">Summary</h2>

        <p>
          {book.summary ? book.summary : <span className="text-muted-foreground">No summary</span>}
        </p>
      </section>

      <h2 className="text-xl font-semibold">{book.chapters.length} chapters</h2>

      {book.chapters.length === 0 && (
        <div className="text-muted-foreground mt-3">
          It looks like author did not write anything yet
        </div>
      )}
      <div className="space-y-2 mt-4">
        {book.chapters.map((chapter) => {
          return <ChapterCard key={chapter.id} bookId={book.id} chapter={chapter} />
        })}
      </div>
    </section>
  )
}

function QuickEditSection({ bookId }: { bookId: string }) {
  return (
    <section className="absolute right-0 top-10">
      <div className="flex gap-2">
        <NavLink to={`/manager/book/${bookId}`}>
          <Button variant="outline">
            <LayoutDashboard /> Go to book manager
          </Button>
        </NavLink>
      </div>
    </section>
  )
}

function BookAdultIndicator() {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <AdultIndicator className="mr-3 relative -top-[0.2em]" />
      </TooltipTrigger>
      <TooltipContent className="max-w-64 font-text font-normal">
        This book's rating indicates it contains some degree of adult content that may not be
        suitable for children.
      </TooltipContent>
    </Tooltip>
  )
}
