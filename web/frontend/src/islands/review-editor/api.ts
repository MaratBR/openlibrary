import { KyResponse } from 'ky'
import { httpClient, OLAPIResponse } from '@/http-client'
import type { RatingAndReview, RatingValue, ReviewDto } from '@/backend-types'

export type CreateReviewRequest = {
  rating: RatingValue
  content: string
}

export function httpUpdateReview(bookId: string, request: CreateReviewRequest): Promise<ReviewDto> {
  return httpClient
    .post(`/_api/reviews/${bookId}`, {
      json: request,
    })
    .then((r) => r.json<ReviewDto>())
}

export async function httpDeleteReview(bookId: string): Promise<KyResponse> {
  return await httpClient.delete(`/_api/reviews/${bookId}`)
}

export async function httpGetReview(bookId: string) {
  return httpClient
    .get(`/_api/reviews/${bookId}`)
    .then((r) => OLAPIResponse.create<RatingAndReview>(r))
}

export type { RatingValue, ReviewDto } from '@/backend-types'
