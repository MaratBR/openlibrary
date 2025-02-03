/* eslint-disable no-unused-vars */
import { z } from 'zod'
import { PreactIslandProps } from '../common'
import { useEffect, useMemo, useRef, useState } from 'preact/hooks'
import { createPortal, CSSProperties } from 'preact/compat'
import BookCardPreviewContent from './BookCardPreviewContent'

const WIDTH = 400

export default function BookCardPreview({ data }: PreactIslandProps) {
  const params = useMemo(() => dataSchema.parse(data), [data])

  const ref = useRef<HTMLDivElement | null>(null)

  const [target, setTarget] = useState<
    | (BookCardPreviewTarget & {
        targetBounds: DOMRect
        style: CSSProperties
      })
    | null
  >(null)
  const targetRef = useRef(target)
  targetRef.current = target

  useEffect(() => {
    const element = params.selector ? document.querySelector(params.selector) : document.body

    if (!element) {
      console.error('[BookCardPreview] cannot find element', params.selector)
      return
    }

    const dispose = initBookCardPreview(element, (target) => {
      if (!target) {
        setTarget(null)
        return
      }

      window.requestAnimationFrame(() => {
        const rect = target.element.getBoundingClientRect()

        const style: CSSProperties = {
          top: window.scrollY + rect.top,
        }

        const OFFSET = 16

        if (rect.right + OFFSET + WIDTH > window.innerWidth - 16) {
          style.right = window.innerWidth - rect.left + 15
        } else {
          style.left = rect.right + 15
        }

        setTarget({
          ...target,
          targetBounds: rect,
          style,
        })
      })
    })
    return dispose
  }, [params.selector])

  if (!target) return

  return createPortal(
    <div
      ref={ref}
      style={{
        height: 300,
        width: WIDTH,
        ...target.style,
      }}
      class="absolute p-3 bg-background shadow-lg z-10 border rounded-xl pointer-events-none"
    >
      <BookCardPreviewContent bookId={target.bookId} />
    </div>,
    document.body,
  )
}

const dataSchema = z.object({
  selector: z.string().nullable().optional(),
})

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function debounce<T extends any[]>(
  fn: (...args: T) => void,
  delay: number,
): [(...args: T) => void, () => void] {
  let timeout: number | null = null
  const cancel = () => {
    if (timeout) {
      clearTimeout(timeout)
      timeout = null
    }
  }
  return [
    (...args: T) => {
      if (timeout) {
        clearTimeout(timeout)
      }
      timeout = window.setTimeout(() => {
        fn(...args)
      }, delay)
    },
    cancel,
  ]
}

type BookCardPreviewTarget = {
  element: HTMLElement
  bookId: string
}

function initBookCardPreview(
  element: Element,

  callback: (target: BookCardPreviewTarget | null) => void,
) {
  const cards = element.querySelectorAll('[data-book-card-preview]')

  const [debouncedCallback, cancel] = debounce(callback, 500)

  const onMouseEnter = (e: MouseEvent) => {
    cancel()
    if (e.target instanceof HTMLElement) {
      const bookId = e.target.getAttribute('data-book-card-preview')
      if (bookId) {
        callback({
          bookId,
          element: e.target,
        })
      }
    }
  }

  const onMouseLeave = () => {
    debouncedCallback(null)
  }

  cards.forEach((element) => {
    if (element instanceof HTMLElement) {
      element.addEventListener('mouseenter', onMouseEnter)
      element.addEventListener('mouseleave', onMouseLeave)
    }
  })

  return () => {
    cards.forEach((element) => {
      if (element instanceof HTMLElement) {
        element.removeEventListener('mouseenter', onMouseEnter)
        element.removeEventListener('mouseleave', onMouseLeave)
      }
    })
  }
}
