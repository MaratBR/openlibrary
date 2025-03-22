import { httpClient, OLAPIResponse } from '@/http-client'
import { z } from 'zod'
import { definedTagDtoSchema } from '../search-filters/api'

export type UpdateBookRequest = {
  rating: string
  summary: string
  tags: string[]
  name: string
  isPubliclyVisible: boolean
}

const bookCollectionDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  position: z.number(),
  size: z.number(),
})

const bookChapterDtoSchema = z.object({
  id: z.string(),
  order: z.number().min(0).int(),
  name: z.string(),
  words: z.number(),
  createdAt: z.string(),
  summary: z.string(),
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
  chapters: z.array(bookChapterDtoSchema),
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

export type ManagerBookDetails = z.infer<typeof managerBookDetailsSchema>

export function httpUpdateBook(id: string, request: UpdateBookRequest) {
  return httpClient
    .post(`/_api/books-manager/book/${id}`, {
      json: request,
    })
    .then((r) => OLAPIResponse.create(r, managerBookDetailsSchema))
}

export type UploadCoverRequest = {
  file: File
  clientCropped: boolean
  bookId: string
}

const uploadCoverResponseSchema = z.string()

export type UploadCoverResponse = z.infer<typeof uploadCoverResponseSchema>

export function httpUploadCover(req: UploadCoverRequest): Promise<UploadCoverResponse> {
  const body = new FormData()
  body.append('file', req.file)
  body.append('clientCropped', req.clientCropped.toString())

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
