import { DEFAULT_SEARCH_DEBOUNCE } from '@/config'
import { httpClient } from '@/http-client'
import type { DefinedTagDto } from '@/backend-types'
import { useQuery } from '@tanstack/react-query'
import { useEffect, useState } from 'react'
type UserDto = {
  id: string
  name: string
  avatar: string
}

type NumberRange = { min: number | null; max: number | null }

export type DetailedBookSearchQuery = {
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
    pageSize: 10,
  }
}

export function getQueryParams(query: DetailedBookSearchQuery): URLSearchParams {
  const params = new URLSearchParams()

  if (query.words.min !== null) params.set('w.min', query.words.min.toString())
  if (query.words.max !== null) params.set('w.max', query.words.max.toString())
  if (query.chapters.min !== null) params.set('c.min', query.chapters.min.toString())
  if (query.chapters.max !== null) params.set('c.max', query.chapters.max.toString())
  if (query.wordsPerChapter.min !== null) params.set('wc.min', query.wordsPerChapter.min.toString())
  if (query.wordsPerChapter.max !== null) params.set('wc.max', query.wordsPerChapter.max.toString())
  if (query.includeTags.length > 0) params.set('it', query.includeTags.map((x) => x.id).join(','))
  if (query.excludeTags.length > 0) params.set('et', query.excludeTags.map((x) => x.id).join(','))
  if (query.includeUsers.length > 0) params.set('iu', query.includeUsers.map((x) => x.id).join(','))
  if (query.excludeUsers.length > 0) params.set('eu', query.excludeUsers.map((x) => x.id).join(','))

  if (query.page > 1) params.set('page', query.page.toString())
  if (query.pageSize !== 20) params.set('pageSize', query.pageSize.toString())

  return params
}

export async function searchTags(query: string): Promise<DefinedTagDto[]> {
  const response = await httpClient.get('/_api/tags', { searchParams: { q: query } })
  return response.json<DefinedTagDto[]>()
}

export type { DefinedTagDto, TagsCategory } from '@/backend-types'

export type TagsSearchOptions = {
  query: string
  fetchDefault?: boolean
  debounceTimeout?: number
  enabled?: boolean
}

export function useTagsSearch({
  query,
  fetchDefault = false,
  debounceTimeout = DEFAULT_SEARCH_DEBOUNCE,
  enabled = true,
}: TagsSearchOptions) {
  query = query.trim().toLowerCase()

  const [realQuery, setRealQuery] = useState(query)

  useEffect(() => {
    const t = setTimeout(() => {
      setRealQuery(query)
    }, debounceTimeout)

    return () => clearInterval(t)
  }, [query, debounceTimeout])

  const queryObject = useQuery({
    queryKey: ['tags-search', realQuery],
    enabled: enabled && (realQuery.length > 0 || fetchDefault),
    queryFn: () => searchTags(realQuery),
  })

  return queryObject
}
