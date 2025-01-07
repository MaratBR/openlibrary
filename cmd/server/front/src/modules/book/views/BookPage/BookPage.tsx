import { useParams } from 'react-router'
import { DefinedTagDto, useBookQuery } from '../../api/api'
import AdultIndicator from '@/components/adult-indicator'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { NavLink } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { LayoutDashboard } from 'lucide-react'
import SanitizeHtml from '@/components/sanitizer-html'
import React, { useMemo, useState } from 'react'

import './BookPage.css'
import BookCard from './BookCard'
import Tag from '../Tag'
import BookRatingCard from './BookRatingCard'
import StartReading from './StartReading'
import { Separator } from '@/components/ui/separator'
import { useTranslation } from 'react-i18next'
import BookComments from './BookComments'

export default function BookPage() {
  const { t } = useTranslation()
  const { id } = useParams<{ id: string }>()

  const { data } = useBookQuery(id)

  if (!data) {
    return null
  }

  return (
    <main className="book-page">
      <div className="book-page-grid">
        <div className="book-page-grid__lcolumn">
          <BookCard book={data} />
        </div>
        <div className="book-page-grid__rcolumn">
          <div className="book-page-header">
            <h1 className="book-title">{data.name}</h1>
            <div className="book-author">
              <NavLink className="link-default" to={`/user/${data.author.id}`}>
                {data.author.name}
              </NavLink>
            </div>
            <BookRatingCard
              bookId={data.id}
              rating={data.rating ? data.rating / 2 : null}
              votes={data.votes}
              reviews={data.reviews}
            />
          </div>

          <div className="book-summary user-content">
            <SanitizeHtml html={data.summary} />
          </div>

          <BookTags tags={data.tags} />

          <div className="book-metadata">
            {t('book.stats.short', {
              words: data.words + '',
              wordsPerChapter: data.wordsPerChapter + '',
              chapters: data.chapters.length + '',
            })}
          </div>

          <StartReading book={data} />
          <Separator className="my-4" />
          <BookComments book={data} />
        </div>
      </div>
    </main>
  )
}

function BookTags({ tags }: { tags: DefinedTagDto[] }) {
  const { t } = useTranslation()
  const [expanded, setExpanded] = useState(false)

  const shownTags = useMemo(() => (expanded ? tags : tags.slice(0, 9)), [expanded, tags])
  const canExpand = !expanded && shownTags.length < tags.length

  return (
    <ul className="book-tags">
      {shownTags.map((tag) => {
        return <Tag key={tag.id} tag={tag} />
      })}
      {canExpand && (
        <button onClick={() => setExpanded(true)} className="book-tags__more">
          {t('tags.more')}
        </button>
      )}
    </ul>
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
