import AgeRatingBadge from '@/components/age-rating-badge'
import { Card } from '@/components/ui/card'
import { ManagerAuthorBookDto } from '@/modules/book/api/api'
import { NavLink } from 'react-router-dom'
import BookIsHiddenIndicator from '../common/book-is-hidden-indicator'
import BookIsBannedIndicator from '../common/book-is-banned-indicator'
import GoToBookPage from '../common/go-to-book-page'
import SanitizeHtml from '@/components/sanitizer-html'

export default function AuthorBookCard({ book }: { book: ManagerAuthorBookDto }) {
  return (
    <Card className="rounded-sm p-2">
      <header>
        <AgeRatingBadge value={book.ageRating} />
        &nbsp;&nbsp;
        <span className="text-muted-foreground">
          <NavLink className="link-default" to={`/manager/book/${book.id}`}>
            {book.name}
          </NavLink>
          &nbsp;&bull; {book.words} words
        </span>
      </header>

      <p className="pt-3 text-sm">
        {book.summary ? (
          <SanitizeHtml html={book.summary} />
        ) : (
          <span className="text-muted-foreground">No summary</span>
        )}
      </p>

      <div className="flex gap-2 mt-3">
        <GoToBookPage bookId={book.id} />
        {!book.isPubliclyVisible && <BookIsHiddenIndicator />}
        {book.isBanned && <BookIsBannedIndicator bookId={book.id} />}
      </div>
    </Card>
  )
}
