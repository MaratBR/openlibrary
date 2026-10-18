import { useEffect, useState } from 'react'
import { notificationStatus, type NotificationStatus } from './api'

// Reject stale status responses and share read acknowledgements with other
// mounted inboxes. The server owns read state across tabs and devices.
export function useNotificationStatus() {
  const [status, setStatus] = useState<NotificationStatus>()
  const [failed, setFailed] = useState(false)

  function accept(next: NotificationStatus) {
    setStatus((previous) => {
      if (previous && new Date(previous.observedAt) > new Date(next.observedAt)) return previous
      if (
        previous?.lastReadAt &&
        (!next.lastReadAt || new Date(previous.lastReadAt) > new Date(next.lastReadAt))
      )
        return previous
      return next
    })
    setFailed(false)
  }

  useEffect(() => {
    const controller = new AbortController()
    let timer: ReturnType<typeof setTimeout>
    const refresh = () => {
      clearTimeout(timer)
      timer = setTimeout(() => {
        void notificationStatus(controller.signal)
          .then(accept)
          .catch(() => {
            if (!controller.signal.aborted) setFailed(true)
          })
      }, 100)
    }
    const onRead = (event: Event) => accept((event as CustomEvent<NotificationStatus>).detail)
    refresh()
    window.addEventListener('focus', refresh)
    window.addEventListener('ol:notification', refresh)
    window.addEventListener('ol:notifications-read', onRead)
    return () => {
      clearTimeout(timer)
      controller.abort()
      window.removeEventListener('focus', refresh)
      window.removeEventListener('ol:notification', refresh)
      window.removeEventListener('ol:notifications-read', onRead)
    }
  }, [])

  return { status, accept, statusFailed: failed }
}
