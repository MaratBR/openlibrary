import { QueryClient, useQuery } from '@tanstack/react-query'
import {
  getPreloadedData,
  httpClient,
  stringArrayToQueryParameterValue,
  withPreloadCache,
} from '../../common/api'
import { z } from 'zod'
import { useNotificationsSlot } from '../../notifications/state'
import { GenericNotification } from '../../notifications'
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
