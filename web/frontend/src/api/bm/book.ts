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

export const BookCollectionDtoSchema = z.object({
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
  collections: BookCollectionDtoSchema.array(),
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
  page: z.number().int(),
})

export type ApiResponseGetBooks = z.infer<typeof ApiResponseGetBooksSchema>

export type ApiPayloadTrashBook = {
  id: string
  trash: boolean
}

export const ManagerBookChapterDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  createdAt: z.string(),
  words: z.number().int(),
  summary: z.string(),
  order: z.number(),
  isAdultOverride: z.boolean(),
  isPubliclyVisible: z.boolean(),
  draftId: z.string().nullable(),
})

export const BookDetailsAuthorDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
})

export const ManagerBookDetailsDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  ageRating: AgeRatingSchema,
  adult: z.boolean(),
  tags: z.array(DefinedTagDtoSchema),
  words: z.number().int(),
  wordsPerChapter: z.number().int(),
  createdAt: z.string(),
  collections: z.array(BookCollectionDtoSchema),
  chapters: z.array(ManagerBookChapterDtoSchema),
  author: BookDetailsAuthorDtoSchema,
  summary: z.string(),
  isPubliclyVisible: z.boolean(),
  isBanned: z.boolean(),
  cover: BookCoverSchema,
})

export type ManagerBookDetailsDto = z.infer<typeof ManagerBookDetailsDtoSchema>

export class BMBookAPI {
  private static _instance = new BMBookAPI()

  public static getInstance() {
    return this._instance
  }

  getBook(id: string) {
    return httpClient
      .get(`/_api/books-manager/books/${id}`)
      .then((r) => OLAPIResponse.create(r, ManagerBookDetailsDtoSchema))
  }

  trashBook(payload: ApiPayloadTrashBook) {
    return httpClient
      .post('/_api/books-manager/books/trash', { searchParams: payload })
      .then((r) => OLAPIResponse.createNoBody(r))
  }

  getBooks(payload: ApiPayloadGetBooks): Promise<OLAPIResponse<ApiResponseGetBooks>> {
    return httpClient
      .get('/_api/books-manager/books', { searchParams: payload })
      .then((r) => OLAPIResponse.create(r, ApiResponseGetBooksSchema))
  }

  normalizeChapterName(name: string) {
    name = name.trim()

    const valid = name.length <= 255 && name.length > 0

    return {
      value: name,
      valid,
    }
  }

  createChapter(
    bookId: string,
    request: {
      name: string
      summary: string
      isAdultOverride: boolean
      content: string
    },
  ) {
    return httpClient
      .post(`/_api/books-manager/book/${bookId}/create-chapter`, {
        body: JSON.stringify(request),
      })
      .then((r) => OLAPIResponse.create(r, z.string()))
  }
}
