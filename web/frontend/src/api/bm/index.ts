import { httpClient, OLAPIResponse } from '@/features/http-client'
import { Schema } from 'effect'
import type { DraftDto, ManagerBookChapterDto } from '@/backend-types'

export function httpUpdateDraft(
  bookId: string,
  chapterId: string,
  draftId: string,
  content: string,
) {
  return httpClient
    .post(`/_api/books-manager/book/${bookId}/${chapterId}/${draftId}`, {
      json: { content },
    })
    .then((r) => OLAPIResponse.create<DraftDto>(r))
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
      json: { content },
      searchParams: {
        makePublic,
      },
    })
    .then((r) => OLAPIResponse.create<DraftDto>(r))
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

export function httpScheduleDraft(
  bookId: string,
  chapterId: string,
  draftId: string,
  scheduledAt: string,
) {
  return httpClient
    .post(`/_api/books-manager/book/${bookId}/${chapterId}/${draftId}/schedule`, {
      json: { scheduledAt },
    })
    .then((r) => OLAPIResponse.create<DraftDto>(r))
}

export type UploadCoverRequest = {
  file: File
  clientCropped: boolean
  bookId: string
}

const uploadCoverResponseSchema = Schema.String

export type UploadCoverResponse = Schema.Schema.Type<typeof uploadCoverResponseSchema>

export function httpUploadCover(req: UploadCoverRequest): Promise<UploadCoverResponse> {
  const body = new FormData()
  body.append('file', req.file)
  body.append('clientCropped', req.clientCropped.toString())

  return httpClient
    .post(`/_api/books-manager/book/${req.bookId}/cover`, {
      body,
    })
    .then((r) => r.json())
    .then(Schema.decodeUnknownSync(uploadCoverResponseSchema))
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
    .then((r) => OLAPIResponse.create<ManagerBookChapterDto[]>(r))
}

export type { ManagerBookChapterDto } from '@/backend-types'
