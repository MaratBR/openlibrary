import { useEffect, useId, useRef, useState } from 'react'
import { Popover } from 'radix-ui'
import {
  markNotificationsRead,
  notificationHistory,
  notificationStatus,
  type NotificationEvent,
} from './api'
import { NotificationRecord } from './Record'
import { useNotificationStatus } from './use-status'
import '@/components/header-popover.scss'
import './menu.scss'

const t = (key: string) => window._(`notifications.${key}`)
const merge = (current: NotificationEvent[], incoming: NotificationEvent[]) =>
  [...new Map([...current, ...incoming].map((event) => [event.id, event])).values()].sort((a, b) =>
    BigInt(a.id) > BigInt(b.id) ? -1 : BigInt(a.id) < BigInt(b.id) ? 1 : 0,
  )

function Bell() {
  return (
    <svg
      width="22"
      height="22"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.75"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      <path d="M18 8a6 6 0 0 0-12 0c0 7-3 7-3 9h18c0-2-3-2-3-9" />
      <path d="M10 21h4" />
    </svg>
  )
}

export function NotificationsMenu() {
  const simulateNoNotifications =
    new URLSearchParams(window.location.search).get('debug.simulateNoNotifications') === '1'
  const titleId = useId()
  const [open, setOpen] = useState(false)
  const [events, setEvents] = useState<NotificationEvent[]>([])
  const { status, accept, statusFailed } = useNotificationStatus()
  const [readFailed, setReadFailed] = useState(false)
  const readThrough = useRef<string | undefined>(undefined)
  const [loaded, setLoaded] = useState(simulateNoNotifications)
  const [pending, setPending] = useState(false)
  const [failed, setFailed] = useState(false)
  const busy = useRef(false)
  const controller = useRef<AbortController | null>(null)

  useEffect(() => () => controller.current?.abort(), [])
  useEffect(() => {
    if (simulateNoNotifications) return
    const receive = (message: Event) => {
      const event = (message as CustomEvent<NotificationEvent>).detail
      setEvents((previous) => merge(previous, [event]))
    }
    window.addEventListener('ol:notification', receive)
    return () => window.removeEventListener('ol:notification', receive)
  }, [simulateNoNotifications])

  async function acknowledge(signal?: AbortSignal) {
    if (!readThrough.current) return
    setReadFailed(false)
    try {
      accept(await markNotificationsRead(readThrough.current, signal))
    } catch {
      if (!signal?.aborted) setReadFailed(true)
    }
  }

  async function load() {
    if (simulateNoNotifications) return
    if (busy.current && !controller.current?.signal.aborted) return
    busy.current = true
    const request = new AbortController()
    controller.current = request
    setPending(true)
    setFailed(false)
    try {
      const snapshot = await notificationStatus(request.signal)
      if (snapshot) accept(snapshot)
      const page = await notificationHistory('0', request.signal)
      if (request.signal.aborted) return
      setEvents((previous) => merge(previous, page))
      setLoaded(true)
      if (snapshot) {
        readThrough.current = snapshot.observedAt
        await acknowledge(request.signal)
      }
    } catch {
      if (!request.signal.aborted) setFailed(true)
    } finally {
      if (controller.current === request) {
        busy.current = false
        setPending(false)
      }
    }
  }

  return (
    <Popover.Root
      open={open}
      onOpenChange={(value) => {
        setOpen(value)
        if (value) {
          void load()
        } else {
          controller.current?.abort()
        }
      }}
    >
      <Popover.Trigger
        className="HeaderPopover-trigger NotificationMenu-trigger"
        aria-label={
          !simulateNoNotifications && status?.unreadCount
            ? `${t('title')}: ${status.unreadCount} ${t('unread')}`
            : t('title')
        }
      >
        <Bell />
        {!simulateNoNotifications && !!status?.unreadCount && (
          <span className="NotificationMenu-badge" aria-hidden="true">
            {status.unreadCount}
          </span>
        )}
      </Popover.Trigger>
      <Popover.Portal>
        <Popover.Content
          className="HeaderPopover NotificationMenu"
          align="end"
          sideOffset={10}
          collisionPadding={12}
          aria-labelledby={titleId}
        >
          <header className="NotificationMenu-header">
            <h2 id={titleId}>{t('title')}</h2>
            <Popover.Close className="HeaderPopover-close" aria-label={t('close')}>
              <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="1.75"
                strokeLinecap="round"
                aria-hidden="true"
              >
                <path d="m6 6 12 12M6 18 18 6" />
              </svg>
            </Popover.Close>
          </header>
          <div className="NotificationMenu-body" aria-busy={pending}>
            {loaded && !pending && !failed && !events.length && (
              <div className="NotificationMenu-empty" role="status">
                <Bell />
                <p>{t('empty')}</p>
                <span>{t('emptyHint')}</span>
              </div>
            )}
            <ol className="NotificationMenu-list">
              {events.map((event) => (
                <li className="NotificationMenu-item" key={event.id}>
                  <NotificationRecord event={event} />
                </li>
              ))}
            </ol>
            {!simulateNoNotifications && (readFailed || statusFailed) && (
              <div className="NotificationMenu-status" role="alert">
                <p>{t(readFailed ? 'readError' : 'statusError')}</p>
                <button
                  className="btn"
                  disabled={pending}
                  onClick={() => (readFailed ? void acknowledge() : void load())}
                >
                  {t('retry')}
                </button>
              </div>
            )}
            {pending && (
              <p className="NotificationMenu-status" role="status">
                {t('loading')}
              </p>
            )}
            {failed && (
              <div className="NotificationMenu-status" role="alert">
                <p>{t('error')}</p>
                <button className="btn" onClick={() => void load()}>
                  {t('retry')}
                </button>
              </div>
            )}
          </div>
          <a className="NotificationMenu-more" href="/notifications">
            {t('center')}
          </a>
        </Popover.Content>
      </Popover.Portal>
    </Popover.Root>
  )
}
