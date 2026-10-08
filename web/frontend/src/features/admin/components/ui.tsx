import { DashboardLoader } from '@/components/dashboard-loader'
import { FormEvent, ReactNode, useEffect, useLayoutEffect, useRef, useState } from 'react'
import { createPortal } from 'react-dom'
import { Link, useBlocker, useRevalidator, useRouteError, useSearchParams } from 'react-router'
import { AdminAPIError } from '@/features/admin/api'

export const t = (key: string, args?: Record<string, string>) => window._(key, args || {})
export function PageHeader({
  title,
  description,
  actions,
}: {
  title: string
  description?: string
  actions?: ReactNode
}) {
  useEffect(() => {
    document.title = `${title} · ${t('admin.ui.title')}`
  }, [title])
  return (
    <header className="Admin-pageHeader">
      <div>
        <h1 tabIndex={-1}>{title}</h1>
        {description && <p>{description}</p>}
      </div>
      {actions && <div className="Admin-actions">{actions}</div>}
    </header>
  )
}
export function Card({
  title,
  description,
  children,
  className = '',
}: {
  title?: string
  description?: string
  children: ReactNode
  className?: string
}) {
  return (
    <section className={`admin-card ${className}`}>
      {title && (
        <header className="admin-card-header">
          <h2>{title}</h2>
          {description && <p>{description}</p>}
        </header>
      )}
      {children}
    </section>
  )
}
export function Field({
  id,
  label,
  children,
  help,
}: {
  id: string
  label: string
  children: ReactNode
  help?: string
}) {
  return (
    <div className="Admin-field">
      <label htmlFor={id}>{label}</label>
      {children}
      {help && (
        <p id={`${id}-help`} className="Admin-help">
          {help}
        </p>
      )}
    </div>
  )
}
export function Notice({ children, error = false }: { children: ReactNode; error?: boolean }) {
  return (
    <p
      className={`Admin-notice ${error ? 'Admin-notice--error' : ''}`}
      role={error ? 'alert' : 'status'}
    >
      {children}
    </p>
  )
}
export function errorMessage(error: unknown) {
  if (error instanceof AdminAPIError) {
    if (error.status === 401) return t('admin.ui.sessionExpired')
    if (error.status === 403) return t('admin.ui.forbidden')
    if (error.status === 400) return t('admin.ui.validationError')
    if (error.status === 404) return t('admin.ui.notFoundDescription')
  }
  return t('admin.ui.requestError')
}
export function RouteError() {
  const error = useRouteError()
  const revalidator = useRevalidator()
  return (
    <>
      <PageHeader
        title={t(
          error instanceof AdminAPIError && error.status === 404
            ? 'admin.ui.notFound'
            : 'admin.ui.unavailable',
        )}
      />
      <Card>
        <Notice error>
          {error instanceof AdminAPIError ? errorMessage(error) : t('admin.ui.loadError')}
        </Notice>
        <div className="Admin-actions">
          {error instanceof AdminAPIError && error.status === 401 ? (
            <a className="Btn Btn--primary" href="/admin/login">
              {t('admin.ui.signIn')}
            </a>
          ) : (
            <button
              className="Btn Btn--primary"
              disabled={revalidator.state !== 'idle'}
              onClick={() => revalidator.revalidate()}
            >
              {t('admin.ui.retry')}
            </button>
          )}
          <Link to="/" className="Btn Btn--outline">
            {t('admin.ui.backHome')}
          </Link>
        </div>
      </Card>
    </>
  )
}
export function Loading() {
  return <DashboardLoader label={t('admin.ui.loading')} />
}

export function Pager({ page, total }: { page: number; total: number }) {
  const [params] = useSearchParams()
  if (total <= 1) return null
  const href = (p: number) => {
    const q = new URLSearchParams(params)
    q.set('p', String(p))
    return `?${q}`
  }
  return (
    <nav
      className="Admin-pager"
      aria-label={t('admin.ui.page', { page: String(page), total: String(total) })}
    >
      {page > 1 ? (
        <Link className="Btn Btn--outline" to={href(page - 1)}>
          {t('admin.ui.previous')}
        </Link>
      ) : (
        <span />
      )}
      <span>{t('admin.ui.page', { page: String(page), total: String(total) })}</span>
      {page < total ? (
        <Link className="Btn Btn--outline" to={href(page + 1)}>
          {t('admin.ui.next')}
        </Link>
      ) : (
        <span />
      )}
    </nav>
  )
}
export function Empty({ kind }: { kind: 'users' | 'tags' }) {
  return (
    <div className="Admin-empty">
      <i className={`fa-solid ${kind === 'users' ? 'fa-users' : 'fa-tags'}`} aria-hidden="true" />
      <h2>{t(kind === 'users' ? 'admin.ui.emptyUsers' : 'admin.ui.emptyTags')}</h2>
      <Link to={`/${kind}`} className="Btn Btn--outline">
        {t('admin.ui.clearFilters')}
      </Link>
    </div>
  )
}
export function Flag({ value }: { value: boolean }) {
  return (
    <span className={`Admin-badge ${value ? 'Admin-badge--accent' : ''}`}>
      {t(value ? 'admin.ui.yes' : 'admin.ui.no')}
    </span>
  )
}
export function Dialog({
  title,
  children,
  onClose,
}: {
  title: string
  children: ReactNode
  onClose: () => void
}) {
  const ref = useRef<HTMLDialogElement>(null)
  useLayoutEffect(() => {
    const dialog = ref.current!
    const trigger = document.activeElement as HTMLElement | null
    dialog.showModal()
    return () => {
      dialog.close()
      trigger?.focus()
    }
  }, [])
  return createPortal(
    <dialog
      ref={ref}
      className="Admin-dialog"
      aria-label={title}
      onCancel={(e) => {
        e.preventDefault()
        onClose()
      }}
    >
      <header>
        <h2>{title}</h2>
        <button
          type="button"
          className="Btn Btn--ghost Btn--icon"
          aria-label={t('common.close')}
          onClick={onClose}
        >
          ×
        </button>
      </header>
      {children}
    </dialog>,
    document.body,
  )
}
export function useDirtyForm(dirty: boolean, pending: boolean) {
  const blocker = useBlocker(dirty || pending)
  const [external, setExternal] = useState<string | null>(null)
  const allowUnload = useRef(false)
  useEffect(() => {
    if (!dirty && !pending) return
    const beforeUnload = (event: BeforeUnloadEvent) => {
      if (!allowUnload.current) {
        event.preventDefault()
        event.returnValue = ''
      }
    }
    const click = (event: MouseEvent) => {
      if (
        event.defaultPrevented ||
        event.button !== 0 ||
        event.metaKey ||
        event.ctrlKey ||
        event.shiftKey ||
        event.altKey
      )
        return
      const anchor =
        event.target instanceof Element
          ? event.target.closest<HTMLAnchorElement>('a[data-admin-external]')
          : null
      if (anchor) {
        event.preventDefault()
        if (!pending) setExternal(anchor.href)
      }
    }
    window.addEventListener('beforeunload', beforeUnload)
    document.addEventListener('click', click)
    return () => {
      window.removeEventListener('beforeunload', beforeUnload)
      document.removeEventListener('click', click)
    }
  }, [dirty, pending])
  function stay() {
    setExternal(null)
    if (blocker.state === 'blocked') blocker.reset()
  }
  return blocker.state === 'blocked' || external ? (
    <Dialog title={t('admin.ui.unsaved')} onClose={stay}>
      <p>{t(pending ? 'admin.ui.saving' : 'admin.ui.discardDescription')}</p>
      <div className="Admin-actions">
        <button className="Btn Btn--primary" type="button" onClick={stay}>
          {t('admin.ui.stay')}
        </button>
        <button
          className="Btn Btn--outline"
          type="button"
          disabled={pending}
          onClick={() => {
            if (blocker.state === 'blocked') blocker.proceed()
            if (external) {
              allowUnload.current = true
              window.location.assign(external)
            }
          }}
        >
          {t('admin.ui.discard')}
        </button>
      </div>
    </Dialog>
  ) : null
}

// Keep errors/input local to the form and guard even same-tick repeated submits.
export function useSave() {
  const [pending, setPending] = useState(false)
  const [error, setError] = useState<unknown>()
  const [saved, setSaved] = useState(false)
  const locked = useRef(false)
  const mounted = useRef(true)
  useEffect(() => {
    mounted.current = true
    return () => {
      mounted.current = false
    }
  }, [])
  async function submit(event: FormEvent, operation: () => Promise<void>) {
    event.preventDefault()
    if (locked.current) return
    locked.current = true
    setPending(true)
    setSaved(false)
    setError(undefined)
    try {
      await operation()
      if (mounted.current) setSaved(true)
    } catch (e) {
      if (mounted.current) setError(e)
    } finally {
      locked.current = false
      if (mounted.current) setPending(false)
    }
  }
  return {
    pending,
    error,
    saved,
    submit,
    changed: () => {
      setSaved(false)
      setError(undefined)
    },
  }
}
export function SaveBar({
  state,
  cancelTo,
  cancelState,
}: {
  state: ReturnType<typeof useSave>
  cancelTo: string
  cancelState?: { back: string }
}) {
  return (
    <footer className="Admin-saveBar">
      <div aria-live="polite">
        {state.error ? (
          <Notice error>{errorMessage(state.error)}</Notice>
        ) : state.saved ? (
          <Notice>{t('admin.ui.saved')}</Notice>
        ) : null}
      </div>
      <div className="Admin-actions">
        <Link to={cancelTo} state={cancelState} className="Btn Btn--outline">
          {t('common.cancel')}
        </Link>
        <button className="Btn Btn--primary" type="submit" disabled={state.pending}>
          {t(state.pending ? 'admin.ui.saving' : 'common.save')}
        </button>
      </div>
    </footer>
  )
}
