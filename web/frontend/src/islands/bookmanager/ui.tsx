import { ReactNode, useEffect, useLayoutEffect, useRef, useState } from 'react'
import { createPortal } from 'react-dom'
import { NavLink, useBlocker } from 'react-router'
import { BookCover } from '@/components/BookCover'
import type { BookCover as Cover } from '@/api/common'

export function usePageTitle(title: string) {
  useEffect(() => {
    document.title = `${title} · ${window._('bookManager.title')}`
  }, [title])
}

export function BookNavigation({
  id,
  active,
}: {
  id: string
  active: 'general' | 'details' | 'chapters'
}) {
  return (
    <nav className="BM-tabs" aria-label={window._('bookManager.ui.bookNavigation')}>
      {(['general', 'details', 'chapters'] as const).map((tab) => (
        <NavLink
          key={tab}
          to={tab === 'details' ? `/books/${id}/edit` : `/books/${id}?t=${tab}`}
          aria-current={active === tab ? 'page' : false}
        >
          {window._(`bookManager.ui.${tab}`)}
        </NavLink>
      ))}
    </nav>
  )
}

export function BackToBooks() {
  return (
    <NavLink to="/books" className="Link inline-flex my-4">
      ← {window._('bookManager.ui.back')}
    </NavLink>
  )
}

export function CoverImage({ cover }: { cover: Cover }) {
  const [failedUrl, setFailedUrl] = useState<string | null>(null)
  return (
    <div className="BM-cover" aria-hidden="true" onErrorCapture={() => setFailedUrl(cover?.url)}>
      {cover?.url && failedUrl !== cover.url ? (
        <BookCover cover={cover} />
      ) : (
        <i className="fa-solid fa-book" />
      )}
    </div>
  )
}

export function Counts({ chapters, words }: { chapters: number; words: number }) {
  return (
    <span>
      {window._(chapters === 1 ? 'bookManager.ui.chapterCountOne' : 'bookManager.ui.chapterCount', {
        count: chapters.toLocaleString(),
      })}{' '}
      ·{' '}
      {window._(words === 1 ? 'bookManager.ui.wordCountOne' : 'bookManager.ui.wordCount', {
        count: words.toLocaleString(),
      })}
    </span>
  )
}

// Native modal dialogs provide focus containment, inert background, Escape and focus restoration.
export function ManagerDialog({
  title,
  children,
  onClose,
  panel = false,
}: {
  title: string
  children: ReactNode
  onClose: () => void
  panel?: boolean
}) {
  const ref = useRef<HTMLDialogElement>(null)
  useLayoutEffect(() => {
    const dialog = ref.current!
    const trigger = document.activeElement as HTMLElement | null
    dialog.showModal()
    dialog.querySelector<HTMLElement>('[data-dialog-focus]')?.focus()
    return () => {
      dialog.close()
      trigger?.focus()
    }
  }, [])
  return createPortal(
    <dialog
      ref={ref}
      className={`BM-dialog ${panel ? 'BM-dialog--panel' : ''}`}
      aria-label={title}
      onCancel={(e) => {
        e.preventDefault()
        onClose()
      }}
      onClick={(e) => {
        if (e.target === e.currentTarget) {
          const r = e.currentTarget.getBoundingClientRect()
          if (
            e.clientX < r.left ||
            e.clientX > r.right ||
            e.clientY < r.top ||
            e.clientY > r.bottom
          )
            onClose()
        }
      }}
    >
      {children}
    </dialog>,
    document.body,
  )
}

export function useUnsavedChanges(dirty: boolean, pending = false) {
  const allowUnload = useRef(false)
  const blocker = useBlocker(dirty || pending)
  const [action, setAction] = useState<(() => void) | null>(null)
  useEffect(() => {
    if (!dirty && !pending) return
    const prevent = (event: BeforeUnloadEvent) => {
      if (!allowUnload.current) {
        event.preventDefault()
        event.returnValue = ''
      }
    }
    window.addEventListener('beforeunload', prevent)
    return () => window.removeEventListener('beforeunload', prevent)
  }, [dirty, pending])
  // Match guarded editor links: modified primary clicks also use Stay/Discard.
  // Capture only the shell's document link; router links keep using useBlocker.
  useEffect(() => {
    if (!dirty && !pending) return
    const followSiteLink = (event: MouseEvent) => {
      if (event.defaultPrevented || event.button !== 0) return
      const anchor =
        event.target instanceof Element
          ? event.target.closest<HTMLAnchorElement>('a[data-manager-site-link]')
          : null
      if (!anchor) return
      event.preventDefault()
      if (!pending)
        setAction(() => () => {
          window.location.href = anchor.href
        })
    }
    document.addEventListener('click', followSiteLink)
    return () => document.removeEventListener('click', followSiteLink)
  }, [dirty, pending])
  const stay = () => {
    setAction(null)
    if (blocker.state === 'blocked') blocker.reset()
  }
  return {
    allowUnload: () => {
      allowUnload.current = true
    },
    leave: (next: () => void) => {
      if (pending) return
      if (dirty) setAction(() => next)
      else next()
    },
    prompt:
      action || blocker.state === 'blocked' ? (
        <ManagerDialog title={window._('bookManager.ui.unsaved')} onClose={stay}>
          <h2>{window._('bookManager.ui.unsaved')}</h2>
          <p className="my-4">
            {window._(pending ? 'bookManager.ui.saving' : 'bookManager.ui.discardDescription')}
          </p>
          <div className="flex gap-2">
            <button type="button" className="Btn Btn--primary" onClick={stay}>
              {window._('bookManager.ui.stay')}
            </button>
            <button
              type="button"
              className="Btn Btn--outline"
              disabled={pending}
              onClick={() => {
                const next = action
                setAction(null)
                if (blocker.state === 'blocked') blocker.proceed()
                if (next) {
                  allowUnload.current = true
                  next()
                }
              }}
            >
              {window._('bookManager.ui.discard')}
            </button>
          </div>
        </ManagerDialog>
      ) : null,
  }
}

export function RequestError() {
  return (
    <p role="alert" className="my-3 text-destructive">
      {window._('bookManager.ui.requestError')}
    </p>
  )
}

// Scroll focused fields above the sticky bar without moving the buttons during clicks.
export function keepFormFocusVisible(element: HTMLElement | null) {
  requestAnimationFrame(() => {
    if (!element) return
    const bar = element.closest('form')?.querySelector<HTMLElement>('.BM-actionBar')
    if (!bar || bar.contains(element)) return
    const field = element.getBoundingClientRect()
    const footer = bar.getBoundingClientRect()
    if (field.bottom > footer.top && field.top < footer.bottom) {
      window.scrollBy({ top: field.bottom - footer.top + 24 })
    }
  })
}
