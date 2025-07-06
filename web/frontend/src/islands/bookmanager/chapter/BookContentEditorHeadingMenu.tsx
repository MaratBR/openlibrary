import { PopupController } from '@/lib/popup-anchor'
import { createPortal } from 'preact/compat'
import { useEffect, useRef, useState } from 'preact/hooks'

type Heading = number | null

const HEADINGS: {
  key: string
  size: string
}[] = [
  {
    key: 'editor.p',
    size: '1em',
  },
  {
    key: 'editor.h1',
    size: '2em',
  },
  {
    key: 'editor.h2',
    size: '1.5em',
  },
  {
    key: 'editor.h3',
    size: '1.17em',
  },
  {
    key: 'editor.h4',
    size: '1em',
  },
  {
    key: 'editor.h5',
    size: '0.83em',
  },
  {
    key: 'editor.h6',
    size: '0.67em',
  },
]

export default function BookContentEditorHeadingMenu({
  heading,
  onChange,
}: {
  heading: Heading
  onChange: (heading: Heading) => void
}) {
  const [open, setOpen] = useState(false)

  const popupRef = useRef<HTMLDivElement | null>(null)
  const buttonRef = useRef<HTMLButtonElement | null>(null)

  const popupController = useRef<PopupController | null>(null)

  useEffect(() => {
    if (!popupRef.current) throw new Error('popupRef is null')
    if (!buttonRef.current) throw new Error('buttonRef is null')
    popupController.current = new PopupController(buttonRef.current, popupRef.current, {
      anchorOrigin: { vertical: 'bottom', horizontal: 'left' },
      targetOrigin: { vertical: 'top', horizontal: 'left' },
    })
    return () => {
      popupController.current?.dispose()
    }
  }, [])

  useEffect(() => {
    if (!popupController.current || !open) return

    popupController.current.update()
  }, [open])

  useEffect(() => {
    if (!open) return

    const callback = (event: MouseEvent) => {
      window.requestAnimationFrame(() => {
        if (
          event.target instanceof Element &&
          !event.target.closest('.ol-card') &&
          // only close if we clicked at element that currently exists in DOM
          document.body.contains(event.target)
        ) {
          setOpen(false)
        }
      })
    }

    document.addEventListener('click', callback)

    return () => {
      document.removeEventListener('click', callback)
    }
  }, [open])

  const items = HEADINGS.map((h, i) => (
    <li
      class="px-2 min-h-8 text-left flex items-center hover:bg-highlight cursor-pointer font-book"
      style={{ height: '1.5em', fontSize: h.size }}
      role="listitem"
      key={i}
      onClick={() => {
        onChange(i === 0 ? null : i)
        setOpen(false)
      }}
    >
      {window._(h.key)}
    </li>
  ))

  const displayed = heading === null ? HEADINGS[0] : HEADINGS[heading]

  return (
    <button
      ref={buttonRef}
      onClick={() => setOpen((x) => !x)}
      class="btn btn--ghost font-book text-lg relative text-foreground"
    >
      <span>{window._(displayed.key)}</span>
      {createPortal(
        <div
          ref={popupRef}
          style={{ display: open ? 'block' : 'none' }}
          class="ol-card border-none rounded-none shadow-md p-0 fixed top-full left-0 z-10"
        >
          <ul>{items}</ul>
        </div>,
        document.body,
      )}
    </button>
  )
}
