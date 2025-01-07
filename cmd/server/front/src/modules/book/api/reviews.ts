import { httpClient, withPreloadCache } from '@/modules/common/api'
import { KyResponse } from 'ky'
import { z } from 'zod'

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

export const reviewsResponseSchema = z.object({
  reviews: z.array(reviewDtoSchema),
  pagination: z.object({
    page: z.number(),
    pageSize: z.number(),
  }),
})

export type ReviewsResponse = z.infer<typeof reviewsResponseSchema>

export function getPreloadedReviews(bookId: string): ReviewsResponse | undefined {
  const key = `/api/reviews/${bookId}`
  if (__server__._preload && __server__._preload[key]) {
    const value = __server__._preload[key]
    delete __server__._preload[key]
    const result = reviewsResponseSchema.safeParse(value)
    return result.success ? result.data : undefined
  }
}

export async function httpGetReviews(bookId: string): Promise<ReviewsResponse> {
  return httpClient
    .get(`/api/reviews/${bookId}`)
    .then((r) => r.json())
    .then((r) => reviewsResponseSchema.parse(r))
}

export type CreateReviewRequest = {
  rating: RatingValue
  content: string
}

export function httpUpdateReview(bookId: string, request: CreateReviewRequest): Promise<ReviewDto> {
  return httpClient
    .post(`/api/reviews/${bookId}`, {
      json: request,
    })
    .then((r) => r.json())
    .then((r) => reviewDtoSchema.parse(r))
}

export async function httpDeleteReview(bookId: string): Promise<KyResponse> {
  return await httpClient.delete(`/api/reviews/${bookId}`)
}

export async function httpGetMyReview(bookId: string): Promise<ReviewDto | null> {
  return withPreloadCache(`/api/reviews/${bookId}/my`, () =>
    httpClient
      .get(`/api/reviews/${bookId}/my`)
      .then((r) => r.json())
      .then((r) => (r === null ? null : reviewDtoSchema.parse(r))),
  )
}
