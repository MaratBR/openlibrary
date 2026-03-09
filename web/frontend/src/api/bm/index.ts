import { httpClient, OLAPIResponse } from '@/http-client'
import { z } from 'zod'
import { DraftDtoSchema } from '@/block-editor/contracts'

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
    .then((r) => OLAPIResponse.create(r, DraftDtoSchema))
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
    .then((r) => OLAPIResponse.create(r, DraftDtoSchema))
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

export type ManagerBookChapterDto = z.infer<typeof managerBookChapterDtoSchema>

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
