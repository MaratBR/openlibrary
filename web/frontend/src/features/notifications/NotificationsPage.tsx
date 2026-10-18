import { useEffect, useState } from 'react'
import { useLoaderData, type LoaderFunctionArgs } from 'react-router'
import { notificationHistory, type NotificationEvent } from './api'
import { WorkspacePage, t, EmptyState, date } from '@/features/moderation/components/workspace/ui'

export function notificationsLoader({ request }: LoaderFunctionArgs) {
  return notificationHistory('0', request.signal)
}

export default function NotificationsPage() {
  const initial = useLoaderData<typeof notificationsLoader>()
  const [events, setEvents] = useState(initial)
  const [cursor, setCursor] = useState(initial.at(-1)?.id ?? '0')
  const [more, setMore] = useState(initial.length === 100)
  const [pending, setPending] = useState(false)
  const [failed, setFailed] = useState(false)
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
            <li key={event.id} className="py-5">
              <h2 className="text-lg font-semibold">{event.title}</h2>
              <time className="text-sm text-secondary-foreground" dateTime={event.createdAt}>
                {date(event.createdAt)}
              </time>
              <p className="mt-2 whitespace-pre-wrap break-words">{event.content}</p>
            </li>
          ))}
      </ol>
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
