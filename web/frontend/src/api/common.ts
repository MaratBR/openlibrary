import z from 'zod'

export const AgeRatingSchema = z.enum(['?', 'G', 'PG', 'PG-13', 'R', 'NC-17'])

export type AgeRating = z.infer<typeof AgeRatingSchema>

export const BookCoverSchema = z.object({
  url: z.string(),
})

export type BookCover = z.infer<typeof BookCoverSchema>
