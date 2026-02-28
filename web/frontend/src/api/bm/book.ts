import z from 'zod'
import { AgeRatingSchema, BookCoverSchema } from '@/api/common'
import { DefinedTagDtoSchema } from '../search'
import { ViewsSchema } from '../analytics'
import { httpClient, OLAPIResponse } from '@/http-client'

export type ApiPayloadGetBooks = {
  page: number
  size: number
  search: string
}

export const BookCollectionDto = z.object({
  id: z.string(),
  name: z.string(),
  pos: z.number().int(),
  size: z.number().int(),
})

export const ManagerBookDtoSchema = z.object({
  id: z.string(),
  slug: z.string(),
  name: z.string(),
  createdAt: z.string(),
  ageRating: AgeRatingSchema,
  tags: DefinedTagDtoSchema.array(),
  words: z.number().int(),
  wordsPerChapter: z.number().int(),
  chapters: z.number().int(),
  collections: BookCollectionDto.array(),
  isPubliclyVisible: z.boolean(),
  isBanned: z.boolean(),
  isTrashed: z.boolean(),
  summary: z.string(),
  cover: BookCoverSchema,
  stats: z.object({
    views: ViewsSchema,
    reviews: z.number().int(),
    ratings: z.number().int(),
  }),
})

export type ManagerBookDto = z.infer<typeof ManagerBookDtoSchema>

export const ApiResponseGetBooksSchema = z.object({
  books: ManagerBookDtoSchema.array(),
  totalPages: z.number().int(),
})

export type ApiResponseGetBooks = z.infer<typeof ApiResponseGetBooksSchema>

export function httpBmGetBooks(
  payload: ApiPayloadGetBooks,
): Promise<OLAPIResponse<ApiResponseGetBooks>> {
  return httpClient
    .get('/_api/books-manager/books', { searchParams: payload })
    .then((r) => OLAPIResponse.create(r, ApiResponseGetBooksSchema))
}

export type ApiPayloadTrashBook = {
  id: string
  trash: boolean
}

export function httpBmTrashBook(payload: ApiPayloadTrashBook) {
  return httpClient
    .post('/_api/books-manager/books/trash', { searchParams: payload })
    .then((r) => OLAPIResponse.createNoBody(r))
}
