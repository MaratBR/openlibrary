import type { DefinedTagDto } from '@/backend-types'

type UserDto = {
  id: string
  name: string
  avatar: string
}

type NumberRange = { min: number | null; max: number | null }

export type BookSearchSort =
  | ''
  | 'chapters'
  | 'words'
  | 'words-per-chapter'
  | 'last-updated'
  | 'created-at'
  | 'reviews'
  | 'readers'
  | 'weighted-rating'
  | 'random'

export type DetailedBookSearchQuery = {
  rating: string
  createdAt: { min: string; max: string }
  updatedAt: { min: string; max: string }
  views: NumberRange
  readers: NumberRange
  reviews: NumberRange
  popularity: NumberRange
  weightedRating: NumberRange
  sort: BookSearchSort
  words: NumberRange
  chapters: NumberRange
  wordsPerChapter: NumberRange
  includeTags: DefinedTagDto[]
  excludeTags: DefinedTagDto[]
  includeUsers: UserDto[]
  excludeUsers: UserDto[]
  includeBanned: boolean
  includeHidden: boolean
  includeEmpty: boolean
  page: number
  pageSize: number
}

export function getDefaultDetailedBookSearchQuery(): DetailedBookSearchQuery {
  return {
    rating: '',
    createdAt: { min: '', max: '' },
    updatedAt: { min: '', max: '' },
    views: { min: null, max: null },
    readers: { min: null, max: null },
    reviews: { min: null, max: null },
    popularity: { min: null, max: null },
    weightedRating: { min: null, max: null },
    sort: '',
    words: { min: null, max: null },
    chapters: { min: null, max: null },
    wordsPerChapter: { min: null, max: null },
    includeTags: [],
    excludeTags: [],
    includeUsers: [],
    excludeUsers: [],
    includeBanned: false,
    includeHidden: false,
    includeEmpty: false,
    page: 1,
    pageSize: 20,
  }
}

export function getQueryParams(query: DetailedBookSearchQuery): URLSearchParams {
  const params = new URLSearchParams()

  if (query.rating) params.set('rating', query.rating)
  for (const key of [
    'views',
    'readers',
    'reviews',
    'popularity',
    'weightedRating',
    'createdAt',
    'updatedAt',
  ] as const) {
    for (const bound of ['min', 'max'] as const) {
      const value = query[key][bound]
      if (value !== null && value !== '') params.set(`${key}.${bound}`, String(value))
    }
  }
  if (query.sort) params.set('sort', query.sort)

  if (query.words.min !== null) params.set('w.min', query.words.min.toString())
  if (query.words.max !== null) params.set('w.max', query.words.max.toString())
  if (query.chapters.min !== null) params.set('c.min', query.chapters.min.toString())
  if (query.chapters.max !== null) params.set('c.max', query.chapters.max.toString())
  if (query.wordsPerChapter.min !== null) params.set('wc.min', query.wordsPerChapter.min.toString())
  if (query.wordsPerChapter.max !== null) params.set('wc.max', query.wordsPerChapter.max.toString())
  if (query.includeTags.length > 0)
    params.set('it', query.includeTags.map((tag) => tag.id).join(','))
  if (query.excludeTags.length > 0)
    params.set('et', query.excludeTags.map((tag) => tag.id).join(','))
  if (query.includeUsers.length > 0)
    params.set('iu', query.includeUsers.map((user) => user.id).join(','))
  if (query.excludeUsers.length > 0)
    params.set('eu', query.excludeUsers.map((user) => user.id).join(','))

  if (query.page > 1) params.set('p', query.page.toString())
  if (query.pageSize !== 20) params.set('ps', query.pageSize.toString())

  return params
}
