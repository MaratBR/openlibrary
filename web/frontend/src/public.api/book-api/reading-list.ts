import { httpClient } from '@/features/http-client'
import type { BookReadingListDto, ReadingListStatus } from '@/backend-types'

/** @deprecated Move reading-list interactions to a feature-local API. */
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
