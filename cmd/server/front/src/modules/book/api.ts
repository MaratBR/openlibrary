import { useQuery } from '@tanstack/react-query'
import { httpClient } from '../common/api'
import { z } from 'zod'
import { useNotificationsSlot } from '../notifications/state'
import { GenericNotification } from '../notifications'

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
  description: z.string(),
  isAdult: z.boolean(),
  isSpoiler: z.boolean(),
  category: tagCategorySchema,
})

export type DefinedTagDto = z.infer<typeof definedTagDtoSchema>

export type SearchTagsResponse = {
  query: string
  tags: DefinedTagDto[]
}

export function httpTagsSearch(q: string): Promise<SearchTagsResponse> {
  return httpClient.get('/api/tags/search', { searchParams: { q } }).then((r) => r.json())
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
  notifications: GenericNotification[]
}

export type GetBookResponse = BookDetailsDto

export function httpGetBook(id: string): Promise<GetBookResponse> {
  return httpClient.get(`/api/books/${id}`).then((r) => r.json())
}

export function useBookQuery(bookId: string | undefined) {
  const setNotifications = useNotificationsSlot()

  return useQuery({
    queryKey: ['book', bookId],
    enabled: !!bookId,
    queryFn: () =>
      httpGetBook(bookId!).then((resp) => {
        if (resp.notifications) setNotifications(resp.notifications)
        return resp
      }),
    staleTime: 0,
    gcTime: 60000,
  })
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
    staleTime: 0,
    gcTime: 0,
  })
}
