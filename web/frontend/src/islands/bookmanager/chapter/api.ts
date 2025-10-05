import { httpClient, OLAPIResponse } from '@/http-client'
import { z } from 'zod'

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
    .post(`/_api/books-manager/book/${bookId}/createChapter`, {
      body: JSON.stringify(request),
    })
    .then((r) => OLAPIResponse.create(r, z.string()))
}
