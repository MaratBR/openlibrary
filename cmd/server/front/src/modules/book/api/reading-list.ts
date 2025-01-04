import { httpClient } from '@/modules/common/api'
import { ReadingListDto, readingListDtoSchema, ReadingListStatus } from './api'

export async function httpUpdateReadingListStatus(
  bookId: string,
  status: ReadingListStatus,
): Promise<ReadingListDto> {
  const response = await httpClient.post('/api/reading-list/status', {
    searchParams: { bookId, status },
  })
  if (!response.ok) {
    throw new Error('unexpected non-ok status code ' + response.status)
  }
  const json = await response.json()
  return readingListDtoSchema.parse(json)
}

export async function httpUpdateReadingListStartReading(
  bookId: string,
  chapterId: string,
): Promise<ReadingListDto> {
  const response = await httpClient.post('/api/reading-list/start-reading', {
    searchParams: { bookId, chapterId },
  })
  if (!response.ok) {
    throw new Error('unexpected non-ok status code ' + response.status)
  }
  const json = await response.json()
  return readingListDtoSchema.parse(json)
}
