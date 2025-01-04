import './BookReview.css'
import { ReviewDto } from '../../api'
import { NavLink } from 'react-router-dom'
import SanitizeHtml from '@/components/sanitizer-html'
import { useRef, useState } from 'react'
import { useResizeObserver } from 'usehooks-ts'
import { useTranslation } from 'react-i18next'
import clsx from 'clsx'
import StarRating from '@/components/star-rating'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'

export default function BookReview({ review, isAuthor }: { review: ReviewDto; isAuthor: boolean }) {
  return (
    <div className="book-review">
      <div className="book-review__user">
        <div className="book-review__user__avatar">
          <div className="inline-block overflow-hidden rounded-full bg-muted">
            <img className="w-21 h-21" src={review.user.avatar} aria-hidden="true" />
          </div>
        </div>
      </div>

      <div className="book-review__content">
        <div className="book-review__username">
          <NavLink className="link-default text-lg" to={`/user/${review.user.id}`}>
            {review.user.name}
          </NavLink>
          {isAuthor && <AuthorQuill />}
        </div>
        <StarRating className="mb-2" value={review.rating / 2} />
        <ReviewContent content={review.content} />
      </div>
    </div>
  )
}

function AuthorQuill() {
  const { t } = useTranslation()

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
          <g
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="1.5"
            color="currentColor"
          >
            <path d="M5.076 17C4.089 4.545 12.912 1.012 19.973 2.224c.286 4.128-1.734 5.673-5.58 6.387c.742.776 2.055 1.753 1.913 2.974c-.1.868-.69 1.295-1.87 2.147C11.85 15.6 8.854 16.78 5.076 17" />
            <path d="M4 22c0-6.5 3.848-9.818 6.5-12" />
          </g>
        </svg>
      </TooltipTrigger>
      <TooltipContent>{t('common.authorQuill.tooltip')}</TooltipContent>
    </Tooltip>
  )
}

function ReviewContent({ content }: { content: string }) {
  const { t } = useTranslation()
  const [expanded, setExpanded] = useState(false)
  const [canBeExpanded, setCaBeExpanded] = useState(false)
  const rootEl = useRef<HTMLDivElement | null>(null)

  useResizeObserver({
    ref: rootEl,
    onResize: (entry) => {
      if (entry.height) setCaBeExpanded(entry.height > 100)
    },
  })

  return (
    <div>
      <div
        className={clsx({
          contents: expanded,
          'overflow-hidden max-h-[100px]': !expanded,
        })}
      >
        <SanitizeHtml ref={rootEl} html={content} />
      </div>
      {canBeExpanded && (
        <button
          onClick={() => setExpanded((x) => !x)}
          className="font-[500] hover:underline underline-offset-2"
        >
          {expanded ? t('common.less') : t('common.more')}
        </button>
      )}
    </div>
  )
}
