import { z } from 'zod'
import { ageRatingSchema } from '../api'

export const savedBookSchema = z.object({
  _id: z.string(),
  _ts: z.number(),
  name: z.string(),
  authorId: z.string(),
  authorName: z.string(),
  cover: z.string().nullable(),
  summary: z.string().nullable(),
  chapters: z.number(),
  words: z.number(),
  wordsPerChapter: z.number(),
  ageRating: ageRatingSchema,
  createdAt: z.string(),
  updatedAt: z.string(),
  isFavorite: z.boolean(),
})

export type SavedBook = z.infer<typeof savedBookSchema>

export const savedChapterSchema = z.object({
  _id: z.string(),
  _ts: z.number(),
  bookId: z.string(),
  name: z.string(),
  summary: z.string().nullable(),
  createdAt: z.string(),
  updatedAt: z.string(),
  content: z.string(),
})

export type SavedChapter = z.infer<typeof savedChapterSchema>
