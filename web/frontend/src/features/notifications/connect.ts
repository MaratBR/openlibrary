import { subscribeUserEvent } from '@/features/sse/api'
import type { NotificationEvent } from './api'

subscribeUserEvent<NotificationEvent>('notification', (event) => {
  if (!event || typeof event.title !== 'string' || typeof event.content !== 'string') return
  window.toast({ title: event.title, text: event.content, type: event.type })
  window.dispatchEvent(new CustomEvent('ol:notification', { detail: event }))
})
