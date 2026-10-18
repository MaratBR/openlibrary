import { httpClient } from '@/features/http-client'

export type NotificationEvent = {
  id: string
  key: string
  type: string
  title: string
  content: string
  category: 'system' | 'moderation' | 'admin'
  urgent: boolean
  icon?: { type: 0 | 1 | 2; payload: string }
  createdAt: string
}

export async function notificationHistory(before = '0', signal?: AbortSignal) {
  return httpClient
    .get('/_api/notifications', { searchParams: { before }, signal })
    .json<NotificationEvent[]>()
}

export type NotificationStatus = {
  unreadCount: number
  lastReadAt: string | null
  observedAt: string
}

export async function notificationStatus(signal?: AbortSignal) {
  return httpClient.get('/_api/notifications/status', { signal }).json<NotificationStatus>()
}

export async function markNotificationsRead(readThrough: string, signal?: AbortSignal) {
  const status = await httpClient
    .post('/_api/notifications/read', {
      json: { readThrough },
      signal,
      retry: 0,
    })
    .json<NotificationStatus>()
  window.dispatchEvent(new CustomEvent('ol:notifications-read', { detail: status }))
  return status
}
