import { httpClient } from '@/http-client';
import { z } from 'zod';

export const readingListStatusSchema = z.enum(['dnf', 'paused', 'read', 'reading', 'want_to_read'])

export type ReadingListStatus = z.infer<typeof readingListStatusSchema>

export const readingListDtoSchema = z.object({
  lastUpdatedAt: z.string(),
  chapterId: z.string().nullable(),
  status: readingListStatusSchema,
})

export type ReadingListDto = z.infer<typeof readingListDtoSchema>


export async function updateReadingListStatus(
  bookId: string,
  status: ReadingListStatus,
): Promise<ReadingListDto> {
  const response = await httpClient.post('/_api/reading-list/status', {
    searchParams: { bookId, status },
  })
  if (!response.ok) {
    throw new Error('unexpected non-ok status code ' + response.status)
  }
  const json = await response.json()
  return readingListDtoSchema.parse(json)
}

export async function updateReadingListStartReading(
  bookId: string,
  chapterId: string,
): Promise<ReadingListDto> {
  const response = await httpClient.post('/_api/reading-list/start-reading', {
    searchParams: { bookId, chapterId },
  })
  if (!response.ok) {
    throw new Error('unexpected non-ok status code ' + response.status)
  }
  const json = await response.json()
  return readingListDtoSchema.parse(json)
}
