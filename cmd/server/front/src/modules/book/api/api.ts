import { QueryClient, useQuery } from '@tanstack/react-query'
import {
  getPreloadedData,
  httpClient,
  stringArrayToQueryParameterValue,
  withPreloadCache,
} from '../../common/api'
import { z } from 'zod'
import { useNotificationsSlot } from '../../notifications/state'
import { genericNotificationSchema } from '../../notifications'
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

export const bookCollectionDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  position: z.number(),
  size: z.number(),
})

export type BookCollectionDto = z.infer<typeof bookCollectionDtoSchema>

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

export const ageRatingSchema = z.enum(['?', 'G', 'PG', 'PG-13', 'R', 'NC-17'])

export type AgeRating = z.infer<typeof ageRatingSchema>

export const AGE_RATINGS_LIST: AgeRating[] = ['?', 'G', 'PG', 'PG-13', 'R', 'NC-17']

export const bookChapterDtoSchema = z.object({
  id: z.string(),
  order: z.number().min(0).int(),
  name: z.string(),
  words: z.number(),
  createdAt: z.string(),
  summary: z.string(),
})

export type BookChapterDto = z.infer<typeof bookChapterDtoSchema>

export const readingListStatusSchema = z.enum(['dnf', 'paused', 'read', 'reading', 'want_to_read'])

export type ReadingListStatus = z.infer<typeof readingListStatusSchema>

export const readingListDtoSchema = z.object({
  lastUpdatedAt: z.string(),
  chapterId: z.string().nullable(),
  status: readingListStatusSchema,
})

export type ReadingListDto = z.infer<typeof readingListDtoSchema>

export const bookDetailsDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  ageRating: ageRatingSchema,
  tags: z.array(definedTagDtoSchema),
  words: z.number(),
  wordsPerChapter: z.number(),
  collections: z.array(bookCollectionDtoSchema),
  chapters: z.array(bookChapterDtoSchema),
  createdAt: z.string(),
  author: z.object({
    id: z.string(),
    name: z.string(),
  }),
  permissions: z.object({
    canEdit: z.boolean(),
  }),
  summary: z.string(),
  favorites: z.number(),
  isFavorite: z.boolean(),
  notifications: z.array(genericNotificationSchema).optional(),
  cover: z.string(),
  rating: z.number().nullable(),
  readingList: readingListDtoSchema.nullable(),
})

export type BookDetailsDto = z.infer<typeof bookDetailsDtoSchema>

export async function httpGetBook(id: string): Promise<BookDetailsDto> {
  const result = await withPreloadCache(`/api/books/${id}`, () =>
    httpClient.get(`/api/books/${id}`).then((r) => r.json()),
  )
  return bookDetailsDtoSchema.parse(result)
}

export function getPreloadedBookResult(id: string) {
  return getPreloadedData<BookDetailsDto>(`/api/books/${id}`)
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

export function httpFavoriteBook(id: string, isFavorite: boolean): Promise<BookDetailsDto> {
  return httpClient
    .post(`/api/favorite`, { searchParams: { bookId: id, isFavorite } })
    .then((r) => r.json())
    .then((r) => bookDetailsDtoSchema.parse(r))
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
