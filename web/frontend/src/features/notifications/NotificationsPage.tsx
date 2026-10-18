import { useEffect, useState } from 'react'
import { useLoaderData, type LoaderFunctionArgs } from 'react-router'
import {
  markNotificationsRead,
  notificationHistory,
  notificationStatus,
  type NotificationEvent,
} from './api'
import { NotificationRecord } from './Record'
import { WorkspacePage, t, EmptyState } from '@/features/moderation/components/workspace/ui'

export async function notificationsLoader({ request }: LoaderFunctionArgs) {
  const status = await notificationStatus(request.signal)
  const events = await notificationHistory('0', request.signal)
  return { status, events }
}

export default function NotificationsPage() {
  const initial = useLoaderData<typeof notificationsLoader>()
  const [events, setEvents] = useState(initial.events)
  const [cursor, setCursor] = useState(initial.events.at(-1)?.id ?? '0')
  const [more, setMore] = useState(initial.events.length === 100)
  const [pending, setPending] = useState(false)
  const [failed, setFailed] = useState(false)
  const [readFailed, setReadFailed] = useState(false)
  async function acknowledge() {
    try {
      await markNotificationsRead(initial.status.observedAt)
      setReadFailed(false)
    } catch {
      setReadFailed(true)
    }
  }
  useEffect(() => {
    void acknowledge()
  }, [initial.status.observedAt])
  useEffect(() => {
    const receive = (message: Event) => {
      const event = (message as CustomEvent<NotificationEvent>).detail
      setEvents((previous) =>
        previous.some((item) => item.id === event.id) ? previous : [...previous, event],
      )
    }
    window.addEventListener('ol:notification', receive)
    return () => window.removeEventListener('ol:notification', receive)
  }, [])
  async function loadMore() {
    setPending(true)
    setFailed(false)
    try {
      const page = await notificationHistory(cursor)
      setEvents((previous) => [
        ...previous,
        ...page.filter((event) => !previous.some((item) => item.id === event.id)),
      ])
      setCursor(page.at(-1)?.id ?? cursor)
      setMore(page.length === 100)
    } catch {
      setFailed(true)
    } finally {
      setPending(false)
    }
  }
  return (
    <WorkspacePage title={t('notifications')}>
      <p className="text-secondary-foreground">{t('notificationsLive')}</p>
      {!events.length && <EmptyState>{t('notificationsEmpty')}</EmptyState>}
      <ol className="divide-y divide-border">
        {[...events]
          .sort((a, b) => (BigInt(a.id) > BigInt(b.id) ? -1 : 1))
          .map((event) => (
            <li key={event.id}>
              <NotificationRecord event={event} />
            </li>
          ))}
      </ol>
      {readFailed && (
        <div role="alert" className="Notice Notice--error">
          <p>{window._('notifications.readError')}</p>
          <button className="Btn Btn--outline" onClick={() => void acknowledge()}>
            {window._('notifications.retry')}
          </button>
        </div>
      )}
      {failed && (
        <p role="alert" className="Notice Notice--error">
          {t('notificationsError')}
        </p>
      )}
      {more && (
        <button className="Btn Btn--outline" disabled={pending} onClick={loadMore}>
          {t('notificationsMore')}
        </button>
      )}
    </WorkspacePage>
  )
}
