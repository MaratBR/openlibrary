import { httpClient, OLAPIResponse } from '@/http-client'
import type {
  BookModerationLogResponse,
  ModerationBookResponse,
  ModerationReasonRequest,
} from '@/backend-types'

export type BookModerationAction =
  | 'ban'
  | 'unban'
  | 'shadow-ban'
  | 'unshadow-ban'
  | 'permanent-remove'

export function getModerationBook(bookId: string) {
  return httpClient
    .get(`/_api/moderation/books/${encodeURIComponent(bookId)}`)
    .then((response) => OLAPIResponse.create<ModerationBookResponse>(response))
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
    .then((response) => OLAPIResponse.create<BookModerationLogResponse>(response))
}

export function performBookModerationAction(
  bookId: string,
  action: BookModerationAction,
  reason: string,
) {
  const body: ModerationReasonRequest = { reason }
  return httpClient
    .post(`/_api/moderation/books/${encodeURIComponent(bookId)}/actions/${action}`, {
      json: body,
    })
    .then((response) => OLAPIResponse.createNoBody(response))
}
