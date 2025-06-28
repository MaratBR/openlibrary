import { z } from 'zod'
import { KyResponse } from 'ky'
import { httpClient, OLAPIResponse } from '@/http-client'

export const ratingSchema = z.union([
  z.literal(1),
  z.literal(2),
  z.literal(3),
  z.literal(4),
  z.literal(5),
  z.literal(6),
  z.literal(7),
  z.literal(8),
  z.literal(9),
  z.literal(10),
])

export type RatingValue = z.infer<typeof ratingSchema>

export const reviewDtoSchema = z.object({
  user: z.object({
    id: z.string(),
    name: z.string(),
    avatar: z.string(),
  }),
  rating: ratingSchema,
  content: z.string(),
  createdAt: z.string(),
  updatedAt: z.string().nullable(),
})

export type ReviewDto = z.infer<typeof reviewDtoSchema>

export const ratingAndReviewSchema = z.object({
  rating: ratingSchema.nullable(),
  review: reviewDtoSchema.nullable(),
})

export type CreateReviewRequest = {
  rating: RatingValue
  content: string
}

export function httpUpdateReview(bookId: string, request: CreateReviewRequest): Promise<ReviewDto> {
  return httpClient
    .post(`/_api/reviews/${bookId}`, {
      json: request,
    })
    .then((r) => r.json())
    .then((r) => reviewDtoSchema.parse(r))
}

export async function httpDeleteReview(bookId: string): Promise<KyResponse> {
  return await httpClient.delete(`/_api/reviews/${bookId}`)
}

export async function httpGetReview(bookId: string) {
  return httpClient
    .get(`/_api/reviews/${bookId}`)
    .then((r) => OLAPIResponse.create(r, ratingAndReviewSchema))
}
