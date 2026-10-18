import { httpClient } from '@/features/http-client'

export type NotificationEvent = {
  id: string
  key: string
  type: string
  title: string
  content: string
  createdAt: string
}

export async function notificationHistory(before = '0', signal?: AbortSignal) {
  return httpClient
    .get('/_api/notifications', { searchParams: { before }, signal })
    .json<NotificationEvent[]>()
}
