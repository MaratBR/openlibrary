import { subscribeUserEvent } from '@/features/sse/api'
import type { NotificationEvent } from './api'

subscribeUserEvent<NotificationEvent>('notification', (event) => {
  if (!event || typeof event.title !== 'string' || typeof event.content !== 'string') return
  window.dispatchEvent(new CustomEvent('ol:notification', { detail: event }))
})
