import { httpClient } from '@/http-client'
import type { BookReadingListDto, ReadingListStatus } from '@/backend-types'

export async function updateReadingListStatus(
  bookId: string,
  status: ReadingListStatus,
): Promise<BookReadingListDto> {
  const response = await httpClient.post('/_api/reading-list/status', {
    searchParams: { bookId, status },
  })
  if (!response.ok) {
    throw new Error(`unexpected non-ok status code ${response.status}`)
  }
  return response.json<BookReadingListDto>()
}

export async function updateReadingListStartReading(
  bookId: string,
  chapterId: string,
): Promise<BookReadingListDto> {
  const response = await httpClient.post('/_api/reading-list/start-reading', {
    searchParams: { bookId, chapterId },
  })
  if (!response.ok) {
    throw new Error(`unexpected non-ok status code ${response.status}`)
  }
  return response.json<BookReadingListDto>()
}

export type { BookReadingListDto as ReadingListDto, ReadingListStatus } from '@/backend-types'
