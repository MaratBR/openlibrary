import { useRef, useState } from 'react'
import { Link, LoaderFunctionArgs, useLoaderData } from 'react-router'
import { loadDebug, save } from '@/features/admin/api'
import {
  Card,
  Dialog,
  Field,
  Notice,
  PageHeader,
  errorMessage,
  t,
} from '@/features/admin/components/ui'

const sections = [
  {
    path: '/users',
    title: 'admin.sidebar.users',
    description: 'admin.ui.usersDescription',
    icon: 'fa-users',
  },
  {
    path: '/tags',
    title: 'admin.sidebar.tags',
    description: 'admin.ui.tagsDescription',
    icon: 'fa-tags',
  },
  {
    path: '/books',
    title: 'admin.sidebar.books',
    description: 'admin.ui.booksDescription',
    icon: 'fa-book-open',
  },
  {
    path: '/debug',
    title: 'admin.sidebar.debugActions',
    description: 'admin.ui.debugDescription',
    icon: 'fa-wrench',
  },
]
export function Home() {
  return (
    <>
      <PageHeader title={t('admin.ui.homeTitle')} description={t('admin.ui.homeDescription')} />
      <div className="Admin-homeGrid">
        {sections.map((section) => (
          <Link
            className="admin-card admin-card--interactive Admin-sectionCard"
            key={section.path}
            to={section.path}
          >
            <span className="Admin-sectionIcon">
              <i className={`fa-solid ${section.icon}`} aria-hidden="true" />
            </span>
            <h2>{t(section.title)}</h2>
            <p>{t(section.description)}</p>
            <span className="Admin-sectionLink">
              {t('admin.ui.openSection')} <span aria-hidden="true">→</span>
            </span>
          </Link>
        ))}
      </div>
    </>
  )
}
export function Books() {
  const [id, setID] = useState('')
  const [error, setError] = useState(false)
  return (
    <>
      <PageHeader title={t('admin.books.title')} description={t('admin.ui.booksDescription')} />
      <Card title={t('admin.ui.bookLookup')} className="Admin-lookupCard">
        <form
          onSubmit={(event) => {
            event.preventDefault()
            if (!/^[1-9]\d*$/.test(id) || BigInt(id) > 9223372036854775807n) {
              setError(true)
              return
            }
            window.location.assign(`/book/${id}`)
          }}
        >
          <Field id="book-id" label={t('admin.ui.bookID')} help={t('admin.ui.bookIDHelp')}>
            <input
              className="input"
              id="book-id"
              inputMode="numeric"
              required
              value={id}
              aria-invalid={error}
              aria-describedby={error ? 'book-id-error book-id-help' : 'book-id-help'}
              onChange={(event) => {
                setID(event.target.value)
                setError(false)
              }}
            />
          </Field>
          {error && (
            <p id="book-id-error" role="alert" className="Admin-notice Admin-notice--error">
              {t('admin.ui.bookIDInvalid')}
            </p>
          )}
          <button type="submit" className="Btn Btn--primary">
            {t('common.open')}
          </button>
        </form>
      </Card>
    </>
  )
}
export function debugLoader({ request }: LoaderFunctionArgs) {
  return loadDebug(request.signal)
}
export function Debug() {
  const actions = useLoaderData<typeof debugLoader>()
  const [confirm, setConfirm] = useState(false)
  const [pending, setPending] = useState(false)
  const [scheduled, setScheduled] = useState(false)
  const [error, setError] = useState<unknown>()
  const locked = useRef(false)
  async function run() {
    if (locked.current) return
    locked.current = true
    setPending(true)
    setError(undefined)
    setScheduled(false)
    try {
      await save('/debug', { act: 'books:elastic:reindex' })
      setConfirm(false)
      setScheduled(true)
    } catch (e) {
      setError(e)
    } finally {
      locked.current = false
      setPending(false)
    }
  }
  return (
    <>
      <PageHeader
        title={t('admin.sidebar.debugActions')}
        description={t('admin.ui.debugDescription')}
      />
      {actions.includes('books:elastic:reindex') ? (
        <Card title={t('admin.ui.reindex')} description={t('admin.ui.reindexDescription')}>
          <button
            className="Btn Btn--outline"
            type="button"
            disabled={pending}
            onClick={() => {
              setConfirm(true)
              setError(undefined)
            }}
          >
            {t('admin.ui.run')}
          </button>
          {scheduled && <Notice>{t('admin.ui.scheduled')}</Notice>}
        </Card>
      ) : (
        <Card>
          <p>{t('admin.ui.noActions')}</p>
        </Card>
      )}
      {confirm && (
        <Dialog
          title={t('admin.ui.confirmAction')}
          onClose={() => {
            if (!pending) setConfirm(false)
          }}
        >
          <p>{t('admin.ui.confirmDescription')}</p>
          {!!error && <Notice error>{errorMessage(error)}</Notice>}
          <div className="Admin-actions">
            <button
              className="Btn Btn--outline"
              disabled={pending}
              onClick={() => setConfirm(false)}
            >
              {t('common.cancel')}
            </button>
            <button className="Btn Btn--primary" disabled={pending} onClick={run}>
              {t(pending ? 'admin.ui.running' : 'admin.ui.run')}
            </button>
          </div>
        </Dialog>
      )}
    </>
  )
}
