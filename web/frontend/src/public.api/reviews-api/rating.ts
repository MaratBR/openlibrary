import { httpClient } from "@/http-client";

export async function updateRating(bookId: string, rating: number): Promise<void> {
  await httpClient.post('/_api/reviews/rating', { searchParams: { bookId, rating } })
}