import AgeRatingBadge from '@/components/age-rating-badge'
import { Card } from '@/components/ui/card'
import { preloadBookQuery } from '@/modules/book/api/api'
import { NavLink } from 'react-router-dom'
import SanitizeHtml from '@/components/sanitizer-html'
import Tag from '../Tag'
import './SearchBookCard.css'
import React from 'react'
import { useQueryClient } from '@tanstack/react-query'
import BookCover from '@/modules/common/components/book-cover'
import { cn } from '@/lib/utils'
import { Heart } from 'lucide-react'
import { useCurrentUser } from '@/modules/auth/state'
import { isAgeRatingAdult } from '../../utils'
import { CensorMode } from '@/modules/account/api'
import { Button } from '@/components/ui/button'
import { BookSearchItemState } from './search-params'

export default function SearchBookCard({ book }: { book: BookSearchItemState }) {
  const [showAdultContent, setShowAdultContent] = React.useState(false)

  const currentUser = useCurrentUser()
  const queryClient = useQueryClient()

  const censorMode: CensorMode = React.useMemo(() => {
    if (showAdultContent) return 'none'

    const isAdult = isAgeRatingAdult(book.ageRating)

    if (!currentUser) {
      return isAdult ? 'censor' : 'none'
    }

    if (currentUser.bookCensoredTags.some((tag) => book.tags.some((t) => t.name === tag))) {
      return currentUser.bookCensoringMode
    }

    if (!currentUser.showAdultContent && isAdult) {
      return currentUser.bookCensoringMode
    }

    return 'none'
  }, [currentUser, book, showAdultContent])

  function handleMouseHover() {
    preloadBookQuery(queryClient, book.id)
  }

  function handleShowClick() {
    setShowAdultContent(true)
  }

  return (
    <li id={`book${book.id}`}>
      <Card
        data-censor-mode={censorMode === 'none' ? undefined : censorMode}
        className={cn('search-book-card', {
          'search-book-card--has-cover': !!book.cover,
          'search-book-card--censored': censorMode !== 'none',
        })}
      >
        {censorMode === 'censor' && (
          <div className="censor-overlay">
            <p className="text-lg">
              Censored according to your{' '}
              <NavLink className="link-default text-primary" to="/account/settings?tab=moderation">
                account settings
              </NavLink>
            </p>

            <Button className="mt-4" variant="outline2" onClick={handleShowClick}>
              Show
            </Button>
          </div>
        )}

        <div className="contents" aria-hidden={censorMode !== 'none' && !showAdultContent}>
          <div
            className={cn({
              'grid grid-cols-[300px_1fr] gap-3': !!book.cover,
            })}
          >
            {book.cover && (
              <div className="search-book-card__cover">
                <BookCover size="sm" url={book.cover} />
              </div>
            )}

            <header className="search-book-card__header">
              <AgeRatingBadge value={book.ageRating} />
              &nbsp;&nbsp;
              <span>
                <NavLink
                  className="link-default"
                  to={`/book/${book.id}`}
                  onMouseEnter={handleMouseHover}
                >
                  {book.name}
                </NavLink>{' '}
                by{' '}
                <NavLink className="link-default" to={`/user/${book.author.id}`}>
                  {book.author.name}
                </NavLink>
              </span>
              {book.tags.length > 0 && (
                <div className="search-book-card__tags">
                  {book.tags.map((tag) => (
                    <Tag key={tag.id} tag={tag} />
                  ))}
                </div>
              )}
            </header>
          </div>

          <p className="pt-3 text-sm search-book-card__summary">
            {book.summary ? (
              <SanitizeHtml html={book.summary} />
            ) : (
              <span className="text-muted-foreground">No summary</span>
            )}
          </p>

          <p className="mt-3 text-muted-foreground text-sm">
            {book.words} words &bull; {book.chapters} chapters &bull; {book.favorites}{' '}
            <Heart className="inline" size="1em" />
          </p>
        </div>
      </Card>
    </li>
  )
}
