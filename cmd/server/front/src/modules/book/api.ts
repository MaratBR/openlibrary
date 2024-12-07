import { QueryClient, useQuery } from '@tanstack/react-query'
import {
  getPreloadedData,
  httpClient,
  parseQueryStringArray,
  stringArrayToQueryParameterValue,
  withPreloadCache,
} from '../common/api'
import { z } from 'zod'
import { useNotificationsSlot } from '../notifications/state'
import { GenericNotification } from '../notifications'
import React from 'react'

export type ManagerAuthorBookDto = {
  id: string
  name: string
  createdAt: string
  ageRating: AgeRating
  words: number
  wordsPerChapter: number
  chapters: number
  tags: DefinedTagDto[]
  collections: BookCollectionDto[]
  isPubliclyVisible: boolean
  isBanned: boolean
  summary: string
}

export type BookCollectionDto = {
  id: string
  name: string
  position: number
  size: number
}

const tagCategorySchema = z.enum(['other', 'warning', 'fandom', 'rel', 'reltype', 'unknown'])

export type TagsCategory = z.infer<typeof tagCategorySchema>

export const definedTagDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  desc: z.string(),
  adult: z.boolean(),
  spoiler: z.boolean(),
  cat: tagCategorySchema,
})

export type DefinedTagDto = z.infer<typeof definedTagDtoSchema>

export type SearchTagsResponse = {
  query: string
  tags: DefinedTagDto[]
}

export function httpTagsSearch(q: string): Promise<SearchTagsResponse> {
  return httpClient
    .get('/api/tags/search-tags', { searchParams: q ? { q } : undefined })
    .then((r) => r.json())
}

export type AgeRating = '?' | 'G' | 'PG' | 'PG-13' | 'R' | 'NC-17'

export const AGE_RATINGS_LIST: AgeRating[] = ['?', 'G', 'PG', 'PG-13', 'R', 'NC-17']

export type BookChapterDto = {
  id: string
  order: number
  name: string
  words: number
  createdAt: string
  summary: string
}

export type BookDetailsDto = {
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
  summary: string
  favorites: number
  isFavorite: boolean
  notifications: GenericNotification[]
  cover: string
}

export type GetBookResponse = BookDetailsDto

export function httpGetBook(id: string): Promise<GetBookResponse> {
  return withPreloadCache(`/api/books/${id}`, () =>
    httpClient.get(`/api/books/${id}`).then((r) => r.json()),
  )
}

export function getPreloadedBookResult(id: string) {
  return getPreloadedData<GetBookResponse>(`/api/books/${id}`)
}

export function preloadBookQuery(queryClient: QueryClient, bookId: string) {
  if (!__server__.clientPreload) return
  queryClient.prefetchQuery({
    queryKey: ['book', bookId],
    queryFn: () => httpGetBook(bookId),
    staleTime: 30000,
    gcTime: 60000,
  })
}

export function useBookQuery(bookId: string | undefined) {
  const setNotifications = useNotificationsSlot()

  const query = useQuery({
    queryKey: ['book', bookId],
    enabled: !!bookId,
    queryFn: () => httpGetBook(bookId!),
    initialData: bookId ? getPreloadedBookResult(bookId) : undefined,
    staleTime: 10000,
    gcTime: 60000,
  })

  const { data } = query

  React.useEffect(() => {
    if (!data) return
    if (data.notifications) setNotifications(data.notifications)
  }, [data, setNotifications])

  return query
}

export type ChapterPrevNextDto = {
  id: string
  name: string
  order: number
}

export type ChapterDto = {
  id: string
  name: string
  words: number
  content: string
  isAdultOverride: boolean
  createdAt: string
  order: number
  summary: string
  nextChapter: ChapterPrevNextDto | null
  prevChapter: ChapterPrevNextDto | null
}

export type GetBookChapterResponse = {
  chapter: ChapterDto
}

export function httpGetBookChapter(
  bookId: string,
  chapterId: string,
): Promise<GetBookChapterResponse> {
  return httpClient.get(`/api/books/${bookId}/chapters/${chapterId}`).then((r) => r.json())
}

export function useBookChapterQuery(bookId: string | undefined, chapterId: string | undefined) {
  return useQuery({
    queryKey: ['book', bookId, 'chapter', chapterId],
    enabled: !!bookId && !!chapterId,
    queryFn: () => httpGetBookChapter(bookId!, chapterId!),
    staleTime: 5000,
    gcTime: 60000,
  })
}

export function preloadBookChapterQuery(
  queryClient: QueryClient,
  bookId: string,
  chapterId: string,
) {
  if (!__server__.clientPreload) return
  queryClient.prefetchQuery({
    queryKey: ['book', bookId, 'chapter', chapterId],
    queryFn: () => httpGetBookChapter(bookId, chapterId),
    staleTime: 60000,
    gcTime: 60000,
  })
}

export function httpFavoriteBook(id: string, isFavorite: boolean): Promise<GetBookResponse> {
  return httpClient
    .post(`/api/favorite`, { searchParams: { bookId: id, isFavorite } })
    .then((r) => r.json())
}

export type BookSearchItem = {
  id: string
  name: string
  createdAt: string
  ageRating: AgeRating
  words: number
  wordsPerChapter: number
  chapters: number
  favorites: number
  summary: string
  author: BookDetailsDto['author']
  tags: string[]
  cover: string
}

export type SearchBooksResponse = {
  booksMeta: {
    cacheHit: boolean
    cacheKey: string
    cacheTook: number
  }
  booksTook: number
  books: BookSearchItem[]
  tags: DefinedTagDto[]
}

export type SearchBooksRequest = {
  'w.min'?: string
  'w.max'?: string
  'c.min'?: string
  'c.max'?: string
  'wc.min'?: string
  'wc.max'?: string
  'f.min'?: string
  'f.max'?: string
  it?: string[]
  et?: string[]
  iu?: string[]
  eu?: string[]
}

export function isSearchBooksRequestEqual(req1: SearchBooksRequest, req2: SearchBooksRequest) {
  return (
    searchBooksRequestToURLSearchParams(req1).toString() ===
    searchBooksRequestToURLSearchParams(req2).toString()
  )
}

export function parseSearchBooksRequest(sp: URLSearchParams): SearchBooksRequest {
  return {
    'w.min': sp.get('w.min') || undefined,
    'w.max': sp.get('w.max') || undefined,
    'c.min': sp.get('c.min') || undefined,
    'c.max': sp.get('c.max') || undefined,
    'wc.min': sp.get('wc.min') || undefined,
    'wc.max': sp.get('wc.max') || undefined,
    'f.min': sp.get('f.min') || undefined,
    'f.max': sp.get('f.max') || undefined,
    it: parseQueryStringArray(sp.get('it')),
    et: parseQueryStringArray(sp.get('et')),
    iu: parseQueryStringArray(sp.get('iu')),
  }
}

export function searchBooksRequestToURLSearchParams(query: SearchBooksRequest): URLSearchParams {
  const urlSp = new URLSearchParams()

  if (query['w.max']) urlSp.set('w.max', query['w.max'])
  if (query['w.min']) urlSp.set('w.min', query['w.min'])
  if (query['c.max']) urlSp.set('c.max', query['c.max'])
  if (query['c.min']) urlSp.set('c.min', query['c.min'])
  if (query['wc.max']) urlSp.set('wc.max', query['wc.max'])
  if (query['wc.min']) urlSp.set('wc.min', query['wc.min'])
  if (query['f.max']) urlSp.set('f.max', query['f.max'])
  if (query['f.min']) urlSp.set('f.min', query['f.min'])
  if (query.it && query.it.length) urlSp.set('it', stringArrayToQueryParameterValue(query.it) || '')
  if (query.et && query.et.length) urlSp.set('et', stringArrayToQueryParameterValue(query.et) || '')
  if (query.iu && query.iu.length) urlSp.set('iu', stringArrayToQueryParameterValue(query.iu) || '')
  if (query.eu && query.eu.length) urlSp.set('eu', stringArrayToQueryParameterValue(query.eu) || '')

  return urlSp
}

export async function httpSearchBooks(query: SearchBooksRequest): Promise<SearchBooksResponse> {
  const sp = searchBooksRequestToURLSearchParams(query)

  return await withPreloadCache(`/api/search?${sp.toString()}`, () =>
    httpClient
      .get('/api/search', {
        searchParams: sp,
      })
      .then((r) => r.json()),
  )
}

export function getPreloadedBookSearchResult(query: SearchBooksRequest) {
  return getPreloadedData<SearchBooksResponse>(
    `/api/search?${searchBooksRequestToURLSearchParams(query).toString()}`,
  )
}

export type BookExtremes = {
  words: {
    min: number
    max: number
  }
  chapters: {
    min: number
    max: number
  }
  wordsPerChapter: {
    min: number
    max: number
  }
  favorites: {
    min: number
    max: number
  }
}

export function httpGetBookExtremes(): Promise<BookExtremes> {
  return withPreloadCache('/api/search/book-extremes', () =>
    httpClient.get('/api/search/book-extremes').then((r) => r.json()),
  )
}

export function httpTagsGetByIds(ids: string[]): Promise<DefinedTagDto[]> {
  const q = stringArrayToQueryParameterValue(ids)
  const searchParams = q ? new URLSearchParams({ q }) : undefined
  return withPreloadCache(
    `/api/tags/lookup` + (searchParams ? `?${searchParams.toString()}` : ''),
    () => httpClient.get('/api/tags/lookup', { searchParams }).then((r) => r.json()),
  )
}
