import { Schema } from 'effect'
import type { AgeRating, ApiResponseGetBooks, ManagerBookDetailsDto } from '@/backend-types'
import { httpClient, OLAPIResponse } from '@/features/http-client'

export type ApiPayloadGetBooks = {
  page: number
  size: number
  search: string
}

export type ApiPayloadTrashBook = {
  id: string
  trash: boolean
}

export type ApiPayloadBookDirectUpdate = {
  name: string
  ageRating: AgeRating
  summary: string
  isAdult: boolean
  isPubliclyVisible: boolean
  tags: string[]
}

export class BMBookAPI {
  private static _instance = new BMBookAPI()

  public static getInstance() {
    return this._instance
  }

  getBook(id: string) {
    return httpClient
      .get(`/_api/books-manager/books/${id}`)
      .then((r) => OLAPIResponse.create<ManagerBookDetailsDto>(r))
  }

  trashBook(payload: ApiPayloadTrashBook) {
    return httpClient
      .post('/_api/books-manager/books/trash', { searchParams: payload })
      .then((r) => OLAPIResponse.createNoBody(r))
  }

  getBooks(payload: ApiPayloadGetBooks): Promise<OLAPIResponse<ApiResponseGetBooks>> {
    return httpClient
      .get('/_api/books-manager/books', { searchParams: payload })
      .then((r) => OLAPIResponse.create<ApiResponseGetBooks>(r))
  }

  normalizeChapterName(name: string) {
    name = name.trim()

    const valid = Array.from(name).length <= 70 && name.length > 0

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
      content: string
    },
  ) {
    return httpClient
      .post(`/_api/books-manager/book/${bookId}/create-chapter`, {
        body: JSON.stringify(request),
      })
      .then((r) => OLAPIResponse.create(r, Schema.String))
  }

  updateBook(bookId: string, body: ApiPayloadBookDirectUpdate) {
    return httpClient
      .post(`/_api/books-manager/book/${bookId}/direct-update`, { json: body })
      .then((r) => OLAPIResponse.create<ManagerBookDetailsDto>(r))
  }

  updateChapter(
    bookId: string,
    chapterId: string,
    body: { name: string; summary: string; isPubliclyVisible: boolean },
  ) {
    return httpClient
      .post(`/_api/books-manager/book/${bookId}/chapter/${chapterId}/direct-update`, { json: body })
      .then((r) => OLAPIResponse.createNoBody(r))
  }
}

export type { ApiResponseGetBooks, ManagerBookDetailsDto, ManagerBookDto } from '@/backend-types'
