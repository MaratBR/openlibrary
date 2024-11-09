import { Button } from '@/components/ui/button'
import { useBookManager, useBookManagerChaptersQuery } from './book-manager-context'
import { ListEnd, PenIcon } from 'lucide-react'
import Spinner from '@/components/spinner'
import { NavLink } from 'react-router-dom'
import BookChaptersList from './BookChaptersList'

export default function BookChapters() {
  const { book } = useBookManager()
  const { data, isLoading } = useBookManagerChaptersQuery()

  return (
    <section className="page-section">
      <div className="space-x-2">
        <NavLink to={`/manager/book/${book.id}/new-chapter`}>
          <Button variant="outline">
            <PenIcon />
            Write new chapter
          </Button>
        </NavLink>

        <NavLink to={`/manager/book/${book.id}/reorder-chapters`}>
          <Button variant="ghost">
            <ListEnd /> Reorder
          </Button>
        </NavLink>
      </div>

      {isLoading && <Spinner />}

      {data && (
        <>
          {data.length === 0 && (
            <p className="my-4 text-muted-foreground">
              No chapters yet. Write a first one by clicking above.
            </p>
          )}
          <div className="space-y-2 mt-2">
            <BookChaptersList chapters={data} />
          </div>
        </>
      )}
    </section>
  )
}
