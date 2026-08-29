import { httpClient } from '@/features/http-client'

/** @deprecated Use the review-editor island API for new review interactions. */
export async function updateRating(bookId: string, rating: number): Promise<void> {
  await httpClient.post('/_api/reviews/rating', { searchParams: { bookId, rating } })
}
