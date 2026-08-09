import { httpGetReview, httpUpdateReview, ReviewDto } from './api'
import { MouseEvent, useEffect, useRef, useState } from 'react'
import { ReactIslandProps } from '../common/react-island'
import { RichTextInput, useOLEditor } from '@/components/rte'

export default function ReviewEditor({ rootElement }: ReactIslandProps) {
  const [review] = useState(getExistingReviewData)
  const editor = useOLEditor({ content: review?.content ?? '' })
  const [rating, setRating] = useState(review?.rating ?? 0)
  const [saving, setSaving] = useState(false)
  const bookId = getBookId()

  useEffect(() => {
    httpGetReview(bookId).then((resp) => {
      if (resp.data.rating) setRating(resp.data.rating)
    })
  }, [bookId])

  const formValid = rating !== 0

  function handleSave() {
    if (saving || !formValid || !editor) return
    setSaving(true)

    httpUpdateReview(bookId, {
      content: editor.getHTML(),
      rating,
    })
      .then((review) => {
        rootElement.dispatchEvent(new CustomEvent('review:updated', { detail: review }))
        window.toast({ title: window._('reviews.updated') })
      })
      .finally(() => {
        setSaving(false)
      })
  }

  return (
    <div>
      <RatingInput scale={0.5} value={rating} onInput={setRating} />

      {editor && <RichTextInput editor={editor} className="ol-simple-editor--review mt-4" />}

      <button
        className="btn btn--lg btn--primary mt-3"
        onClick={() => handleSave()}
        disabled={!formValid}
      >
        {window._('common.save')}
      </button>
    </div>
  )
}

function RatingInput({
  scale = 1,
  value,
  onInput,
}: {
  scale?: number
  value: number

  onInput: (value: number) => void
}) {
  const rootElement = useRef<HTMLDivElement | null>(null)
  const disableHalfPoints = true

  function handleClick(event: MouseEvent) {
    if (!rootElement) return

    const rect = rootElement.current!.getBoundingClientRect()
    const x = event.clientX - rect.left
    const width = rect.width
    let newValue = Math.max(Math.min(Math.ceil((x / width) * 10), 10), 1)

    if (disableHalfPoints && newValue % 2 === 1) {
      newValue += 1
    }

    if (newValue !== value) {
      onInput(newValue)
    }
  }

  return (
    <div
      ref={rootElement}
      className="relative cursor-pointer"
      onClick={handleClick}
      style={{ width: 540 * scale, height: 100 * scale }}
    >
      <div className="star-background h-full w-full opacity-15" />
      <div
        className="absolute left-0 top-0 star-background star-background--filled h-full"
        style={{ width: `${calcPerc(value)}%`, backgroundSize: `auto ${scale * 100}px` }}
      />
    </div>
  )
}

function calcPerc(value: number): number {
  return ((500 * (value / 10) + Math.floor(value / 2) * 10) / 540) * 100
}

/**
 * Finds a hidden element containing review data in JSON format and parses it
 * into a ReviewDto object. If the element is not found,
 * null is returned.
 *
 * This function is used to pre-fill the review text area with existing review
 * data when the review editor is rehydrated on the server.
 *
 * @returns {ReviewDto | null}
 */
function getExistingReviewData(): ReviewDto | null {
  const el = document.getElementById('island-review-editor-data')

  if (el instanceof HTMLTemplateElement) {
    return JSON.parse(el.content.textContent || '') as ReviewDto
  }

  return null
}

function getBookId() {
  const v = window.__server__?.bookId
  if (typeof v === 'string' && v) {
    return v
  }

  throw new Error('could not find bookId, __server__.bookId is not set')
}
