import { z } from "zod"

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
