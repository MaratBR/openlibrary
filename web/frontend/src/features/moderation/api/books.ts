import { httpClient, OLAPIResponse } from '@/http-client'
import type {
  BookModerationLogResponse,
  ModerationBookResponse,
  ModerationReasonRequest,
  ModerationValueRequest,
  ModerationBookChapterResponse,
  ModerationBooksPageResponse,
} from '@/backend-types'

export function searchModerationBooks(
  search = '',
  exact = false,
  includeBanned = false,
  includeDeleted = false,
  page = 1,
  pageSize = 20,
) {
  return httpClient
    .get('/_api/moderation/books', {
      searchParams: { search, exact, includeBanned, includeDeleted, page, pageSize },
    })
    .then((response) => OLAPIResponse.create<ModerationBooksPageResponse>(response))
}

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

export function getBookModerationChapters(bookId: string) {
  return httpClient
    .get(`/_api/moderation/books/${encodeURIComponent(bookId)}/chapters`)
    .then((response) => OLAPIResponse.create<ModerationBookChapterResponse[]>(response))
}

export function changeBookAgeRating(bookId: string, value: string, reason: string) {
  return bookValueAction(bookId, 'change-age-rating', { value, reason })
}

export function changeBookSummary(bookId: string, value: string, reason: string) {
  return bookValueAction(bookId, 'change-summary', { value, reason })
}

function bookValueAction(bookId: string, action: string, body: ModerationValueRequest) {
  return httpClient
    .post(`/_api/moderation/books/${encodeURIComponent(bookId)}/actions/${action}`, { json: body })
    .then((response) => OLAPIResponse.createNoBody(response))
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
