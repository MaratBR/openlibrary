import { httpClient, OLAPIResponse } from '@/http-client'
import { BookModerationLogSchema, ModerationBookSchema } from './schemas'

export type BookModerationAction =
  | 'ban'
  | 'unban'
  | 'shadow-ban'
  | 'unshadow-ban'
  | 'permanent-remove'

export function getModerationBook(bookId: string) {
  return httpClient
    .get(`/_api/moderation/books/${encodeURIComponent(bookId)}`)
    .then((response) => OLAPIResponse.create(response, ModerationBookSchema))
}

export function getBookModerationLog(
  bookId: string,
  options: { page?: number; pageSize?: number } = {},
) {
  return httpClient
    .get(`/_api/moderation/books/${encodeURIComponent(bookId)}/log`, {
      searchParams: {
        page: options.page ?? 1,
        pageSize: options.pageSize ?? 25,
      },
    })
    .then((response) => OLAPIResponse.create(response, BookModerationLogSchema))
}

export function performBookModerationAction(
  bookId: string,
  action: BookModerationAction,
  reason: string,
) {
  return httpClient
    .post(`/_api/moderation/books/${encodeURIComponent(bookId)}/actions/${action}`, {
      json: { reason },
    })
    .then((response) => OLAPIResponse.createNoBody(response))
}
