import { httpClient } from '@/features/http-client'

export type UserEvent = { type: string; data: unknown }

export async function getSSECapabilities() {
  return httpClient.get('/_api/sse/capabilities').json<{ sseAllowed: boolean; userId?: string }>()
}

export function subscribeUserEvent<T>(type: string, handler: (data: T) => void) {
  const receive = (message: Event) => {
    const event = (message as CustomEvent<UserEvent>).detail
    if (event.type === type) handler(event.data as T)
  }
  window.addEventListener('ol:sse', receive)
  return () => window.removeEventListener('ol:sse', receive)
}

let source: EventSource | undefined
let stopped = false

export function disconnectUserStream() {
  stopped = true
  clearTimeout(retry)
  source?.close()
  source = undefined
}

export function resumeUserStream() {
  stopped = false
  void connectUserStream()
}

let connecting = false
let retry: ReturnType<typeof setTimeout> | undefined
let after = '0'
let recipient = ''

export async function connectUserStream() {
  if (stopped || source || connecting) return
  connecting = true
  try {
    const { sseAllowed, userId } = await getSSECapabilities()
    if (!sseAllowed || !userId || stopped) return
    const storageKey = `ol.sse.${userId}`
    if (recipient !== userId) {
      recipient = userId
      after = '0'
      try {
        after = sessionStorage.getItem(storageKey) ?? '0'
      } catch {
        /* Storage may be disabled. */
      }
      if (!/^[A-Za-z0-9_-]+$/.test(after)) after = '0'
    }
    source = new EventSource(`/_api/sse/stream?after=${encodeURIComponent(after)}`)
    source.addEventListener('user-event', (message) => {
      try {
        const event: UserEvent = JSON.parse(message.data)
        if (typeof event.type !== 'string' || !('data' in event)) return
        if (message.lastEventId) after = message.lastEventId
        try {
          sessionStorage.setItem(storageKey, after)
        } catch {
          /* Keep the cursor in memory. */
        }
        window.dispatchEvent(new CustomEvent('ol:sse', { detail: event }))
      } catch {
        /* Ignore malformed messages without interrupting the stream. */
      }
    })
    source.onerror = () => {
      source?.close()
      source = undefined
      if (!stopped) retry = setTimeout(connectUserStream, 3000)
    }
  } catch {
    if (!stopped) retry = setTimeout(connectUserStream, 10000)
  } finally {
    connecting = false
  }
}
