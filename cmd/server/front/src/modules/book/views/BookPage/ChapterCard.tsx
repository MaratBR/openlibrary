import { BookChapterDto } from '../../api/api'
import { NavLink } from 'react-router-dom'
import { useChapterName } from '../../utils'
import React from 'react'
import { cn } from '@/lib/utils'

export default function ChapterCard({
  chapter,
  bookId,
  onClick,
  className,
}: {
  chapter: BookChapterDto
  bookId: string
  onClick?: React.MouseEventHandler<HTMLAnchorElement>
  className?: string
}) {
  const chapterName = useChapterName(chapter.name, chapter.order)

  return (
    <NavLink
      to={`/book/${bookId}/chapters/${chapter.id}`}
      className={cn(
        'p-2 bg-card text-card-foreground block w-full hover:bg-muted group',
        className,
      )}
      onClick={onClick}
    >
      <span className="font-[500] group-hover:underline underline-offset-4">{chapterName}</span>
      &nbsp;&nbsp;&bull;&nbsp;&nbsp;
      <span className="text-sm text-muted-foreground">
        {chapter.words} words &nbsp;&nbsp;&bull;&nbsp;&nbsp; published{' '}
        {new Date(chapter.createdAt).toLocaleDateString('en-US')}
      </span>
      <p className="text-sm pt-2">
        {chapter.summary ? (
          chapter.summary
        ) : (
          <span className="text-muted-foreground">No summary</span>
        )}
      </p>
    </NavLink>
  )
}
