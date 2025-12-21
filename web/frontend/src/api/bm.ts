import { httpClient, OLAPIResponse } from '@/http-client'
import { z } from 'zod'
import { definedTagDtoSchema } from './search'

export function httpUpdateDraft(
  bookId: string,
  chapterId: string,
  draftId: string,
  content: string,
) {
  return httpClient
    .post(`/_api/books-manager/book/${bookId}/${chapterId}/${draftId}`, {
      body: content,
      headers: {
        'Content-Type': 'text/plain',
      },
    })
    .then((r) => OLAPIResponse.createNoBody(r))
}

export function httpUpdateAndPublishDraft(
  bookId: string,
  chapterId: string,
  draftId: string,
  content: string,
  makePublic: boolean,
) {
  return httpClient
    .post(`/_api/books-manager/book/${bookId}/${chapterId}/${draftId}/publish`, {
      body: content,
      headers: {
        'Content-Type': 'text/plain',
      },
      searchParams: {
        makePublic,
      },
    })
    .then((r) => OLAPIResponse.createNoBody(r))
}

export function httpUpdateDraftChapterName(
  bookId: string,
  chapterId: string,
  draftId: string,
  chapterName: string,
) {
  return httpClient
    .post(`/_api/books-manager/book/${bookId}/${chapterId}/${draftId}/chapterName`, {
      body: chapterName,
      headers: {
        'Content-Type': 'text/plain',
      },
    })
    .then((r) => OLAPIResponse.createNoBody(r))
}

export function httpCreateChapter(
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

export function httpUpdateChaptersOrder(
  bookId: string,
  payload: {
    modifications: {
      chapterId: string
      newIndex: number
    }[]
  },
) {
  return httpClient.post(`/_api/books-manager/book/${bookId}/chapters-order`, {
    json: payload,
  })
}

export function httpGetBookChapters(bookId: string) {
  return httpClient
    .get(`/_api/books-manager/book/${bookId}/chapters`)
    .then((r) => OLAPIResponse.create(r, z.array(managerBookChapterDtoSchema)))
}
