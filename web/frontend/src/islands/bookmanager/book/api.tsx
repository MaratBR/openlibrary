import { definedTagDtoSchema } from '@/api/search'
import { httpClient } from '@/http-client'
import { z } from 'zod'

const bookCollectionDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  position: z.number(),
  size: z.number(),
})

const managerBookChapterDtoSchema = z.object({
  id: z.string(),
  order: z.number().min(0).int(),
  name: z.string(),
  words: z.number(),
  createdAt: z.string(),
  summary: z.string(),
  isAdultOverride: z.boolean(),
  isPubliclyVisible: z.boolean(),
  draftId: z.string().nullable(),
})

export const managerBookDetailsSchema = z.object({
  id: z.string(),
  name: z.string(),
  ageRating: z.string(),
  adult: z.boolean(),
  tags: z.array(definedTagDtoSchema),
  words: z.number(),
  wordsPerChapter: z.number(),
  collections: z.array(bookCollectionDtoSchema),
  chapters: z.array(managerBookChapterDtoSchema),
  createdAt: z.string(),
  author: z.object({
    id: z.string(),
    name: z.string(),
  }),
  summary: z.string(),
  isPubliclyVisible: z.boolean(),
  isBanned: z.boolean(),
  cover: z.string(),
})

export type ManagerBookDetailsDto = z.infer<typeof managerBookDetailsSchema>

export type UploadCoverRequest = {
  file: File
  bookId: string
}

const uploadCoverResponseSchema = z.string()

export type UploadCoverResponse = z.infer<typeof uploadCoverResponseSchema>

export function httpUploadCover(req: UploadCoverRequest): Promise<UploadCoverResponse> {
  const body = new FormData()
  body.append('file', req.file)

  return httpClient
    .post(`/_api/books-manager/book/${req.bookId}/cover`, {
      body,
    })
    .then((r) => r.json())
    .then(uploadCoverResponseSchema.parse)
}

export function httpUpdateChaptersOrder(bookId: string, order: string[]) {
  return httpClient.post(`/_api/books-manager/book/${bookId}/chapters-order`, {
    json: order,
  })
}
