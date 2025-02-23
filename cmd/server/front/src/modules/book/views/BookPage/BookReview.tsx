import './BookReview.css'
import { httpDeleteReview, ReviewDto } from '../../api'
import { NavLink } from 'react-router-dom'
import SanitizeHtml from '@/components/sanitizer-html'
import { useRef, useState } from 'react'
import { useResizeObserver } from 'usehooks-ts'
import { useTranslation } from 'react-i18next'
import clsx from 'clsx'
import StarRating from '@/components/star-rating'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { CheckCircle, PenBox, Trash2 } from 'lucide-react'
import { useCurrentUser } from '@/modules/auth/state'
import ReviewEditor, { ReviewData } from './ReviewEditor'
import { useUpdateReviewMutation } from './WriteReview'
import { cn, delayMs } from '@/lib/utils'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Button } from '@/components/ui/button'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { getErrorMessage } from '@/lib/errors'

type BookReviewProps = {
  review: ReviewDto
  isAuthor: boolean
  bookId: string
  isEditable?: boolean
}

export default function BookReview({ review, isAuthor, isEditable, bookId }: BookReviewProps) {
  const [isEditing, setIsEditing] = useState(false)

  const isEditableDefault = useCurrentUser()?.id === review.user.id
  if (isEditable === undefined) isEditable = isEditableDefault

  const updateReview = useUpdateReviewMutation(bookId)

  function handleStartEditing() {
    if (!isEditable) return
    setIsEditing(true)
  }

  function handleDelete() {}

  function handleStopEditing() {
    setIsEditing(false)
  }

  async function handleUpdate(_reviewData: ReviewData, review: ReviewDto | null) {
    if (!review) return
    await updateReview.mutateAsync({
      content: review.content,
      rating: review.rating,
    })
    setIsEditing(false)
  }

  return (
    <div
      className={cn('book-review', {
        'book-review--editing': isEditing,
      })}
    >
      <div className="book-review__user">
        <NavLink to={`/users/${review.user.id}`} className="size-[84px]">
          <div className="inline-block overflow-hidden rounded-full bg-muted">
            <img className="w-21 h-21" src={review.user.avatar} aria-hidden="true" />
          </div>
        </NavLink>
      </div>

      <div className="book-review__content">
        <div className="book-review__username">
          <NavLink className="link-default text-lg" to={`/users/${review.user.id}`}>
            {review.user.name}
          </NavLink>
          {isAuthor && <AuthorQuill />}
        </div>
        {isEditing ? (
          <ReviewEditor onClose={handleStopEditing} onUpdated={handleUpdate} review={review} />
        ) : (
          <>
            <StarRating className="mb-2" size="2rem" value={review.rating / 2} />
            <ReviewContent content={review.content} />
          </>
        )}
      </div>

      {!isEditing && (
        <div className="book-review__actions">
          {isEditable && (
            <>
              <button onClick={handleStartEditing} className="book-review__action">
                <PenBox size={16} />
              </button>
              <DeleteReview bookId={bookId} />
            </>
          )}
        </div>
      )}
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
    <div className="text-sm">
      <div
        className={clsx('__user-content', {
          contents: expanded,
          'overflow-hidden max-h-[100px]': !expanded,
        })}
      >
        <SanitizeHtml ref={rootEl} html={content} />
      </div>
      {canBeExpanded && (
        <button
          onClick={() => setExpanded((x) => !x)}
          className="font-[500] underline-offset-2 block hover:underline relative after:bg-highlight after:inset-[-8px] after:rounded-lg after:absolute after:hidden hover:after:block"
        >
          {expanded ? t('common.less') : t('common.more')}
        </button>
      )}
    </div>
  )
}

function DeleteReview({ bookId }: { bookId: string }) {
  const [open, setOpen] = useState(false)
  const { t } = useTranslation()

  const queryClient = useQueryClient()

  const deleteReview = useMutation({
    mutationFn: async () => {
      const delayPromise = delayMs(350)
      try {
        await Promise.all([httpDeleteReview(bookId), delayPromise])
        toast(
          <div className="flex items-center gap-2">
            <CheckCircle size={16} />
            <span>{t('book.review.delete.deleted')}</span>
          </div>,
        )
        queryClient.setQueryData(['book', bookId, 'reviews', 'my'], null)
      } catch (e: unknown) {
        await delayPromise
        throw e
      }
    },
    onError(error) {
      toast(
        <div className="flex items-center gap-2">
          <span>{getErrorMessage(error)}</span>
        </div>,
      )
    },
  })

  return (
    <Popover modal open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <button
          onClick={() => setOpen(true)}
          className="book-review__action book-review__action--delete"
        >
          <Trash2 size={16} />
        </button>
      </PopoverTrigger>
      <PopoverContent align="end" className="border-destructive/50 border-2 pt-3">
        <h3 className="text-lg font-[500]">{t('book.review.delete.title')}</h3>
        <p className="text-sm">{t('book.review.delete.description')}</p>
        <Button
          disabled={deleteReview.isPending}
          onClick={() => deleteReview.mutate()}
          variant="ghost"
          className="rounded-lg -ml-1 mt-3 p-1 h-auto bg-destructive/10 hover:text-destructive hover:bg-destructive/20"
        >
          <Trash2 />
          {deleteReview.isPending
            ? t('book.review.delete.deleting')
            : t('book.review.delete.delete')}
        </Button>
      </PopoverContent>
    </Popover>
  )
}
