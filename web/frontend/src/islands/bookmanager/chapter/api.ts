import { httpClient, OLAPIResponse } from '@/http-client'

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
) {
  return httpClient
    .post(`/_api/books-manager/book/${bookId}/${chapterId}/${draftId}/publish`, {
      body: content,
      headers: {
        'Content-Type': 'text/plain',
      },
    })
    .then((r) => OLAPIResponse.createNoBody(r))
}
