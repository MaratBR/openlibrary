import {
  AgeRating,
  ManagerAuthorBookDto,
  BookChapterDto,
  BookCollectionDto,
  DefinedTagDto,
} from '@/modules/book/api'
import { httpClient } from '@/modules/common/api'
import { z } from 'zod'

export * from './cover'

export type CreateBookRequest = {
  name: string
  ageRating: AgeRating
  tags: string[]
  summary: string
  isPubliclyVisible: boolean
}

export type CreateBookResponse = {
  id: string
}

export function httpCreateBook(req: CreateBookRequest): Promise<CreateBookResponse> {
  return httpClient.post('/api/manager/books', { json: req }).then((r) => r.json())
}

export type ImportFromAo3Request = {
  id: string
}

export type ImportFromAo3Response = {
  id: string
}

export function httpImportBookFromAo3(req: ImportFromAo3Request): Promise<ImportFromAo3Response> {
  return httpClient.post('/api/manager/books/ao3-import', { json: req }).then((r) => r.json())
}

export type CreateBookChapterRequest = {
  content: string
  isAdultOverride: boolean
  name: string
  summary: string
}

export const createBookChapterRequestSchema = z.object({
  content: z.string(),
  isAdultOverride: z.boolean(),
  name: z.string(),
  summary: z.string(),
})

export type CreateBookChapterResponse = z.infer<typeof createBookChapterRequestSchema>

export function httpCreateBookChapter(
  bookId: string,
  req: CreateBookChapterRequest,
): Promise<CreateBookChapterResponse> {
  return httpClient
    .post(`/api/manager/books/${bookId}/chapters`, { json: req })
    .then((r) => r.json())
}

export const updateBookChapterRequestSchema = z.object({
  content: z.string(),
  isAdultOverride: z.boolean(),
  name: z.string(),
  summary: z.string(),
})

export type UpdateBookChapterRequest = z.infer<typeof updateBookChapterRequestSchema>

export function httpUpdateBookChapter(
  bookId: string,
  chapterId: string,
  req: UpdateBookChapterRequest,
) {
  return httpClient.post(`/api/manager/books/${bookId}/chapters/${chapterId}`, { json: req })
}

export type UpdateBookRequest = {
  name: string
  ageRating: AgeRating
  tags: string[]
  summary: string
  isPubliclyVisible: boolean
}

export type UpdateBookResponse = {
  book: ManagerBookDetailsDto
}

export function httpUpdateBook(id: string, req: UpdateBookRequest): Promise<UpdateBookResponse> {
  return httpClient.post(`/api/manager/books/${id}`, { json: req }).then((r) => r.json())
}

export type ManagerBookDetailsDto = {
  id: string
  name: string
  ageRating: AgeRating
  isAdult: boolean
  tags: DefinedTagDto[]
  words: number
  wordsPerChapter: number
  collections: BookCollectionDto[]
  chapters: BookChapterDto[]
  createdAt: string
  author: {
    id: string
    name: string
  }
  permissions: {
    canEdit: boolean
  }
  isPubliclyVisible: boolean
  isBanned: boolean
  cover: string
}

export type ManagerGetBookResponse = ManagerBookDetailsDto

export function httpManagerGetBook(id: string): Promise<ManagerGetBookResponse> {
  return httpClient.get(`/api/manager/books/${id}`).then((r) => r.json())
}

export const managerChapterDetailsDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  createdAt: z.string().datetime(),
  words: z.number().nonnegative(),
  summary: z.string(),
  isAdultOverride: z.boolean(),
  order: z.number(),
  isPubliclyVisible: z.boolean(),
  content: z.string(),
})

export type ManagerChapterDetailsDto = z.infer<typeof managerChapterDetailsDtoSchema>

export function httpManagerGetBookChapter(
  bookId: string,
  id: string,
): Promise<ManagerChapterDetailsDto> {
  return httpClient.get(`/api/manager/books/${bookId}/chapters/${id}`).then((r) => r.json())
}

export const managerBookChapterDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  createdAt: z.string(),
  words: z.number().nonnegative(),
  summary: z.string(),
  isAdultOverride: z.boolean(),
  order: z.number(),
})

export type ManagerBookChapterDto = z.infer<typeof managerBookChapterDtoSchema>

export function httpManagerGetChapters(id: string): Promise<ManagerBookChapterDto[]> {
  return httpClient.get(`/api/manager/books/${id}/chapters`).then((r) => r.json())
}

export function httpGetMyBooks(): Promise<{ books: ManagerAuthorBookDto[] }> {
  return httpClient.get('/api/manager/books/my-books').then((r) => r.json())
}

export type ReorderChaptersRequest = {
  sequence: string[]
}

export function httpReorderChapters(bookId: string, chapterIds: string[]) {
  return httpClient.post(`/api/manager/books/${bookId}/chapters/reorder`, {
    json: { sequence: chapterIds } satisfies ReorderChaptersRequest,
  })
}
