import { censorModeSchema } from '@/modules/account/api'
import { AgeRating, BookCollectionDto, DefinedTagDto } from '@/modules/book/api'
import { httpClient, withPreloadCache } from '@/modules/common/api'
import { QueryClient, useQuery } from '@tanstack/react-query'
import { z } from 'zod'

export type UserDetailsDto = {
  id: string
  name: string
  email: string
  avatar: {
    lg: string
    md: string
  }
  favorites: number
  following: number
  followers: number
  booksTotal: number
  joinedAt: string
  isBlocked: boolean
  isAdmin: boolean
  hasCustomTheme: boolean
  about: {
    bio: string
    gender: string
  }
  books: AuthorBookDto[]

  hideEmail: boolean
  hideStats: boolean
  hideFavorites: boolean

  isFollowing: boolean
}

export type AuthorBookDto = {
  id: string
  name: string
  createdAt: string
  ageRating: AgeRating
  tags: DefinedTagDto[]
  words: number
  wordsPerChapter: number
  chapters: number
  collections: BookCollectionDto[]
  isPinned: boolean
}

export function httpGetUser(id: string): Promise<UserDetailsDto> {
  return withPreloadCache(`/api/users/${id}`, () =>
    httpClient.get(`/api/users/${id}`).then((r) => r.json()),
  )
}

export const userRoleSchema = z.enum(['user', 'admin', 'moderator', 'system'])

export type UserRole = z.infer<typeof userRoleSchema>

export type SelfUserDto = z.infer<typeof selfUserDtoSchema>

export const selfUserDtoSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  role: userRoleSchema,
  avatar: z.object({
    lg: z.string(),
    md: z.string(),
  }),
  joinedAt: z.string(),
  isBlocked: z.boolean(),
  preferredTheme: z.string(),
  showAdultContent: z.boolean(),
  bookCensoredTags: z.array(z.string()),
  bookCensoringMode: censorModeSchema,
})

export type WhoamiResponse = {
  user: SelfUserDto | null
}

export function httpWhoami(): Promise<WhoamiResponse> {
  return withPreloadCache('/api/users/whoami', () =>
    httpClient.get('/api/users/whoami').then((r) => r.json()),
  )
}

export function preloadWhoami(queryClient: QueryClient) {
  if (!__server__.clientPreload) return
  queryClient.prefetchQuery({
    queryKey: ['whoami'],
    queryFn: () => httpWhoami(),
    staleTime: 20000,
    gcTime: Infinity,
  })
}

export function useWhoamiQuery() {
  return useQuery({
    queryKey: ['whoami'],
    queryFn: () => httpWhoami(),
    staleTime: 10000,
    gcTime: 10000,
  })
}

export function httpFollowUser(userId: string) {
  return httpClient.post('/api/users/follow', { searchParams: { userId } })
}

export function httpUnfollowUser(userId: string) {
  return httpClient.post('/api/users/unfollow', { searchParams: { userId } })
}
