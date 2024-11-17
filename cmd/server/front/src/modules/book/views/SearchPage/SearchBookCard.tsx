import AgeRatingBadge from '@/components/age-rating-badge'
import { Card } from '@/components/ui/card'
import { BookSearchItem, preloadBookQuery } from '@/modules/book/api'
import { NavLink } from 'react-router-dom'
import SanitizeHtml from '@/components/sanitizer-html'
import Tag from '../Tag'
import './SearchBookCard.css'
import React from 'react'
import { useQueryClient } from '@tanstack/react-query'

export default function SearchBookCard({ book }: { book: BookSearchItem }) {
  const queryClient = useQueryClient()

  function handleMouseHover() {
    preloadBookQuery(queryClient, book.id)
  }

  return (
    <Card className="search-book-card">
      <header>
        <AgeRatingBadge value={book.ageRating} />
        &nbsp;&nbsp;
        <span>
          <NavLink className="link-default" to={`/book/${book.id}`} onMouseEnter={handleMouseHover}>
            {book.name}
          </NavLink>{' '}
          by{' '}
          <NavLink className="link-default" to={`/user/${book.author.id}`}>
            {book.author.name}
          </NavLink>
        </span>
      </header>

      {book.tags.length > 0 && (
        <div className="search-book-card__tags">
          <span className="search-book-card__tags-label">Tags:</span>
          {book.tags.map((tag) => (
            <Tag key={tag.id} tag={tag} />
          ))}
        </div>
      )}

      <p className="pt-3 text-sm search-book-card__summary">
        {book.summary ? (
          <SanitizeHtml html={book.summary} />
        ) : (
          <span className="text-muted-foreground">No summary</span>
        )}
      </p>

      <p className="mt-3 text-muted-foreground text-sm">
        {book.words} words &bull; {book.chapters} chapters &bull; {book.favorites} favorites
      </p>
    </Card>
  )
}
