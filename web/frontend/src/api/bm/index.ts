import { httpClient, OLAPIResponse } from '@/features/http-client'
import { Schema } from 'effect'
import type { ManagerBookChapterDto } from '@/backend-types'

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
